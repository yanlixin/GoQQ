package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"strings"

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
	fmt.Sprintf("")
	return http.Client{
		nil,
		nil,
		jar,
		t * time.Millisecond,
	}
}

func HttpGet(url string, refer string) (string, error) {

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {

		return "", err
	}
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.125 Safari/537.36")
	req.Header.Add("Host", "ssl.ptlogin2.qq.com")
	req.Header.Add("referer", refer)
	res, err := client.Do(req)
	if err != nil {

		return "", err
	}
	defer res.Body.Close()
	body := ReadString(res.Body)
	return body, nil
}
func HttpDown(url string, path string, refer string) error {

	out, err := os.Create(path)
	defer out.Close()

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		ColorLog("[ERRO] Download fail ,Error:%+v\n", err)
		return err
	}

	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.125 Safari/537.36")
	req.Header.Add("Host", "ssl.ptlogin2.qq.com")
	req.Header.Add("referer", refer)

	res, err := client.Do(req)

	if err != nil {

		return err
	}

	defer res.Body.Close()

	none, err := io.Copy(out, res.Body)
	fmt.Sprintf("%d", none)
	if nil != err {

		return err
	}
	return nil

}
func HttpPost(url string, data string, refer string) (string, error) {
	/*
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()*/
	DebugLog("\nPOST: URL: %v\nDATA: %v", url, data)
	//req, err := http.NewRequest("POST", u, strings.NewReader(data.Encode()))
	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	//ErrHandle(err, `p`)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if client.Jar != nil {
		for _, cookie := range client.Jar.Cookies(req.URL) {
			req.AddCookie(cookie)
		}
	}

	req.Header.Add(`referer`, `http://d.web2.qq.com/proxy.html?v=20110331002&callback=2&id=3`)

	res, err := client.Do(req)
	if nil !=err{
		return "",err
	}
	//ErrHandle(err, `p`)
	defer res.Body.Close()
	body := ReadString(res.Body)
	client.Jar.SetCookies(req.URL, res.Cookies())
	DebugLog("%+v", res)
	return body, nil
}
