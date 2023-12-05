package main

type CommandFunc func() string

var CommandMap = map[string]CommandFunc{
	"sayhi": func() string {
		return "Hi!"
	},
	"jellystatus": func() string {
		return GetSessions()
	},
}
