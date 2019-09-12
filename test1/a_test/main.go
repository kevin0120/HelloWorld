package test1

import (
	"encoding/binary"
	"fmt"
)

func main() {
	var testInt int32 = 0x01020304 // 十六进制表示
	fmt.Printf("%d use big endian: \n", testInt)

	var testBytes []byte = make([]byte, 4)
	binary.BigEndian.PutUint32(testBytes, uint32(testInt)) //大端序模式
	fmt.Println("int32 to bytes:", testBytes)
}