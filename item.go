package appcast

type Item struct {
    Title       string        `xml:"title"`
    Link        string        `xml:"link"`
    Comments    string        `xml:"comments"`
    PubDate     Date          `xml:"pubDate"`
    GUID        string        `xml:"guid"`
    Category    []string      `xml:"category"`
    Enclosure   ItemEnclosure `xml:"enclosure"`
    Description string        `xml:"description"`
    Content     string        `xml:"content"`
	ReleaseNotesLink string	  `xml:"releaseNotesLink"`
}

func (item * Item) SetEnclosure(enclosure * ItemEnclosure) {
	item.Enclosure = *enclosure
}

func (item * Item) AddCategory(category string) {
	item.Category = append(item.Category, category)
}


