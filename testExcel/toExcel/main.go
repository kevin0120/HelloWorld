package main

import (
	"fmt"
	"github.com/tealeg/xlsx"
	"os"
)

func main() {
	dir, _ := os.Getwd()
	//D:\Code\lianxi\HelloWorld
	fmt.Println("当前路径：", dir)
	excelFileName := "./testExcel/toExcel/1.xlsx"
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		panic(err)
	}
	for s, sheet := range xlFile.Sheets {
		if s == 1 {
			break
		}
		for _, row := range sheet.Rows {
			for j, cell := range row.Cells {
				if j == 0 {
					fmt.Printf("\t\n")
				}
				switch cell.Type() {
				case xlsx.CellTypeString:
					fmt.Printf("str %d %s\t", j, cell.String())
					break
				case xlsx.CellTypeStringFormula:
					fmt.Printf("formula %d %s\t", j, cell.Formula())
					break
				case xlsx.CellTypeNumeric:
					x, _ := cell.Int64()
					fmt.Printf("int %d %d\t", j, x)
					break
				case xlsx.CellTypeBool:
					fmt.Printf("bool %d %v\t", j, cell.Bool())
					break
				case xlsx.CellTypeDate:
					t, _ := cell.GetTime(false)
					fmt.Printf("date %d %v\t", j, t)
					break
				}
			}
			row.AddCell().Value = "hello world!"
		}
	}
	err = xlFile.Save("./testExcel/toExcel/1.xlsx")
	if err != nil {
		fmt.Printf(err.Error())
	}
}
