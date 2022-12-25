package main

import (
	"bytes"
	"fmt"
	"log"
	"syscall"
	"time"
)

const (
	OK = 200
)

var contentType = []byte(`Content-Type: text/html; charset=utf-8`)

func http1_0() {

	t := time.Now()

	socket := makeServerSocket()
	defer syscall.Close(socket)

	for {
		b := bytes.Buffer{}
		nfd, _, err := syscall.Accept(socket)
		if err != nil {
			log.Fatal(err)
		}
		recvbuf := make([]byte, 1500)
		n, clientsa, err := syscall.Recvfrom(nfd, recvbuf, 0)
		fmt.Printf("%s\n", string(recvbuf[:n]))

		// HTTPレスポンスの作成
		b.Write([]byte(fmt.Sprintf("HTTP/1.0 %d OK\r\n", OK)))
		b.Write([]byte(fmt.Sprintf("Date: %s\r\n", t.Format(time.RFC1123))))
		b.Write([]byte(fmt.Sprintf("Content-Length %d\r\n", len(hello))))
		b.Write(contentType)
		b.Write([]byte("\r\n\r\n"))
		b.Write(hello)

		//syscall.Write(nfd, b.Bytes())

		syscall.Sendmsg(nfd, b.Bytes(), nil, clientsa, 0)
		syscall.Close(nfd)
	}

}
