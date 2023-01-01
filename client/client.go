package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"net/textproto"
	"net/url"
	"os"
	"strings"
)

var server = "http://127.0.0.1:18888"
var githuburl = "http://github.com"
var multipartdata stringFlags

type stringFlags []string

func (v *stringFlags) String() string {
	return fmt.Sprintf("%v", multipartdata)
}

func (s *stringFlags) Set(v string) error {
	*s = append(*s, v)
	return nil
}

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

// 例3-10 Go言語でマルチパートフォームをPOST送信
func postMultipart(data []string) {
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)
	for _, v := range data {
		tmp := strings.Split(v, "=")
		if tmp[0] == "name" {
			writer.WriteField(tmp[0], tmp[1])
		} else {
			fileWriter, err := writer.CreateFormFile(tmp[0], tmp[1])
			if err != nil {
				log.Fatal(err)
			}
			readFile, err := os.Open(tmp[1])
			if err != nil {
				log.Fatal(err)
			}
			defer readFile.Close()
			io.Copy(fileWriter, readFile)
			writer.Close()
		}
	}
	resp, err := http.Post(server, writer.FormDataContentType(), &buffer)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Status:", resp.Status)
}

// 例3-11　送信するファイルに任意のMIMEタイプを設定
func postMIME(data []string) {
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	for _, v := range data {
		tmp := strings.Split(v, "=")
		if tmp[0] == "thumbnaill" {
			part := make(textproto.MIMEHeader)
			part.Set("Content-Type", "image/jpeg")
			part.Set("Content-Disposition", `form-data; name="thumbnaill"; filename="photo.jpg"`)
			fileWriter, err := writer.CreatePart(part)
			if err != nil {
				log.Fatal(err)
			}
			readFile, err := os.Open("photo.jpg")
			if err != nil {
				log.Fatal(err)
			}
			io.Copy(fileWriter, readFile)
			writer.Close()
		}
	}

	resp, err := http.Post(server, writer.FormDataContentType(), &buffer)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Status:", resp.Status)
}

// 例3-12 http.Client構造体を使用して、クッキーの送受信を行う
func sendCookie() {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}
	client := http.Client{
		Jar: jar,
	}
	for i := 0; i < 2; i++ {
		resp, err := client.Get(server + "/cookie")
		if err != nil {
			log.Fatal(resp)
		}
		dump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(string(dump))
	}
}

func proxy(proxyto string) {
	proxyUrl, err := url.Parse(proxyto)
	if err != nil {
		log.Fatal(err)
	}
	client := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		},
	}
	resp, err := client.Get(githuburl)
	if err != nil {
		log.Fatal(err)
	}
	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(dump))
}

func main() {
	var ishead bool
	var dataUrlEncode string
	var postValue string
	var postContent string
	var cookie string
	var proxyurl string

	flag.BoolVar(&ishead, "head", false, "use HTTP HEAD")
	flag.StringVar(&dataUrlEncode, "data-urlencode", "", "data for URL Encode")
	flag.StringVar(&postValue, "d", "", "data for POST")
	flag.StringVar(&postContent, "T", "", "file data for POST")
	flag.Var(&multipartdata, "F", "send multipart form data")
	flag.StringVar(&cookie, "b", "", "cookie value")
	flag.StringVar(&proxyurl, "x", "", "set proxy URL")
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
		} else if len(multipartdata) != 0 {
			postMultipart(multipartdata)
			postMIME(multipartdata)
		} else if len(cookie) != 0 {
			sendCookie()
		} else if len(proxyurl) != 0 {
			proxy(proxyurl)
		}
	}

}
