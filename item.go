package appcast

import "github.com/c9s/rss"

// Item from rss
type Item struct {

//	HTable string `xml:"http://www.w3.org/TR/html4/ table,attr"`

	rss.Item
    Enclosure ItemEnclosure `xml:"enclosure"`

	SparkleReleaseNotesLink string `xml:"http://www.andymatuschak.org/xml-namespaces/sparkle releaseNotesLink,omitempty"`

	// XXX: support localization
	// <sparkle:releaseNotesLink xml:lang="de">http://you.com/app/2.0_German.html</sparkle:releaseNotesLink>

	// <sparkle:minimumSystemVersion>10.7.1</sparkle:minimumSystemVersion>
	SparkleMinimumSystemVersion string `xml:"http://www.andymatuschak.org/xml-namespaces/sparkle minimumSystemVersion,omitempty"`
}

