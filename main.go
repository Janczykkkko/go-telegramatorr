package main

import (
	"log"
	"time"
)

var (
	jellyfinAddress    string
	apiKey             string
	chatID             int64
	botToken           string
	silentNotification bool
	pollingInterval    = 30 * time.Second
)

func main() {
	if err := LoadConfig(); err != nil {
		log.Fatal("Error loading configuration:", err)
	}

	SendMessage(chatID, botToken, "hello", silentNotification)
}
