package main

import (
	"os"
	"telegramatorr/bot"
)

// Config struct holds all the environment variables
var (
	jellyfinAddress string = os.Getenv("JELLYFIN_ADDRESS")
	jellyfinApiKey  string = os.Getenv("JELLYFIN_APIKEY")
	plexAddress     string = os.Getenv("PLEX_ADDRESS")
	plexApiKey      string = os.Getenv("PLEX_APIKEY")
	telegramApiKey  string = os.Getenv("TELEGRAM_APIKEY")
	telegramChatId  string = os.Getenv("TELEGRAM_CHATID")
)

func main() {
	bot.Init(
		jellyfinAddress,
		jellyfinApiKey,
		plexAddress,
		plexApiKey,
		telegramApiKey,
		telegramChatId,
	)
}
