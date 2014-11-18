package urinceptiontest

import (
	"fmt"
	"github.com/heroku/urinception"
	web "github.com/heroku/urinception/cmd/urinception-web"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
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

	portRegex := regexp.MustCompile(".*?(\\d+)$")
	port = portRegex.FindStringSubmatch(l.Addr().String())[1]
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
