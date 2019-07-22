package testmyticker

import (
	"fmt"
	"reflect"
	"runtime"
	"time"
)

type f func(a ...interface{}) (error error)

type Ticker struct {
	T          *time.Ticker
	CountFirst bool
	NeedHandle bool
	Par        []interface{}
	Func       string
}

func New(method f, d time.Duration) (a *Ticker) {

	return &Ticker{T: time.NewTicker(d),
		CountFirst: true,
		NeedHandle: false,
		Func:       runtime.FuncForPC(reflect.ValueOf(method).Pointer()).Name(),
	}
}

func (b *Ticker) Handle(method f) {
	for {
		select {
		case <-b.T.C:
			if !b.CountFirst && b.NeedHandle {
				method(b.Par[0:]...)
				b.NeedHandle = false
			} else {
				break
			}
		default:
			if b.CountFirst && b.NeedHandle {
				method(b.Par[0:]...)
				b.NeedHandle = false
				b.CountFirst = false
			}
			break
		}
	}

}

func (b *Ticker) Flutter(method f, a ...interface{}) (error error) {
	//fmt.Println(&b.Func,&method)
	funcname := runtime.FuncForPC(reflect.ValueOf(method).Pointer()).Name()
	if funcname == b.Func {
		b.Par = a[0:]
		if !b.NeedHandle {
			go b.Handle(method)
		}
		b.NeedHandle = true
		return nil
	} else {
		return fmt.Errorf("此次加入节流器的函数与节流器初始化时的不一致!!!")
	}
}
