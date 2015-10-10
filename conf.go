package main

import (
	"encoding/json"
	"os"
)

const CONF_VER = 0

var defaultConf = `{
	"version":0,
	"gopm":{"enable":false,
		"install":false
	}
	"cmmd_args":[],
	"evns":[],
	"ConnectReferer":"http://d.web2.qq.com/proxy.html?v=20030916001&callback=1&id=2",
	"SmartQQUrl":"http://w.qq.com/login.html",
	"QRCodePath":"./v.jpg",
	"TucaoPath":"./data/tucao_save/",


}
`
var conf struct {
	Version int `json:"version"`
	Gopm    struct {
		Enable  bool
		Install bool
	}
	CmdArgs        []string `json:"cmd_args"`
	Envs           []string
	ConnectReferer string `json:"ConnectReferer"`
	SmartQQUrl     string `json:"SmartQQUrl"`
	QRCodePath     string `json:"QRCodePath"`
	TucaoPath      string `json:"TucaoPath"`
}

func loadConfig() error {
	f, err := os.Open("goqq.json")
	if nil != err {
		err = json.Unmarshal([]byte(defaultConf), &conf)
		if nil != err {
			return err
		}
	} else {
		defer f.Close()
		ColorLog("[INFO] Detected goqq.json\n")
		d := json.NewDecoder(f)
		err = d.Decode(&conf)
		if nil != err {
			return err
		}
	}
	if CONF_VER != conf.Version {
		ColorLog("[WARN] Your goqq.json is out-of-date,please update!\n")
		ColorLog("[HINT] Compare goqq.json under goqq source code path and yours\n")
	}
	return nil
}
