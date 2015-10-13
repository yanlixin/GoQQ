package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"
	"time"
)

var client http.Client

func init() {
	client = newClient(time.Duration(10000))
}

//https://github.com/Yinzo/SmartQQBot/blob/master/QQLogin.py
//https://github.com/doomsplayer/gowebQQ/blob/master/test/test.go
func newClient(t time.Duration) http.Client {
	jar, err := cookiejar.New(nil)
	//ErrHandle(err, `x`, `obtain_cookiejar`)
	if err != nil {
		ColorLog("[ERRO] Login fail ,Error:%+v\n", err)
	}

	return http.Client{
		nil,
		nil,
		jar,
		t * time.Millisecond,
	}
}

func HttpGet(url string, refer string) (string, error) {

	//req, err := http.NewRequest("GET", conf.SmartQQUrl, nil)
	req, err := http.NewRequest("GET", url, nil)

	//fmt.Printf("%s\n", conf.SmartQQUrl)
	if err != nil {
		ColorLog("[ERRO] Login fail ,Error:%+v\n", err)
		return "", err
	}

	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.125 Safari/537.36")
	req.Header.Add("Host", "ssl.ptlogin2.qq.com")
	req.Header.Add("referer", refer)
	response, err := client.Do(req)
	//defer response.Body.Close()
	if err != nil {
		ColorLog("[ERRO] Login fail ,Error:%+v\n", err)
		return "", err
	}
	client.Jar.SetCookies(req.URL, response.Cookies())
	//fmt.Printf("%+v", response)
	defer response.Body.Close()

	body := ReadString(response.Body)
	return body, nil
}
func HttpDown(url string, path string, refer string) error {

	out, err := os.Create(path)
	defer out.Close()
	//req, err := http.NewRequest("GET", conf.SmartQQUrl, nil)
	req, err := http.NewRequest("GET", url, nil)

	//fmt.Printf("%s\n", conf.SmartQQUrl)
	if err != nil {
		ColorLog("[ERRO] Login fail ,Error:%+v\n", err)
		return err
	}
	req.Header.Add("Accept", "application/javascript, */*;q=0.8")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.125 Safari/537.36")
	req.Header.Add("referer", refer)

	//req.Header.Add(`referer`, conf.ConnectReferer)

	response, err := client.Do(req)

	if err != nil {
		ColorLog("[ERRO] Download QRCode fail ,Error:%+v\n", err)
		return err
	}
	//fmt.Printf("%+v", response)
	defer response.Body.Close()

	none, err := io.Copy(out, response.Body)
	fmt.Sprintf("%+v", none)
	if nil != err {
		ColorLog("[ERRO] Save QRCode fail ,Error:%+v\n", err)
		return err
	}
	return nil

}
