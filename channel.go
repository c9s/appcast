package appcast

type Channel struct {
    Title         string `xml:"title"`
    Link          string `xml:"link"`
    Description   string `xml:"description"`
    Language      string `xml:"language"`
    LastBuildDate Date   `xml:"lastBuildDate"`
    Item          []SparkleItem `xml:"item"`
}

func (channel * Channel) AddItem( item * SparkleItem ) {
	channel.Item = append(channel.Item, *item)
}

