package request

import (
	"errors"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
	Headers     map[string]string
	Body        []byte
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, 4096)

	line := ""

	for {
		n, err := reader.Read(buf)

		lines := strings.Split(string(buf[:n]), "\r\n")

		if len(lines) > 1 {
			line += lines[0]
			break
		} else {
			line += lines[0]
		}
		if err == io.EOF {
			break
		}

	}

	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return nil, errors.New("invalid request line: " + line + "")
	}
	if parts[2] != "HTTP/1.1" || parts[1][0] != '/' {
		return nil, errors.New("invalid request line: " + line + "")
	}
	request := &Request{
		Headers: make(map[string]string),
		Body:    []byte{},
		RequestLine: RequestLine{
			HttpVersion:   strings.Split(parts[2], "/")[1],
			RequestTarget: parts[1],
			Method:        parts[0],
		},
	}

	return request, nil
}
