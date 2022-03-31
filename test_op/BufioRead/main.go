package main

import (
	"bufio"
	"bytes"
	"fmt"
)

var s = "hello\nworld\nhhhh\ndd"

var s01 = "hello         worldhhhh                                                                                  r00                 dd\n"
var s02 = "hello         worldhhhh                                                                                  r00                 dd\n"
var s03 = "hello         worldhhhh                                                                                  r00                 dd\n"
var s04 = "hello         worldhhhh                                                                                  r00                 dd\n"
var s05 = "hello         worldhhhh                                                                                  r00                 dd\n"
var s06 = "hello         worldhhhh                                                                                  r00                 dd\n"
var s07 = "hello         worldhhhh                                                                                  r00                 dd\n"
var s08 = "hello         worldhhhh                                                                                  r00                 dd\n"
var s09 = "hello         worldhhhh                                                                                  r00                 dd\n"
var s10 = "hello         worldhhhh                                                                                  r00                 dd\n"

func main() {
	s = s01 + s02 + s03 + s04 + s05 + s06 + s07 + s08 + s09 + s10
	s += s
	s += s
	s += s

	fmt.Println(len(s))

	//testBigString(s)

	//B:= []byte{0x4A, 0x01, 0x5D, 0x1a, 0x01, 0x01, 0x4c, 0x01, 0x50, 0xfa, 0x01, 0x01}
	scanner := bufio.NewScanner(bytes.NewReader([]byte(s)))

	splitFunc := func(data []byte, atEOF bool) (advance int, token []byte, err error) {

		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		if i := bytes.IndexByte(data, '\n'); i >= 0 {
			// We have a full newline-terminated line.
			return i + 1, dropCR(data[0:i]), nil
		}
		// If we're at EOF, we have a final, non-terminated line. Return it.
		if atEOF {
			return len(data), dropCR(data), nil
		}
		// Request more data.
		return 0, nil, nil
	}

	scanner.Split(splitFunc)
	var i int
	for scanner.Scan() {
		_ = scanner.Bytes()
		i++
		fmt.Println(i)

	}

}

func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}

func testBigString(s string) {
	newBuf := make([]byte, 9)
	scanner := bytes.NewReader([]byte(s))
	for {

		n, err := scanner.Read(newBuf)

		fmt.Println(n, "@@@@@@@@@@@@@", newBuf[0:n])
		if err != nil {
			fmt.Println(err)
			break
		}

	}

}
