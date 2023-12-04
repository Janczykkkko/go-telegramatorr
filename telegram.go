package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	msg tgbotapi.MessageConfig
)

func botInit() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {

		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "")

		// Irrelevant scenario
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

	// Check if the received message contains a known command
	command := update.Message.Command()
	reply, found := CommandMap[command]
	if found {
		replyStr := reply()
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, replyStr)
	} else {
		// If the command is unknown, send a default reply
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "I don't know that command")
		// Send the default reply message using your Telegram bot API client
		// bot.Send(msg)
	}
	return msg
}
