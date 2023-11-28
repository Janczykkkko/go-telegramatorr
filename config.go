package main

import (
	"errors"
	"os"
	"strconv"
	"time"
)

var (
	chatIDStr             string
	silentNotificationStr string
)

func LoadConfig() error {
	jellyfinAddress = os.Getenv("JELLYFIN_ADDRESS")
	apiKey = os.Getenv("API_KEY")
	chatIDStr = os.Getenv("CHAT_ID")
	botToken = os.Getenv("BOT_TOKEN")
	silentNotificationStr = os.Getenv("SILENT_NOTIFICATION")

	if chatIDStr == "" || botToken == "" {
		return errors.New("missing required envvars")
	}

	if err := ParseSilentNotification(); err != nil {
		return err
	}

	if err := ParseChatID(); err != nil {
		return err
	}

	if err := ParsePollingInterval(); err != nil {
		return err
	}

	return nil
}

func ParseSilentNotification() error {
	if silentNotificationStr != "" {
		parsedSilentNotification, err := strconv.ParseBool(silentNotificationStr)
		if err != nil {
			return err
		}
		silentNotification = parsedSilentNotification
	}
	return nil
}

func ParseChatID() error {
	if chatIDStr != "" {
		parsedChatID, err := strconv.ParseInt(chatIDStr, 10, 64)
		if err != nil {
			return err
		}
		chatID = parsedChatID
	}
	return nil
}

func ParsePollingInterval() error {
	pollingIntervalStr := os.Getenv("POLLING_INTERVAL")
	if pollingIntervalStr != "" {
		interval, err := strconv.Atoi(pollingIntervalStr)
		if err != nil {
			return err
		}
		pollingInterval = time.Duration(interval) * time.Second
	}
	return nil
}
