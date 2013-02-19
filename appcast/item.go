package appcast

import "github.com/c9s/go-rss/rss"

// Item from go-rss
type Item struct {
	rss.Item
    Enclosure SparkleItemEnclosure `xml:"enclosure"`
	SparkleReleaseNotesLink string `xml:"sparkle:releaseNotesLink"`
}

func (item * Item) SetEnclosure(enclosure * SparkleItemEnclosure) {
	item.Enclosure = *enclosure
}

func (item * Item) AddCategory(category string) {
	item.Category = append(item.Category, category)
}


