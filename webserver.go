package main

import (
	"bytes"
	"code.cloudfoundry.org/bytefmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type file struct {
	bytesReader *bytes.Reader
	modTime     time.Time
}

var files map[string]file

func fileHandler(w http.ResponseWriter, r *http.Request) {
	if file, found := files[r.URL.Path]; found {
		http.ServeContent(w, r, r.URL.Path, file.modTime, file.bytesReader)
	} else {
		if r.URL.Path == "/" {
			w.WriteHeader(http.StatusNoContent)
		} else {
			http.NotFound(w, r)
		}
	}
}

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	files = make(map[string]file)
	_ = filepath.Walk("/serve", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() == false {
			filePath := strings.ReplaceAll(path, "/serve", "")
			log.Infof("Found file: %s, size: %s, modTime: %s", filePath, bytefmt.ByteSize(uint64(info.Size())), info.ModTime())
			fileBytes, err := ioutil.ReadFile(path)
			if err != nil {
				log.Fatalf("Couldn't read file: %s: %s", path, err)
			}

			files[filePath] = file{
				bytesReader: bytes.NewReader(fileBytes),
				modTime:     info.ModTime(),
			}
		}
		return nil
	})

	if root, isSet := os.LookupEnv("ROOT"); isSet {
		if file, exists := files[root]; exists {
			log.Infof("Parameter ROOT found. Additionally routing %s under /", root)
			files["/"] = file
		} else {
			log.Fatalf("Parameter ROOT (%s) found, but no corresponding file found.", root)
		}
	}

	http.HandleFunc("/", fileHandler)
	log.Info("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
