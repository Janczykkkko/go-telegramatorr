package bot

import (
	"fmt"
	"log"
	"strconv"
)

func generateHelpText() string {
	helpText := "I understand: "
	for key := range CommandMap {
		helpText += fmt.Sprintf("/%s ", key)
	}
	return helpText
}

func checkAssignEnv(JellyfinAddress, JellyfinApiKey, PlexAddress, PlexApiKey, TelegramChatId string) (monitor bool, chatID int64) {
	//glorified printer
	sources := 2
	monitor = true
	if JellyfinAddress == "" || JellyfinApiKey == "" {
		log.Println("Jellyfin env vars not specified, disabling gatherer...")
		sources--
	}
	if PlexAddress == "" || PlexApiKey == "" {
		log.Println("Plex env vars not specified, disabling gatherer...")
		sources--
	}
	chatID, err := strconv.ParseInt(TelegramChatId, 10, 64)
	if err != nil {
		log.Printf("Telegram chat id not specified or wrong format, disabling monitor...")
		monitor = false
	}
	jellyfinAddress = JellyfinAddress
	jellyfinApiKey = JellyfinApiKey
	plexAddress = PlexAddress
	plexApiKey = PlexApiKey
	if sources == 0 {
		log.Fatalln("No sources enabled! Nothing will work bruh...")
	}
	return monitor, chatID
}
