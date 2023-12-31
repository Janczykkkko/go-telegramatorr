package bot

import (
	"fmt"
	"log"
	"math"
	"time"

	gatherers "github.com/Janczykkkko/jellyplexgatherer"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/xid"
)

type ActiveSession struct {
	UserName     string
	Name         string
	Bitrate      string
	PlayMethod   string
	SubStream    string
	DeviceName   string
	Service      string
	StartTime    time.Time
	EndTime      time.Time //used by db
	EndTimeStr   string    //used by db
	StartTimeStr string    //used by db
	Duration     string    //used by db
	ID           string
}

var (
	sessionStore []ActiveSession
	persist      bool = true
)

func botMonitorAndInform(bot *tgbotapi.BotAPI, chatID int64, dblocation string) {
	if !DbExistsAndWorks(dblocation) {
		err := CreateDb(dblocation)
		if err != nil {
			log.Printf("Failed to create db, disabling reporting: %s", err)
			persist = false
		}
	}
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		sessions, errors := gatherers.GetAllSessions(jellyfinAddress, jellyfinApiKey, plexAddress, plexApiKey)
		if errors != "" {
			log.Printf("Encountered errors. Some data might be missing:\n%s", errors)
		} else {
			log.Println("All stream data collected succesfully")
		}

		processSessions(sessions, chatID, bot)

		if len(sessionStore) > 0 {
			fmt.Println("Active Sessions monitored:")
			for _, s := range sessionStore {
				fmt.Printf(
					"%s (%s) - %s running on %s for %.0f minutes\n",
					s.UserName,
					s.DeviceName,
					s.Name,
					s.Service,
					math.Round(time.Since(s.StartTime).Seconds())/60)
			}
		} else {
			fmt.Println("No sessions are monitored.")
		}
	}
}

func processSessions(currentSessions []gatherers.SessionData, chatID int64, bot *tgbotapi.BotAPI) {
	//get new sessions
	for _, c := range currentSessions {
		var count int
		for _, s := range sessionStore {
			if c.UserName == s.UserName &&
				c.Name == s.Name &&
				c.DeviceName == s.DeviceName {
				//update if changed
				s.SubStream = c.SubStream
				s.Bitrate = c.Bitrate
				s.PlayMethod = c.PlayMethod
				break //found monitored session in currently active streams
			}
			count++
		}
		if count == len(sessionStore) {
			id := xid.New()
			sessionStore = append(sessionStore, ActiveSession{
				UserName:   c.UserName,
				Name:       c.Name,
				Bitrate:    c.Bitrate,
				PlayMethod: c.PlayMethod,
				SubStream:  c.SubStream,
				DeviceName: c.DeviceName,
				Service:    c.Service,
				StartTime:  time.Now(),
				ID:         id.String(),
			})
			log.Printf("Stream registered: %s - %s on %s", c.UserName, c.Name, c.Service)
		}
	}
	//check for stopped and lingering sessions
	for _, s := range sessionStore {
		var count int
		for _, c := range currentSessions {
			if c.UserName == s.UserName &&
				c.Name == s.Name &&
				c.DeviceName == s.DeviceName {
				break
			}
			count++
		}
		if count == len(currentSessions) {
			//session ended, persist in db for reporting
			if persist {
				err := InsertDataToDb(s, time.Now(), dblocation)
				if err != nil {
					log.Println("Error persisting session in db", err)
				}
			}
			removeSession(s.ID)
			log.Printf("Deregistered finished stream: %s - %s on %s after %.0f minutes\n", s.UserName, s.Name, s.Service, math.Round(time.Since(s.StartTime).Seconds())/60)
			msgStr := fmt.Sprintf(
				"User %s (%s) was playing %s on %s for %.0f minutes\nmethod: %s\nbitrate: %s Mbps\nsubs: %s",
				s.UserName,
				s.DeviceName,
				s.Name,
				s.Service,
				math.Round(time.Since(s.StartTime).Seconds())/60,
				s.PlayMethod,
				s.Bitrate,
				s.SubStream)
			msg := tgbotapi.NewMessage(chatID, msgStr)
			msg.DisableNotification = true
			if _, err := bot.Send(msg); err != nil {
				log.Printf("Error sending message: %s", err)
			}
		}
		if math.Round(time.Since(s.StartTime).Seconds())/60 > 180 {
			log.Printf("Deregistered lingering stream: %s - %s on %s after %.0f minutes\n",
				s.UserName,
				s.Name,
				s.Service,
				math.Round(time.Since(s.StartTime).Seconds())/60)
			removeSession(s.ID)
		}
	}
}

func removeSession(ID string) {
	var tmpActiveStreams []ActiveSession
	for i, session := range sessionStore {
		if session.ID == ID {
			tmpActiveStreams = append(sessionStore[:i], sessionStore[i+1:]...)
		}
	}
	sessionStore = tmpActiveStreams
}
