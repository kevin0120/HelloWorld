package StandardKey

import "fmt"

type Standard struct {
	Name string

}


func (s1 *Standard)Connect(){
fmt.Printf("%s 连接成功\n",s1.Name)
}

//扫码枪一定要有回车，且光标要移动到命令窗口相应位置
func (s1 *Standard)Read() (string, error) {
	var s string
	_, err := fmt.Scan(&s)
	if err != nil {
		return "", err
	}
	return s, nil
}
