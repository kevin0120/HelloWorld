package first

import (
	"fmt"
	_ "github.com/kevin0120/HelloWorld/go_import/first/second"
	"time"
)

var _ = func() error {
	fmt.Printf("first2.go下第一个全局变量运行时间%s\n", time.Now())
	return nil
}()

func init() {
	fmt.Printf("first2.go下init时间%s\n", time.Now())
}

var _ = func() error {
	fmt.Printf("first2.go下第二个全局变量运行时间%s\n", time.Now())
	return nil
}()
