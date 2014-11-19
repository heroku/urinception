/*
	Package `urinceptiontest` provides other Go applications
	a way to run URInception locally for use in tests.
	Importing `urinceptiontest` automatically starts an HTTP
	server on a random port on the local machine. Several
	methods are provided to create URI fixtures that can be
	be passed the system under test which will call the
	local server.

	For example, if a test that wants to assert `http.Get` works,
	`urinceptiontest` can be used to create a URI that will return
	a given response body:

		import "github.com/heroku/urinception/urinceptiontest"

		// create the URI fixture
		txt := "hello world"
		uri := urinceptiontest.StringUri(txt)

		// pass the URI to the system under test
		res, _ := http.Get(uri)
		defer res.Body.Close()
		bytes, _ := ioutil.ReadAll(res.Body)

		// assert the result
		obtained := string(bytes)
		expected := txt
		if obtained != expected {
			t.Errorf("Obtained: '%v'; Expected: '%v'", obtained, expected)
		}

*/
package urinceptiontest

import (
	"fmt"
	"github.com/heroku/urinception"
	web "github.com/heroku/urinception/cmd/urinception-web"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"sync"
)

var (
	scheme = "http"
	host   = "localhost"
	port   = "0" // automatically choose an available port
	start  = &sync.Once{}
)

func init() {
	start.Do(func() {
		setAvailablePort()
		go web.Start(port)
	})
}

func setAvailablePort() {
	addr, err := net.ResolveTCPAddr("tcp", ":"+port)
	if err != nil {
		panic(err)
	}

	l, err := net.ListenTCP("tcp", addr)
	if l != nil {
		defer l.Close()
	}
	if err != nil {
		panic(err)
	}

	_, port, err = net.SplitHostPort(l.Addr().String())
	if err != nil {
		panic(err)
	}
}

// Create a URI returning a null character
// Not too useful on its own, but can be used as a dummy URI
// or with a special status or path
func NullUri() string {
	return StringUri("\x00")
}

// Create a URI returning the string provided
func StringUri(data string) string {
	return BytesUri([]byte(data))
}

// Create a URI returning the contents of the file provided
func FileUri(filename string) string {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return WithStatus("", http.StatusInternalServerError)
	} else {
		return BytesUri(bytes)
	}
}

// Create a URI returning the bytes provided
func BytesUri(data []byte) string {
	return urinception.CreateUri(scheme, host+":"+port, "/", http.DetectContentType(data), data)
}

// Transform an existing URI to return the HTTP status provided
// If no URI is provided, a dummy value will be used
func WithStatus(uri string, statusCode int) string {
	if uri == "" {
		uri = NullUri()
	}
	return fmt.Sprintf("%s&status=%d", uri, statusCode)
}

// Transform an existing URI to include the path provided
// If no URI is provided, a dummy value will be used
func WithPath(uri string, path string) string {
	if uri == "" {
		uri = NullUri()
	}
	parts := strings.SplitN(uri, "/?", 2)
	return parts[0] + path + parts[1]
}
