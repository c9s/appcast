package appcast

import "github.com/c9s/rss"
import "mime"
import "path"
import "os"

func init() {
	var err = mime.AddExtensionType(".zip", "application/octet-stream")
	if err != nil {
		panic(err)
	}
}

type ItemEnclosure struct {
	rss.ItemEnclosure
	SparkleVersion            string `xml:"http://www.andymatuschak.org/xml-namespaces/sparkle version,attr"`
	SparkleVersionShortString string `xml:"http://www.andymatuschak.org/xml-namespaces/sparkle versionShortString,attr,omitempty"`
	SparkleDSASignature       string `xml:"http://www.andymatuschak.org/xml-namespaces/sparkle dsaSignature,attr,omitempty"`
}

// Return ItemEnclosure object with Type, Length
func CreateItemEnclosureFromFile(filepath string) (*ItemEnclosure, error) {
	enclosure := ItemEnclosure{}
	mimetype := mime.TypeByExtension(path.Ext(filepath))
	enclosure.Type = mimetype

	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	enclosure.Length = stat.Size()
	defer file.Close()
	return &enclosure, nil
}
