package main

import (
	"github.com/c9s/appcast"
	"github.com/c9s/rss"
	"testing"
)

func TestChannel(t *testing.T) {
	db = ConnectDB(":memory:")
	ch := appcast.Channel{
		rss.Channel{Title: "Testing", Description: "Description"},
		[]appcast.Item{},
	}
	id, err := CreateChannel("testing", &ch)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("created channel record", id)

	ch2 := FindChannelByIdentity("testing")
	if ch2 == nil {
		t.Fatal("testing channel not found.")
	}
}
