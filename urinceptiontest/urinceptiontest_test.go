package urinceptiontest

import (
	"io"
	"net/http"
	"testing"
)

func TestNullUri(t *testing.T) {
	txt := "hello world"
	uri := StringUri(txt)

	res, err := http.Get(uri)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	obtained := string(bytes)
	expected := txt
	if obtained != expected {
		t.Errorf("Obtained: '%v'; Expected: '%v'", obtained, expected)
	}
}
