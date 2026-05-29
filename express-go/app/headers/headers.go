package headers

import "strings"

const (
	UserAgent      = "User-Agent"
	ContentType    = "Content-Type"
	ContentLength  = "Content-Length"
	Host           = "Host"
	AcceptEncoding = "Accept-Encoding"
	Connection     = "Connection"
)

func Parse(lines []string) []string {
	result := []string{}
	for line := 1; line < len(lines); line++ {
		if lines[line] == "" {
			break
		}
		result = append(result, lines[line])
	}
	return result
}

func Get(headers []string, key string) string {
	for _, header := range headers {
		keyVal := strings.SplitN(header, ":", 2)
		if len(keyVal) < 2 {
			continue
		}

		if keyVal[0] == key {
			return strings.TrimSpace(keyVal[1])
		}
	}

	return ""
}
