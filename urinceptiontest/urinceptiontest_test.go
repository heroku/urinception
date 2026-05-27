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
		t.Fatalf("http.Get failed: %v", err)
	}
	defer res.Body.Close()
	bytes, _ := io.ReadAll(res.Body)

	obtained := string(bytes)
	expected := txt
	if obtained != expected {
		t.Errorf("Obtained: '%v'; Expected: '%v'", obtained, expected)
	}
}
