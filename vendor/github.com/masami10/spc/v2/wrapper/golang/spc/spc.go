package spc

/*
#cgo CFLAGS: -I../../../include/
#cgo LDFLAGS: ${SRCDIR}/../libs/libspc_lib.a -lstdc++

#include <stdlib.h>
#include "spc.h"

*/
import "C"
import (
	"errors"
	"unsafe"
)

func Cpk(data []float64, length uint, usl, lsl float64) (float64, error) {
	if len(data) == 0 {
		return 0, errors.New("data length error")
	}
	c := C.cpk((*C.double)(unsafe.Pointer(&data[0])), C.size_t(length), C.double(usl), C.double(lsl))
	return float64(c), nil
}

func Cmk(data []float64, length uint, usl, lsl float64) (float64, error) {
	if len(data) == 0 {
		return 0, errors.New("data length error")
	}
	c := C.cmk((*C.double)(unsafe.Pointer(&data)), C.size_t(length), C.double(usl), C.double(lsl))
	return float64(c), nil
}

func Cr(data []float64, length uint, usl, lsl float64) (float64, error) {
	if len(data) == 0 {
		return 0, errors.New("data length error")
	}
	c := C.cr((*C.double)(unsafe.Pointer(&data)), C.size_t(length), C.double(usl), C.double(lsl))
	return float64(c), nil
}

func Cp(data []float64, length uint, usl, lsl float64) (float64, error) {
	if len(data) == 0 {
		return 0, errors.New("data length error")
	}
	c := C.cp((*C.double)(unsafe.Pointer(&data)), C.size_t(length), C.double(usl), C.double(lsl))
	return float64(c), nil
}
