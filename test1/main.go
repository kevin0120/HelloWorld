package main

import (
	"flag"
	"fmt"
)

func main() {
	//fs := flag.NewFlagSet("", flag.ContinueOnError)
	//若加了这一条,则后面fs.Parse(参数),需要加参数.
	infile := flag.String("i", "infile", "File contains values for sorting")
	outfile := flag.String("o", "outfile", "File to receive sorted values")
	algorithm := flag.String("a", "qsort", "Sort algorithm")

	flag.Parse()
	if infile != nil {
		fmt.Println("infile =", *infile, "outfile =", *outfile, "algorithm =", *algorithm)
	}
	var a []byte
	var b []byte
	a = append(a, 0x55)
	b = append(b, 0x00)
	a = append(a, b[0])
	fmt.Println(len(a))
	//	a=append(a, b[0:0]...)

	copy(a, b[0:0])

	fmt.Println(len(a))
	/*
		fmt.Println(string(a))

		c:=strings.Split(string(a),string(b))
		fmt.Println(c[0],len(c[0]))
		fmt.Println(c[1],len(c[1]))
	*/
}
