package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func main() {
	i := 1573530963287321846 - 1573530951241648837
	s := 1573531251396910951 - 1573531237316635691

	i0 := 1573532014619150152.0 - 1573532014604790592.0
	s0 := 1573531440705412959.0 - 1573531440690693847.0

	i1 := 1573531693944508704.0 - 1573531693928653814.0
	s1 := 1573531794316803309.0 - 1573531794300326803.0

	i2 := 1573545913768123155.0 - 1573545913752886288.0
	fmt.Println(i, s)

	fmt.Println((-i0 + s0) / s0)
	fmt.Println((-i1 + s1) / s1)
	fmt.Println((-i2 + s1) / s1)
	db, _ := sql.Open("postgres", "user=odoo password=odoo dbname=saturm sslmode=disable")
	//fmt.Println(err1)
	rew, _ := db.Query("select (jsons::json#>>'{product,url}')::text from testjson where (jsons::json#>>'{code}')::text = 'MO0002'")
	//fmt.Println(err2)
	//rew, _ := db.Query("SELECT json FROM  testjson ")
	var a string
	rew.Next()
	rew.Scan(&a)

	fmt.Println(a)
}
