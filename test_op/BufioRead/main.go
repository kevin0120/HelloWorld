package main

import (
	"bufio"
	"bytes"
	"fmt"
)

var s="hello\nworld\nhhhh\ndd"

func main()  {
	fmt.Println(len(s))
	//B:= []byte{0x4A, 0x01, 0x5D, 0x1a, 0x01, 0x01, 0x4c, 0x01, 0x50, 0xfa, 0x01, 0x01}
	scanner:=bufio.NewScanner( bytes.NewReader([]byte(s)))


	splitFunc:= func(data []byte, atEOF bool) (advance int, token []byte, err error) {

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

	for scanner.Scan(){
		h:=scanner.Bytes()
		fmt.Println(string(h))

	}


}



func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}