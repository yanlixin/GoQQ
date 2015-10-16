package main

import (
//	"os"
//	"unsafe"
)

var cmdSend = &Command{
	UsageLine: "Send Message",
	Short:     "Send a QQ message",
	Long: `
	Run comand will Send All friend QQ message `,
}

//var mainFiles ListOpts
//var downdoc docValue
//var gendoc docValue
//var excludedPaths strFlags

func init() {
	cmdSend.Run = SendMsg
	//	cmdRun.Flag.Var(&mainFiles, "main", "specify main go files")
	//	cmdRun.Flag.Var(&gendoc, "gendoc", "auto generate the docs")
	//	cmdRun.Flag.Var(&downdoc, "downdoc", "auto download swagger file when not exist")
	//	cmdRun.Flag.Var(&excludedPaths, "e", "Excluded paths[].")
}

func SendMsg(cmd *Command, args []string) int {
	ColorLog("[INFO] 发送消息: %s\r\n", args)
	/*
		c, _ = f.Read(b)
		bb = b[:c-1]
		str = *(*string)(unsafe.Pointer(&bb))
		ColorLog("[INFO] 输入的消息为: %s\r\n", str)
		if "y" == str {
			ColorLog("[INFO] 消息发送成功 !%s\r\n")
		}

		//curpath, _ := os.Getwd()
		//fmt.Println(curpath)
		//ColorLog("[INFO] Uses '%s' as 'appname'\n", appname)
		status, err := SendMessage()
		if nil != err {

			ColorLog("[ERRO] Send fail ,Error:%+v\n", err)
			return 2
		}
		if status == 0 {
			ColorLog("[INFO] Send succeed\n")

		} else {

			ColorLog("[INFO] Send fail,Code:%d\n", status)
		} //ips := []string{"a", "b", "c"}
		//writeFile(ips)
		//fmt.Println(len(iparr))
		//	ips := TestIP(ipMap)
	*/
	//writeFile(ips)
	ColorLog("[INFO] 消息发送成功 !\r\n")
	return 0
}
