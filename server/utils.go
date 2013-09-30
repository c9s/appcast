package main

import "path"
import "mime"
import "os"

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
