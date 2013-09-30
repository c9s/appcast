package main

import (
	"database/sql"
	// "github.com/c9s/appcast"
	// "github.com/c9s/rss"
	"github.com/c9s/gatsby"
	"testing"
)

func createTestDb() *sql.DB {
	var db = ConnectDB(":memory:")
	gatsby.SetupConnection(db, gatsby.DriverSqlite)
	return db
}

func TestChannel(t *testing.T) {
	db = createTestDb()

	var newChannel = Channel{
		Title:       "Testing",
		Description: "Description",
		Identity:    "testing",
	}
	newChannel.Init()
	var res = newChannel.Create()
	if res.Error != nil {
		t.Fatal(res.Error)
	}
	ch2 := FindChannelByIdentity("testing")
	if ch2 == nil {
		t.Fatal("testing channel not found.")
	}
}
