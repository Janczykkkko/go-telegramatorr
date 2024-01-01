package main

import (
	"os"
	"telegramatorr/bot"
)

func main() {
	//init bot and its services
	bot.Init(
		os.Getenv("JELLYFIN_ADDRESS"),
		os.Getenv("JELLYFIN_APIKEY"),
		os.Getenv("PLEX_ADDRESS"),
		os.Getenv("PLEX_APIKEY"),
		os.Getenv("TELEGRAM_APIKEY"),
		os.Getenv("TELEGRAM_CHATID"),
		os.Getenv("ENABLE_REPORTS"),
	)
}
