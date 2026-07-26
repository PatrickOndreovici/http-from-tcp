package main

import (
	"fmt"
	"httpfromtcp/internal/request"
	"io"
	"net"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {

	ch := make(chan string)

	fn := func() {
		data := make([]byte, 8)
		var currentLine string = ""

		for {
			n, err := f.Read(data)
			if err != nil {
				if currentLine != "" {
					ch <- currentLine + "\n"
				}
				break
			}

			str := string(data[:n])

			parts := strings.Split(str, "\n")

			for i := 0; i < len(parts)-1; i++ {
				ch <- (currentLine + parts[i] + "\n")
				currentLine = ""
			}
			currentLine += parts[len(parts)-1]
		}
		close(ch)
	}

	go fn()

	return ch
}

func main() {
	ln, err := net.Listen("tcp", ":42069")

	if err != nil {
		panic(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			println("error accepting connection: " + err.Error())
			continue
		}
		println("accepted connection")

		req, err := request.RequestFromReader(conn)
		if err != nil {
			println("error parsing request: " + err.Error())
			conn.Close()
			continue
		}

		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", req.RequestLine.Method)
		fmt.Printf("- Target: %s\n", req.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", req.RequestLine.HttpVersion)

		conn.Close()

	}

	defer ln.Close()
}
