package main

import (
	"fmt"
	"HelloWorld/myticker/testmyticker"
	"time"
)

func Test(a ...interface{}) (eorr error) {

	fmt.Println("现在时间是:", a[0], a[1])
	fmt.Println(time.Now())
	return fmt.Errorf("这是现在时间")
}
func Test1(a ...interface{}) (eorr error) {

	fmt.Println("Test1现在时间是:", a[0], a[1])
	fmt.Println(time.Now())
	return fmt.Errorf("这是现在时间")
}
func main() {
	tick := testmyticker.New(Test, 2*time.Second)
	//	go tick.Handle(test)
	time.Sleep(1580 * time.Millisecond)
	for i := 0; i < 100; i++ {
		error := tick.Flutter(Test1, i, "hhh")
		if error != nil {
			fmt.Println(error)
			break
		}

		time.Sleep(200 * time.Millisecond)
	}

	time.Sleep(1 * time.Second)
}
