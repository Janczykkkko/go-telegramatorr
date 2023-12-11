package main

import (
	"fmt"
	"strconv"
	"strings"
)

func GetSessionsStr() string {
	var (
		formattedSessions string
		sessionStrings    []string
	)
	genericInfo := "Here's an activity report from Jellyfin: \n\n"

	jellyJSON, err := GetSessions()
	if err != nil {
		errMsg := "Error: " + err.Error()
		return errMsg
	}

	for _, obj := range jellyJSON {
		var sessionString string
		if len(obj.NowPlayingQueueFullItems) > 0 &&
			obj.PlayState.PlayMethod != "" {
			var state string
			var bitrate float64
			var substream string
			if obj.PlayState.IsPaused {
				state = "paused"
			} else {
				state = "in progress"
			}
			bitrate = float64(obj.NowPlayingQueueFullItems[0].MediaSources[0].Bitrate) / 1000000.0
			name := obj.NowPlayingQueueFullItems[0].MediaSources[0].Name

			SubtitleStreamIndex := obj.PlayState.SubtitleStreamIndex
			if SubtitleStreamIndex >= 0 && SubtitleStreamIndex < len(obj.NowPlayingQueueFullItems[0].MediaStreams) {
				substream = obj.NowPlayingQueueFullItems[0].MediaStreams[obj.PlayState.SubtitleStreamIndex].DisplayTitle
			} else {
				substream = "None"
			}

			sessionString = fmt.Sprintf("%s is playing (%s): %s\nPlayback: %s\nBitrate: %.2f Mbps\nSubtitles: %s\nDevice: %s\n", obj.UserName, state, name, obj.PlayState.PlayMethod, bitrate, substream, obj.DeviceName)

		} else if len(obj.FullNowPlayingItem.Container) > 0 && //mobile not showing fix
			obj.PlayState.PlayMethod != "" &&
			!obj.PlayState.IsPaused {
			var bitrateData int
			var substream string
			var state string
			if obj.PlayState.IsPaused {
				state = "paused"
			} else {
				state = "in progress"
			}
			for _, stream := range obj.NowPlayingItem.MediaStreams {
				if stream.Type == "Video" {
					bitrateData = stream.BitRate
					break
				}
			}
			bitrateFloat := float64(bitrateData) / 1000000.0
			bitrate := strconv.FormatFloat(bitrateFloat, 'f', -1, 64)
			SubtitleStreamIndex := obj.PlayState.SubtitleStreamIndex
			if SubtitleStreamIndex >= 0 && SubtitleStreamIndex < len(obj.NowPlayingItem.MediaStreams) {
				substream = obj.NowPlayingItem.MediaStreams[obj.PlayState.SubtitleStreamIndex].DisplayTitle
			} else {
				substream = "None"
			}
			sessionString = fmt.Sprintf("%s is playing(%s): %s\nPlayback: %s\nBitrate: %s Mbps\nSubtitles: %s\nDevice: %s\n", obj.UserName, state, obj.NowPlayingItem.Name, obj.PlayState.PlayMethod, bitrate, substream, obj.DeviceName)
		} else {
			continue
		}
		sessionStrings = append(sessionStrings, sessionString)
	}
	if len(strings.Join(sessionStrings, "\n")) != 0 {
		formattedSessions = genericInfo + strings.Join(sessionStrings, "\n")
	} else {
		formattedSessions = "Nothing is playing."
	}
	return formattedSessions
}
