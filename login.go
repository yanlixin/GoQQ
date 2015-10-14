package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

var UserInfo struct {
	UserName   string
	Account    string
	ClientId   int
	PTWebQQ    string
	FriendList map[string]string
	VFWebQQ    string
	PSessionId string
}

//https://github.com/Yinzo/SmartQQBot/blob/master/QQLogin.py
//https://github.com/doomsplayer/gowebQQ/blob/master/test/test.go
func getRevalue(ex string, body string) (string, error) {
	result := regexp.MustCompile(ex).FindStringSubmatch(body)
	if len(result) > 1 {
		return result[1], nil
	} else if len(result) > 0 {
		return result[0], nil
	} else {
		return "", nil
	}
}

func LoginByQRCode() (int, error) {

	ColorLog("[INFO] Requesting the login pages... \n")

	refer := conf.ConnectReferer
	sBody, err := HttpGet(conf.SmartQQUrl, refer)
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
	var ret []string
	done := make(chan bool, 1)
	go func() {

		err := HttpDown(qr_url, conf.QRCodePath, refer)
		if nil != err {
			ColorLog("[ERRO] DownLoad the QRCode faild,%+v \n", err)
			done <- false
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
			ret = strings.Split(html, "'")
			//fmt.Printf("%+v \r\n", ret[1])
			if ret[1] == "65" {
				fmt.Println("")
				done <- false
				// 65: QRCode 失效, 0: 验证成功, 66: 未失效, 67: 验证中
				break
			}
			if ret[1] == "0" {
				fmt.Println("")
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

		// 删除QRCode文件
		//if _, err := os.Stat(conf.QRCodePath); os.IsNotExist(err) {
		//	ColorLog("[INFO] No such file or directory: %s", conf.QRCodePath)
		//
		//}

		os.Remove(conf.QRCodePath)
		userinfo := &UserInfo
		// 记录登陆账号的昵称

		userinfo.UserName = ret[11]
		DebugLog(userinfo.UserName)
		DebugLog("rev[5]:%s\r\n", ret[5])
		html, err := HttpGet(ret[5], refer)
		if nil != err {
			ColorLog("[ERRO] Get Login Status faild,%+v \n", err)
			//return 10010,err
		}
		DebugLog("mibao_res html:  %s\r\n", html)

		url1, err := getRevalue(` src="(.+?)"`, html)
		if nil != err {
			ColorLog("[ERRO] Get Login Status faild,%+v \n", err)
			return 10010, err
		}
		DebugLog("url=%s\r\n", url1)

		if url1 != "" {
			html, err := HttpGet(strings.Replace(url1, "&amp;", "&", 0), refer)
			if nil != err {
				ColorLog("[ERRO] Get Login Status faild,%+v \n", err)
				return 10010, err
			}
			url1, err = getRevalue(`location\.href=""(.+?)""`, html)
			if nil != err {
				ColorLog("[ERRO] Get Login Status faild,%+v \n", err)
				return 10010, err
			}
			DebugLog("url=%s\r\n", url1)
			none, err := HttpGet(url1, refer)
			if nil != err {
				ColorLog("[ERRO] Get Login Status faild,%+v \n", err)
				return 10010, err
			}
			DebugLog(none)
		}
		u, _ := url.Parse(fmt.Sprintf("%v", ret[5]))
		client.Jar.Cookies(u)
		//self.ptwebqq = self.req.getCookie('ptwebqq')
		login_error := 1
		var ret map[string]interface{}

		for login_error > 0 {
			r := fmt.Sprintf(`{"ptwebqq":"%s","clientid":%d,"psessionid":"%s","status":"online"}`,
				userinfo.PTWebQQ,
				userinfo.ClientId,
				userinfo.PSessionId)
			data := fmt.Sprintf(`{"r":%s}`, r)
			html, err := HttpPost("http://d.web2.qq.com/channel/login2", data, refer)
			if nil != err {

				ColorLog("[ERRO] Post User Login fail,%+v \n", err)

			}
			DebugLog("login html: %s ", html)
			byt := []byte(html)
			err1 := json.Unmarshal(byt, &ret)
			if nil != err1 {
				login_error += 1
				ColorLog("[ERRO] Get login fail, retrying...,%+v \n", err1)

			} else {
				login_error = 0
			}

			if ret["retcode"] != "0" {
				DebugLog("%+v", ret)
				ColorLog("[ERRO] return code:,%+v \n", ret["retcode"])

				return 10011, nil
			}
		}
		//vfwebqq := ret["result"]["vfwebqq"]
		//psessionid := ret["result"]["psessionid"]
		//account := ret["result"]["uin"]

	} else {
		ColorLog("[INFO] QRCode 失效.")
	}

	return 0, nil

}
