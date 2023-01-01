package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	OK = 200
)

var contentType = []byte(`Content-Type: text/html; charset=utf-8`)

type httpHeader struct {
	method         string
	path           string
	version        string
	host           string
	authorization  map[string]string
	userAgent      string
	accept         string
	acceptEncoding string
	contentLength  int
	contentType    string
	content        string
}

func splitHeaderStr(headerstr string) string {
	tmp := strings.Split(headerstr, ": ")
	return tmp[1]
}

func authBasicHeader(basicstr string) {
	tmp := strings.Split(basicstr, " ")
	decode, _ := base64.StdEncoding.DecodeString(tmp[1])
	tmp = strings.Split(string(decode), ":")
	fmt.Printf("encoding is user is %s, pass is %s\n", tmp[0], tmp[1])
}

func parseBuf(recv string) httpHeader {
	var header httpHeader
	http := strings.Split(recv, "\r\n")
	// fmt.Printf("recv http is %+v\n", http)

	for _, v := range http {
		if strings.Contains(v, "HTTP/") {
			tmp := strings.Split(v, " ")
			header.method = tmp[0]
			header.path = tmp[1]
			header.version = tmp[2]
		} else if strings.Contains(v, "Host") {
			header.host = splitHeaderStr(v)
		} else if strings.Contains(v, "Authorization") {
			authBasicHeader(splitHeaderStr(v))
		} else if strings.Contains(v, "User-Agent") {
			header.userAgent = splitHeaderStr(v)
		} else if strings.Contains(v, "Accept:") {
			header.accept = splitHeaderStr(v)
		} else if strings.Contains(v, "Accept-Encoding:") {
			header.acceptEncoding = splitHeaderStr(v)
		} else if strings.Contains(v, "Content-Length") {
			length := splitHeaderStr(v)
			header.contentLength, _ = strconv.Atoi(length)
		} else if strings.Contains(v, "Content-Type") {
			header.contentType = splitHeaderStr(v)
		} else {
			if header.contentLength == len(v) {
				header.content = v
			}
		}
	}
	fmt.Printf("header is %+v\n", header)
	return header
}

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

		httpcontent := parseBuf(string(recvbuf[:n]))

		// HTTPレスポンスの作成
		b.Write([]byte(fmt.Sprintf("HTTP/1.0 %d OK\r\n", OK)))
		if httpcontent.path == "/cookie" {
			b.Write([]byte(fmt.Sprintf("Set-Cookie: LAST_ACCESS_TIME=%s\r\n", t.Format("03:04:05"))))
		}
		b.Write([]byte(fmt.Sprintf("Date: %s\r\n", t.Format(time.RFC1123))))
		b.Write([]byte(fmt.Sprintf("Content-Length: %d\r\n", len(hello))))
		b.Write(contentType)
		b.Write([]byte("\r\n\r\n"))
		b.Write(hello)

		syscall.Sendmsg(nfd, b.Bytes(), nil, clientsa, 0)
		syscall.Close(nfd)
	}

}
