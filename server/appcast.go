package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"text/template"
)
import _ "github.com/c9s/appcast"
import _ "github.com/c9s/appcast/server/uploader"

const UPLOAD_DIR = "uploads"

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
	}

	templates, err := template.ParseFiles("templates/upload.html")
	if err != nil {
		panic(err)
	}
	t := templates.Lookup("upload.html")
	if t != nil {
		err := t.Execute(w, nil)
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
	http.HandleFunc("/upload", UploadPageHandler)
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	log.Println("Listening http://localhsot:8080 ...")
	// http.HandleFunc("/upload", UploadNewReleaseHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
