package main

import (
    "fmt"
    "net/http"
    "os"
    "regexp"
)

func main() {
    r, _ := regexp.Compile("data:(.*?)?(;base64)?,(.+)")
    http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
        path := req.URL.Path
        match := r.FindStringSubmatch(path)
        if len(match) == 0 {
          http.Error(res, "Must be in RFC 2397 form", http.StatusBadRequest)
          return
        }
        mimeType := match[1]
        isBase64 := match[2] != ""
        data := match[3]
        fmt.Fprintln(res, match)
        fmt.Fprintln(res, mimeType)
        fmt.Fprintln(res, isBase64)
        fmt.Fprintln(res, data)
    })

    fmt.Println("listening...")
    err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
    if err != nil {
        panic(err)
    }
}
