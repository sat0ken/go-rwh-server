package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	var ishttp0_9 bool
	var ishttp1_0 bool
	flag.BoolVar(&ishttp0_9, "http0.9", false, "listen server HTTP/0.9")
	flag.BoolVar(&ishttp1_0, "http1.0", false, "listen server HTTP/1.0")
	flag.Parse()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if ishttp0_9 {
			http0_9()
		} else if ishttp1_0 {
			http1_0()
		}
	}()

	<-sigs
	os.Exit(0)

}
