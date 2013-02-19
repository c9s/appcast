package appcast
import "mime"
import "os"

type ItemEnclosure struct {
    URL					string `xml:"url,attr"`
    Type				string `xml:"type,attr"`
    Length				int64 `xml:"length,attr"`
}

type SparkleItemEnclosure struct {
	ItemEnclosure
	SparkleVersion			   string `xml:"sparkle:version,attr"`
	SparkleVersionShortString  string `xml:"sparkle:versionShortString,attr"`
	SparkleDSASignature		   string `xml:"sparkle:dsaSignature,attr"`
}

// Return SparkleItemEnclosure object with Type, Length
func CreateItemEnclosureFromFile(path string) (*SparkleItemEnclosure, error) {
	enclosure := SparkleItemEnclosure{}
	mimetype := mime.TypeByExtension(path)
	enclosure.Type = mimetype

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	stat ,err := file.Stat()
	if err != nil {
		return nil, err
	}
	enclosure.Length = stat.Size()
	defer file.Close()
	return &enclosure, nil
}

