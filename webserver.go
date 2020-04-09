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
	bytes   []byte
	modTime time.Time
}

var files map[string]file

func fileHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		w.WriteHeader(http.StatusNoContent)
	} else if foundFile, ok := files[r.URL.Path]; ok {
		//log.Infof("Serving %s", "/serve"+r.URL.Path)
		http.ServeContent(w, r, r.URL.Path, foundFile.modTime, bytes.NewReader(foundFile.bytes))
	} else {
		//log.Infof("Not found: %s", r.URL.Path)
		http.NotFound(w, r)
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
				//log.Errorf("Couldn't read file: %s: %s", path, err)
				return nil
			}

			files[filePath] = file{
				bytes:   fileBytes,
				modTime: info.ModTime(),
			}
		}
		return nil
	})

	http.HandleFunc("/", fileHandler)
	log.Info("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
