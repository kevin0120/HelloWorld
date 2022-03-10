//go:generate stringer -type=Pill

package main

import "fmt"

type Pill int

const (
	Placebo Pill = iota
	Aspirin
	Ibuprofen
	Paracetamol
	Acetaminophen = Paracetamol
)

func main()  {
	fmt.Println(Placebo.String())
	fmt.Println(Aspirin.String())
	fmt.Println(Ibuprofen.String())
	fmt.Println(Paracetamol.String())
}