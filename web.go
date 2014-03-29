package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"regexp"
)

func main() {
	r, _ := regexp.Compile("^data:(.*?)?(;base64)?,(.+)$")
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		dataurl := req.URL.RawQuery
		match := r.FindStringSubmatch(dataurl)
		if len(match) == 0 {
			fmt.Println("match.error.input:", dataurl)
			http.Error(res, "Path must be in RFC 2397 form", http.StatusBadRequest)
			return
		}

		contentType := match[1]
		isBase64 := match[2] != ""
		data := match[3]
		fmt.Println("request.type:", contentType, "request.base64:", isBase64)

		res.Header().Set("Content-Type", contentType)
		if isBase64 {
			decoded, err := base64.StdEncoding.DecodeString(data)
			if err != nil {
				fmt.Println("base64.decode.error:", err, "dataurl:", dataurl)
				http.Error(res, "Error decoding base64: "+err.Error(), http.StatusBadRequest)
				return
			}
			res.Write(decoded)
		} else {
			fmt.Fprintln(res, data)
		}
	})

	fmt.Println("listening:true port:", os.Getenv("PORT"))
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		panic(err)
	}
}
