package main

import (
	"crypto/sha1"
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
	"text/template"
	"time"
)

import (
	"github.com/c9s/appcast"
	_ "github.com/c9s/appcast/server/uploader"
	"github.com/c9s/jsondata"
	"github.com/c9s/rss"
	_ "github.com/mattn/go-sqlite3"
)

const UPLOAD_DIR = "uploads"
const SQLITEDB = "appcast.db"

var ErrFileIsRequired = errors.New("file is required.")
var ErrReleaseInsertFailed = errors.New("release insert failed.")

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
		createAccountTable(db)
		createReleaseTable(db)
		createChannelTable(db)
	}
	return db
}

func createAccountTable(db *sql.DB) {
	if _, err := db.Exec(`create table account(
		id integer auto_increment,
		account varchar,
		token varchar
	);`); err != nil {
		log.Fatal(err)
	}
}

func createChannelTable(db *sql.DB) {
	if _, err := db.Exec(`create table channels(
		id integer auto_increment,
		title varchar,
		description varchar,
		identity varchar
	);`); err != nil {
		log.Fatal(err)
	}
}

func createReleaseTable(db *sql.DB) {
	if _, err := db.Exec(`create table releases(
		id integer auto_increment,
		title varchar,
		desc text,
		releaseNotesLink varchar,
		pubDate datetime default current_timestamp,
		filename varchar,
		channelId int,
		length int,
		mimetype varchar,
		url varchar,
		version varchar,
		shortVersionString varchar,
		dsaSignature varchar,
		token varchar
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

func UploadReleaseHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	item, err := CreateNewReleaseFromRequest(r)
	if err != nil {
		var msg = jsondata.Map{"error": err}
		msg.WriteTo(w)
	}
	_ = item

}

func CreateNewReleaseFromRequest(r *http.Request) (*appcast.Item, error) {
	file, fileReader, err := r.FormFile("file")
	if err == http.ErrMissingFile {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	dstFilePath := path.Join(UPLOAD_DIR, fileReader.Filename)
	if err = ioutil.WriteFile(dstFilePath, data, 0777); err != nil {
		return nil, err
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

	h := sha1.New()
	h.Write([]byte(title))
	h.Write([]byte(version))
	h.Write([]byte(shortVersionString))
	h.Write(data)
	token := fmt.Sprintf("% x", h.Sum(nil))
	_ = token

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

	result, err := db.Exec(`INSERT INTO releases 
		(title, desc, pubDate, version, shortVersionString, releaseNotesLink, dsaSignature, filename, length, mimetype)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		title,
		desc,
		pubDate,
		version,
		shortVersionString,
		releaseNotesLink,
		dsaSignature,
		fileReader.Filename,
		length, mimetype)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	log.Println("Record created", id)

	log.Println("New Release Uploaded", title, version, shortVersionString, desc, pubDate, dsaSignature, length, mimetype)
	return &newItem, nil
}

func UploadPageHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	if r.Method == "POST" {
		if _, err := CreateNewReleaseFromRequest(r); err != nil {
			panic(err)
		}
	}

	templates, err := template.ParseFiles("templates/upload.html")
	if err != nil {
		panic(err)
	}
	t := templates.Lookup("upload.html")
	if t != nil {
		channel := GetChannel("gotray")
		err := t.Execute(w, channel)
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

func GetChannel(identity string) *appcast.Channel {
	if channel, ok := channels[identity]; ok {
		return &channel
	}
	return nil
}

func AppcastXmlHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/xml; charset=UTF-8")

	rows, err := QueryReleases()
	if err != nil {
		log.Fatal("Query failed:", err)
	}
	defer rows.Close()

	channel := GetChannel("gotray")

	appcastRss := appcast.New()
	appcastRss.Channel.Title = channel.Title
	appcastRss.Channel.Description = channel.Description
	appcastRss.Channel.Link = channel.Link
	appcastRss.Channel.Language = channel.Language

	for rows.Next() {
		if item, err := ScanRowToAppcastItem(rows); err == nil {
			appcastRss.Channel.AddItem(item)
		}
	}
	appcastRss.WriteTo(w)
}

/*
For route: /download/gotray/{token}

/download/gotray/be24d1c54d0ba415b8897b02f0c38d89
*/
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
	http.HandleFunc("/=/upload", UploadReleaseHandler)
	http.HandleFunc("/appcast.xml", AppcastXmlHandler)
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	log.Println("Listening http://localhost:8080 ...")
	// http.HandleFunc("/upload", UploadNewReleaseHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
