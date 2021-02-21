package main

import (
	"fmt"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
	"time"
)

//參考文檔
//https://github.com/go-xorm/xorm/blob/master/README_CN.md
//http://gobook.io/read/github.com/go-xorm/manual-zh-CN/chapter-01/
type User struct {
	Id      int64 //唯一约束Id
	Name    string
	Salt    string
	Age     int
	Passwd  string    `xorm:"varchar(200)"`
	Created time.Time `xorm:"created"`
	Updated time.Time `xorm:"updated"`
}

//`xorm:"varchar(25) notnull unique 'usr_name'"`
func main() {

	engine, _ := xorm.NewEngine("postgres", "user=odoo password=odoo dbname=test sslmode=disable")
	//则会在控制台打印出生成的SQL语句；
	//

	//如果User为大写则不存在，若ｕｓｅｒ为小写则存在．
	//exist, err := engine.IsTableExist("User")
	exist, err := engine.IsTableExist("user")
	if err == nil {
		if !exist {
			if err := engine.Sync2(new(User)); err != nil {
				fmt.Println("you BUG")
			}
		}
	}

	//engine.ShowSQL(true)
	users1 := User{
		//Id:   2,
		Name: "FJLDJFLDFLD",
		Age:  33,
	}
	affected, err := engine.Insert(&users1)
	// INSERT INTO struct () values ()
	fmt.Println(affected, err)
	//engine.Ping()
	//有的数据库超时断开ping可以重连。可以通过起一个定期Ping的Go程来保持连接鲜活。

	a, _ := engine.DBMetas()
	fmt.Println(a)

	engine.ShowSQL(true)
	var users2 []User
	ss := engine.Alias("r").Where("r.name = ?", "FJLDJFLDFLD")

	ss.Find(&users2)

	fmt.Println(users2)

}
