package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	jellyfinAddress string
	jellyfinApiKey  string
	plexAddress     string
	plexApiKey      string
)

func Init(jellyfinAddress, jellyfinApiKey, plexAddress, plexApiKey, telegramApiKey, telegramChatId string) {

	bot, err := tgbotapi.NewBotAPI(telegramApiKey)
	if err != nil {
		log.Fatal("Error connecting to bot, is the apikey correct?", err)
	}

	monitor, chatID := checkAssignEnv(jellyfinAddress, jellyfinApiKey, plexAddress, plexApiKey, telegramChatId)

	bot.Debug = true
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	log.Printf("Authorized on account %s", bot.Self.UserName)

	go botWatch(u, bot)

	if monitor {
		go botMonitorAndInform(bot, chatID)
	}

	select {}
}
