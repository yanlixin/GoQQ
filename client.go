package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
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

func HttpGet(url string, refer string) (res *http.Response, err error) {

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {

		return nil, err
	}
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.125 Safari/537.36")
	req.Header.Add("Host", "ssl.ptlogin2.qq.com")
	req.Header.Add("referer", refer)
	res, err = client.Do(req)
	if err != nil {

		return nil, err
	}

	return
}

func HttpDown(url string, path string, refer string) error {

	out, err := os.Create(path)
	defer out.Close()
	res, err := HttpGet(url, refer)
	if err != nil {

		return nil
	}
	defer res.Body.Close()

	none, err := io.Copy(out, res.Body)
	fmt.Sprintf("%d", none)
	if nil != err {

		return err
	}
	return nil

}
func HttpPost(u string, data url.Values) (re *http.Response, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	//lg.Trace("\nPOST: URL: %v\nDATA: %v", u, data.Encode())
	req, err := http.NewRequest("POST", u, strings.NewReader(data.Encode()))
	if nil != err {
		return nil, err
	}
	//ErrHandle(err, `p`)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if client.Jar != nil {
		for _, cookie := range client.Jar.Cookies(req.URL) {
			req.AddCookie(cookie)
		}
	}

	req.Header.Add(`referer`, `http://d.web2.qq.com/proxy.html?v=20110331002&callback=2&id=3`)

	re, err = client.Do(req)
	if nil != err {
		return nil, err
	}

	//ErrHandle(err, `p`)

	client.Jar.SetCookies(req.URL, re.Cookies())
	return re, nil
}
