package main

import (
	"database/sql"
	"fmt"
	"github.com/weigj/go-odbc"
	_ "github.com/weigj/go-odbc/driver"
)

func main() {
	fmt.Printf("%s\n", "创建数据库链接")
	//conn, _ := odbc.Connect("DSN=hvb")
	conn, _ := odbc.Connect("driver={Microsoft Access Driver (*.mdb,*.accdb)};dbq=./data1.mdb")
	stmt, _ := conn.Prepare("select top 10 * from data")
	_ = stmt.Execute()
	rows, err := stmt.FetchAll()
	if err != nil {
		fmt.Println(err)
		return
	}
	for i, row := range rows {
		println(i, row)
	}
	stmt.Close()
	conn.Close()
	return

	//access()
}

func access() {

	conn, err := sql.Open("odbc", "driver={Microsoft Access Driver (*.mdb,*.accdb)};dbq=./data1.mdb")
	if err != nil {
		fmt.Println("Connecting Error")
		return
	}
	fmt.Println(conn)
	defer conn.Close()
	a, b := conn.Exec("select * from data")

	println(a,b)
	stmt, err := conn.Prepare("select * from data") //ALTER TABLE tb ALTER COLUMN aa Long
	if err != nil {
		fmt.Println("Query Error1")
		return
	}
	defer stmt.Close()
	row, err := stmt.Query()
	if err != nil {
		fmt.Print(err)
		fmt.Println("Query Error")
		return
	}
	defer row.Close()
	for row.Next() {
		var ID string
		var SequenceNumber int
		var ValueCode string
		if err := row.Scan(&ID, &SequenceNumber, &ValueCode); err == nil {
			fmt.Println(ID, SequenceNumber, ValueCode)
		}
	}
	fmt.Printf("%s\n", "finish")
	return

}
