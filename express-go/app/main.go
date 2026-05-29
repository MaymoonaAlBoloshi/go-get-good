package main

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/app/headers"
	"github.com/codecrafters-io/http-server-starter-go/app/response"
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
	defer connection.Close()

	buf := make([]byte, 1024)
	n, err := connection.Read(buf)

	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		return
	}

	req := string(buf[:n])
	lines := strings.Split(req, CRLF)

	path := strings.Split(lines[0], " ")[1]

	method := strings.Split(lines[0], " ")[0]

	pathParts := strings.Split(path, "/")

	headerLines := headers.Parse(lines)

	fmt.Println(method)

	if path == "/" {
		response.Write(connection, 200, "", "")

	} else if pathParts[1] == "echo" {
		echoStr := pathParts[2]
		response.Write(connection, 200, echoStr, response.Text)

	} else if pathParts[1] == "user-agent" {
		userAgent := headers.Get(headerLines, headers.UserAgent)
		response.Write(connection, 200, userAgent, response.Text)

	} else if pathParts[1] == "files" {
		fileName := pathParts[2]
		if method == "POST" {
			contentLength := headers.Get(headerLines, headers.ContentLength)

			size, err := strconv.Atoi(contentLength)
			if err != nil || size < 0 {
				response.Write(connection, 400, "", "")
				return
			}
			content := make([]byte, size)

			filePath, isFiles := strings.CutPrefix(path, "/files/")
			if isFiles {
				os.WriteFile(*dirFlag+string(filePath), content, 0644)
				response.Write(connection, 201, string(content), response.File)
				// response.Write(connection, 200, string(content), response.File)
				return
			}

		} else {
			content, err := os.ReadFile(filepath.Join(*dirFlag, fileName))

			if err != nil {
				if errors.Is(err, os.ErrNotExist) {
					response.Write(connection, 404, "", "")
				} else {
					response.Write(connection, 500, "", "")
				}
			} else {
				response.Write(connection, 200, string(content), response.File)
			}
		}

	} else {
		response.Write(connection, 404, "", "")

	}
}
