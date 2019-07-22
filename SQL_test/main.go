package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func main() {
	db, _ := sql.Open("postgres", "user=odoo password=odoo dbname=test sslmode=disable")
	//fmt.Println(err1)
	rew, _ := db.Query("SELECT code FROM  mrp_bom")
	//fmt.Println(err2)
	var a string
	rew.Next()
	rew.Scan(&a)

	fmt.Println(a)
}
