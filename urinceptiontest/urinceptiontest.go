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

func EmptyUri() string {
	return BytesUri([]byte{})
}

func StringUri(data string) string {
	return BytesUri([]byte(data))
}

func FileUri(filename string) string {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return WithStatus("", http.StatusInternalServerError)
	} else {
		return BytesUri(bytes)
	}
}

func BytesUri(data []byte) string {
	return urinception.CreateUri(scheme, host+":"+port, "/", http.DetectContentType(data), data)
}

func WithStatus(uri string, statusCode int) string {
	if uri == "" {
		uri = EmptyUri()
	}
	return fmt.Sprintf("%s&status=%d", uri, statusCode)
}

func WithPath(uri string, path string) string {
	if uri == "" {
		uri = EmptyUri()
	}
	parts := strings.SplitN(uri, "/?", 2)
	return parts[0] + path + parts[1]
}
