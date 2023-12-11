package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func generateHelpText() string {
	helpText := "I understand: "
	for key := range CommandMap {
		helpText += fmt.Sprintf("/%s ", key)
	}
	return helpText
}

func GetSessions() (JellyJSON []JellySession, err error) {
	url := jellyfinAddress + "/Sessions?api_key=" + jellyfinApiKey
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	log.Printf("API request to %s completed with status code: %d", jellyfinAddress, resp.StatusCode)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &JellyJSON)
	if err != nil {
		return nil, err
	}
	return JellyJSON, nil
}

func RemoveActiveSession(ID string, activeStreams []ActiveSession) []ActiveSession {
	var tmpActiveStreams []ActiveSession
	for i, session := range activeStreams {
		if session.MediaID == ID {
			tmpActiveStreams = append(activeStreams[:i], activeStreams[i+1:]...)
		}
	}
	return tmpActiveStreams
}

func AppendMessage(msg string, additive string) string {
	if msg != "" {
		msg += "\n" + additive
	} else {
		msg = additive
	}
	return msg
}
