package appcast
import "testing"

func TestItemEnclosure(t * testing.T) {
	enclosure, err := CreateItemEnclosureFromFile("tests/tests.zip")
	if err != nil {
		t.Error(err)
	}
	if enclosure == nil {
		t.Error("Enclosure is empty")
	}

// 	a := map[string] string {
// 		"sparkle:version":"220",
// 	}
// 	_ = a


}

