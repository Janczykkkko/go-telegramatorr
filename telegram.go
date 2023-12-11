package main

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func botInit() {
	bot, err := tgbotapi.NewBotAPI(telegramApiKey)
	if err != nil {
		log.Fatal("Error connecting to bot, is the apikey correct?", err)
	}

	chatID, err := strconv.ParseInt(telegramChatId, 10, 64)
	if err != nil {
		log.Fatal("Error parsing chat id", err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	// watch for commands
	go botWatch(u, bot)
	// start monitoring and updating on user playback if enabled
	if botMonitor {
		go botMonitorAndInform(bot, chatID)
	}
	// all goroutines are meant to run idefinitely
	select {}
}

func botWatch(u tgbotapi.UpdateConfig, bot *tgbotapi.BotAPI) {
	var msg tgbotapi.MessageConfig
	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		} else if update.Message.Command() == "help" {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, generateHelpText())
		} else if update.Message.Command() != "" {
			msg = botObey(update)
		} else {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, `
			Message not a command! Good luck talking to yourself! Try /help for list of available comands :)
			`)
		}
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Error sending message: %s", err)
		}
	}
}

func botObey(update tgbotapi.Update) (msg tgbotapi.MessageConfig) {
	command := update.Message.Command()
	reply, found := CommandMap[command]
	if found {
		replyStr := reply()
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, replyStr)
	} else {
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "I don't know that command")
	}
	return msg
}

func botMonitorAndInform(bot *tgbotapi.BotAPI, chatID int64) {
	var msg tgbotapi.MessageConfig
	var activeStreams []ActiveSession
	var msgStr string
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		msgStr = ""
		var skipOnErr bool = false
		jellyJSON, err := GetSessions()
		if err != nil {
			errMsg := "Error: " + err.Error()
			fmt.Println(errMsg)
			skipOnErr = true
			continue
		} else {
			fmt.Println("Stream data collected succesfully")
		}
		if !skipOnErr {
			var filteredJellyJSON []JellySession
			//eliminate non-playback sessions
			for _, obj := range jellyJSON {
				if obj.PlayState.PlayMethod != "" &&
					!obj.PlayState.IsPaused {
					filteredJellyJSON = append(filteredJellyJSON, obj)
				} else {
					continue
				}
			}

			//Catch new / continuing streams
			for _, obj := range filteredJellyJSON {
				var count int
				if len(activeStreams) > 0 {
					for _, data := range activeStreams {
						if data.UserName == obj.UserName && data.MediaID == obj.PlayState.MediaSourceID {
							break
						} else {
							count++
						}
					}
					if count == len(activeStreams) {
						activeStreams = append(activeStreams, ActiveSession{
							UserName:  obj.UserName,
							MediaID:   obj.PlayState.MediaSourceID,
							MediaName: obj.NowPlayingItem.Name,
							StartTime: time.Now(),
						})
						fmt.Println("Stream registered: ", obj.UserName, "-", obj.NowPlayingItem.Name)
					}
				} else {
					activeStreams = append(activeStreams, ActiveSession{
						UserName:  obj.UserName,
						MediaID:   obj.PlayState.MediaSourceID,
						MediaName: obj.NowPlayingItem.Name,
						StartTime: time.Now(),
					})
					fmt.Println("Stream registered: ", obj.UserName, "-", obj.NowPlayingItem.Name)
				}
			}
			//Catch stopped streams
			var tmpActiveStreams []ActiveSession
			if len(activeStreams) > 0 {
				for _, data := range activeStreams {
					var count int
					for _, obj := range filteredJellyJSON {
						if data.UserName == obj.UserName && data.MediaID == obj.PlayState.MediaSourceID && data.MediaName == obj.NowPlayingItem.Name {
							tmpActiveStreams = activeStreams
							break
						} else {
							count++
						}
					}
					if count == len(filteredJellyJSON) {
						tmpActiveStreams = RemoveActiveSession(data.MediaID, activeStreams)
						fmt.Printf("Deregistered finished stream: %s (%s) - %s after %.0f minutes\n", data.UserName, data.MediaName, data.MediaID, math.Round(time.Since(data.StartTime).Seconds())/60)
						additive := fmt.Sprintf("User %s was playing %s for %.0f minutes - finished.", data.UserName, data.MediaName, math.Round(time.Since(data.StartTime).Seconds())/60)
						msgStr = AppendMessage(msgStr, additive)
					}
				}
				activeStreams = tmpActiveStreams
			}
		}
		if len(activeStreams) > 0 {
			fmt.Println("Active Sessions monitored:")
			for _, data := range activeStreams {
				fmt.Printf("%s (%s) - %s running for %.0f minutes\n", data.UserName, data.MediaName, data.MediaID, math.Round(time.Since(data.StartTime).Seconds())/60)
			}
		} else {
			fmt.Println("No sessions are monitored.")
		}
		// account for bugs - remove lingering sessions
		var tmpActiveStreams []ActiveSession
		for _, data := range activeStreams {
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
