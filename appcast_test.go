package appcast

import "testing"

func TestEmptyRSSXML(t *testing.T) {
	rss := RSS{}
	item := Item{}
	rss.Channel.AddItem(&item)
	if len(rss.Channel.Items) != 1 {
		t.Error("AddItem failed.")
	}
	_ = rss
}

func TestRSSXML(t *testing.T) {
	rss, err := ReadFile("tests/appcast.xml");
	if err != nil {
		t.Errorf("RSS read fail.")
	}
	if rss == nil {
		t.Errorf("RSS is empty.")
	}

	// fmt.Printf("rss type: %T\n",rss)
	// fmt.Printf("channel type: %T\n",rss.Channel)

	if len(rss.Channel.Items) == 0 {
		t.Errorf("Items length is zero.")
	}

	for _ , item := range rss.Channel.Items {
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

	/*
	err = WriteFile("tests/appcast-out.xml", rss)
	if err != nil {
		t.Error("Can not write xml file.")
	}
	*/
}

