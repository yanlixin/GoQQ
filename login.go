package main

import (
	"fmt"
	"regexp"
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
	sBody, _ := HttpGet("http://w.qq.com/login.html", refer)
	//fmt.Printf("ptui_checkVC is %v\n", sBody)

	src, _ := getRevalue(`\.src = "(.+?)"`, sBody)

	html, _ := HttpGet(src+"0", refer)

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
	//date_to_millis(datetime.datetime.utcnow())

	error_times := 0
	//var ret []string

	qr_url := fmt.Sprintf("https://ssl.ptlogin2.qq.com/ptqrshow?appid={0}&e=0&l=L&s=8&d=72&v=4", appid)

	//out, err := os.Create(conf.QRCodePath)
	QRCodePath := "./v.jpg"
	HttpDown(qr_url, QRCodePath, refer)
	fmt.Println("Please scan the downloaded QRCode")
	//bufio.NewReader(os.Stdin)

	done := make(chan bool, 1)
	go func() {
		for {
			callbackUrl := fmt.Sprintf(`https://ssl.ptlogin2.qq.com/ptqrlogin?webqq_type=10&remember_uin=1&login2qq=1&aid=%s&u1=http%%3A%%2F%%2Fw.qq.com%%2Fproxy.html%%3Flogin2qq%%3D1%%26webqq_type%%3D10&ptredirect=0&ptlang=2052&daid=164&from_ui=1&pttype=1&dumy=&fp=loginerroralert&action=0-0-%d&mibao_css=%s&t=undefined&g=1&js_type=0&js_ver=%s&login_sig=%s`, appid, time.Now().Unix()*1000-star_time, mibao_css, js_ver, sign)
			html, err := HttpGet( callbackUrl,refer)
			if nil != err {
				ColorLog("[ERRO] QRCode check faild,error: %s\r\n  ", err)
			}
			//fmt.Println(html)
			ColorLog("[INFO] QRCode check html:   " + html)
			time.Sleep(time.Second)
			
			ret = strings.Split(html,"'")
			fmt.Printf("%+v",ret)
			if ret[1] in ('0', '65'){  
				// 65: QRCode 失效, 0: 验证成功, 66: 未失效, 67: 验证中
			    break
			 }
			if ret[1] == '0' || error_times > 10{
				done <- true
			}
			error_times += 1
		}
	}()
	status := <-done
	fmt.Println(status)

	//if ret[1] != "0" {
	//	return 10002, nil
	//}
	ColorLog("[INFO] QRCode scaned, now logging in.")

	return 0, nil

}
func worker(done chan bool) {
	fmt.Print("working...")

	fmt.Println("done")

}
