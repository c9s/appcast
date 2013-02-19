package appcast

import "github.com/c9s/go-rss/rss"

// Item from go-rss
type Item struct {
	rss.Item
    Enclosure ItemEnclosure `xml:"enclosure"`


	SparkleReleaseNotesLink string `xml:"sparkle:releaseNotesLink"`

	// XXX: support localization
	// <sparkle:releaseNotesLink xml:lang="de">http://you.com/app/2.0_German.html</sparkle:releaseNotesLink>

	SparkleMinimumSystemVersion string `xml:"sparkle:minimumSystemVersion"`

	// <sparkle:minimumSystemVersion>10.7.1</sparkle:minimumSystemVersion>
}

func (item * Item) SetEnclosure(enclosure * ItemEnclosure) {
	item.Enclosure = *enclosure
}

func (item * Item) AddCategory(category string) {
	item.Category = append(item.Category, category)
}


