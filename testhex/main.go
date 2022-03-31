package main

import "encoding/hex"

func main() {
	// 转换的用的 byte数据
	byteData := []byte(`测试数据`)
	// 将 byte 装换为 16进制的字符串
	hexStringData := hex.EncodeToString(byteData)
	// byte 转 16进制 的结果
	println(hexStringData)

	/* ====== 分割线 ====== */

	// 将 16进制的字符串 转换 byte
	hexData, _ := hex.DecodeString(hexStringData)
	// 将 byte 转换 为字符串 输出结果
	println(string(hexData))
}
