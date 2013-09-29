package main

import (
	"errors"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
	"text/template"
	"time"
)
import "database/sql"
import _ "github.com/mattn/go-sqlite3"
import "github.com/c9s/appcast"
import "github.com/c9s/rss"
import _ "github.com/c9s/appcast/server/uploader"

const UPLOAD_DIR = "uploads"
const SQLITEDB = "appcast.db"

var ErrFileIsRequired = errors.New("file is required.")
var ErrReleaseInsertFailed = errors.New("release insert failed.")

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
		createReleaseTable(db)
	}
	return db
}

func createReleaseTable(db *sql.DB) {
	if _, err := db.Exec(`create table releases(
		id integer auto_increment,
		title varchar,
		desc text,
		releaseNotesLink varchar,
		pubDate datetime default current_timestamp,
		filename varchar,
		length int,
		mimetype varchar,
		url varchar,
		version varchar,
		shortVersionString varchar,
		dsaSignature varchar
	);`); err != nil {
		log.Fatal(err)
	}
}

func GetMimeTypeByFilename(filename string) string {
	return mime.TypeByExtension(path.Ext(filename))
}

func GetFileLength(filepath string) (int64, error) {
	stat, err := os.Stat(filepath)
	if err != nil {
		return 0, err
	}
	return stat.Size(), nil
}

func UploadNewReleaseFromRequest(r *http.Request) error {
	file, fileReader, err := r.FormFile("file")
	if err == http.ErrMissingFile {
		return err
	}
	if err != nil {
		return err
	}

	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	dstFilePath := path.Join(UPLOAD_DIR, fileReader.Filename)

	if err = ioutil.WriteFile(dstFilePath, data, 0777); err != nil {
		return err
	}

	length, _ := GetFileLength(dstFilePath)
	mimetype := GetMimeTypeByFilename(fileReader.Filename)

	title := r.FormValue("title")
	desc := r.FormValue("desc")
	pubDate := r.FormValue("pubDate")
	version := r.FormValue("version")
	shortVersionString := r.FormValue("shortVersionString")
	releaseNotesLink := r.FormValue("releaseNotesLink")
	dsaSignature := r.FormValue("dsaSignature")

	var newItem = appcast.Item{}
	newItem.Title = title
	newItem.Description = desc
	if pubDate != "" {
		newItem.PubDate = rss.Date(pubDate)
	}
	newItem.Enclosure.SparkleVersion = version
	newItem.Enclosure.SparkleVersionShortString = shortVersionString
	newItem.Enclosure.SparkleDSASignature = dsaSignature
	newItem.SparkleReleaseNotesLink = releaseNotesLink
	_ = newItem

	result, err := db.Exec(`INSERT INTO releases 
		(title, desc, pubDate, version, shortVersionString, releaseNotesLink, dsaSignature, filename, length, mimetype)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		title, desc, pubDate, version, shortVersionString, releaseNotesLink, dsaSignature, fileReader.Filename, length, mimetype)
	if err != nil {
		return err
	}

	if id, err := result.LastInsertId(); err == nil {
		log.Println("Record created", id)
	}

	log.Println("New Release Uploaded", title, version, shortVersionString, desc, pubDate, dsaSignature, length, mimetype)
	return nil
}

func UploadPageHandler(w http.ResponseWriter, r *http.Request) {

	var err error
	if r.Method == "POST" {
		if err := UploadNewReleaseFromRequest(r); err != nil {
			panic(err)
		}
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

func ScanRowToAppcastItem(rows *sql.Rows) (*appcast.Item, error) {
	var title, desc, version, shortVersionString, filename, mimetype, dsaSignature string
	var pubDate time.Time
	var length int64
	var err = rows.Scan(&title, &desc, &pubDate, &version, &shortVersionString, &filename, &mimetype, &length, &dsaSignature)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var item = appcast.Item{}
	item.Title = title
	item.Description = desc
	item.PubDate = rss.Date(pubDate.Format(time.RFC822Z))
	item.Enclosure.Length = length
	item.Enclosure.Type = mimetype
	item.Enclosure.SparkleVersion = version
	item.Enclosure.SparkleVersionShortString = shortVersionString
	item.Enclosure.SparkleDSASignature = dsaSignature
	return &item, nil
}

func QueryReleases() (*sql.Rows, error) {
	return db.Query(`SELECT 
		title, desc, pubDate, version, shortVersionString, filename, mimetype, length, dsaSignature
		FROM releases ORDER BY pubDate DESC`)
}

func AppcastXmlHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/xml; charset=UTF-8")

	rows, err := QueryReleases()
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
		if item, err := ScanRowToAppcastItem(rows); err == nil {
			appcastRss.Channel.AddItem(item)
		}
	}
	appcastRss.WriteTo(w)
}

func DownloadFileHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.RequestURI, r.URL)
	/*
		log.Println(r.URL.Path)
		log.Println(r.URL.Opaque)
		log.Println(r.URL.Fragment)
	*/
}

func main() {
	db = ConnectDB()
	defer db.Close()

	http.HandleFunc("/download/", DownloadFileHandler)
	http.HandleFunc("/upload", UploadPageHandler)
	http.HandleFunc("/appcast.xml", AppcastXmlHandler)

	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	log.Println("Listening http://localhost:8080 ...")
	// http.HandleFunc("/upload", UploadNewReleaseHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
