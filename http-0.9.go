package main

import (
	"bytes"
	"fmt"
	"log"
	"syscall"
)

const (
	OK = 200
)

var hello = []byte(`<html><body>hello</body></html>`)
var contentType = []byte(`Content-Type: text/html; charset=utf-8`)

func main() {

	// ソケット生成
	socket, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		log.Fatal(err)
	}

	// ソケットをbind
	err = syscall.Bind(socket, &syscall.SockaddrInet4{
		Port: 18888,
		Addr: [4]byte{127, 0, 0, 1},
	})
	if err != nil {
		log.Fatal(err)
	}

	// 接続を待つ
	err = syscall.Listen(socket, syscall.SOMAXCONN)
	if err != nil {
		log.Fatal(err)
	}

	b := bytes.Buffer{}

	for {
		nfd, _, err := syscall.Accept(socket)
		if err != nil {
			log.Fatal(err)
		}
		recvbuf := make([]byte, 1500)
		n, _, err := syscall.Recvfrom(nfd, recvbuf, 0)
		fmt.Printf("%s\n", string(recvbuf[:n]))

		// HTTPレスポンスの作成
		b.Write([]byte(fmt.Sprintf("HTTP/1.0 %d OK\r\n", OK)))
		b.Write([]byte(fmt.Sprintf("Content-Length %d\r\n", len(hello))))
		b.Write(contentType)
		b.Write([]byte("\r\n\r\n"))
		b.Write(hello)
		syscall.Write(nfd, b.Bytes())
		syscall.Close(nfd)
	}
}
