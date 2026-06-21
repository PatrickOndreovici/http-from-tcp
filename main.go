package main

import (
	"fmt"
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
		for line := range getLinesChannel(conn) {
			fmt.Println(line)
		}

	}

	defer ln.Close()
}
