package main

import (
	"encoding/json"
	"flag"
	"net/http"
	"os"
	"regexp"
)

type assets struct {
	Url  string `json:"url"`
	Name string `json:"name"`
}

type rel struct {
	TagName    string   `json:"tag_name"`
	PreRelease bool     `json:"prerelease"`
	Assets     []assets `json:"assets"`
}

var listen string
var url string
var dir string

func main() {
	flag.StringVar(&listen, "listen", "0.0.0.0:8080", "http listen")
	flag.StringVar(&url, "url", "http://127.0.0.1:8080", "download url")
	flag.StringVar(&dir, "dir", "./dl", "file dir")
	flag.Parse()
	run()
}

func run() {
	http.Handle("/dl/", http.StripPrefix("/dl/", http.FileServer(http.Dir(dir))))
	http.HandleFunc("/meta.json", func(writer http.ResponseWriter, request *http.Request) {
		body := getJson()
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = writer.Write(body)
	})
	_ = http.ListenAndServe(listen, nil)
}

func getJson() []byte {

	list := map[string][]string{}

	dirEntry, err := os.ReadDir(dir)
	if err != nil {
		return []byte("[]")
	}
	for _, file := range dirEntry {
		if file.IsDir() == false {

			var ver string

			name := file.Name()

			re, _ := regexp.Compile(`syncthing-(.*)-(.*)-(v.*)(\.tar\.gz|\.zip)`)
			stList := re.FindStringSubmatch(name)
			if len(stList) == 5 {
				ver = stList[3]
			} else {
				continue
			}

			if _, ok := list[ver]; ok {
				list[ver] = append(list[ver], name)
			} else {
				list[ver] = []string{name}
			}

		}
	}

	r := make([]rel, 0)

	for ver, item := range list {

		itemLen := len(item)
		as := make([]assets, itemLen+2)
		as[0].Url = url + "/dl/sha1sum.txt.asc"
		as[0].Name = "sha1sum.txt.asc"
		as[1].Url = url + "/dl/sha256sum.txt.asc"
		as[1].Name = "sha256sum.txt.asc"

		for j, fileName := range item {
			as[j+2].Url = url + "/dl/" + fileName
			as[j+2].Name = fileName
		}

		r = append(r, rel{Assets: as, PreRelease: false, TagName: ver})
	}
	b, _ := json.Marshal(&r)

	return b
}
