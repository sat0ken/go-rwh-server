package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var server = "http://localhost:18888"

func parsedataUrlEncode(urlstr string) url.Values {
	tmp := strings.Split(urlstr, "=")
	value := url.Values{}
	value.Add(tmp[0], tmp[1])

	return value
}

// 例3-6 HEADメソッドでヘッダーを取得
func head() {
	resp, err := http.Head(server)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Status :", resp.Status)
	log.Println("Header :", resp.Header)
}

// 例3-4 GETメソッドでクエリーを送信
func get(values url.Values) {
	resp, err := http.Get(server + "?" + values.Encode())
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

// 例3-7 x-www-form-urlencoded形式のPOSTメソッドの送信
func post(values url.Values) {
	resp, err := http.PostForm(server, values)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Status :", resp.Status)
}

// 例3-8 Go言語で任意のボディをPOST送信
func contentpost(file string) {
	content, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := http.Post(server, "text/plain", content)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Status :", resp.Status)

	// 例3-9 Go言語で任意の文字列をPOST送信
	resp, err = http.Post(server, "text/plain", strings.NewReader("テキスト"))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Status :", resp.Status)
}

func main() {
	var ishead bool
	var dataUrlEncode string
	var postValue string
	var postContent string
	flag.BoolVar(&ishead, "head", false, "use HTTP HEAD")
	flag.StringVar(&dataUrlEncode, "data-urlencode", "", "data for URL Encode")
	flag.StringVar(&postValue, "d", "", "data for POST")
	flag.StringVar(&postContent, "T", "", "file data for POST")
	flag.Parse()

	if ishead {
		head()
	} else {
		if len(dataUrlEncode) != 0 {
			values := parsedataUrlEncode(dataUrlEncode)
			get(values)
		} else if len(postValue) != 0 {
			values := parsedataUrlEncode(postValue)
			post(values)
		} else if len(postContent) != 0 {
			contentpost(postContent)
		}
	}

}
