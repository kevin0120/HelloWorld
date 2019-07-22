package testmyticker

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

/**

  每个测试文件必须以 _test.go 结尾，不然 go test 不能发现测试文件
  每个测试文件必须导入 testing 包
  功能测试函数必须以 Test 开头，然后一般接测试函数的名字，这个不强求

*/
func test(a ...interface{}) (eorr error) {

	fmt.Println("现在时间是:", a[0], a[1])
	fmt.Println(time.Now())
	return fmt.Errorf("这是现在时间")
}
func test1(a ...interface{}) (eorr error) {

	fmt.Println("test1现在时间是:", a[0], a[1])
	fmt.Println(time.Now())
	return fmt.Errorf("这是现在时间")
}
func TestFlutter(t *testing.T) {
	tick := New(test, 2*time.Second)
	//	go tick.Handle(test)
	//fmt.Println(&tick.Func)
	time.Sleep(1580 * time.Millisecond)
	for i := 0; i < 10; i++ {
		error := tick.Flutter(test, i, "hhh")
		assert.Nil(t, error)
		if error != nil {
			break
		}

		time.Sleep(200 * time.Millisecond)
	}

	time.Sleep(10 * time.Second)
}
