package bot

import (
	"fmt"
	"log"
	"math"
	"telegramatorr/gatherers"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/xid"
)

type ActiveSession struct {
	UserName   string
	Name       string
	Bitrate    string
	PlayMethod string
	SubStream  string
	DeviceName string
	Service    string
	StartTime  time.Time
	ID         string
}

var sessionStore []ActiveSession

func botMonitorAndInform(bot *tgbotapi.BotAPI, chatID int64) {
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
		if count == len(sessionStore) {
			removeSession(s.ID)
			log.Printf("Deregistered finished stream: %s - %s on %s after %.0f minutes\n", s.UserName, s.Name, s.Service, math.Round(time.Since(s.StartTime).Seconds())/60)
			msgStr := fmt.Sprintf(
				"User %s (%s) was playing %s on %s for %.0f minutes\nmethod: %s\nbitrate: %s\nsubs: %s",
				s.UserName,
				s.DeviceName,
				s.Name,
				s.Service,
				math.Round(time.Since(s.StartTime).Seconds())/60,
				s.PlayMethod,
				s.Bitrate,
				s.SubStream)
			msg := tgbotapi.NewMessage(chatID, msgStr)
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
