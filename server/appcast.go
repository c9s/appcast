package main

import (
	"crypto/sha1"
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

const BIND = ":5000"
const BASEURL = "http://localhost:5000"
const UPLOAD_DIR = "uploads"
const SQLITEDB = "appcast.db"

var ErrFileIsRequired = errors.New("file is required.")
var ErrReleaseInsertFailed = errors.New("release insert failed.")

var db *sql.DB
var templates = template.Must(template.ParseFiles("templates/upload.html"))

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
		(channel, title, desc, pubDate, version, shortVersionString, releaseNotes, dsaSignature, filename, length, mimetype, token)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		channelIdentity,
		title,
		desc,
		pubDate,
		version,
		shortVersionString,
		releaseNotes,
		dsaSignature,
		fileReader.Filename,
		length, mimetype, token)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	log.Println("New Release Uploaded", id, title, version, shortVersionString, desc, pubDate, dsaSignature, length, mimetype)
	return &newItem, nil
}

func UploadPageHandler(w http.ResponseWriter, r *http.Request) {
	uploadPageRegExp := regexp.MustCompile("/release/upload/([^/]+)/([^/]+)")
	submatches := uploadPageRegExp.FindStringSubmatch(r.URL.Path)
	if len(submatches) != 3 {
		ForbiddenHandler(w, r)
		return
	}

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

func ScanRowToAppcastItem(rows *sql.Rows, channelIdentity, channelToken string) (*appcast.Item, error) {
	var title, desc, version, shortVersionString, filename, mimetype, dsaSignature, token string
	var pubDate time.Time
	var length int64
	var err = rows.Scan(&title, &desc, &pubDate, &version, &shortVersionString, &filename, &mimetype, &length, &dsaSignature, &token)
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
	item.Enclosure.URL = BASEURL + "/release/download/" + channelIdentity + "/" + channelToken + "/" + token
	return &item, nil
}

func ForbiddenHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte("Forbidden"))
	return
}

/*
For route: /download/gotray/{token}

/download/gotray/be24d1c54d0ba415b8897b02f0c38d89
*/
func DownloadFileHandler(w http.ResponseWriter, r *http.Request) {
	downloadRegExp := regexp.MustCompile("/release/download/([^/]+)/([^/]+)/([^/]+)")
	submatches := downloadRegExp.FindStringSubmatch(r.URL.Path)

	if len(submatches) != 4 {
		ForbiddenHandler(w, r)
		return
	}

	channelIdentity := submatches[1]
	channelToken := submatches[2]
	releaseToken := submatches[3]

	if channel := FindChannelByIdentity(channelIdentity, channelToken); channel != nil {
		if release := LoadReleaseByChannelAndToken(channelIdentity, releaseToken); release != nil {
			log.Println(r.URL.Path, release.Filename, release.Mimetype)
			w.Header().Set("Content-Type", release.Mimetype)
			w.Header().Set("Content-Disposition", "inline; filename=\""+release.Filename+"\"")

			data, err := ioutil.ReadFile(path.Join(UPLOAD_DIR, release.Filename))
			if err != nil {
				panic(err)
			}
			w.Write(data)
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Release not found"))
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Channel not found"))
	}
}

func main() {
	db = ConnectDB(SQLITEDB)
	gatsby.SetupConnection(db, gatsby.DriverSqlite)
	defer db.Close()

	/*
		/release/download/{channel identity}/{channel token}/{release token}
		/release/upload/{channel identity}/{channel token}
		/release/new/{channel identity}
		/appcast/{channel identity}.xml
	*/
	http.HandleFunc("/release/download/", DownloadFileHandler)
	http.HandleFunc("/release/upload/", UploadPageHandler)
	http.HandleFunc("/release/new/", UploadReleaseHandler)
	http.HandleFunc("/appcast/", AppcastXmlHandler)
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	log.Println("Listening  " + BASEURL + " ...")
	if err := http.ListenAndServe(BIND, nil); err != nil {
		panic(err)
	}
}
