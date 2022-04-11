package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"time"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func main() {
	sConn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable connect_timeout=10",
		"localhost",
		"5432",
		"odoo",
		"saturm",
		"odoo")

	db, err := gorm.Open("postgres", sConn)

	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&Product{})

	// Create
	db.Create(&Product{Code: "L1212", Price: 1000})

	// Read
	for {
		fmt.Println(db.DB().Ping())
		product := Product{}
		db.First(&product, 1) // find product with id 1
		db.First(&product, "code = ?", "L1212")
		fmt.Println(&product, product)
		time.Sleep(5 * time.Second)
		// find product with code l1212
	}
	//// Update - update product's price to 2000
	//db.Model(&product).Update("Price", 2000)
	//
	//// Delete - delete product
	//db.Delete(&product)
}
