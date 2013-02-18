package appcast

import (
    "encoding/xml"
    "io/ioutil"
    "net/http"
    "time"
)

type Date string

func (self Date) Parse() (time.Time, error) {
    t, err := time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", string(self)) // Wordpress format
    if err != nil {
        t, err = time.Parse(time.RFC822, string(self)) // RSS 2.0 spec
    }
    return t, err
}

func (self Date) Format(format string) (string, error) {
    t, err := self.Parse()
    if err != nil {
        return "", err
    }
    return t.Format(format), nil
}

func (self Date) MustFormat(format string) string {
    s, err := self.Format(format)
    if err != nil {
        return err.Error()
    }
    return s
}

func ParseContent(text []byte) (*Channel, error) {
    var rss struct {
        Channel Channel `xml:"channel"`
    }
    err := xml.Unmarshal(text, &rss)
    if err != nil {
        return nil, err
    }
    return &rss.Channel, nil
}

func WriteFile(path string, channel * Channel) (error) {
    content, err := xml.Marshal(channel)
    if err != nil {
        return err
    }
	err = ioutil.WriteFile(path, content, 0666)
	if err != nil {
		return err
	}
	return nil
}

func ReadFile(file string) (*Channel, error) {
    text, err := ioutil.ReadFile(file)
    if err != nil {
        return nil, err
    }
    return ParseContent(text)
}

func ReadUrl(url string) (*Channel, error) {
    response, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer response.Body.Close()
    text, err := ioutil.ReadAll(response.Body)
    if err != nil {
        return nil, err
    }
    return ParseContent(text)
}


