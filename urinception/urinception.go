package urinception

import (
	"encoding/base64"
	"net/http"
	"net/url"
)

func CreateUri(scheme, host, path, contentType string, data []byte) string {
	if contentType == "" || contentType == "application/x-www-form-urlencoded" {
		contentType = http.DetectContentType(data)
	}

	base64 := base64.StdEncoding.EncodeToString(data)
	datauri := "data:" + contentType + ";base64," + base64
	return scheme + "://" + host + path + "?uri=" + url.QueryEscape(datauri)
}
