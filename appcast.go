package appcast

import (
	"encoding/xml"
	"github.com/c9s/rss"
	"io"
	"io/ioutil"
	"net/http"
)

// XXX: better solution? use lower-case, because we need to encode it with lowercase
type Appcast struct {
	XMLName      xml.Name `xml:"rss"`
	XmlNSSparkle string   `xml:"http://www.andymatuschak.org/xml-namespaces/sparkle sparkle,attr"`
	XmlNSDC      string   `xml:"http://purl.org/dc/elements/1.1 dc,attr"`
	Channel      Channel  `xml:"channel"`
	rss.RSS
	/*
		<rss version="2.0"
			xmlns:sparkle="http://www.andymatuschak.org/xml-namespaces/sparkle"
			xmlns:dc="http://purl.org/dc/elements/1.1/">
	*/
}

func (self *Appcast) MarshalIndent() ([]byte, error) {
	content, err := xml.MarshalIndent(self, "", "  ")
	if err != nil {
		return nil, err
	}
	return content, nil
}

func (self *Appcast) WriteTo(w io.Writer) {
	content, err := self.MarshalIndent()
	if err == nil {
		w.Write(content)
	}
}

/*
Write appcast XML content to file.
*/
func (self *Appcast) WriteFile(path string) error {
	content, err := self.MarshalIndent()
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(path, content, 0666); err != nil {
		return err
	}
	return nil
}

func New() *Appcast {
	appcast := Appcast{}
	appcast.Version = "2.0"
	appcast.XmlNSSparkle = "http://www.andymatuschak.org/xml-namespaces/sparkle"
	appcast.XmlNSDC = "http://purl.org/dc/elements/1.1/"
	return &appcast
}

/*
Parse appcast XML content from bytes
*/
func ParseContentString(text string) (*Appcast, error) {
	var appcast = New()
	if err := xml.Unmarshal([]byte(text), appcast); err != nil {
		return nil, err
	}
	return appcast, nil
}

/*
Parse appcast XML content from bytes
*/
func ParseContent(text []byte) (*Appcast, error) {
	var appcast = New()
	err := xml.Unmarshal(text, appcast)
	if err != nil {
		return nil, err
	}
	return appcast, nil
}

func ParseFile(file string) (*Appcast, error) {
	text, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return ParseContent(text)
}

func ParseContentFromUrl(url string) (*Appcast, error) {
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
