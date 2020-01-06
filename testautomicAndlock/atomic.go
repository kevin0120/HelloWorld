package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

/*原子操作---增或减、比较并交换、载入、存储和交换
在这些操作的过程中保证不被其他协程调用,但不能保证变量读出后其他协程不再访问,
	if *value == 1 {
			fmt.Println(name, "value==1", *value)
		}
如上例:value值在满足if判断条件后可能发生改变!!!!!!
*/
var ConfValue atomic.Value

func main() {
	var value int32

	var value1 int32
	var value2 int32
	fmt.Println("初始值:", value1, value2)
	/*存储*******原子操作**********/
	atomic.StoreInt32(&value1, 3)
	atomic.StoreInt32(&value2, 5)
	fmt.Println("原子操作存储后:", value1, value2)
	/*交换*******原子操作***不用判断旧值是否正确返,回值为旧值,*******/
	fmt.Println(atomic.SwapInt32(&value1, 4))
	fmt.Println("原子操作value1交换后:", value1, value2)

	/*比较并交换*******原子操作***判断旧值是否正确,正确才会替换原值,返回bool*******/
	fmt.Println(atomic.CompareAndSwapInt32(&value1, 4, 3),
		atomic.CompareAndSwapInt32(&value2, 4, 100))
	fmt.Println("原子操作比较并交换后:", value1, value2)

	/*增或减*******原子操作**********/
	fmt.Println("原子操作增或减:", atomic.AddInt32(&value1, -1),
		atomic.AddInt32(&value2, 1))

	///原子操作赋值
	var c int32
	c = 45
	ConfValue.Store(c)
	fmt.Println(ConfValue.Load().(int32))

	/*载入*******原子操作**********/
	/*注意，虽然我们在这里使用atomic.LoadInt32函数原子的载入value的值，
	但是其后面的CAS操作仍然是有必要的。因为，那条赋值语句和if语句并不会被原子的执行。
	在它们被执行期间，CPU仍然可能进行其它的针对value的值的读或写操作。也就是说，
	value的值仍然有可能被并发的改变。*/

	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("222222222222222222222222")
		fmt.Println(value1)
		atomic.AddInt32(&value1, 1)
		fmt.Println(value1)
	}()

	fmt.Println(time.Now())
	fmt.Println(time.Now().UTC())
	fmt.Println(time.Now().UnixNano())

	func(delta int32) {
		for {
			fmt.Println("1111111111111111111111111111")
			v := atomic.LoadInt32(&value1)
			fmt.Println(v)
			time.Sleep(3 * time.Second)
			if atomic.CompareAndSwapInt32(&value1, v, (v + delta)) {
				break
			}
		}
	}(10)
	fmt.Println("原子操作载入并+10后:", value1, value2)

	fmt.Println("origin value:", value)

	go entry1("1", &value)

	entry2("2", &value)

	//	time.Sleep(time.Second)
}

func entry2(name string, value *int32) {

	for i := 0; i < 10; i++ {
		swapFlag := atomic.CompareAndSwapInt32(value, 0, 1)
		fmt.Println(name, swapFlag)
		if *value == 1 {
			fmt.Println(name, "value==1", *value)
		} else {
			fmt.Println(name, "value==0", *value)
		}
		time.Sleep(1 * time.Second)

	}

}
func entry1(name string, value *int32) {

	for {
		swapFlag := atomic.CompareAndSwapInt32(value, 1, 0)
		fmt.Println(name, swapFlag)
		if *value == 1 {
			fmt.Println(name, "value==1", *value)
		} else {
			time.Sleep(3 * time.Second)
			fmt.Println(name, "value==0", *value)
		}
		time.Sleep(3 * time.Second)

	}

}
