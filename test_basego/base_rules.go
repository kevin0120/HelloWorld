package main

import (
	"fmt"
)

func testIf() {
	/* 定义局部变量 */
	var a = 20
	/* 使用 if 语句判断布尔表达式 */
	if a < 20 {
		/* 如果条件为 true 则执行以下语句 */
		fmt.Printf("a 小于 20\n")
	} else if a > 20 {
		fmt.Printf("a 大于 20\n")
	} else {
		fmt.Printf("a 大于等于 20\n")
	}
	fmt.Printf("a 的值为 : %d\n", a)

}

func testSwitch() {
	var b = 2
	switch b {
	case 1:
		fmt.Printf("hello %v", 1)
		break
	case 2:
		fmt.Printf("hello %v", 2)
		fallthrough
	default:
		fmt.Printf("hello %v", "这是一个异常的数字")
	}
}

func testFor() {
	var a int
	var b int
	for {

		fmt.Printf("请输入密码：   \n")
		_, _ = fmt.Scan(&a)
		fmt.Println(a)
		if a == 5211314 {
			_, _ = fmt.Scan(&b)
			if b == 5211314 {
				fmt.Printf("密码正确，门锁已打开")
				break
			} else {
				fmt.Printf("非法入侵，已自动报警")
			}
		} else {
			fmt.Printf("非法入侵，已自动报警")
		}
	}
}

func testForrange() {
	strings := make([]string, 10, 19)
	fmt.Println(len(strings))
	strings = []string{"google", "baidu"}
	fmt.Println(len(strings))
	strings = append(strings, "souhu")
	fmt.Println(len(strings))
	for i, s := range strings {
		fmt.Println(i, s, len(strings))
	}

	numbers := [6]int{1, 2, 3, 5}
	for i, x := range numbers {
		fmt.Println(i, x, len(numbers))
	}

	for i, s := range "hello world!" {
		fmt.Println(i, string(s))
	}

	var map1 map[string]string
	map1 = make(map[string]string)
	fmt.Println(len(map1))
	/* map插入key - value对,各个国家对应的首都 */
	map1["France"] = "巴黎"
	map1["Italy"] = "罗马"
	map1["Japan"] = "东京"
	map1["India"] = "新德里"
	fmt.Println(map1)
	fmt.Println(len(map1))
	d, y := map1["china"]
	fmt.Println(d, "hello china", y)
	d, y = map1["India"]
	fmt.Println(d, "hello India", y)

	fmt.Println(map1["china"])
	for k, v := range map1 {
		fmt.Println(k, v)
	}
}

func main() {
	fmt.Println("####################test_if#######################")
	testIf()
	fmt.Println("####################test_switch#######################")
	testSwitch()
	fmt.Println("####################testForrange#######################")
	//map 和slice是不同的make出来的len（）
	testForrange()
	fmt.Println("####################testFor#######################")
	testFor()
}
