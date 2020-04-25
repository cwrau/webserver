package main

import (
	"bytes"
	"code.cloudfoundry.org/bytefmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	path2 "path"
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
		_, _ = file.bytesReader.Seek(0, io.SeekStart)
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
	indexIsSet := strings.EqualFold(os.Getenv("INDEX"), "true")
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
			if indexIsSet && strings.HasSuffix(filePath, "index.html") {
				files[path2.Dir(filePath)+"/"] = files[filePath]
			}
		}
		return nil
	})

	usedMemory := 0
	for _, file := range files {
		usedMemory += file.bytesReader.Len()
	}
	log.Infof("Using %s memory", bytefmt.ByteSize(uint64(usedMemory)))

	if indexIsSet {
		log.Info("Parameter 'INDEX' set to 'true'. Routing all '*/' folders to their corresponding '*/index.html', if they exist")
	}

	http.HandleFunc("/", fileHandler)
	log.Info("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
