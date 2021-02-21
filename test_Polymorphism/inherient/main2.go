package main

import (
	"fmt"
	"reflect"
)

type AnimalIF interface {
	Sleep()
}
type Animal struct {
	MaxAge int
}

func (this *Animal) Sleep() {
	fmt.Println("Animal need sleep")
}

type Dog struct {
	Animal
}

//func (this *Dog) Sleep() {
//	fmt.Println("Dog need sleep")
//}

type Cat struct {
	Animal
}

func (this *Cat) Sleep() {
	fmt.Println("Cat need sleep")
}
func (this *Cat) Say() {
	fmt.Println("HELLO Cat!")
}

func Factory(name string) AnimalIF {
	switch name {
	case "dog":
		return &Dog{Animal{MaxAge: 20}}
	case "cat":
		return &Cat{Animal{MaxAge: 10}}
	default:
		panic("No such animal")
	}
}

func main() {
	animal := Factory("dog")
	animal.Sleep()

	animal = Factory("cat")
	animal.Sleep()

	fmt.Println(reflect.TypeOf(animal).String())
	animal.(*Cat).Say()
}
