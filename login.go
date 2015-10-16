package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var qq *UserInfo

func init() {
	qq = &UserInfo{}
}

type UserInfo struct {
	username      string
	account       string
	clientid      string
	ptwebqq       string
	friendlist    map[string]string
	vfwebqq       string
	psessionid    string
	verifysession string
	msgid         int64
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
	res, err := HttpGet(conf.SmartQQUrl, refer)
	if nil != err {
		ColorLog("[ERRO] Requesting the login faild,%+v \n", err)
		return 10001, err
	}
	defer res.Body.Close()
	body := ReadString(res.Body)

	src, _ := getRevalue(`\.src = "(.+?)"`, body)
	//fmt.Println(src)
	res, err = HttpGet(src+"0", refer)
	if nil != err {
		ColorLog("[ERRO] Requesting the QRCode login faild,%+v \n", err)
		return 10002, err
	}
	defer res.Body.Close()
	body = ReadString(res.Body)
	appid, _ := getRevalue(`var g_appid =encodeURIComponent\("(\d+)"\);`, body)
	sign, _ := getRevalue(`var g_login_sig=encodeURIComponent\("(.*?)"\);`, body)
	js_ver, _ := getRevalue(`var g_pt_version=encodeURIComponent\("(\d+)"\);`, body)
	mibao_css, _ := getRevalue(`var g_mibao_css=encodeURIComponent\("(.+?)"\);`, body)

	star_time := time.Now().Unix() * 1000
	error_times := 0
	qr_url := fmt.Sprintf("https://ssl.ptlogin2.qq.com/ptqrshow?appid=%s&e=0&l=L&s=8&d=72&v=4", appid)
	var ret []string
	done := make(chan string, 1)
	go func() {

		err := HttpDown(qr_url, conf.QRCodePath, refer)
		if nil != err {
			ColorLog("[ERRO] DownLoad the QRCode faild,%+v \n", err)
			done <- ""
		}
		fmt.Println("Please scan the downloaded QRCode")

		for {
			checkStatusUrl := fmt.Sprintf(`https://ssl.ptlogin2.qq.com/ptqrlogin?webqq_type=10&remember_uin=1&login2qq=1&aid=%s&u1=http%%3A%%2F%%2Fw.qq.com%%2Fproxy.html%%3Flogin2qq%%3D1%%26webqq_type%%3D10&ptredirect=0&ptlang=2052&daid=164&from_ui=1&pttype=1&dumy=&fp=loginerroralert&action=0-0-%d&mibao_css=%s&t=undefined&g=1&js_type=0&js_ver=%s&login_sig=%s`,
				appid, time.Now().Unix()*1000-star_time, mibao_css, js_ver, sign)

			res, err := HttpGet(checkStatusUrl, qr_url)
			if nil != err {
				ColorLog("[ERRO] Check the QRCode Login Status faild,%+v \n", err)
			}
			defer res.Body.Close()
			body := ReadString(res.Body)
			fmt.Printf(".")

			time.Sleep(time.Second)
			ret = strings.Split(body, "'")
			//fmt.Printf("%+v \r\n", ret[1])
			if ret[1] == "65" {
				fmt.Println("")
				done <- ""
				// 65: QRCode 失效, 0: 验证成功, 66: 未失效, 67: 验证中
				break
			}
			if ret[1] == "0" {
				fmt.Println("")
				done <- body
			}

		}
		if error_times > 10 {
			fmt.Printf("%s \r\n", "Done")
			done <- ""
		}
		error_times += 1
	}()
	loginStatus := <-done
	if len(loginStatus) > 0 {
		ColorLog("[INFO] QRCode scaned, now logging in.")
		DebugLog("LoginStatus:%s", loginStatus)
		os.Remove(conf.QRCodePath)
		statusCode, err := setLoginStatus(loginStatus, ret)
		if nil != err {
			ColorLog("[ERRO] Check the QRCode Login Status faild,%+v \n", err)
			return statusCode, err

		}
	} else {
		ColorLog("[INFO] QRCode 失效.")
	}

	return 0, nil

}
func setLoginStatus(sBody string, ret []string) (statusCode int, err error) {
	rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	clientid := strconv.Itoa(rd.Intn(90000000) + 10000000)

	// 记录登陆账号的昵称

	qq.username = ret[11]
	DebugLog(qq.username)

	reg := regexp.MustCompile(`ptuiCB\('0','0','(.*)','0','登录成功！', '.*'\);`)

	if !reg.MatchString(sBody) {

		panic(errors.New(`第一次握手失败（密码错误？）`))
	}

	ssl := reg.FindStringSubmatch(sBody)
	re, err := HttpGet(ssl[1], conf.ConnectReferer)
	if nil != err {
		return -1, err
	}
	defer re.Body.Close()

	v := url.Values{}
	v.Set(`clientid`, clientid)
	v.Set(`psessionid`, `null`)

	c, err := json.Marshal(
		map[string]interface{}{
			`status`:     `online`,
			`ptwebqq`:    qq.ptwebqq,
			`passwd_sig`: ``,
			`clientid`:   qq.clientid,
			`psessionid`: nil})
	v.Set(`r`, string(c))
	re, err = HttpPost(`http://d.web2.qq.com/channel/login2`, v)

	if err != nil {
		return -2, err
	}
	defer re.Body.Close()

	retb := ReadByte(re.Body)
	var js map[string]interface{}
	json.Unmarshal(retb, &js)
	//DebugLog("%+v", js)

	if js[`retcode`] != float64(0) {
		DebugLog("第二次握手失败,错误码:%v", js[`retcode`])
		return -3, errors.New(fmt.Sprintf("第二次握手失败,错误码:%v", js[`retcode`]))
	}
	//var result map[string]interface{}
	result := js[`result`].(map[string]interface{})
	//DebugLog("%+v", result)

	qq.vfwebqq = result[`vfwebqq`].(string)
	qq.psessionid = result[`psessionid`].(string)
	DebugLog("%+v", qq)
	return 0, nil

}
