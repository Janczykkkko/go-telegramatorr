package main

import (
	"log"
	"os"
	"telegramatorr/bot"
)

// Config struct holds all the environment variables
var env = map[string]string{
	"jellyfinAddress": os.Getenv("JELLYFIN_ADDRESS"),
	"jellyfinAPIKey":  os.Getenv("JELLYFIN_APIKEY"),
	"plexAddress":     os.Getenv("PLEX_ADDRESS"),
	"plexAPIKey":      os.Getenv("PLEX_APIKEY"),
	"telegramAPIKey":  os.Getenv("TELEGRAM_APIKEY"),
	"telegramChatID":  os.Getenv("TELEGRAM_CHATID"),
}

func CheckConfig() {
	for key, val := range env {
		if val == "" {
			log.Printf("%s variable not provided, disabling related features...", key)
		}
	}
}

func main() {
	CheckConfig()
	bot.Init(env)
}
