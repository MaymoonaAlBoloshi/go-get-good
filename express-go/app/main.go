package main

import (
	"bytes"
	"compress/gzip"
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
	"github.com/codecrafters-io/http-server-starter-go/app/request"
	"github.com/codecrafters-io/http-server-starter-go/app/response"
)

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

	req, ok := request.Parse(string(buf[:n]))
	if !ok {
		response.Write(connection, response.Response{
			StatusCode: 400,
		})
		return
	}

	fmt.Println(req.Method)

	if req.Path == "/" {
		response.Write(connection, response.Response{
			StatusCode: 200,
		})

	} else if req.PathParts[1] == "echo" {
		echoStr := req.PathParts[2]

		acceptEncoding := headers.Get(req.Headers, headers.AcceptEncoding)
		encodings := []string{}

		if strings.Contains(acceptEncoding, " ") {
			encodings = strings.Split(acceptEncoding, ",")
			for encoding := 0; encoding < len(encodings); encoding++ {
				encodings[encoding] = strings.TrimSpace(encodings[encoding])
			}
		}

		res := response.Response{
			StatusCode:  200,
			Body:        echoStr,
			ContentType: response.Text,
		}

		if slices.Contains(encodings, "gzip") || acceptEncoding == "gzip" {

			compressed, err := compressGzip([]byte(echoStr))
			if err != nil {
				response.Write(connection, response.Response{
					StatusCode: 500,
				})
				return
			}
			res.Body = string(compressed)
			res.ContentEncoding = response.Gzip

		}
		response.Write(connection, res)

	} else if req.PathParts[1] == "user-agent" {
		userAgent := headers.Get(req.Headers, headers.UserAgent)
		response.Write(connection, response.Response{
			StatusCode:  200,
			Body:        userAgent,
			ContentType: response.Text,
		})

	} else if req.PathParts[1] == "files" {
		fileName := req.PathParts[2]
		if req.Method == "POST" {
			contentLength := headers.Get(req.Headers, headers.ContentLength)
			size, err := strconv.Atoi(contentLength)
			if err != nil || size < 0 {
				response.Write(connection, response.Response{
					StatusCode: 400,
				})
				return
			}

			content := []byte(req.Body)
			if len(content) > size {
				content = content[:size]
			}

			filePath, isFiles := strings.CutPrefix(req.Path, "/files/")
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

func compressGzip(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)
	if _, err := writer.Write(data); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
