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
	"github.com/codecrafters-io/http-server-starter-go/app/http"
	"github.com/codecrafters-io/http-server-starter-go/app/request"
	"github.com/codecrafters-io/http-server-starter-go/app/response"
)

var dirFlag = flag.String("directory", "", "")

func main() {
	flag.Parse()
	fmt.Println("Logs from your program will appear here!")

	server := http.New()
	server.Get("/", handleRoot)
	server.Get("/echo/*", handleEcho)
	server.Get("/user-agent", handleUserAgent)
	server.Get("/files/*", handleGetFile)
	server.Post("/files/*", handlePostFile)

	listener, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		connection, err := listener.Accept()

		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleConnection(connection, server)

	}
}

func handleConnection(connection net.Conn, server *http.App) {
	defer connection.Close()

	for {

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

		connectionHeader := headers.Get(req.Headers, headers.Connection)
		shouldClose := strings.EqualFold(connectionHeader, "close")

		res := server.Handle(req)
		if shouldClose {
			res.Connection = "close"
		}

		response.Write(connection, res)

		if shouldClose {
			return
		}
	}
}

func handleRoot(req request.Request) response.Response {
	return response.Response{
		StatusCode: 200,
	}
}

func handleEcho(req request.Request) response.Response {
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
			return response.Response{
				StatusCode: 500,
			}
		}
		res.Body = string(compressed)
		res.ContentEncoding = response.Gzip
	}

	return res
}

func handleUserAgent(req request.Request) response.Response {
	userAgent := headers.Get(req.Headers, headers.UserAgent)
	return response.Response{
		StatusCode:  200,
		Body:        userAgent,
		ContentType: response.Text,
	}
}

func handleGetFile(req request.Request) response.Response {
	fileName := req.PathParts[2]
	content, err := os.ReadFile(filepath.Join(*dirFlag, fileName))

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return response.Response{
				StatusCode: 404,
			}
		}
		return response.Response{
			StatusCode: 500,
		}
	}

	return response.Response{
		StatusCode:  200,
		Body:        string(content),
		ContentType: response.File,
	}
}

func handlePostFile(req request.Request) response.Response {
	contentLength := headers.Get(req.Headers, headers.ContentLength)
	size, err := strconv.Atoi(contentLength)
	if err != nil || size < 0 {
		return response.Response{
			StatusCode: 400,
		}
	}

	content := []byte(req.Body)
	if len(content) > size {
		content = content[:size]
	}

	filePath, isFiles := strings.CutPrefix(req.Path, "/files/")
	if !isFiles {
		return response.Response{
			StatusCode: 404,
		}
	}

	os.WriteFile(*dirFlag+string(filePath), content, 0644)
	return response.Response{
		StatusCode:  201,
		Body:        string(content),
		ContentType: response.File,
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
