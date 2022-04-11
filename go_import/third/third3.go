package third

import (
	"fmt"
	"time"
)

var _ = func() error {
	fmt.Printf("third3.go下第一个全局变量运行时间%s\n", time.Now())
	return nil
}()

func init() {
	fmt.Printf("third3.go下init时间%s\n", time.Now())
}

var _ = func() error {
	fmt.Printf("third3.go下第二个全局变量运行时间%s\n", time.Now())
	return nil
}()
