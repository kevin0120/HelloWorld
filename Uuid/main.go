package main

import (
	"encoding/base64"
	"fmt"
	"github.com/satori/go.uuid"
	"time"
)

func GenerateID() string {
	u4 := uuid.NewV4()
	//fmt.Println(u4.String())
	return base64.RawURLEncoding.EncodeToString(u4.Bytes())
}

func main() {
	for i:=0;i<132 ;i++ {
		fmt.Println(i+1,".",GenerateID())
		time.Sleep(15 * time.Millisecond)
	}
}
