package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Message represents the structure of the message to be sent
type Message struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

// SendMessage sends a message via Telegram
func SendMessage(chatId int64, botToken string, text string, silent bool) {
	message := Message{
		ChatID: chatId,
		Text:   text,
	}

	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)

	if silent {
		apiURL += "?disable_notification=true" // Add the parameter to send the message silently
	}

	requestBody, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error encoding message: %s\nMessage content: %+v\n", err, message.Text)
		return
	}

	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Printf("Error sending message: %s\nMessage content: %+v\n", err, message.Text)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Unexpected status code: %d\nMessage content: %+v\n", resp.StatusCode, message.Text)
		return
	}

	log.Printf("Message: \"%s\" sent successfully!", message.Text)
}
