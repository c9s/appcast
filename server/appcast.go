package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"text/template"
)
import "database/sql"
import _ "github.com/mattn/go-sqlite3"
import _ "github.com/c9s/appcast"
import _ "github.com/c9s/appcast/server/uploader"
import "os"

const UPLOAD_DIR = "uploads"
const SQLITEDB = "appcast.db"

var channelmeta = map[string]interface{}{
	"title":         "GoTray Appcast",
	"link":          "http://gotray.extremedev.org/appcast.xml",
	"description":   "Most recent changes with links to updates.",
	"language":      "en",
	"lastBuildDate": nil,
}

var db *sql.DB

func ConnectDB() *sql.DB {
	var initDB = false
	_, err := os.Stat(SQLITEDB)
	if os.IsNotExist(err) {
		// init db
		initDB = true
	}

	// os.Remove("./appcast.db")
	db, err := sql.Open("sqlite3", SQLITEDB)
	if err != nil {
		log.Fatal(err)
	}

	if initDB {
		log.Println("Initializing database schema...")

		/*
			<title>Version 1.4.10 (2 bugs fixed)</title>
			<pubDate>Sun, 17 Feb 2013 21:28:09 +0800</pubDate>
			<description>Fix auto-update</description>
			<enclosure
				url="http://gotray.extremedev.org/app/GoTray-1.4.10.zip"
				type="application/octet-stream"
				length="7292169"
				sparkle:version="105"
				sparkle:dsaSignature="MC0CFQCHnbi7kJ7C5wAA+QLu52NvFim4ZQIUdgxVJatWmwbWGWXrNGZJc2sDKjk=">
			</enclosure>
		*/

		if _, err := db.Exec(`create table releases(
			id integer auto_increment,
			title varchar,
			description text,
			pubDate datetime default current_timestamp,
			length int,
			type varchar,
			url varchar,
			version varchar,
			dsaSignature varchar
		);`); err != nil {
			log.Println(err)
		}
	}
	return db
}

func init() {
}

func UploadPageHandler(w http.ResponseWriter, r *http.Request) {

	var err error
	if r.Method == "POST" {
		file, fileReader, err := r.FormFile("file")
		if err != nil {
			log.Println("FormFile", err)
		}
		defer file.Close()
		data, err := ioutil.ReadAll(file)
		if err != nil {
			log.Println("ReadAll", err)
		}
		if err = ioutil.WriteFile(path.Join(UPLOAD_DIR, fileReader.Filename), data, 0777); err != nil {
			log.Println(err)
		}

		title := r.FormValue("title")
		desc := r.FormValue("desc")
		pubDate := r.FormValue("pubDate")
		version := r.FormValue("version")
		dsaSignature := r.FormValue("dsaSignature")
		log.Println("New Upload", title, version, desc, pubDate, dsaSignature)
	}

	templates, err := template.ParseFiles("templates/upload.html")
	if err != nil {
		panic(err)
	}
	t := templates.Lookup("upload.html")
	if t != nil {
		err := t.Execute(w, channelmeta)
		if err != nil {
			panic(err)
		}
	}
}

func AppcastXmlHandler(w http.ResponseWriter, r *http.Request) {
	// appcast := appcast.New()
	// appcast.Write()
}

func main() {
	db = ConnectDB()
	defer db.Close()

	http.HandleFunc("/upload", UploadPageHandler)
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	log.Println("Listening http://localhost:8080 ...")
	// http.HandleFunc("/upload", UploadNewReleaseHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
