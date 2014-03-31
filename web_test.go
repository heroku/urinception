package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const base = "http://example.com"
const data = "hello, world"
const uri = base + "?uri=data%3Atext%2Fplain%3B+charset%3Dutf-8%3Bbase64%2CaGVsbG8sIHdvcmxk"

func TestHandleGet(t *testing.T) {
	req, _ := http.NewRequest("GET", uri, nil)
	res := httptest.NewRecorder()
	compileDatauriPattern()
	handleGet(res, req)

	assertEquals(t, res.Body.String(), data)
	assertEquals(t, res.Header().Get("Content-Type"), "text/plain; charset=utf-8")
	assertEquals(t, res.Code, 200)
}

func TestHandlePost(t *testing.T) {
	req, _ := http.NewRequest("POST", base, strings.NewReader(data))
	res := httptest.NewRecorder()
	handlePost(res, req)

	assertEquals(t, res.Body.String(), uri+"\n")
	assertEquals(t, res.Header().Get("Content-Type"), "text/uri-list; charset=utf-8")
	assertEquals(t, res.Code, 200)
}

func assertEquals(t *testing.T, actual interface{}, expected interface{}) {
	if actual != expected {
		t.Errorf("Actual: '%v'; Expected: '%v'", actual, expected)
	}
}
