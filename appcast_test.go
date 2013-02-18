package rss

import "testing"

func TestRSSXML(t *testing.T) {
	channel, err := ReadFile("tests/appcast.xml");

	if err != nil {
		t.Errorf("RSS read fail.")
	}
	if channel == nil {
		t.Errorf("Channel is empty.")
	}

	if len(channel.Item) == 0 {
		t.Errorf("Item length is zero")
	}

	for _ , item := range channel.Item {
		if len(item.Title) == 0 {
			t.Errorf("Item Title is empty")
		}
	}
	_ = channel
	_ = err
}

