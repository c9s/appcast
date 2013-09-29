package main

import (
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"path"
	"text/template"
	"time"
)
import "database/sql"
import _ "github.com/mattn/go-sqlite3"
import "github.com/c9s/appcast"
import "github.com/c9s/rss"
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
			desc text,
			releaseNote text,
			pubDate datetime default current_timestamp,
			filename varchar,
			length int,
			mimetype varchar,
			url varchar,
			version varchar,
			shortVersion varchar,
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
		if err == http.ErrMissingFile {
			log.Println("Missing file", err)
		}
		if err != nil {
			log.Println("FormFile", err)
		}
		defer file.Close()
		data, err := ioutil.ReadAll(file)
		if err != nil {
			log.Println("ReadAll", err)
		}

		dstFilePath := path.Join(UPLOAD_DIR, fileReader.Filename)

		if err = ioutil.WriteFile(dstFilePath, data, 0777); err != nil {
			log.Println(err)
		}

		/*
			filename varchar,
			length int,
			mimetype varchar,
			url varchar,
			version varchar,
			dsaSignature varchar
		*/
		stat, err := os.Stat(dstFilePath)
		if err != nil {
			log.Println(err)
		}
		length := stat.Size()
		mimetype := mime.TypeByExtension(path.Ext(fileReader.Filename))

		title := r.FormValue("title")
		desc := r.FormValue("desc")
		pubDate := r.FormValue("pubDate")
		version := r.FormValue("version")
		shortVersion := r.FormValue("shortVersion")
		releaseNote := r.FormValue("releaseNote")
		dsaSignature := r.FormValue("dsaSignature")

		result, err := db.Exec(`INSERT INTO releases 
			(title, desc, pubDate, version, shortVersion, releaseNote, dsaSignature, filename, length, mimetype)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			title, desc, pubDate, version, shortVersion, releaseNote, dsaSignature, fileReader.Filename, length, mimetype)
		if err != nil {
			panic(err)
		}

		if id, err := result.LastInsertId(); err == nil {
			log.Println("Record created", id)
		}

		log.Println("New Release Uploaded", title, version, shortVersion, desc, pubDate, dsaSignature, length, mimetype)
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
	w.Header().Set("Content-Type", "text/xml; charset=UTF-8")

	rows, err := db.Query(`SELECT 
		title, desc, pubDate, version, shortVersion, filename, mimetype, length 
		FROM releases ORDER BY pubDate DESC`)
	if err != nil {
		log.Fatal("Query failed:", err)
	}
	defer rows.Close()

	appcastRss := appcast.New()
	appcastRss.Channel.Title = channelmeta["title"].(string)
	appcastRss.Channel.Description = channelmeta["description"].(string)
	appcastRss.Channel.Link = channelmeta["link"].(string)
	appcastRss.Channel.Language = channelmeta["language"].(string)

	for rows.Next() {
		var title, desc, version, shortVersion, filename, mimetype string
		var pubDate time.Time
		var length int64
		err = rows.Scan(&title, &desc, &pubDate, &version, &shortVersion, &filename, &mimetype, &length)
		if err != nil {
			log.Println(err)
			continue
		}

		var item = appcast.Item{}
		item.Title = title
		item.Description = desc
		// item.PubDate = rss.Date(time.Unix(pubDate, 0).Format(time.RFC822Z))
		item.PubDate = rss.Date(pubDate.Format(time.RFC822Z))
		item.Enclosure.Length = length
		item.Enclosure.Type = mimetype
		item.Enclosure.SparkleVersion = version
		item.Enclosure.SparkleVersionShortString = shortVersion
		// item.ImportFile(filename)

		appcastRss.Channel.AddItem(&item)
	}
	appcastRss.WriteTo(w)
}

func main() {
	db = ConnectDB()
	defer db.Close()

	http.HandleFunc("/upload", UploadPageHandler)
	http.HandleFunc("/appcast.xml", AppcastXmlHandler)

	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	log.Println("Listening http://localhost:8080 ...")
	// http.HandleFunc("/upload", UploadNewReleaseHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
