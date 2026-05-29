package request

import (
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/app/headers"
)

const CRLF = "\r\n"

type Request struct {
	Raw       string
	Method    string
	Path      string
	PathParts []string
	Headers   []string
	Body      string
}

func Parse(raw string) (Request, bool) {
	lines := strings.Split(raw, CRLF)
	if len(lines) == 0 {
		return Request{}, false
	}

	requestLine := strings.Split(lines[0], " ")
	if len(requestLine) < 2 {
		return Request{}, false
	}

	req := Request{
		Raw:       raw,
		Method:    requestLine[0],
		Path:      requestLine[1],
		PathParts: strings.Split(requestLine[1], "/"),
		Headers:   headers.Parse(lines),
	}

	parts := strings.SplitN(raw, CRLF+CRLF, 2)
	if len(parts) == 2 {
		req.Body = parts[1]
	}

	return req, true
}
