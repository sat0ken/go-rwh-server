package main

import (
	"bytes"
	"fmt"
	"log"
	"syscall"
)

var hello = []byte(`<html><body>hello</body></html>`)

func http0_9() {

	socket := makeServerSocket()

	for {
		b := bytes.Buffer{}
		nfd, _, err := syscall.Accept(socket)
		if err != nil {
			log.Fatal(err)
		}
		recvbuf := make([]byte, 1500)
		n, clientsa, err := syscall.Recvfrom(nfd, recvbuf, 0)
		fmt.Printf("%s\n", string(recvbuf[:n]))

		b.Write(hello)

		syscall.Sendmsg(nfd, b.Bytes(), nil, clientsa, 0)
		syscall.Close(nfd)
	}

}
