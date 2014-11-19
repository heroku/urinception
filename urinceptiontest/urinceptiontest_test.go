package urinceptiontest

import (
	"io/ioutil"
	"net/http"
	"testing"
)

func TestNullUri(t *testing.T) {
	txt := "hello world"
	uri := StringUri(txt)

	res, _ := http.Get(uri)
	defer res.Body.Close()
	bytes, _ := ioutil.ReadAll(res.Body)

	obtained := string(bytes)
	expected := txt
	if obtained != expected {
		t.Errorf("Obtained: '%v'; Expected: '%v'", obtained, expected)
	}
}
