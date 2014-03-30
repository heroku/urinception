package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
)

var dataUrlPattern *regexp.Regexp

func handleGet(res http.ResponseWriter, req *http.Request) {
	dataUrl := req.URL.Query().Get("url")
	match := dataUrlPattern.FindStringSubmatch(dataUrl)
	if len(match) == 0 {
		log.Println("get.error.url:", dataUrl)
		http.Error(res, "Parameter 'url' must be present and in RFC 2397 form", http.StatusBadRequest)
		return
	}

	contentType := match[1]
	isBase64 := match[2] != ""
	data := match[3]

	res.Header().Set("Content-Type", contentType)
	if isBase64 {
		decoded, err := base64.StdEncoding.DecodeString(data)
		if err != nil {
			log.Println("get.error.base64.decode:", err)
			http.Error(res, "Error decoding base64: "+err.Error(), http.StatusBadRequest)
			return
		}
		res.Write(decoded)
	} else {
		fmt.Fprintln(res, data)
	}
}

func handlePost(res http.ResponseWriter, req *http.Request) {
	scheme := "http"
	if req.TLS != nil || req.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}

	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println("post.error.body:", err)
		http.Error(res, "Error reading request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	contentType := req.Header.Get("Content-Type")
	if contentType == "" || contentType == "application/x-www-form-urlencoded" {
		contentType = http.DetectContentType(data)
	}

	base64 := base64.StdEncoding.EncodeToString(data)
	dataUrl := "data:" + contentType + ";base64," + base64
	url := scheme + "://" + req.Host + req.URL.Path + "?url=" + url.QueryEscape(dataUrl)

	fmt.Fprint(res, url)
}

func main() {
	dataUrlPattern, _ = regexp.Compile("^data:(.*?)?(;base64)?,(.+)$")

	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case "GET":
			handleGet(res, req)
		case "POST":
			handlePost(res, req)
		default:
			http.Error(res, "Only GET and POST supported", http.StatusMethodNotAllowed)
		}
	})

	log.Println("listening:true port:", os.Getenv("PORT"))
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		panic(err)
	}
}
