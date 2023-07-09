package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func readRawRequest(conn net.Conn) (string, error) {

	// Read the raw request in a small trunk at a time.
	rawRequest := []byte{}
	bufferSize := 128
	buffer := make([]byte, bufferSize)
	totalReceived := 0
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			return "", err
		}
		rawRequest = append(rawRequest, buffer[:n]...)
		totalReceived += n
		fmt.Printf("Read %d bytes.\n", n)
		if n < bufferSize {
			break
		}
	}
	fmt.Printf("Received %d bytes.\n", totalReceived)
	return string(rawRequest), nil
}

func parseRequest(request string) (string, string, string, map[string]string) {
	requestLines := strings.Split(request, "\r\n")
	x := strings.Split(requestLines[0], " ")
	method := x[0]
	path := x[1]
	version := x[2]
	headers := map[string]string{}
	for i := 1; i < len(requestLines); i++ {
		headerLine := strings.TrimSpace(requestLines[i])
		x = strings.Split(headerLine, ": ")
		if len(x) == 2 {
			headers[x[0]] = x[1]
		}
		if requestLines[i] == "" {
			break
		}
	}
	return method, path, version, headers
}

func readResource(path string) (string, []byte, error) {
	filePath := path
	if strings.HasPrefix(path, "/") {
		filePath = path[1:]
	}
	var contentType string
	if filePath == "" {
		filePath = "index.html"
	}
	if strings.HasSuffix(filePath, ".html") {
		contentType = "text/html; charset=utf-8"
	}
	if strings.HasSuffix(filePath, ".css") {
		contentType = "text/css; charset=utf-8"
	}
	if strings.HasSuffix(filePath, ".js") {
		contentType = "application/javascript; charset=utf-8"
	}
	if strings.HasSuffix(filePath, ".ico") {
		contentType = "image/x-icon"
	}
	if strings.HasSuffix(filePath, ".png") {
		contentType = "image/png"
	}
	if strings.HasSuffix(filePath, ".jpeg") {
		contentType = "image/jpeg"
	}
	fmt.Println("File path:", filePath)
	fmt.Println("Content type:", contentType)

	file, err := os.Open(filePath)
	if err != nil {
		return "", nil, err
	}
	defer file.Close()

	// Read file content.
	fileContent := []byte{}
	bufferSize := 128
	buffer := make([]byte, bufferSize)
	totalRead := 0
	for {
		n, err := file.Read(buffer)
		if err != nil {
			return "", nil, err
		}
		fileContent = append(fileContent, buffer[:n]...)
		totalRead += n
		fmt.Printf("Read %d bytes.\n", n)
		if n < bufferSize {
			break
		}
	}
	fmt.Printf("Total read %d bytes.\n", totalRead)

	return contentType, fileContent, nil
}

func writeResponse(conn net.Conn, statusCode int, contentType string, content []byte) error {
	fmt.Println("Send response.")

	httpStatusCodes := map[int]string{
		200: "OK",
		404: "Not Found",
	}

	_, err := conn.Write([]byte(fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, httpStatusCodes[statusCode])))
	if err != nil {
		return err
	}
	if contentType != "" {
		_, err = conn.Write([]byte(fmt.Sprintf("Content-Type: %s\r\n", contentType)))
		if err != nil {
			return err
		}
	}
	if len(content) > 0 {
		_, err = conn.Write([]byte(fmt.Sprintf("Content-Length: %d\r\n", content)))
		if err != nil {
			return err
		}
	}
	_, err = conn.Write([]byte("\r\n"))
	if err != nil {
		return err
	}
	_, err = conn.Write(content)
	if err != nil {
		return err
	}

	return nil
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Println("New connection accepted.")

	// Read the raw request.
	request, err := readRawRequest(conn)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("------------------------------------------------")
	fmt.Println(request)
	fmt.Println("------------------------------------------------")

	// Parse the request.
	method, path, version, headers := parseRequest(request)
	fmt.Println("Method:", method)
	fmt.Println("Path:", path)
	fmt.Println("Version:", version)
	fmt.Println("Accept:", headers["Accept"])

	// Read the resource.
	contentType, fileContent, err := readResource(path)
	if err == nil {
		err = writeResponse(conn, 200, contentType, fileContent)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		fmt.Println(err)
		err = writeResponse(conn, 404, "text/html; charset=utf-8", []byte("404 Not Found"))
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func main() {
	fmt.Println("Server is listening on port 8080...")

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println("Have a connection...")
		go handleConnection(conn)
		fmt.Println("Connection handled.")
	}
}
