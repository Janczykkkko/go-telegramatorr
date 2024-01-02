package bot

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	gatherers "github.com/Janczykkkko/jellyplexgatherer"
)

func generateHelpText() string {
	helpText := "I understand:\n"
	for key := range CommandMap {
		helpText += fmt.Sprintf("/%s\n", key)
	}
	return helpText
}

func checkAssignEnv(JellyfinAddress, JellyfinApiKey, PlexAddress, PlexApiKey, TelegramChatId, enableReport string) (monitor bool, chatID int64, reports bool) {
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

	reportBool, err := strconv.ParseBool(enableReport)
	if err != nil {
		fmt.Println("Reports variable not assigned or wrong format, disabling...", err)
		reports = false
	}
	reports = reportBool

	return monitor, chatID, reports
}

func GetAllSessionsStr(jellyfinAddress, jellyfinApiKey, plexAddress, plexApiKey string) string {
	response := []string{"Here's a report from your player(s):"}
	sessions, errors := gatherers.GetAllSessions(jellyfinAddress, jellyfinApiKey, plexAddress, plexApiKey)
	if errors != "" {
		return errors
	}
	if len(sessions) == 0 {
		return "Nothing is playing."
	}
	for _, session := range sessions {
		response = append(response, fmt.Sprintf(
			"%s is playing(%s) on %s: %s\nBitrate: %s Mbps\nDevice: %s\nSubs: %s",
			session.UserName,
			session.PlayMethod,
			session.Service,
			session.Name,
			session.Bitrate,
			session.DeviceName,
			session.SubStream,
		))
	}
	return strings.Join(response, "\n\n")
}

func FormatTimeToNiceString(t time.Time, showDays bool) string {
	// Format the time to display "Mon-Sun HH:mm"
	formattedTime := t.Format("15:04")
	if showDays {
		formattedTime = t.Format("Mon 15:04")
	}
	return formattedTime
}
