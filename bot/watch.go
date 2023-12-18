package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

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
	msg = tgbotapi.NewMessage(update.Message.Chat.ID, "I don't know that command")
	command := update.Message.Command()
	reply, found := CommandMap[command]
	if found {
		replyStr := reply()
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, replyStr)
	}
	return msg
}
