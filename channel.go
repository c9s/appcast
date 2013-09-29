package appcast

import "github.com/c9s/rss"

type Channel struct {
	rss.Channel
	Items []Item `xml:"item"`
}

func (c *Channel) AddItem(item *Item) {
	c.Items = append(c.Items, *item)
}

func (c *Channel) Len() int {
	return len(c.Items)
}
