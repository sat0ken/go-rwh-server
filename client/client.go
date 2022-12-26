package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func parsedataUrlEncode(urlstr string) url.Values {
	tmp := strings.Split(urlstr, "=")
	value := url.Values{}
	value.Add(tmp[0], tmp[1])

	return value
}

func head() {
	resp, err := http.Head("http://localhost:18888")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Status :", resp.Status)
	log.Println("Header :", resp.Header)
}

func get(values url.Values) {
	resp, err := http.Get("http://localhost:18888" + "?" + values.Encode())
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(body))
	log.Println("Status :", resp.Status)
	log.Println("StatusCode :", resp.StatusCode)
	log.Println("Header :", resp.Header)
}

func main() {
	var dataUrlEncode string
	var ishead bool
	flag.BoolVar(&ishead, "head", false, "use HTTP HEAD")
	flag.StringVar(&dataUrlEncode, "data-urlencode", "", "data for URL Encode")
	flag.Parse()

	if ishead {
		head()
	} else {
		values := parsedataUrlEncode(dataUrlEncode)
		get(values)
	}

}
