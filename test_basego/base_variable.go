package main

import "fmt"

func main() {
	const LENGTH int = 10
	const WIDTH int = 5
	var area int
	const a1, b1, c1 = 1, false, "str" //多重赋值

	area = LENGTH * WIDTH
	fmt.Printf("面积为 : %d", area)
	fmt.Println()
	hello := "hello world"
	fmt.Println(a1, b1, c1, hello)

	var a = 21
	var b = 10
	var c int

	c = a + b
	fmt.Printf("第一行 - c 的值为 %d\n", c)
	c = a - b
	fmt.Printf("第二行 - c 的值为 %d\n", c)
	c = a * b
	fmt.Printf("第三行 - c 的值为 %d\n", c)
	c = a / b
	fmt.Printf("第四行 - c 的值为 %d\n", c)
	c = a % b
	fmt.Printf("第五行 - c 的值为 %d\n", c)
	a++
	fmt.Printf("第六行 - a 的值为 %d\n", a)
	a = 21 // 为了方便测试，a 这里重新赋值为 21
	a--
	fmt.Printf("第七行 - a 的值为 %d\n", a)
}
