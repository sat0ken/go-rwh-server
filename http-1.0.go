package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	OK = 200
)

var contentType = []byte(`Content-Type: text/html; charset=utf-8`)
var html401 = []byte(`<html><body>secret page</body></html>`)

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

func authBasicHeader(basicstr string) error {
	tmp := strings.Split(basicstr, " ")
	decode, _ := base64.StdEncoding.DecodeString(tmp[1])
	tmp = strings.Split(string(decode), ":")
	// fmt.Printf("encoding is user is %s, pass is %s\n", tmp[0], tmp[1])
	if tmp[0] != "user" {
		return errors.New("user is not correct")
	} else if tmp[1] != "pass" {
		return errors.New("password is not correct")
	}
	return nil
}

func parseBuf(recv string) (httpHeader, error) {
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
			err := authBasicHeader(splitHeaderStr(v))
			if err != nil {
				return header, err
			}
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
	return header, nil
}

func doproxy(hosturl string) []byte {
	client := http.Client{}
	resp, err := client.Get(hosturl)
	if err != nil {
		log.Fatal(err)
	}
	dump, err := httputil.DumpResponse(resp, true)
	return dump
}

func http1_0() {

	t := time.Now()
	socket := makeServerSocket()

	for {
		b := bytes.Buffer{}
		nfd, _, err := syscall.Accept(socket)
		if err != nil {
			log.Fatal(err)
		}
		recvbuf := make([]byte, 1500)
		n, clientsa, err := syscall.Recvfrom(nfd, recvbuf, 0)

		httpcontent, err := parseBuf(string(recvbuf[:n]))
		// basic認証失敗
		if err != nil {
			log.Println(err)
			b.Write([]byte(fmt.Sprintf("HTTP/1.1 %d Unauthorized\r\n", 401)))
			b.Write([]byte(fmt.Sprintf("Date: %s\r\n", t.Format(time.RFC1123))))
			b.Write([]byte("\r\n"))
			b.Write(html401)
			syscall.Sendmsg(nfd, b.Bytes(), nil, clientsa, 0)
			syscall.Close(nfd)
		}

		// HTTPレスポンスの作成
		b.Write([]byte(fmt.Sprintf("HTTP/1.0 %d OK\r\n", OK)))
		if httpcontent.path == "/cookie" {
			b.Write([]byte(fmt.Sprintf("Set-Cookie: LAST_ACCESS_TIME=%s\r\n", t.Format("03:04:05"))))
		}
		// proxy対応
		if strings.HasPrefix(httpcontent.path, "http") {
			b.Write(doproxy(httpcontent.path))
		} else {
			b.Write([]byte(fmt.Sprintf("Date: %s\r\n", t.Format(time.RFC1123))))
			b.Write([]byte(fmt.Sprintf("Content-Length: %d\r\n", len(hello))))
			b.Write(contentType)
			b.Write([]byte("\r\n\r\n"))
			b.Write(hello)
		}
		syscall.Sendmsg(nfd, b.Bytes(), nil, clientsa, 0)
		syscall.Close(nfd)
	}

}
