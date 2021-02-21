package main

import (
	"fmt"
)

type animal interface {
	Say()
	Jump()
}

type human struct {
	age int
}

func (h *human) Say() {
	fmt.Println("我是一个人")
}

type man struct {
	human
	hith int
}

func (m *man) Jump() {
	fmt.Println("我跳的 好高")
}

/*
func (m man) Say() {
	fmt.Println("我是一个男人")
}
*/
func main() {

	s := []animal{}

	h := &human{12}
	m := &man{*h,
		30,
	}

	m.Jump()
	m.Say()

	//var a animal
	//a=m
	//a.Jump()
	//a.Say()

	s = append(s, m)
	s[0].Jump()
	s[0].Say()
}
