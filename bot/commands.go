package bot

type CommandFunc func() string

var CommandMap = map[string]CommandFunc{
	"sayhi": func() string {
		return "Hi!"
	},
	"playstatus": func() string {
		return GetAllSessionsStr(jellyfinAddress, jellyfinApiKey, plexAddress, plexApiKey)
	},
}
