package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Init(env map[string]string) {
	bot, err := tgbotapi.NewBotAPI(env["telegramApiKey"])
	if err != nil {
		log.Fatal("Error connecting to bot, is the apikey correct?", err)
	}
	/*
		chatID, err := strconv.ParseInt(telegramChatId, 10, 64)
		if err != nil {
			log.Fatal("Error parsing chat id", err)
		}
	*/
	bot.Debug = true
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	log.Printf("Authorized on account %s", bot.Self.UserName)

	//check
}

func Obey(update tgbotapi.Update) (msg tgbotapi.MessageConfig) {
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
