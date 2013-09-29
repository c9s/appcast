package main

import (
	// "github.com/c9s/appcast"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
)

func UploadNewReleaseHandler(w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Println(err)
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
	}
	if err = ioutil.WriteFile(handler.Filename, data, 0777); err != nil {
		log.Println(err)
	}
}

func UploadPageHandler(w http.ResponseWriter, r *http.Request) {
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
	// http.HandleFunc("/upload", UploadNewReleaseHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
