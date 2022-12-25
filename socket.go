package main

import (
	"log"
	"syscall"
)

func makeServerSocket() int {
	// ソケット生成
	socket, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		log.Fatal(err)
	}

	// SO_REUSEADDRを設定
	// https://www.geekpage.jp/programming/winsock/so_reuseaddr.php
	syscall.SetsockoptInt(socket, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)

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

	return socket
}
