package main

type CommandFunc func() string

var CommandMap = map[string]CommandFunc{
	"sayhi": func() string {
		return "a"
	},
	"jellystatus": func() string {
		return "b"
	},
}
