package appcast

type ItemEnclosure struct {
    URL					string `xml:"url,attr"`
    Type				string `xml:"type,attr"`
    Length				string `xml:"length,attr"`
    Version				string `xml:"version,attr"`
    VersionShortString  string `xml:"versionShortString,attr"`
	DSASignature		string `xml:"dsaSignature,attr"`
}
