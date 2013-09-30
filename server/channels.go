package main

import (
	"crypto/sha1"
	"database/sql"
	"fmt"
	"github.com/c9s/appcast"
	"github.com/c9s/gatsby"
	"github.com/c9s/rss"
	"time"
)

// Channel Record
type Channel struct {
	Title       string `field:"title"`
	Description string `field:"description"`
	Identity    string `field:"identity"`
	Token       string `field:"token"`
	gatsby.BaseRecord
}

func (self *Channel) RegenerateToken(secret string) string {
	h := sha1.New()
	h.Write([]byte(self.Title + self.Description + self.Identity))
	h.Write([]byte(time.Now().Format(time.RFC822Z)))
	self.Token = fmt.Sprintf("%x", h.Sum(nil))
	return self.Token
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

func FindChannelByIdentity(identity string, token string) *Channel {
	row := db.QueryRow(`SELECT id, title, description FROM channels WHERE identity = ? AND token = ?`, identity, token)
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
		Token:       token,
	}
	return &channel
}
