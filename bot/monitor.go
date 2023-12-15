package bot

import (
	"fmt"
	"log"
	"math"
	"telegramatorr/gatherers"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var sessionStore []gatherers.SessionData

func removeInactiveSession(ID string, activeStreams []gatherers.SessionData) []ActiveSession {
	var tmpActiveStreams []ActiveSession
	for i, session := range activeStreams {
		if session.MediaID == ID {
			tmpActiveStreams = append(activeStreams[:i], activeStreams[i+1:]...)
		}
	}
	return tmpActiveStreams
}

func botMonitorAndInform(bot *tgbotapi.BotAPI, chatID int64) {
	var msg tgbotapi.MessageConfig
	var msgStr string

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {

		msgStr = ""

		sessions, errors := gatherers.GetAllSessions(jellyfinAddress, jellyfinApiKey, plexAddress, plexApiKey)
		if errors != "" {
			log.Printf("Encountered some errors. Some data might be missing:\n%s", errors)
		} else {
			log.Println("Stream data collected succesfully")
		}

		//Catch stopped streams

		if len(sessionStore) > 0 {
			fmt.Println("Active Sessions monitored:")
			for _, data := range sessionStore {
				fmt.Printf("%s (%s) - %s running for %.0f minutes\n", data.UserName, data.MediaName, data.MediaID, math.Round(time.Since(data.StartTime).Seconds())/60)
			}
		} else {
			fmt.Println("No sessions are monitored.")
		}
		// account for bugs - remove lingering sessions
		for _, data := range sessionStore {
			if math.Round(time.Since(data.StartTime).Seconds())/60 > 180 {
				tmpActiveStreams = RemoveActiveSession(data.MediaID, activeStreams)
			} else {
				tmpActiveStreams = activeStreams
			}
		}
		activeStreams = tmpActiveStreams

		if msgStr != "" {
			msg = tgbotapi.NewMessage(chatID, msgStr)
			if _, err := bot.Send(msg); err != nil {
				log.Printf("Error sending message: %s", err)
			}
		}

	}
}

func processSessions(sessions []gatherers.SessionData) {
	//get new sessions
	for _, obj := range sessions {
		var count int
		for _, data := range sessions {
			if data.UserName == obj.UserName && data.MediaID == obj.PlayState.MediaSourceID {
				break
			} else {
				count++
			}
		}
		if count == len(sessions) {
			sessionStore = append(sessionStore, gatherers.SessionData{
				UserName:  obj.UserName,
				MediaID:   obj.PlayState.MediaSourceID,
				MediaName: obj.NowPlayingItem.Name,
				StartTime: time.Now(),
			})
			log.Println("Stream registered: ", obj.UserName, "-", obj.NowPlayingItem.Name)
		}
	}
	//check for stopped sessions
	for _, data := range sessionStore {
		var count int
		for _, obj := range sessions {
			if data.UserName == obj.UserName && data.MediaID == obj.PlayState.MediaSourceID && data.MediaName == obj.NowPlayingItem.Name {
				tmpActiveStreams = activeStreams
				break
			} else {
				count++
			}
		}
		if count == len(sessionStore) {
			tmpActiveStreams = removeInactiveSession(data.MediaID, activeStreams)
			fmt.Printf("Deregistered finished stream: %s (%s) - %s after %.0f minutes\n", data.UserName, data.MediaName, data.MediaID, math.Round(time.Since(data.StartTime).Seconds())/60)
			additive := fmt.Sprintf("User %s was playing %s for %.0f minutes - finished.", data.UserName, data.MediaName, math.Round(time.Since(data.StartTime).Seconds())/60)
			msgStr = appendMessage(msgStr, additive)
		}
	}
}
