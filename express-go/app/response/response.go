package response

import (
	"fmt"
	"net"
)

const (
	File = "file"
	Text = "text"
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

func Write(conn net.Conn, statusCode int, body string, contentType string) {
	status := statusText[statusCode]
	if status == "" {
		status = "Unknown"
	}

	resp := fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, status)

	if contentType != "" {
		resp += fmt.Sprintf("Content-Type: %s\r\n", contentTypeText[contentType])
	}

	if body != "" || contentType != "" {
		resp += fmt.Sprintf("Content-Length: %d\r\n", len(body))
	}

	resp += "\r\n" + body

	conn.Write([]byte(resp))
}
