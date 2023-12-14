package main

import (
	"fmt"
)

func generateHelpText() string {
	helpText := "I understand: "
	for key := range CommandMap {
		helpText += fmt.Sprintf("/%s ", key)
	}
	return helpText
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
