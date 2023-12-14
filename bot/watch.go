package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Watch(u tgbotapi.UpdateConfig, bot *tgbotapi.BotAPI, CommandMap map[string]CommandFunc) {
	var msg tgbotapi.MessageConfig
	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		} else if update.Message.Command() == "help" {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, GenerateHelpText())
		} else if update.Message.Command() != "" {
			msg = Obey(update, CommandMap)
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
