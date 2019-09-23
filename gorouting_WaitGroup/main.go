package main

import "sync"

//rang 里的i 是共享的如果range下多个协程访问i会错乱,例如print(a)改成print(i)
//如果要吧waitGROUP 当参数传入函数,一定要传指针!!!!
func main() {

	wg := sync.WaitGroup{}
	si := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	for i := range si {
		wg.Add(1)
		go func(a int) {
			println(a)
			wg.Done()
		}(i)
	}
	wg.Wait()
}