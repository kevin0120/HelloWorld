package main

import "fmt"

func testMap(m map[string]string) {

	m["china"] = "北京"

}



type HashMap struct {
	key string
	value string
	hashCode int
	next *HashMap
}

var table [16]*HashMap

func initTable() {
	for i := range table{
		table[i] = &HashMap{"","",i,nil}
	}
}

func getInstance() [16]*HashMap {
	if table[0] == nil {
		initTable()
	}
	return table
}

func genHashCode(k string) int{
	if len(k) == 0{
		return 0
	}
	var hashCode int = 0
	var lastIndex int = len(k) - 1
	for i := range k {
		if i == lastIndex {
			hashCode += int(k[i])
			break
		}
		hashCode += (hashCode + int(k[i]))*31
	}
	return hashCode
}

func indexTable(hashCode int) int{
	return hashCode%16
}

func indexNode(hashCode int) int {
	return hashCode>>4
}

func put(k string, v string) string {
	var hashCode = genHashCode(k)
	var thisNode = HashMap{k,v,hashCode,nil}

	var tableIndex = indexTable(hashCode)
	var nodeIndex = indexNode(hashCode)

	var headPtr  = getInstance()
	var headNode = headPtr[tableIndex]

	if (*headNode).key == "" {
		*headNode = thisNode
		return ""
	}

	var lastNode = headNode
	var nextNode = (*headNode).next

	for nextNode != nil && (indexNode((*nextNode).hashCode) < nodeIndex){
		lastNode = nextNode
		nextNode = (*nextNode).next
	}
	if (*lastNode).hashCode == thisNode.hashCode {
		var oldValue = lastNode.value
		lastNode.value = thisNode.value
		return oldValue
	}
	if lastNode.hashCode < thisNode.hashCode {
		lastNode.next = &thisNode
	}
	if nextNode != nil {
		thisNode.next = nextNode
	}
	return ""
}

func get(k string) string {
	var hashCode = genHashCode(k)
	var tableIndex = indexTable(hashCode)

	var headPtr = getInstance()
	var node = headPtr[tableIndex]

	if (*node).key == k{
		return (*node).value
	}

	for (*node).next != nil {
		if k == (*node).key {
			return (*node).value
		}
		node = (*node).next
	}
	return ""
}



func myMap(){
	getInstance()
	put("a","a_put")
	put("b","b_put")
	fmt.Println(get("a"))
	fmt.Println(get("b"))
	put("p","p_put")
	fmt.Println(get("p"))
}
func main() {
	//此处的长度没有什么实际意义 内存分配控件
	v := make(map[string]string,1)
	fmt.Println(len(v))
	v["Franch"] = "巴黎"
	testMap(v)
	//所有的能make出来的都是应用传递的
	fmt.Println(v)
	fmt.Println(len(v))

	fmt.Println("$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$定义一个自己的map$$$$$$$$$$$$$$$$$$$$$$$$$$$$")

	myMap()
}
