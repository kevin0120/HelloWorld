package main


/*
#cgo CFLAGS: -I ${SRCDIR}/include/
#cgo LDFLAGS: ${SRCDIR}/cmake-build-debug/libspc_lib.a
#include "library.h"

void hello1() {
    hello();
}
 */
import "C"

func main()  {
	C.hello1()
}