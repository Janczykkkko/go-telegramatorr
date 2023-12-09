package main

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	msg tgbotapi.MessageConfig
)

func botInit() {
	bot, err := tgbotapi.NewBotAPI(telegramApiKey)
	if err != nil {
		log.Fatal("Error connecting to bot, is the apikey correct?", err)
	}
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	// watch for commands
	go botWatch(u, bot)
	// start monitoring and updating on user playback if enabled
	if botMonitor {
		go botMonitorAndInform(u, bot)
	}
	// all goroutines are meant to run idefinitely
	select {}
}

func botWatch(u tgbotapi.UpdateConfig, bot *tgbotapi.BotAPI) {
	updates := bot.GetUpdatesChan(u)

	for update := range updates {

		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "")
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
			log.Panic(err)
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

func botMonitorAndInform(u tgbotapi.UpdateConfig, bot *tgbotapi.BotAPI) {
	var activeSessions []ActiveSession
	for {
		jellyJSON, err := GetSessions()
		if err != nil {
			errMsg := "Error: " + err.Error()
			fmt.Println(errMsg)
		}
		for _, obj := range jellyJSON {
			if len(obj.NowPlayingQueueFullItems) > 0 ||
				len(obj.FullNowPlayingItem.Container) > 0 &&
					obj.PlayState.PlayMethod != "" &&
					!obj.PlayState.IsPaused {

			} else {
				continue
			}
		}
		// eval if notifying
	}
}
