package bot

import (
	"telegramatorr/gatherers"
)

type CommandFunc func() string

var CommandMap = map[string]CommandFunc{
	"sayhi": func() string {
		return "Hi!"
	},
	"jellystatus": func() string {
		return GetSessionsStr()
	},
}

func GetSessionsStr() string {
	return gatherers.GetAllSessionsStr(jellyfinAddress, jellyfinApiKey, plexAddress, plexApiKey)
}
