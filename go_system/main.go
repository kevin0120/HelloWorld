package main

import (
	"fmt"
	"github.com/kevin0120/HelloWorld/go_system/drivers"
	"os"
)

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func main() {
	//a := "\"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&#43;AAAAAAAAAAAAAAAAAAAAAAAMMD/JQAAAACQkAAAAAAAAAAAAAAAAAAAAAAAAAAAAQBQcmludEFycmF5AAAAAAIAAAACAAAABAAAAAAABwAAAAMAAAAAAAQAAAADAAAAAAAEAAAAAwAudGV4dAAAAAAAAAABAAAAAwAuaWRhdGEkNwAAAAACAAAAAwAuaWRhdGEkNQAAAAADAAAAAwAuaWRhdGEkNAAAAAAEAAAAAwAuaWRhdGEkNgAAAAAFAAAAAwAAAAAABAAAAAAAAAABAAAAAgAAAAAADwAAAAAAAAADAAAAAgAAAAAAIAAAAAAAAAAAAAAAAgA1AAAAUHJpbnRBcnJheQBfX2ltcF9QcmludEFycmF5AF9oZWFkX2xpYnNwY19saWJfZGxsAAo=\""
	//var bb []byte
	//err := json.Unmarshal([]byte(a), &bb)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//c,_:=json.Marshal(bb)
	//d:=string(c)
	//fmt.Println(d)

	drivers.Nkio()
	fmt.Println("hello world!!")
}
