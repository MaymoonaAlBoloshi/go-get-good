package main

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
)

// Ensures gofmt doesn't remove the "net" and "os" imports above (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

const CRLF = "\r\n"

var dirFlag = flag.String("directory", "", "")

func main() {
	flag.Parse()
	fmt.Println("Logs from your program will appear here!")

	app, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		connection, err := app.Accept()

		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleConnection(connection)

	}
}

func handleConnection(connection net.Conn) {
	buf := make([]byte, 1024)
	n, err := connection.Read(buf)

	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
	}

	req := string(buf[:n])
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
		for line := 1; line < len(lines); line++ {
			keyVal := strings.Split(lines[line], ":")
			if keyVal[0] == "User-Agent" {
				userAgent := strings.TrimSpace(keyVal[1])
				res = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(userAgent), userAgent)
			}
			fmt.Println(keyVal)
		}

	} else if pathParts[1] == "file" {
		fileName := pathParts[2]
		content, err := os.ReadFile(filepath.Join(*dirFlag, fileName))

		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				res = "HTTP/1.1 404 Not Found\r\n\r\n"
			} else {
				res = "HTTP/1.1 500 Internal Server Error\r\n\r\n"
			}

			res = "HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: " + fmt.Sprintf("%d", len(content)) + "\r\n\r\n" + string(content)
		}
	} else {
		res = "HTTP/1.1 404 Not Found\r\n\r\n"

	}
	fmt.Println(res)

	connection.Write([]byte(res))

	connection.Close()
}
