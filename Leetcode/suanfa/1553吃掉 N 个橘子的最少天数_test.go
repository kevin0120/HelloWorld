package suanfa

import (
	"fmt"
	"testing"
)

func minDays(n int) int {

	if n == 1 {
		return 1
	}

	if n == 2 {
		return 2
	}

	if n == 3 {
		return 2
	}
	if n >2000000000{
		return 0
	}
	var i, j int

	i = minDays(n/2) + n%2+1

	j = minDays(n/3) + n%3+1

	return min(i, j)
}

func min(i int, j int) int {
	if i < j {
		j = i
	}
	return j
}

func Test_suanfa(t *testing.T) {
	fmt.Println(minDays(10))
	fmt.Println("########################")
	fmt.Println(minDays(55))
	fmt.Println(minDays(28))
	fmt.Println("########################")
	fmt.Println(minDays(54))

	fmt.Println(minDays(27))
	fmt.Println(minDays(14))
	//
	fmt.Println(fmt.Sprintf("吃掉所有的橘子需要 %d 天", minDays(56)))
}
