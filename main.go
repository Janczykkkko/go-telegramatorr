package main

import (
	"os"
)

var (
	jellyfinAddress string
	jellyfinApiKey  string
	telegramApiKey  string
)

func main() {
	jellyfinAddress = os.Getenv("JELLYFIN_ADDRESS")
	jellyfinApiKey = os.Getenv("JELLYFIN_APIKEY")
	telegramApiKey = os.Getenv("TELEGRAM_APIKEY")

	botInit()
}
