package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"time"
)

//參考文檔
//https://github.com/go-xorm/xorm/blob/master/README_CN.md
//http://gobook.io/read/github.com/go-xorm/manual-zh-CN/chapter-01/
type Rush struct {
	eng *xorm.Engine
}

type Workorder struct {
	Id           int64  `xorm:"pk autoincr notnull 'id'" json:"-"`
	Workorder    string `xorm:"text 'workorder'" json:"-"`
	Code         string `xorm:"varchar(128) 'code'"  json:"code"`
	Track_code   string `xorm:"varchar(128) 'track_code'" json:"track_code"`
	Product_code string `xorm:"varchar(128) 'product_code'" json:"product_code"`

	Created time.Time `xorm:"created" json:"-"`
	Updated time.Time `xorm:"updated" json:"-"`
}

type Step struct {
	Id          int64  `xorm:"pk autoincr notnull 'id'" json:"-"`
	WorkorderID int64  `xorm:"bigint 'x_workorder_id'" json:"-"`
	Step        string `xorm:"text 'step'" json:"-"`

	Test_type string `xorm:"varchar(128) 'test_type'" json:"test_type"`

	Code string `xorm:"varchar(128) 'code'" json:"code"`

	Status1 string    `xorm:"varchar(128) 'status1'" json:"status1"`
	Status2 string    `xorm:"varchar(128) 'status2'" json:"status2"`
	Image   string    `xorm:"varchar(128) 'image'" json:"image"`
	Created time.Time `xorm:"created" json:"-"`
	Updated time.Time `xorm:"updated" json:"-"`
}

//`xorm:"varchar(25) notnull unique 'usr_name'"`
func main() {

	file, _ := os.Open("/home/kevin/Downloads/gopath/src/HelloWorld/SQL_test/xorm_test/rush/workorder.json")

	var wor []byte
	wor, _ = ioutil.ReadAll(file)

	fmt.Println(string(wor))

	r := NewRush()

	r.WorkorderIn(wor)
	fmt.Println("###############################################")
	//根据工单号,查询工单
	_, rr := r.WorkorderOut("MO0001")
	fmt.Println(string(rr))

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
	exist, err := engine.IsTableExist("Workorder")
	if err == nil {
		if !exist {
			if err := engine.Sync2(new(Workorder)); err != nil {
				errors.Wrapf(err, "Create Table Workorder fail")
			}
		}
	}

	exist, err = engine.IsTableExist("Step")
	if err == nil {
		if !exist {
			if err := engine.Sync2(new(Step)); err != nil {
				errors.Wrapf(err, "Create Table Step fail")
			}
		}
	}

	r := &Rush{
		eng: engine,
	}

	return r
}

func (s *Rush) WorkorderIn(in []byte) error {

	session := s.eng.NewSession()
	defer session.Close()

	var work Workorder
	err := json.Unmarshal(in, &work)

	workorder1 := Workorder{
		Workorder:    string(in),
		Code:         work.Code,
		Track_code:   work.Track_code,
		Product_code: work.Product_code,
	}

	//var hh map[string]interface{}
	//
	//err = json.Unmarshal(wor, &hh)
	_, err = s.eng.Insert(&workorder1)
	// INSERT INTO struct () values ()
	//engine.Ping()
	//有的数据库超时断开ping可以重连。可以通过起一个定期Ping的Go程来保持连接鲜活。
	if err != nil {
		session.Rollback()
		return errors.Wrapf(err, "store data fail")
	}

	var hh map[string]interface{}
	var step []map[string]interface{}

	err = json.Unmarshal(in, &hh)

	cc, _ := json.Marshal(hh["steps"])

	err = json.Unmarshal(cc, &step)

	for i := 0; i < len(step); i++ {
		a, _ := json.Marshal(step[i])
		var msg Step
		err = json.Unmarshal(a, &msg)

		step := Step{
			WorkorderID: workorder1.Id,
			Step:        string(a),
			Image:       msg.Image,

			Test_type: msg.Test_type,

			Code: msg.Code,
		}

		_, err := s.eng.Insert(&step)
		// INSERT INTO struct () values ()
		if err != nil {
			session.Rollback()
			return errors.Wrapf(err, "store data fail")
		}

	}

	err = session.Commit()
	if err != nil {
		return errors.Wrapf(err, "commit fail")
	}

	return nil

}

func (s *Rush) WorkorderOut(order string) (error, []byte) {

	var workorder Workorder
	ss := s.eng.Alias("r").Where("r.code = ?", order)
	_, e := ss.Get(&workorder)
	if e != nil {
		return e, nil
	}

	var step []Step
	ss = s.eng.Alias("r").Where("r.x_workorder_id = ?", workorder.Id)
	e = ss.Find(&step)
	if e != nil {
		return e, nil
	}

	var steps []map[string]interface{}
	for i := 0; i < len(step); i++ {
		hh := stringtomap(step[i].Step)
		jj, _ := json.Marshal(hh["test_type"])

		if !(string(jj) == `"tightening"`) {
			steps = append(steps, hh)
			continue
		}

		mm := structomap(step[i])
		for k, v := range mm {
			hh[k] = v
		}
		findpicture()
		hh["image"] = "http://127.0.0.1:8080/picture"
		steps = append(steps, hh)
	}

	ww := stringtomap(workorder.Workorder)
	ww["steps"] = steps

	rr, _ := json.Marshal(ww)

	return nil, rr
}

func structomap(in interface{}) (m map[string]interface{}) {
	j, _ := json.Marshal(in)
	json.Unmarshal(j, &m)
	return
}

func stringtomap(in string) (m map[string]interface{}) {
	json.Unmarshal([]byte(in), &m)
	return
}

func findpicture() {
	fmt.Println("查找图片-----")

}
