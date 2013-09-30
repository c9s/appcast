package main

import (
	"database/sql"
	"github.com/c9s/appcast"
	"github.com/c9s/gatsby"
	"github.com/c9s/rss"
)

// Channel Record
type Channel struct {
	Title       string `field:"title"`
	Description string `field:"description"`
	Identity    string `field:"identity"`
	gatsby.BaseRecord
}

func (self *Channel) Init() {
	self.BaseRecord.SetTarget(self)
}

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

func CreateChannel(identity string, ch *appcast.Channel) (int64, error) {
	result, err := db.Exec(`INSERT INTO channels (title,description,identity) VALUES (?,?,?)`, ch.Title, ch.Description, identity)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func FindChannelByIdentity(identity string) *Channel {
	row := db.QueryRow(`SELECT id, title, description FROM channels WHERE identity = ?`, identity)
	var id int64
	var title, description string
	err := row.Scan(&id, &title, &description)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		panic(err)
	}
	channel := Channel{
		Title:       title,
		Description: description,
		Identity:    identity,
	}
	return &channel
}
