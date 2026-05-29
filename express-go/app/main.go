package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "net" and "os" imports above (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

const CRLF = "\r\n"

func main() {
	fmt.Println("Logs from your program will appear here!")

	app, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	connection, err := app.Accept()

	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	buf := make([]byte, 1024)
	_, err = connection.Read(buf)

	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
	}

	req := string(buf)
	lines := strings.Split(req, CRLF)
	path := strings.Split(lines[0], " ")[1]
	pathParts := strings.Split(path, "/")
	fmt.Println(path)

	var res string
	if path == "/" {
		res = "HTTP/1.1 200 OK\r\n\r\n"
	} else if pathParts[1] == "echo" {
		echoStr := pathParts[2]
		res = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(echoStr), echoStr)

	} else if pathParts[1] == "user-agent" {

		fmt.Println("correct spot")
		// var userAgent string

		for l := 1; l < len(lines); l++ {
			keyVal := strings.Split(lines[l], ":")
			if keyVal[0] == "User-Agent" {

				res = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(keyVal[1]), keyVal[1])
			}
			fmt.Println(keyVal)
		}
	} else {
		res = "HTTP/1.1 404 Not Found\r\n\r\n"
	}
	fmt.Println(res)

	connection.Write([]byte(res))

}
