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

func ParseContent(text []byte) (*RSS, error) {
    var rss = RSS{}
    err := xml.Unmarshal(text, &rss)
    if err != nil {
        return nil, err
    }
    return &rss, nil
}

func WriteFile(path string, rss * RSS) (error) {
	rss.Version = "2.0"
	rss.XmlNSSparkle = "http://www.andymatuschak.org/xml-namespaces/sparkle"
	rss.XmlNSDC = "http://purl.org/dc/elements/1.1/"
    content, err := xml.Marshal(rss)
    if err != nil {
        return err
    }
	err = ioutil.WriteFile(path, content, 0666)
	if err != nil {
		return err
	}
	return nil
}

func ReadFile(file string) (*RSS, error) {
    text, err := ioutil.ReadFile(file)
    if err != nil {
        return nil, err
    }
    return ParseContent(text)
}

func ReadUrl(url string) (*RSS, error) {
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


