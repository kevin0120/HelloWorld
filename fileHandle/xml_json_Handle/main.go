package main

import (
	"encoding/json"
	"fmt"
)

type Student struct {
	Name  string `json:"name"` //序列化时将字段变为小写
	Age   int    `json:"age"`
	Score int    `json:"score"`
}

func main() {

	//序列化
	stu1 := &Student{Name: "wd", Age: 22, Score: 100}
	res, err := json.Marshal(stu1)
	if err != nil {
		fmt.Println("json encode error")
		return
	}
	fmt.Printf("json string: %s\n", res)
	fmt.Println("#####################################")
	//反序列化
	data := `{"name":"wd","age":22,"score":100}`
	var stu2 Student
	err1 := json.Unmarshal([]byte(data), &stu2)
	if err1 != nil {
		fmt.Println("json decode error: ", err1)
		return
	}
	fmt.Printf("struct obj is : %s", stu2.Name)
}

// 结果
//json string :{"name":"wd","age":22,"score":100}
