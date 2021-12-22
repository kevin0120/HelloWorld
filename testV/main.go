package main
//https://www.cnblogs.com/xinliangcoder/p/13282964.html
import (
	"encoding/json"
	"fmt"
)


func main ()  {
	V(9999999)
	V(113106)

}
func V(TighteningId int64)  {
	var b interface{}
	jsonStrs, _ := json.Marshal(TighteningId)
	_= json.Unmarshal(jsonStrs, &b)
	fmt.Println(fmt.Sprintf("i am %v",b))
	fmt.Println(fmt.Sprintf("i am %v",int64(b.(float64))))
}
