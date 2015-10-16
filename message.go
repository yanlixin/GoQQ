package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

func HandleMessage() error {
	return nil
}
func SendMessage() (int, error) {
	return 0, nil

}
func buddyMsgStructer(uiuin string, msg_id int64, msg, fontname, fontsize, fontcolor string, fontstyle [3]int) url.Values {
	uin, _ := strconv.Atoi(uiuin)
	v := url.Values{}
	v.Set(`clientid`, qq.clientid)
	v.Set(`psessionid`, qq.psessionid)
	ms := [2]interface{}{
		msg,
		[2]interface{}{
			"font",
			map[string]interface{}{
				"name":  fontname,
				"size":  fontsize,
				"style": fontstyle,
				"color": fontcolor}}}
	byts, _ := json.Marshal(ms)
	m := map[string]interface{}{
		"to":         uin,
		"face":       0,
		"content":    string(byts),
		"msg_id":     msg_id,
		"clientid":   qq.clientid,
		"psessionid": qq.psessionid,
	}

	byts, _ = json.Marshal(m)
	v.Set(`r`, string(byts))
	return v
}

func SendBuddyMsgEasy(uin string, msg_id int64, msg string) (code int, err error) {
	code, err = SendBuddyMsg(uin, msg_id, msg, `宋体`, `15`, `000000`, [3]int{0, 0, 0})
	return
}

func SendBuddyMsg(uin string, msg_id int64, msg, fontname, fontsize, fontcolor string, fontstyle [3]int) (code int, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
		qq.msgid++
	}()
	v := buddyMsgStructer(uin, qq.msgid, msg, fontname, fontsize, fontcolor, fontstyle)
	re, err := HttpPost(`http://d.web2.qq.com/channel/send_buddy_msg2`, v)
	if err != nil {
		return -2, err
	}
	defer re.Body.Close()

	retb := ReadByte(re.Body)
	var js map[string]interface{}
	json.Unmarshal(retb, &js)
	DebugLog("%+v", js)

	if js[`retcode`] != float64(0) {
		DebugLog(fmt.Sprintf("发送个人消息:%v 失败，错误代码:%v", msg, js[`retcode`]))
		return -3, errors.New(fmt.Sprintf("发送个人消息:%v 失败，错误代码:%v", msg, js[`retcode`]))
	}
	return 0, nil
}
func GetFriends() (code string, err error) {

	v := url.Values{}
	m := map[string]interface{}{
		"vfwebqq": qq.vfwebqq,
		"hash":    CalcHash("2951916302", qq.ptwebqq),
	}

	byts, _ := json.Marshal(m)
	v.Set(`r`, string(byts))
	re, err := HttpPost(`http://s.web2.qq.com/api/get_user_friends2`, v)
	if err != nil {
		return "", err
	}
	defer re.Body.Close()

	retb := ReadByte(re.Body)
	var js map[string]interface{}
	json.Unmarshal(retb, &js)
	DebugLog("%+v", js)
	return
}

//#id, ptwebqq

/*
func getHash(uin,ptwebqq)
{
    var x
    x= u(uin,ptwebqq);
	return x;
};

func u(x, K)
{
    x += ""
    for (var N = [], T = 0; T < K.length; T++) {
    	N[T % 4] ^= K.charCodeAt(T);
    }
    var U = [
      'EC',
      'OK'
    ],
    V = [
    ];
    V[0] = x >> 24 & 255 ^ U[0].charCodeAt(0);
    V[1] = x >> 16 & 255 ^ U[0].charCodeAt(1);
    V[2] = x >> 8 & 255 ^ U[1].charCodeAt(0);
    V[3] = x & 255 ^ U[1].charCodeAt(1);
    U = [
    ];
    for (T = 0; T < 8; T++) U[T] = T % 2 == 0 ? N[T >> 1] : V[T >> 1];
    N = [
      '0',
      '1',
      '2',
      '3',
      '4',
      '5',
      '6',
      '7',
      '8',
      '9',
      'A',
      'B',
      'C',
      'D',
      'E',
      'F'
    ];
    V = '';
    for (T = 0; T < U.length; T++) {
      V += N[U[T] >> 4 & 15];
      V += N[U[T] & 15]
    }
    return V
}
*/
