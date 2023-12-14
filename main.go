package main

import (
	"fmt"
	"os"
	"strconv"
)

var (
	jellyfinAddress string
	jellyfinApiKey  string
	plexAddress     string
	plexApiKey      string
	telegramApiKey  string
	botMonitor      bool
	telegramChatId  string
)

func main() {
	jellyfinAddress = os.Getenv("JELLYFIN_ADDRESS")
	jellyfinApiKey = os.Getenv("JELLYFIN_APIKEY")
	plexAddress = os.Getenv("PLEX_ADDRESS")
	plexApiKey = os.Getenv("PLEX_APIKEY")
	telegramApiKey = os.Getenv("TELEGRAM_APIKEY")
	telegramChatId = os.Getenv("TELEGRAM_CHATID")
	var err error
	botMonitor, err = strconv.ParseBool(os.Getenv("BOT_MONITOR"))
	if err != nil {
		fmt.Println("BOT_MONITOR variable unspecified or incorrect. Disabling botMonitor")
		botMonitor = false
	}
	botInit()
}
