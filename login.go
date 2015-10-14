package main

import (
	"fmt"

	"regexp"
	"strings"
	"time"
)

//https://github.com/Yinzo/SmartQQBot/blob/master/QQLogin.py
//https://github.com/doomsplayer/gowebQQ/blob/master/test/test.go
func getRevalue(ex string, body string) (string, error) {
	result := regexp.MustCompile(ex).FindStringSubmatch(body)
	if len(result) > 1 {
		return result[1], nil
	} else {
		return result[0], nil
	}
}

func LoginByQRCode() (int, error) {

	ColorLog("[INFO] Requesting the login pages... \n")

	refer := `http://d.web2.qq.com/proxy.html?v=20110331002&callback=2&id=3`
	sBody, err := HttpGet("http://w.qq.com/login.html", refer)
	if nil != err {
		ColorLog("[ERRO] Requesting the login faild,%+v \n", err)
		return 10001, err
	}

	src, _ := getRevalue(`\.src = "(.+?)"`, sBody)
	//fmt.Println(src)
	html, err := HttpGet(src+"0", refer)
	if nil != err {
		ColorLog("[ERRO] Requesting the QRCode login faild,%+v \n", err)
		return 10002, err
	}

	appid, _ := getRevalue(`var g_appid =encodeURIComponent\("(\d+)"\);`, html)
	sign, _ := getRevalue(`var g_login_sig=encodeURIComponent\("(.*?)"\);`, html)
	js_ver, _ := getRevalue(`var g_pt_version=encodeURIComponent\("(\d+)"\);`, html)
	mibao_css, _ := getRevalue(`var g_mibao_css=encodeURIComponent\("(.+?)"\);`, html)
	/*
		fmt.Printf("ptui_checkVC is %s\n", src)

		fmt.Printf("html is %+v\n", html)
		fmt.Printf("appid is %s\n", appid)
		fmt.Printf("sign is %s\n", sign)
		fmt.Printf("js_ver is %s\n", js_ver)
		fmt.Printf("mibao_css is %s\n", mibao_css)
	*/
	star_time := time.Now().Unix() * 1000
	error_times := 0
	qr_url := fmt.Sprintf("https://ssl.ptlogin2.qq.com/ptqrshow?appid=%s&e=0&l=L&s=8&d=72&v=4", appid)

	done := make(chan bool, 1)
	go func() {
		QRCodePath := "./v.jpg"
		err := HttpDown(qr_url, QRCodePath, refer)
		if nil != err {
			ColorLog("[ERRO] DownLoad the QRCode faild,%+v \n", err)
			//return 10003, err
		}
		fmt.Println("Please scan the downloaded QRCode")

		for {
			checkStatusUrl := fmt.Sprintf(`https://ssl.ptlogin2.qq.com/ptqrlogin?webqq_type=10&remember_uin=1&login2qq=1&aid=%s&u1=http%%3A%%2F%%2Fw.qq.com%%2Fproxy.html%%3Flogin2qq%%3D1%%26webqq_type%%3D10&ptredirect=0&ptlang=2052&daid=164&from_ui=1&pttype=1&dumy=&fp=loginerroralert&action=0-0-%d&mibao_css=%s&t=undefined&g=1&js_type=0&js_ver=%s&login_sig=%s`,
				appid, time.Now().Unix()*1000-star_time, mibao_css, js_ver, sign)

			html, err := HttpGet(checkStatusUrl, qr_url)
			if nil != err {
				ColorLog("[ERRO] Check the QRCode Login Status faild,%+v \n", err)
			}
			fmt.Printf(".")

			time.Sleep(time.Second)
			ret := strings.Split(html, "'")
			//fmt.Printf("%+v \r\n", ret[1])
			if ret[1] == "65" {
				//done <- false
				// 65: QRCode 失效, 0: 验证成功, 66: 未失效, 67: 验证中
				break
			}
			if ret[1] == "0" {
				done <- true
			}

		}
		if error_times > 10 {
			fmt.Printf("%s \r\n", "Done")
			done <- false
		}
		error_times += 1
	}()
	status := <-done
	if status {
		ColorLog("[INFO] QRCode scaned, now logging in.")
	}

	return 0, nil

}
