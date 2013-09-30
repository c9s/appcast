package main

import "github.com/c9s/gatsby"

type Release struct {
	Id                 int64     `field:"id"`
	Title              string    `field:"title"`
	Description        string    `field:"desc"`
	ReleaseNotesLink   string    `field:"releaseNotesLink"`
	PubDate            time.Time `field:"pubDate"`
	Filename           string    `field:"filename"`
	ChannelId          int64     `field:"channelId"`
	Length             int64     `field:"length"`
	Mimetype           string    `field:"mimetype"`
	DSASignature       string    `field:"dsaSignature"`
	Version            string    `field:"version"`
	ShortVersionString string    `field:"shortVersionString"`
	Token              string    `field:"token"`
	gatsby.BaseRecord
}

func FindReleaseByToken(token string) *appcast.Channel {
	row := db.QueryRow(`SELECT 
		id, 
		title, 
		description, 
		version, 
		shortVersionString, 
		filename,
		length,
		mimetype
		WHERE token = ?`, token)

	var (
		id                 int64
		title              string
		description        string
		version            string
		shortVersionString string
		filename           string
		length             int64
		mimetype           string
	)
	err := row.Scan(&id, &title, &description, &version, &shortVersionString, &filename, &length, &mimetype)
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
