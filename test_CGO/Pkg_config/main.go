package main

// #cgo pkg-config: libhello
// #include < stdlib.h >
// #include < hello_world.h >
import "C"
import (
	"unsafe"
)

func main() {
	msg := "Hello, world!"
	cmsg := C.CString(msg)
	C.hello(cmsg)
	C.free(unsafe.Pointer(cmsg))
}

//https://www.cnblogs.com/mokliu/p/5538926.html
