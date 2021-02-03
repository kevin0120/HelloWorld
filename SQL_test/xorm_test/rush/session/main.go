package main

import (
	"fmt"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"os"
	"time"
)

//參考文檔
//https://github.com/go-xorm/xorm/blob/master/README_CN.md
//http://gobook.io/read/github.com/go-xorm/manual-zh-CN/chapter-01/
type Rush struct {
	eng *xorm.Engine
}

type Money struct {
	Id  int64 `xorm:"pk autoincr notnull 'id'" json:"-"`
	RMB int64 `xorm:"bigint 'rmb'" json:"-"`
}

type Payrecord struct {
	Id  int64 `xorm:"pk autoincr notnull 'id'" json:"-"`
	RMB int64 `xorm:"bigint 'rmb'" json:"-"`
}

//`xorm:"varchar(25) notnull unique 'usr_name'"`
func main() {

	r := NewRush()
	fmt.Println("%%%%%%%%%%", time.Now().UnixNano())
	for i := 10; i < 20; i++ {
		go r.Pay(int64(i))
	}

	for {
		time.Sleep(time.Second)
	}
	//results, err := engine.Query("select (jsons::json#>>'{product,url}')::text from testjson where (jsons::json#>>'{code}')::text = 'MO0002'")
	//
	//fmt.Println(string(results[0]["text"]))
	//	//则会在控制台打印出生成的SQL语句；
	//engine.ShowSQL(true)
	//	engine.Sync2(new(User))

}

func NewRush() *Rush {
	info := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		"odoo",
		"odoo",
		"localhost",
		"new_database")
	engine, err := xorm.NewEngine("postgres", info)

	engine.DatabaseTZ = time.UTC
	engine.TZLocation = time.UTC

	if err != nil {
		errors.Wrapf(err, "Create postgres engine fail")
	}
	//_ = engine.Sync2(new(User))

	exist, err := engine.IsTableExist("money")
	if err == nil {
		if !exist {
			if err := engine.Sync2(new(Money)); err != nil {
				errors.Wrapf(err, "Create Table Money fail")
			}

			money := Money{
				RMB: 100,
			}
			_, err = engine.Insert(&money)

		}
	}

	exist, err = engine.IsTableExist("payrecord")
	if err == nil {
		if !exist {
			if err := engine.Sync2(new(Payrecord)); err != nil {
				errors.Wrapf(err, "Create Table Payrecord fail")
			}

		}
	}

	r := &Rush{
		eng: engine,
	}

	return r
}

func (s *Rush) Pay(in int64) error {

	// session
	// select if exist
	// if exist
	// 		update workorders set no = {id} where code = '' and no <= {id}
	// else
	// 		insert
	// end session
	session := s.eng.NewSession().ForUpdate()
	defer session.Close()
	session.Begin()
	//sql1 := fmt.Sprintf("insert into `payrecord` values (%d)", in)
	//_, err2 := session.Exec(sql1)
	//if err2 != nil {
	//	fmt.Println(err2)
	//}
	//fmt.Println(in)
	var m Money
	ss := session.Alias("r").Where("r.id = ?", 1)

	_, e := ss.Get(&m)

	if e != nil {
		fmt.Println(e)
	}
	fmt.Println("b", m.RMB, in)
	//if  m.RMB-in< 0 {
	//	return nil
	//	//var hh map[string]interface{}
	//} //

	//sql := fmt.Sprintf("update  `money` set rmb = %d", m.RMB-in)
	sql := fmt.Sprintf("insert into `payrecord` (rmb) values (%d)", in)
	_, err1 := session.Exec(sql)
	if err1 != nil {
		fmt.Println(err1)
	}
	//手动调用rollback后事物之间的锁就失效了
	//session.Rollback()
	sql1 := fmt.Sprintf("update  `money` set rmb = %d", m.RMB-in)
	//fmt.Println(in)
	_, err1 = session.Exec(sql1)
	if err1 != nil {
		fmt.Println(err1)
	}

	if m.RMB-in < 0 {
		//fmt.Fprintf(os.Stdout, "消费%d元不成功,余额不足.余额:%d元\n", in, m.RMB)
		fmt.Fprintf(os.Stderr, "消费%d元不成功,余额不足.余额:%d元\n", in, m.RMB)
		session.Rollback()
		fmt.Println("%%%%%%%%%%", time.Now().UnixNano())
		return nil
		//var hh map[string]interface{}
	} //

	//time.Sleep(2000 * time.Millisecond)
	err := session.Commit()
	if err != nil {
		fmt.Println(err)
		session.Rollback()
	}
	//va
	fmt.Println("a", m.RMB, in)
	fmt.Println("%%%%%%%%%%", time.Now().UnixNano())
	return nil

}
