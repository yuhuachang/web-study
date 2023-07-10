package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func readResource(path string) (string, []byte, error) {
	filePath := path
	if strings.HasPrefix(path, "/") {
		filePath = path[1:]
	}
	if filePath == "" {
		filePath = "index.html"
	}
	fmt.Println("File path:", filePath)

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return "", nil, err
	}

	contentType := http.DetectContentType(fileContent)
	fmt.Println("Content type:", contentType)

	return contentType, fileContent, nil
}

func handler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Method: ", r.Method)
	fmt.Println("URL: ", r.URL)
	fmt.Println("Proto: ", r.Proto)
	for k, v := range r.Header {
		fmt.Printf("  %s = %s\n", k, v[0])
	}

	fmt.Println("------------------------------------------------")
	fmt.Println(r.Body)
	fmt.Println("------------------------------------------------")

	// type URL struct {
	// 	Scheme      string
	// 	Opaque      string    // encoded opaque data
	// 	User        *Userinfo // username and password information
	// 	Host        string    // host or host:port
	// 	Path        string    // path (relative paths may omit leading slash)
	// 	RawPath     string    // encoded path hint (see EscapedPath method)
	// 	OmitHost    bool      // do not emit empty host (authority)
	// 	ForceQuery  bool      // append a query ('?') even if RawQuery is empty
	// 	RawQuery    string    // encoded query values, without '?'
	// 	Fragment    string    // fragment for references, without '#'
	// 	RawFragment string    // encoded fragment hint (see EscapedFragment method)
	// }

	contentType, fileContent, err := readResource(r.URL.Path)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", contentType)
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(fileContent)))
		n, err := w.Write(fileContent)
		if err == nil {
			fmt.Printf("Write %d bytes.\n", n)
		} else {
			fmt.Println(err)
			return
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func main() {
	certFile := "localhost.pem"
	keyFile := "localhost-key.pem"

	server := http.Server{
		Addr:    ":8080",
		Handler: http.HandlerFunc(handler),
		TLSConfig: &tls.Config{
			MinVersion:   tls.VersionTLS13,
			Certificates: []tls.Certificate{},
		},
	}
	err := server.ListenAndServeTLS(certFile, keyFile)

	if err != nil {
		fmt.Println("Error starting http server:", err)
	}
}
