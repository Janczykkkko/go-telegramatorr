package gatherers

import (
	"fmt"
	"strings"
)

func GetAllSessionsStr(jellyfinAddress, jellyfinApiKey, plexAddress, plexApiKey string) string {
	response := []string{"Here's a report from your player(s):\n"}

	jellySessions, err := GetJellySessions(jellyfinAddress, jellyfinApiKey)
	if err != nil {
		response = append(response, fmt.Sprintf("Error getting Jellyfin sessions: %s", err))
	}
	plexSessions, err := GetPlexSessions(plexAddress, plexApiKey)
	if err != nil {
		response = append(response, fmt.Sprintf("Error getting Plex sessions: %s", err))
	}

	sessions := append(plexSessions, jellySessions...)

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

	if len(response) == 1 {
		return "Nothing is playing."
	}
	return strings.Join(response, "\n")
}
