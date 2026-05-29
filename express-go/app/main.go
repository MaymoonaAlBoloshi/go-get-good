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

		server.Serve(connection, req)

		connectionHeader := headers.Get(req.Headers, headers.Connection)

		if connectionHeader == "close" {
			return
		}
	}
}

func handleRoot(connection net.Conn, req request.Request) {
	response.Write(connection, response.Response{
		StatusCode: 200,
	})
}

func handleEcho(connection net.Conn, req request.Request) {
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
}

func handleUserAgent(connection net.Conn, req request.Request) {
	userAgent := headers.Get(req.Headers, headers.UserAgent)
	response.Write(connection, response.Response{
		StatusCode:  200,
		Body:        userAgent,
		ContentType: response.Text,
	})
}

func handleGetFile(connection net.Conn, req request.Request) {
	fileName := req.PathParts[2]
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
		return
	}

	response.Write(connection, response.Response{
		StatusCode:  200,
		Body:        string(content),
		ContentType: response.File,
	})
}

func handlePostFile(connection net.Conn, req request.Request) {
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
	if !isFiles {
		response.Write(connection, response.Response{
			StatusCode: 404,
		})
		return
	}

	os.WriteFile(*dirFlag+string(filePath), content, 0644)
	response.Write(connection, response.Response{
		StatusCode:  201,
		Body:        string(content),
		ContentType: response.File,
	})
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
