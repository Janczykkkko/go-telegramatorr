package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/PaesslerAG/jsonpath"
)

// GetSessions fetches sessions from Jellyfin
func GetSessions(jellyfinAddress, apiKey string) (int, error) {
	// Construct the URL with provided address and API key
	url := jellyfinAddress + "/Sessions?api_key=" + apiKey

	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	log.Printf("API request to %s completed with status code: %d", jellyfinAddress, resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var jsonData interface{}
	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		return 0, err
	}

	// Use jsonpath to extract the NowPlayingItem count
	result, err := jsonpath.Get("$[*].NowPlayingItem", jsonData)
	if err != nil {
		return 0, err
	}

	// Type assertion to extract the integer value
	count := len(result.([]interface{}))
	return count, nil
}
