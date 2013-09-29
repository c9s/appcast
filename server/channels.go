package main

import "github.com/c9s/appcast"
import "github.com/c9s/rss"
import "database/sql"

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

func FindChannel(identity string) *appcast.Channel {
	row := db.QueryRow(`SELECT id, title, description WHERE identity = ?`, identity)

	var id int64
	var title, description string
	err := row.Scan(&id, &title, &description)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		panic(err)
	}
	channel := appcast.Channel{}
	channel.Title = title
	channel.Description = description
	return &channel
}
