package gatherers

import (
	"fmt"
	"strings"
)

func GetAllSessionsStr(jellyfinAddress, jellyfinApiKey, plexAddress, plexApiKey string) string {
	response := []string{"Here's a report from your player(s):"}
	sessions, errors := GetAllSessions(jellyfinAddress, jellyfinApiKey, plexAddress, plexApiKey)
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

func GetAllSessions(jellyfinAddress, jellyfinApiKey, plexAddress, plexApiKey string) (allSessions []SessionData, errors string) {

	var jellySessions []SessionData
	if jellyfinAddress != "" || jellyfinApiKey != "" {
		sessions, err := GetJellySessions(jellyfinAddress, jellyfinApiKey)
		if err != nil {
			errors = fmt.Sprintf("Error getting Jellyfin sessions: %s", err)
		}
		jellySessions = sessions
	}

	var plexSessions []SessionData
	if plexAddress != "" || plexApiKey != "" {
		sessions, err := GetPlexSessions(plexAddress, plexApiKey)
		if err != nil {
			errors = errors + "\n" + fmt.Sprintf("Error getting Plex sessions: %s", err)
		}
		plexSessions = sessions
	}

	allSessions = append(jellySessions, plexSessions...)

	return allSessions, errors
}
