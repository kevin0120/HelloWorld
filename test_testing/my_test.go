package test_testing

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

//func TestErrorInCode(t *testing.T) {
//	fmt.Println("Start")
//	t.Error("Error")
//	fmt.Println("End")
//	/** 运行结果：
//	  === RUN   TestErrorInCode
//	  Start
//	      TestErrorInCode: functions_test.go:25: Error
//	  End
//	  --- FAIL: TestErrorInCode (0.00s)
//	*/
//}
//
//func TestFatalInCode(t *testing.T) {
//	fmt.Println("Start")
//	t.Fatal("Error")
//	fmt.Println("End")
//	/** 运行结果：
//	  === RUN   TestFatalInCode
//	  Start
//	      TestFatalInCode: functions_test.go:38: Error
//	  --- FAIL: TestFatalInCode (0.00s)
//	*/
//}

func square(op int) int {
	return op * op
}

func TestSquareWithAssert(t *testing.T) {
	inputs := [...]int{1, 2, 3}
	expected := [...]int{1, 4, 9}
	for i := 0; i < len(inputs); i++ {
		ret := square(inputs[i])
		assert.Equal(t, expected[i], ret)
	}
}

// 利用+=连接
