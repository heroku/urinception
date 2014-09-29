package main

import (
	"encoding/base64"
	"fmt"
	"github.com/heroku/urinception"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

var datauriPattern *regexp.Regexp

func init() {
	compileDatauriPattern()
}

func main() {
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case "GET":
			handleGet(res, req)
		case "POST":
			handlePost(res, req)
		case "PUT":
			handlePost(res, req)
		default:
			http.Error(res, "Only GET and POST supported", http.StatusMethodNotAllowed)
		}
	})

	port := os.Getenv("PORT")
	log.Println("listening:true port:", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func compileDatauriPattern() {
	var err error
	datauriPattern, err = regexp.Compile("^data:(.*?)?(;base64)?,(.+)$")
	if err != nil {
		log.Fatal(err)
	}
}

func handleGet(res http.ResponseWriter, req *http.Request) {
	uri := req.URL.Query().Get("uri")
	match := datauriPattern.FindStringSubmatch(uri)
	if len(match) == 0 {
		log.Println("get.error.uri:", uri)
		http.Error(res, "Parameter 'uri' must be present and in RFC 2397 form", http.StatusBadRequest)
		return
	}
	contentType := match[1]
	isBase64 := match[2] != ""
	data := match[3]

	if req.URL.Query().Get("status") != "" {
		statusCode, err := strconv.Atoi(req.URL.Query().Get("status"))
		if err != nil {
			log.Println("get.error.base64.decode:", err)
			http.Error(res, "Error parsing status code to integer", http.StatusBadRequest)
			return
		}
		res.WriteHeader(statusCode)
	}

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
		fmt.Fprint(res, data)
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

	res.Header().Set("Content-Type", "text/uri-list; charset=utf-8")
	contentType := req.Header.Get("Content-Type")
	uri := urinception.CreateUri(scheme, req.Host, req.URL.Path, contentType, data)
	fmt.Fprintln(res, uri)
}
