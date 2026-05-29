package response

import (
	"fmt"
	"net"
)

const (
	File = "file"
	Text = "text"
	Gzip = "gzip"
)

var statusText = map[int]string{
	200: "OK",
	201: "Created",
	400: "Bad Request",
	404: "Not Found",
	500: "Internal Server Error",
}

var contentTypeText = map[string]string{
	File: "application/octet-stream",
	Text: "text/plain",
}

type Response struct {
	StatusCode      int
	Body            string
	ContentType     string
	ContentEncoding string
	Connection      string
}

func Write(conn net.Conn, res Response) {
	status := statusText[res.StatusCode]
	if status == "" {
		status = "Unknown"
	}

	resp := fmt.Sprintf("HTTP/1.1 %d %s\r\n", res.StatusCode, status)

	if res.ContentType != "" {
		resp += fmt.Sprintf("Content-Type: %s\r\n", contentTypeText[res.ContentType])
	}

	if res.ContentEncoding != "" {
		resp += fmt.Sprintf("Content-Encoding: %s\r\n", res.ContentEncoding)
	}

	if res.Connection != "" {
		resp += fmt.Sprintf("Connection: %s\r\n", res.Connection)
	}

	if res.Body != "" || res.ContentType != "" {
		resp += fmt.Sprintf("Content-Length: %d\r\n", len(res.Body))
	}

	resp += "\r\n" + res.Body

	conn.Write([]byte(resp))
}
