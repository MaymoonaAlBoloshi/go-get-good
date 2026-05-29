package main

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"slices"
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
		response.Write(connection, response.Response{
			StatusCode: 200,
		})

	} else if pathParts[1] == "echo" {
		echoStr := pathParts[2]

		acceptEncoding := headers.Get(headerLines, headers.AcceptEncoding)
		encodings := []string{}

		if strings.Contains(acceptEncoding, " ") {
			encodings = strings.Split(acceptEncoding, ",")
			for encoding := 1; encoding < len(encodings); encoding++ {
				strings.TrimSpace(encodings[encoding])
			}
		}

		res := response.Response{
			StatusCode:  200,
			Body:        echoStr,
			ContentType: response.Text,
		}

		if slices.Contains(encodings, "gzip") {
			res.ContentEncoding = response.Gzip
		}
		response.Write(connection, res)

	} else if pathParts[1] == "user-agent" {
		userAgent := headers.Get(headerLines, headers.UserAgent)
		response.Write(connection, response.Response{
			StatusCode:  200,
			Body:        userAgent,
			ContentType: response.Text,
		})

	} else if pathParts[1] == "files" {
		fileName := pathParts[2]
		if method == "POST" {
			contentLength := headers.Get(headerLines, headers.ContentLength)
			parts := strings.SplitN(req, CRLF+CRLF, 2)

			if len(parts) < 2 {
				response.Write(connection, response.Response{
					StatusCode: 400,
				})
				return
			}

			body := parts[1]

			size, err := strconv.Atoi(contentLength)
			if err != nil || size < 0 {
				response.Write(connection, response.Response{
					StatusCode: 400,
				})
				return
			}

			content := []byte(body)
			if len(content) > size {
				content = content[:size]
			}

			filePath, isFiles := strings.CutPrefix(path, "/files/")
			if isFiles {
				os.WriteFile(*dirFlag+string(filePath), content, 0644)
				response.Write(connection, response.Response{
					StatusCode:  201,
					Body:        string(content),
					ContentType: response.File,
				})
				// response.Write(connection, 200, string(content), response.File)
				return
			}

		} else {
			content, err := os.ReadFile(filepath.Join(*dirFlag, fileName))

			if err != nil {
				if errors.Is(err, os.ErrNotExist) {
					response.Write(connection, response.Response{
						StatusCode: 404,
					})
				} else {
					response.Write(connection, response.Response{
						StatusCode: 500,
					})
				}
			} else {
				response.Write(connection, response.Response{
					StatusCode:  200,
					Body:        string(content),
					ContentType: response.File,
				})
			}
		}

	} else {
		response.Write(connection, response.Response{
			StatusCode: 404,
		})

	}
}
