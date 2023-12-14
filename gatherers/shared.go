package gatherers

import (
	"fmt"
	"strings"
)

func GetAllSessionsStr(jellyfinAddress, jellyfinApiKey, plexAddress, plexApiKey string) string {

	response := []string{"Here's a report from your player(s):\n"}

	var jellySessions []SessionData
	if jellyfinAddress != "" || jellyfinApiKey != "" {
		sessions, err := GetJellySessions(jellyfinAddress, jellyfinApiKey)
		if err != nil {
			response = append(response, fmt.Sprintf("Error getting Jellyfin sessions: %s", err))
		}
		jellySessions = sessions
	}

	var plexSessions []SessionData
	if plexAddress != "" || plexApiKey != "" {
		sessions, err := GetPlexSessions(plexAddress, plexApiKey)
		if err != nil {
			response = append(response, fmt.Sprintf("Error getting Plex sessions: %s", err))
		}
		plexSessions = sessions
	}

	sessions := append(plexSessions, jellySessions...)

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

	return strings.Join(response, "\n")
}
