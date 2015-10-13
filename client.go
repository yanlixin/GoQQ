package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"

	"os"
	"strconv"
	"strings"
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

func HttpGet(url string, refer string) (string, error, []*http.Cookie) {

	fmt.Println("#####################\r\n")
	//req, err := http.NewRequest("GET", conf.SmartQQUrl, nil)
	req, err := http.NewRequest("GET", url, nil)

	//fmt.Printf("%s\n", conf.SmartQQUrl)
	if err != nil {
		ColorLog("[ERRO] Login fail ,Error:%+v\n", err)
		return "", err, nil
	}

	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.125 Safari/537.36")
	req.Header.Add("Host", "ssl.ptlogin2.qq.com")
	req.Header.Add("referer", refer)
	response, err := client.Do(req)
	//defer response.Body.Close()
	if err != nil {
		ColorLog("[ERRO] Login fail ,Error:%+v\n", err)
		return "", err, nil
	}
	//for _,h :=range response.Header{
	//	fmt.Printf("%+v \r\n",h)
	//}

	cookies := readSetCookies(response.Header)
	// client.Jar.Cookies(req.URL)
	for _, v := range client.Jar.Cookies(req.URL) {
		println(v.Name + "=" + v.Value)
	}

	//client.Jar.SetCookies(req.URL, cookies)
	//fmt.Printf("%+v", response)
	defer response.Body.Close()

	body := ReadString(response.Body)
	fmt.Println("---------------------\r\n")
	return body, nil, cookies
}
func HttpDown(url string, path string, refer string) error {

	out, err := os.Create(path)
	defer out.Close()
	//req, err := http.NewRequest("GET", conf.SmartQQUrl, nil)
	req, err := http.NewRequest("GET", url, nil)

	//fmt.Printf("%s\n", conf.SmartQQUrl)
	if err != nil {
		ColorLog("[ERRO] Download fail ,Error:%+v\n", err)
		return err
	}

	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.125 Safari/537.36")
	req.Header.Add("Host", "ssl.ptlogin2.qq.com")
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

func parseCookieValue(raw string, allowDoubleQuote bool) (string, bool) {
	// Strip the quotes, if present.
	if allowDoubleQuote && len(raw) > 1 && raw[0] == '"' && raw[len(raw)-1] == '"' {
		raw = raw[1 : len(raw)-1]
	}
	for i := 0; i < len(raw); i++ {
		if !validCookieValueByte(raw[i]) {
			return "", false
		}
	}
	return raw, true
}

func isCookieNameValid(raw string) bool {
	if raw == "" {
		return false
	}
	return true //strings.IndexFunc(raw, isNotToken) < 0
}
func validCookieValueByte(b byte) bool {
	return 0x20 <= b && b < 0x7f && b != '"' && b != ';' && b != '\\'
}
func readSetCookies(h http.Header) []*http.Cookie {
	cookies := []*http.Cookie{}
	for _, line := range h["Set-Cookie"] {
		parts := strings.Split(strings.TrimSpace(line), ";")
		if len(parts) == 1 && parts[0] == "" {
			//	continue
		}
		parts[0] = strings.TrimSpace(parts[0])
		j := strings.Index(parts[0], "=")
		if j < 0 {
			//	continue
		}
		name, value := parts[0][:j], parts[0][j+1:]
		//if !isCookieNameValid(name) {
		//	continue
		//}
		value, success := parseCookieValue(value, true)
		if !success {
			//	continue
		}
		c := &http.Cookie{
			Name:  name,
			Value: value,
			Raw:   line,
		}
		for i := 1; i < len(parts); i++ {
			parts[i] = strings.TrimSpace(parts[i])
			if len(parts[i]) == 0 {
				continue
			}

			attr, val := parts[i], ""
			if j := strings.Index(attr, "="); j >= 0 {
				attr, val = attr[:j], attr[j+1:]
			}
			lowerAttr := strings.ToLower(attr)
			val, success = parseCookieValue(val, false)
			if !success {
				c.Unparsed = append(c.Unparsed, parts[i])
				continue
			}
			switch lowerAttr {
			case "secure":
				c.Secure = true
				continue
			case "httponly":
				c.HttpOnly = true
				continue
			case "domain":
				c.Domain = val
				continue
			case "max-age":
				secs, err := strconv.Atoi(val)
				if err != nil || secs != 0 && val[0] == '0' {
					break
				}
				if secs <= 0 {
					c.MaxAge = -1
				} else {
					c.MaxAge = secs
				}
				continue
			case "expires":
				c.RawExpires = val
				exptime, err := time.Parse(time.RFC1123, val)
				if err != nil {
					exptime, err = time.Parse("Mon, 02-Jan-2006 15:04:05 MST", val)
					if err != nil {
						c.Expires = time.Time{}
						break
					}
				}
				c.Expires = exptime.UTC()
				continue
			case "path":
				c.Path = val
				continue
			}
			c.Unparsed = append(c.Unparsed, parts[i])
		}
		cookies = append(cookies, c)
	}
	return cookies
}

