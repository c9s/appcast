package main

import "github.com/c9s/appcast"
import "github.com/c9s/rss"

var channels = map[string]appcast.Channel{
	"gotray": appcast.Channel{
		rss.Channel{
			Title:         "GoTray Appcast",
			Link:          "http://gotray.extremedev.org/appcast.xml",
			Description:   "Most recent changes with links to updates.",
			Language:      "en",
			LastBuildDate: "",
		},
		[]appcast.Item{},
	},
}
