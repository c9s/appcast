package appcast
import "encoding/xml"
import "github.com/c9s/go-rss/rss"

// XXX: better solution? use lower-case, because we need to encode it with lowercase
type RSS struct {


	XMLName xml.Name `xml:"rss"`

	rss.RSS
	XmlNSSparkle string `xml:"http://www.andymatuschak.org/xml-namespaces/sparkle sparkle,attr"`
	XmlNSDC string `xml:"http://purl.org/dc/elements/1.1 dc,attr"`
	Channel Channel `xml:"channel"`
	/*
	<rss version="2.0" 
		xmlns:sparkle="http://www.andymatuschak.org/xml-namespaces/sparkle"  
		xmlns:dc="http://purl.org/dc/elements/1.1/">
	*/
}

