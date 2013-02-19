package appcast
import "testing"

func TestItemEnclosure(t * testing.T) {
	en, err := CreateItemEnclosureFromFile("tests/tests.zip")
	if err != nil {
		t.Error(err)
	}
	if en == nil {
		t.Error("Enclosure is empty")
	}

	// application/octet-stream
	// application/zip
	if en.Type == "" {
		t.Error("Enclosure Type empty")
	}
	if en.Type != "application/octet-stream" {
		t.Error("enclosure type is not application/octet-stream")
	}
	if en.Length == 0 {
		t.Error("enclosure length is empty")
	}
}

