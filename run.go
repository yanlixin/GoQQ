package main

import (
	"fmt"
	"os"
)

var cmdRun = &Command{
	UsageLine: "run [appname]",
	Short:     "run the appp and  start a tester google app",
	Long: `
	Run comand will `,
}

//var mainFiles ListOpts
//var downdoc docValue
//var gendoc docValue
//var excludedPaths strFlags

func init() {
	cmdRun.Run = runApp
	//	cmdRun.Flag.Var(&mainFiles, "main", "specify main go files")
	//	cmdRun.Flag.Var(&gendoc, "gendoc", "auto generate the docs")
	//	cmdRun.Flag.Var(&downdoc, "downdoc", "auto download swagger file when not exist")
	//	cmdRun.Flag.Var(&excludedPaths, "e", "Excluded paths[].")
}

var appname string

func runApp(cmd *Command, args []string) int {
	curpath, _ := os.Getwd()
	fmt.Println(curpath)
	ColorLog("[INFO] Uses '%s' as 'appname'\n", appname)
	status, err := LoginByQRCode()
	if nil != err {

		ColorLog("[ERRO] Login fail ,Error:%+v\n", err)
		return 2
	}
	if status == 0 {
		ColorLog("[INFO] Login succeed\n")

	} else {

		ColorLog("[INFO] Login fail,Code:%d\n", status)
	} //ips := []string{"a", "b", "c"}
	//writeFile(ips)
	//fmt.Println(len(iparr))
	//	ips := TestIP(ipMap)

	//writeFile(ips)
	return 0
}
