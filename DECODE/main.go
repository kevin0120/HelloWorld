package main

import (
	"encoding/json"
	"fmt"
)

type User struct {
	UserName string `json:"username"`
	NickName string `json:"nickname"`
	Age      int
	Birthday string
	Sex      string
	Email    string
	Phone    string
}

func testStruct() (ret string, err error) {
	user1 := &User{
		UserName: "user1",
		NickName: "上课看似",
		Age:      18,
		Birthday: "2008/8/8",
		Sex:      "男",
		Email:    "mahuateng@qq.com",
		Phone:    "110",
	}

	data, err := json.Marshal(user1)
	if err != nil {
		err = fmt.Errorf("json.marshal failed, err:", err)
		return
	}

	ret = string(data)
	fmt.Println(ret)
	return
}

func test() {
	data, err := testStruct()
	fmt.Printf("Field%v\n", []byte(data))
	fmt.Println(data[0:5])
	if err != nil {
		fmt.Println("test struct failed, ", err)
		return
	}

	//var user1 User

	var user1 User
	err = json.Unmarshal([]byte(data), &user1)
	if err != nil {
		fmt.Println("Unmarshal failed, ", err)
		return
	}
	fmt.Println(user1)

	var user2 map[string]interface{}
	err = json.Unmarshal([]byte(data), &user2)
	if err != nil {
		fmt.Println("Unmarshal failed, ", err)
		return
	}
	fmt.Println(user2)
	/*
		A:=user2.(User)
		fmt.Println(A)
	*/
	json.Marshal()
	fmt.Println(user2["nickname"])

}

func main() {
	test()
	// test2()
}
