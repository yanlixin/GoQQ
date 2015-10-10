package main

import (
	"fmt"
	"testing"
)

func Test_loadConfig(t *testing.T) {
	err := loadConfig()
	fmt.Printf("%+v\n", conf.Version)
	t.Log(conf)
	t.Log(err)
}
