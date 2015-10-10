package main

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"regexp"
	"time"
)
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
func getRevalue(ex string, body string) (string, error) {
	result := regexp.MustCompile(ex).FindStringSubmatch(body)
	if len(result) > 1 {
		return result[1], nil
	} else {
		return result[0], nil
	}
}
func get(url string) (string, error) {
	client := newClient(time.Duration(10000))
	//req, err := http.NewRequest("GET", conf.SmartQQUrl, nil)
	req, err := http.NewRequest("GET", url, nil)
	//fmt.Printf("%s\n", conf.SmartQQUrl)
	if err != nil {
		//ColorLog("[ERRO] Login fail ,Error:%+v\n", err)
		return "", err
	}
	req.Header.Add(`referer`, `http://d.web2.qq.com/proxy.html?v=20110331002&callback=2&id=3`)
	//req.Header.Add(`referer`, conf.ConnectReferer)

	response, err1 := client.Do(req)
	//defer response.Body.Close()
	if err1 != nil {
		ColorLog("[ERRO] Login fail ,Error:%+v\n", err1)
		return "", err1
	}
	//fmt.Printf("%+v", response)
	defer response.Body.Close()

	sBody := ReadString(response.Body)
	return sBody, nil

}
func LoginByQRCode() (int, error) {

	ColorLog("[INFO] Requesting the login pages... \n")
	/*
		client := newClient(time.Duration(10000))
		//req, err := http.NewRequest("GET", conf.SmartQQUrl, nil)
		req, err := http.NewRequest("GET", "http://w.qq.com/login.html", nil)
		//fmt.Printf("%s\n", conf.SmartQQUrl)
		if err != nil {
			//ColorLog("[ERRO] Login fail ,Error:%+v\n", err)
			return 10001, err
		}
		req.Header.Add(`referer`, `http://d.web2.qq.com/proxy.html?v=20110331002&callback=2&id=3`)
		//req.Header.Add(`referer`, conf.ConnectReferer)

		response, err1 := client.Do(req)
		//defer response.Body.Close()
		if err1 != nil {
			ColorLog("[ERRO] Login fail ,Error:%+v\n", err1)
			return 10001, err1
		}
		//fmt.Printf("%+v", response)
		defer response.Body.Close()

		sBody := ReadString(response.Body)
	*/
	sBody, _ := get("http://w.qq.com/login.html")
	//fmt.Printf("ptui_checkVC is %v\n", sBody)

	src, _ := getRevalue(`\.src = "(.+?)"`, sBody)
	
	html, _ := get(src + "0")
	
	appid, _ := getRevalue(`var g_appid =encodeURIComponent\("(\d+)"\);`,html)
	sign, _ := getRevalue(`var g_login_sig=encodeURIComponent\("(.*?)"\);`,html)
	js_ver, _ := getRevalue(`var g_pt_version=encodeURIComponent\("(\d+)"\);`,html)
	mibao_css, _ := getRevalue(`var g_mibao_css=encodeURIComponent\("(.+?)"\);`,html)

	fmt.Printf("ptui_checkVC is %s\n", src)
	fmt.Printf("html is %+v\n", html)
	fmt.Printf("appid is %s\n", appid)
	fmt.Printf("sign is %s\n", sign)
	fmt.Printf("js_ver is %s\n", js_ver)
	fmt.Printf("mibao_css is %s\n", mibao_css)
	




	star_time = date_to_millis(datetime.datetime.utcnow())

        error_times = 0
        ret = []
        while True:
            error_times += 1
            self.req.Download('https://ssl.ptlogin2.qq.com/ptqrshow?appid={0}&e=0&l=L&s=8&d=72&v=4'.format(appid),
                              self.qrcode_path)
            logging.info("Please scan the downloaded QRCode")

            while True:
                html = self.req.Get(
                    'https://ssl.ptlogin2.qq.com/ptqrlogin?webqq_type=10&remember_uin=1&login2qq=1&aid={0}&u1=http%3A%2F%2Fw.qq.com%2Fproxy.html%3Flogin2qq%3D1%26webqq_type%3D10&ptredirect=0&ptlang=2052&daid=164&from_ui=1&pttype=1&dumy=&fp=loginerroralert&action=0-0-{1}&mibao_css={2}&t=undefined&g=1&js_type=0&js_ver={3}&login_sig={4}'.format(
                        appid, date_to_millis(datetime.datetime.utcnow()) - star_time, mibao_css, js_ver, sign),
                    initurl)
                logging.debug("QRCode check html:   " + str(html))
                ret = html.split("'")
                if ret[1] in ('0', '65'):  # 65: QRCode 失效, 0: 验证成功, 66: 未失效, 67: 验证中
                    break
            if ret[1] == '0' or error_times > 10:
                break

        if ret[1] != '0':
            return
        logging.info("QRCode scaned, now logging in.")


	return 0, nil

}
