package main

import (
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"

	"github.com/xuri/excelize/v2"
)

func NewSheet1(f *excelize.File) {

	//更改sheet名
	f.SetSheetName("Sheet1", "Sheet1")
	categories := map[string]string{
		"A2": "Small", "A3": "Normal", "A4": "Large",
		"B1": "Apple", "C1": "Orange", "D1": "Pear"}
	values := map[string]int{
		"B2": 2, "C2": 3, "D2": 3, "B3": 5, "C3": 2, "D3": 4, "B4": 6, "C4": 7, "D4": 8}
	for k, v := range categories {
		f.SetCellValue("Sheet1", k, v)
	}
	for k, v := range values {
		f.SetCellValue("Sheet1", k, v)
	}
	if err := f.AddChart("Sheet1", "E1", `{
        "type": "col3DClustered",
        "series": [
        {
            "name": "Sheet1!$A$2",
            "categories": "Sheet1!$B$1:$D$1",
            "values": "Sheet1!$B$2:$D$2"
        },
        {
            "name": "Sheet1!$A$3",
            "categories": "Sheet1!$B$1:$D$1",
            "values": "Sheet1!$B$3:$D$3"
        },
        {
            "name": "Sheet1!$A$4",
            "categories": "Sheet1!$B$1:$D$1",
            "values": "Sheet1!$B$4:$D$4"
        }],
        "title":
        {
            "name": "Fruit 3D Clustered Column Chart"
        }
    }`); err != nil {
		fmt.Println(err)
		return
	}
}

func NewSheet2(f *excelize.File) {

	// Create a new sheet
	index := f.NewSheet("Sheet2")
	// Set value of a cell.
	f.SetCellValue("Sheet2", "A3", "Hello world.")
	// Set active sheet of the workbook.
	f.SetActiveSheet(index)
}

func NewSheet3(f *excelize.File) {
	// Create a new sheet
	index := f.NewSheet("Sheet3")
	f.SetActiveSheet(index)
	// Insert a picture.
	if err := f.AddPicture("Sheet3", "A2", "./testExcel/fromandtoexcel/write-read/images/excel.png", ""); err != nil {
		fmt.Println(err)
	}
	// Insert a picture to worksheet with scaling.
	if err := f.AddPicture("Sheet3", "D2", "./testExcel/fromandtoexcel/write-read/images/excel.jpg",
		`{"x_scale": 0.5, "y_scale": 0.5}`); err != nil {
		fmt.Println(err)
	}
	// Insert a picture offset in the cell with printing support.
	if err := f.AddPicture("Sheet3", "H2", "./testExcel/fromandtoexcel/write-read/images/excel.gif", `{
        "x_offset": 15,
        "y_offset": 10,
        "print_obj": true,
        "lock_aspect_ratio": false,
        "locked": false
    }`); err != nil {
		fmt.Println(err)
	}
}




func ReadCellValue() {
	f, err := excelize.OpenFile("./testExcel/fromandtoexcel/write-read/Book1.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// Get value from cell by given worksheet name and axis.
	cell, err := f.GetCellValue("Sheet1", "B2")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(cell)
	// Get all the rows in the Sheet1.
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, row := range rows {
		for _, colCell := range row {
			fmt.Print(colCell, "\t")
		}
		fmt.Println()
	}
}



func ReadPicture() {
	f, err := excelize.OpenFile("./testExcel/fromandtoexcel/write-read/Book1.xlsx")
	   if err != nil {
	       fmt.Println(err)
	       return
	   }
	   defer func() {
	       if err := f.Close(); err != nil {
	           fmt.Println(err)
	       }
	   }()
	   file, raw, err := f.GetPicture("Sheet3", "D2")
	   if err != nil {
	       fmt.Println(err)
	       return
	   }
	   if err := ioutil.WriteFile("./testExcel/fromandtoexcel/write-read/"+file, raw, 0644); err != nil {
	       fmt.Println(err)
	   }
}


func main() {
	f := excelize.NewFile()

	//chart
	NewSheet1(f)

	NewSheet2(f)

	//picture
	NewSheet3(f)

	// Save spreadsheet by the given path.
	if err := f.SaveAs("./testExcel/fromandtoexcel/write-read/Book1.xlsx"); err != nil {
		fmt.Println(err)
	}

	ReadCellValue()
	ReadPicture()
}
