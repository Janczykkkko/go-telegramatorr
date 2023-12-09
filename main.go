package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	jellyfinAddress string
	jellyfinApiKey  string
	telegramApiKey  string
	botMonitor      bool
	bot             *tgbotapi.BotAPI
	updates         tgbotapi.UpdatesChannel
	u               tgbotapi.UpdateConfig
)

func main() {
	jellyfinAddress = os.Getenv("JELLYFIN_ADDRESS")
	jellyfinApiKey = os.Getenv("JELLYFIN_APIKEY")
	telegramApiKey = os.Getenv("TELEGRAM_APIKEY")
	var err error
	botMonitor, err = strconv.ParseBool(os.Getenv("BOT_MONITOR"))
	if err != nil {
		fmt.Println("BOT_MONITOR variable unspecified or incorrect.")
		log.Fatal(err)
	}
	botInit()
}
