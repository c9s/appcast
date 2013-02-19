package appcast
import "testing"

func TestRSSXML(t *testing.T) {
	channel, err := ReadFile("tests/appcast.xml");


	if err != nil {
		t.Errorf("RSS read fail.")
	}
	if channel == nil {
		t.Errorf("Channel is empty.")
	}

	if len(channel.Item) == 0 {
		t.Errorf("Item length is zero.")
	}

	for _ , item := range channel.Item {

		if len(item.Title) == 0 {
			t.Errorf("Item Title is empty.")
		}
		if len(item.SparkleReleaseNotesLink) == 0 {
			t.Errorf("sparkle:releaseNotesLink is empty")
		}
		if len(item.Enclosure.SparkleVersion) == 0 {
			t.Errorf("Enclosure sparkle:version not found.")
		}
		if item.Enclosure.Length == 0 {
			t.Errorf("Enclosure length not found.")
		}
		if len(item.Enclosure.SparkleDSASignature) == 0 {
			t.Errorf("Enclosure sparkle:dsaSignature not found.")
		}
	}

	err = WriteFile("tests/appcast-out.xml", channel)
	if err != nil {
		t.Error("Can not write xml file.")
	}

	_ = channel
}

