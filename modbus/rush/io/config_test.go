package io

import (
	"fmt"
	"testing"
	"time"
)

func TestConfig(t *testing.T) {
	//cfg := NewConfig()
	//err := cfg.Validate()
	//assert.Nil(t, err)
	//
	//cfg.IOS[0].Model = "wef"
	//err = cfg.Validate()
	//assert.NotNil(t, err)

	dateString := "2019-10-16T11:20:30+08:00"

	loc, _ := time.LoadLocation("Local")
	dt, _ := time.ParseInLocation(time.RFC3339, dateString, loc)

	fmt.Println(dt.Format(time.RFC3339))
}
