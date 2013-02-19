package appcast

import (
    "encoding/xml"
    "io/ioutil"
    "net/http"
)

func ParseContent(text []byte) (*RSS, error) {
    var rss = RSS{}
    err := xml.Unmarshal(text, &rss)
    if err != nil {
        return nil, err
    }
    return &rss, nil
}

func MarshalIndent(rss * RSS) ([]byte,error) {
    content, err := xml.MarshalIndent(rss,"","  ")
	if err != nil {
		return nil, err
	}
	return content, nil
}

func WriteFile(path string, rss * RSS) (error) {
	rss.Version = "2.0"
	rss.XmlNSSparkle = "http://www.andymatuschak.org/xml-namespaces/sparkle"
	rss.XmlNSDC = "http://purl.org/dc/elements/1.1/"
    content, err := MarshalIndent(rss)
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

