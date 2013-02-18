package appcast
import "mime"
import "os"

type ItemEnclosure struct {
    URL					string `xml:"url,attr"`
    Type				string `xml:"type,attr"`
    Length				int64 `xml:"length,attr"`
    Version				string `xml:"version,attr"`
    VersionShortString  string `xml:"versionShortString,attr"`
	DSASignature		string `xml:"dsaSignature,attr"`
}

// Return ItemEnclosure object with Type, Length
func CreateItemEnclosureFromFile(path string) (*ItemEnclosure, error) {
	enclosure := ItemEnclosure{}
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

