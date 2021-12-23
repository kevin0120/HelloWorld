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
	return base64.StdEncoding.EncodeToString(u4.Bytes())
}

func GenerateID1() string {
	u4 := uuid.NewV4()
	//fmt.Println(u4.String())
	return base64.URLEncoding.EncodeToString(u4.Bytes())
}
func GenerateID2() string {
	u4 := uuid.NewV4()
	//fmt.Println(u4.String())
	return base64.RawStdEncoding.EncodeToString(u4.Bytes())
}
func GenerateID3() string {
	u4 := uuid.NewV4()
	//fmt.Println(u4.String())
	return base64.RawURLEncoding.EncodeToString(u4.Bytes())
}



func main() {
	for i := 0; i < 132; i++ {
		fmt.Println(i+1, ".", GenerateID())
		fmt.Println(i+1, ".", GenerateID1())
		fmt.Println(i+1, ".", GenerateID2())
		fmt.Println(i+1, ".", GenerateID3())
		time.Sleep(15 * time.Millisecond)
	}
}
