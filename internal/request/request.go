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

const CLRF = "\r\n"

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, 8)

	encounterClrf := false
	readIndex := 0
	currentClrfIndex := 0
	lineEnd := -1

	for !encounterClrf {
		n, err := reader.Read(buf[readIndex:])
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		if n == 0 {
			return nil, errors.New("invalid request line: " + string(buf[:n]) + "")
		}
		for i := readIndex; i < n+readIndex; i++ {
			if buf[i] == CLRF[currentClrfIndex] {
				currentClrfIndex++
			} else {
				if buf[i] == CLRF[0] {
					currentClrfIndex = 1
				} else {
					currentClrfIndex = 0
				}
			}
			if currentClrfIndex == len(CLRF) {
				lineEnd = i - 1
				encounterClrf = true
				break
			}
		}
		readIndex += n
		if readIndex >= len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

	}
	line := string(buf[:lineEnd])

	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return nil, errors.New("invalid request line: " + line)
	}

	if parts[2] != "HTTP/1.1" || parts[1][0] != '/' {
		return nil, errors.New("invalid request line: " + line)
	}

	request := &Request{
		Headers: make(map[string]string),
		Body:    []byte{},
		RequestLine: RequestLine{
			Method:        parts[0],
			RequestTarget: parts[1],
			HttpVersion:   "1.1",
		},
	}

	return request, nil
}
