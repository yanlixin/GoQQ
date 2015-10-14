package main

import (
	"fmt"
	"testing"
)

func Test_loadConfig(t *testing.T) {
	err := loadConfig()
	if nil != err {
		t.Errorf("%+v", err)
		fmt.Sprintf("%+v", err)
	}

}
