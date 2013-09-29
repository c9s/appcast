package uploader

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
)

const MIN_FILE_SIZE = 1 // bytes
const MAX_FILE_SIZE = 5000000

type FileInfo struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Size  int64  `json:"size"`
	Error string `json:"error,omitempty"`
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func handleUpload(r *http.Request, p *multipart.Part) (fi *FileInfo) {
	fi = &FileInfo{
		Name: p.FileName(),
		Type: p.Header.Get("Content-Type"),
	}
	// recover function
	defer func() {
		if rec := recover(); rec != nil {
			log.Println(rec)
			fi.Error = rec.(error).Error()
		}
	}()

	// XXX: do copy here.

	/*
			lr := &io.LimitedReader{R: p, N: MAX_FILE_SIZE + 1}
			context := appengine.NewContext(r)
			w, err := blobstore.Create(context, fi.Type)
			defer func() {
				w.Close()
				fi.Size = MAX_FILE_SIZE + 1 - lr.N
				fi.Key, err = w.Key()
				check(err)
				if !fi.ValidateSize() {
					err := blobstore.Delete(context, fi.Key)
					check(err)
					return
				}
				delayedDelete(context, fi)
				fi.CreateUrls(r, context)
			}()
			check(err)
		_, err = io.Copy(w, lr)
	*/
	return
}

func getFormValue(p *multipart.Part) string {
	var b bytes.Buffer
	io.CopyN(&b, p, int64(1<<20)) // Copy max: 1 MiB
	return b.String()
}

func handleUploads(r *http.Request) (fileInfos []*FileInfo) {
	fileInfos = make([]*FileInfo, 0)
	mr, err := r.MultipartReader()
	check(err)
	r.Form, err = url.ParseQuery(r.URL.RawQuery)
	check(err)
	part, err := mr.NextPart()
	for err == nil {
		if name := part.FormName(); name != "" {
			if part.FileName() != "" {
				fileInfos = append(fileInfos, handleUpload(r, part))
			} else {
				r.Form[name] = append(r.Form[name], getFormValue(part))
			}
		}
		part, err = mr.NextPart()
	}
	return
}
