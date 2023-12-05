package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// GetSessions fetches sessions from Jellyfin
func GetSessions() string {
	var (
		JellyJSON         []JellySession
		genericInfo       string
		sessionStrings    []string
		formattedSessions string
	)
	genericInfo = "Here's an activity report from Jellyfin: \n\n"
	url := jellyfinAddress + "/Sessions?api_key=" + jellyfinApiKey
	resp, err := http.Get(url)
	if err != nil {
		formattedSessions = "Error fetching sessions: " + err.Error()
	}
	defer resp.Body.Close()
	log.Printf("API request to %s completed with status code: %d", jellyfinAddress, resp.StatusCode)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		formattedSessions = "Error fetching sessions: " + err.Error()
	}
	err = json.Unmarshal(body, &JellyJSON)
	if err != nil {
		formattedSessions = "Error fetching sessions: " + err.Error()
	}
	for _, obj := range JellyJSON {
		var sessionString string
		if len(obj.NowPlayingQueueFullItems) > 0 &&
			//len(obj.NowPlayingQueueFullItems[0].MediaSources) > 0 &&
			obj.PlayState.PlayMethod != "" {
			var state string

			if !obj.PlayState.IsPaused {
				state = "paused"
			} else {
				state = "in progress"
			}
			bitrate := float64(obj.NowPlayingQueueFullItems[0].MediaSources[0].Bitrate) / 1000000.0
			name := obj.NowPlayingQueueFullItems[0].MediaSources[0].Name
			sessionString = fmt.Sprintf("%s is playing (%s): %s\nPlayback: %s\nBitrate: %.2f Mbps\nDevice: %s\n", obj.UserName, state, name, obj.PlayState.PlayMethod, bitrate, obj.DeviceName)
		} else if len(obj.NowPlayingQueueFullItems) > 0 &&
			//len(obj.NowPlayingQueueFullItems[0].MediaSources) > 0 &&
			obj.PlayState.PlayMethod == "" {
			continue
		} else {
			sessionString = fmt.Sprintf("%s is chilling in the menus\n", obj.UserName)
		}
		sessionStrings = append(sessionStrings, sessionString)
	}
	formattedSessions = genericInfo + strings.Join(sessionStrings, "\n")

	return formattedSessions
}
