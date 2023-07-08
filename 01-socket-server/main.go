package main

import (
	"fmt"
	"net"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)

	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Received %d bytes.\n", n)

	request := string(buffer[:n])
	fmt.Println("------------------------------------------------")
	fmt.Println(request)
	fmt.Println("------------------------------------------------")

	response := `HTTP/1.1 200 OK
Content-Type: text/html; charset=utf-8

<!DOCTYPE html>
<html>
<head>
	<title>Socket Server</title>
</head>
<body>
  <h1>Socket Server</h1>
</body>
</html>
`
	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Println(err)
		return
	}
}

func main() {
	fmt.Println("Hello, World!")

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go handleConnection(conn)
	}
}
