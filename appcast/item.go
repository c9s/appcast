package appcast

// Item from go-rss
type Item struct {
    Title       string        `xml:"title"`
    Link        string        `xml:"link"`
    Comments    string        `xml:"comments"`
    PubDate     Date          `xml:"pubDate"`
    GUID        string        `xml:"guid"`
    Category    []string      `xml:"category"`
    Description string        `xml:"description"`
    Content     string        `xml:"content"`
    Enclosure   SparkleItemEnclosure `xml:"enclosure"`
}

type SparkleItem struct {
	Item
	SparkleReleaseNotesLink string `xml:"sparkle:releaseNotesLink"`
}

func (item * SparkleItem) SetEnclosure(enclosure * SparkleItemEnclosure) {
	item.Enclosure = *enclosure
}

func (item * SparkleItem) AddCategory(category string) {
	item.Category = append(item.Category, category)
}


