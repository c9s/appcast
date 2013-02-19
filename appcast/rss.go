package appcast

// XXX: better solution? use lower-case, because we need to encode it with lowercase
type rss struct {
	Channel Channel `xml:"channel"`
    Version	string `xml:"version,attr"`
	XmlNSSparkle string `xml:"xmlns:sparkle,attr"`
	XmlNSDC string `xml:"xmlns:dc,attr"`
	/*
	<rss version="2.0" 
		xmlns:sparkle="http://www.andymatuschak.org/xml-namespaces/sparkle"  
		xmlns:dc="http://purl.org/dc/elements/1.1/">
	*/
}
