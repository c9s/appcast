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
	"regexp"
	"text/template"
	"time"
)

import (
	"github.com/c9s/appcast"
	_ "github.com/c9s/appcast/server/uploader"
	"github.com/c9s/gatsby"
	// "github.com/c9s/jsondata"
	"github.com/c9s/rss"
	_ "github.com/mattn/go-sqlite3"
)

const UPLOAD_DIR = "uploads"
const SQLITEDB = "appcast.db"

var ErrFileIsRequired = errors.New("file is required.")
var ErrReleaseInsertFailed = errors.New("release insert failed.")

var db *sql.DB
var templates = template.Must(template.ParseFiles("templates/upload.html"))

func ConnectDB(dbname string) *sql.DB {
	var initDB = false
	_, err := os.Stat(dbname)
	if os.IsNotExist(err) {
		// init db
		initDB = true
	}

	// os.Remove("./appcast.db")
	db, err := sql.Open("sqlite3", dbname)
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
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		account varchar,
		token varchar
	);`); err != nil {
		log.Fatal(err)
	}
}

func createChannelTable(db *sql.DB) {
	if _, err := db.Exec(`create table channels(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title varchar,
		description text,
		identity varchar,
		token varchar
	);`); err != nil {
		log.Fatal(err)
	}
	/*
		http://localhost:8080/appcast/gotray/4cbd040533a2f43fc6691d773d510cda70f4126a
	*/
	if _, err := db.Exec(`insert into channels(title,description, identity, token) values (?,?,?,?)`, "GoTray", "Desc", "gotray", "4cbd040533a2f43fc6691d773d510cda70f4126a"); err != nil {
		panic(err)
	}
}

func createReleaseTable(db *sql.DB) {
	if _, err := db.Exec(`create table releases(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title varchar,
		desc text,
		releaseNotes text,
		pubDate datetime default current_timestamp,
		filename varchar,
		channel varchar,
		length integer,
		mimetype varchar,
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
	/*
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		item, err := CreateNewReleaseFromRequest(r)
		if err != nil {
			var msg = jsondata.Map{"error": err}
			msg.WriteTo(w)
		}
		_ = item
	*/
}

func CreateNewReleaseFromRequest(r *http.Request, channelIdentity string) (*appcast.Item, error) {
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
	releaseNotes := r.FormValue("releaseNotes")
	dsaSignature := r.FormValue("dsaSignature")

	h := sha1.New()
	h.Write([]byte(title))
	h.Write([]byte(version))
	h.Write([]byte(shortVersionString))
	h.Write(data)
	token := fmt.Sprintf("%x", h.Sum(nil))
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
	// newItem.SparkleReleaseNotesLink = releaseNotes

	result, err := db.Exec(`INSERT INTO releases 
		(channel, title, desc, pubDate, version, shortVersionString, releaseNotes, dsaSignature, filename, length, mimetype)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		channelIdentity,
		title,
		desc,
		pubDate,
		version,
		shortVersionString,
		releaseNotes,
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
	uploadPageRegExp := regexp.MustCompile("/release/upload/([^/]+)/([^/]+)")
	submatches := uploadPageRegExp.FindStringSubmatch(r.URL.Path)
	channelIdentity := submatches[1]
	channelToken := submatches[2]

	if channel := FindChannelByIdentity(channelIdentity, channelToken); channel != nil {
		if r.Method == "POST" {
			if _, err := CreateNewReleaseFromRequest(r, channelIdentity); err != nil {
				panic(err)
			}
		}

		t := templates.Lookup("upload.html")
		if t != nil {
			err := t.Execute(w, channel)
			if err != nil {
				panic(err)
			}
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Channel not found"))
		return
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

func AppcastXmlHandler(w http.ResponseWriter, r *http.Request) {
	var channelRegExp = regexp.MustCompile("/appcast/([^/]+)/([^/]+)")
	var submatches = channelRegExp.FindStringSubmatch(r.URL.Path)
	var channelIdentity = submatches[1]
	var channelToken = submatches[2]

	log.Println(r.URL)

	if channel := FindChannelByIdentity(channelIdentity, channelToken); channel != nil {
		w.Header().Set("Content-Type", "text/xml; charset=UTF-8")

		rows, err := QueryReleasesByChannel(channelIdentity)
		if err != nil {
			log.Fatal("Query failed:", err)
		}
		defer rows.Close()

		appcastRss := appcast.New()
		appcastRss.Channel.Title = channel.Title
		appcastRss.Channel.Description = channel.Description
		appcastRss.Channel.Link = "http://" + r.Host + "/appcast/" + channelIdentity + ".xml"
		// appcastRss.Channel.Language = channel.Language

		for rows.Next() {
			if item, err := ScanRowToAppcastItem(rows); err == nil {
				appcastRss.Channel.AddItem(item)
			}
		}
		appcastRss.WriteTo(w)
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Channel not found"))
	}
}

/*
For route: /download/gotray/{token}

/download/gotray/be24d1c54d0ba415b8897b02f0c38d89
*/
func DownloadFileHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.RequestURI, r.URL)
	path := r.URL.Path

	downloadRegExp := regexp.MustCompile("/download/([^/]+)/([^/]+)")
	// submatches := downloadRegExp.FindAllStringSubmatch(path)
	submatches := downloadRegExp.FindStringSubmatch(path)
	// log.Println(submatches)
	identity := submatches[1]
	token := submatches[2]

	_ = identity
	_ = token
	/*
		log.Println(r.URL.Opaque)
		log.Println(r.URL.Fragment)
	*/
}

func main() {
	db = ConnectDB(SQLITEDB)
	gatsby.SetupConnection(db, gatsby.DriverSqlite)
	defer db.Close()

	/*
		/release/download/{channel identity}/{release token}/{validation token}
		/release/upload/{channel identity}
		/release/new/{channel identity}
		/appcast/{channel identity}.xml
	*/
	http.HandleFunc("/release/download/", DownloadFileHandler)
	http.HandleFunc("/release/upload/", UploadPageHandler)
	http.HandleFunc("/release/new/", UploadReleaseHandler)
	http.HandleFunc("/appcast/", AppcastXmlHandler)
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	log.Println("Listening http://localhost:5000 ...")
	// http.HandleFunc("/upload", UploadNewReleaseHandler)
	if err := http.ListenAndServe(":5000", nil); err != nil {
		panic(err)
	}
}
