package main

import "testing"

import "time"

func TestRelease(t *testing.T) {
	db = createTestDb()

	now := time.Now()

	r := Release{
		Title:       "release title",
		Description: "description",
		PubDate:     &now,
		Token:       "s3cr3t",
	}
	r.Init()
	res := r.Create()
	if res.Error != nil {
		t.Log(res.Sql)
		t.Fatal(res.Error)
	}

	var r2 = FindReleaseByToken("s3cr3t")
	if r2 == nil {
		t.Fatal("release not found.")
	}

}
