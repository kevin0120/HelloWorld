package main

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	"reflect"
)

type Document struct {
	Metadata postgres.Jsonb
	//Secrets  postgres.Hstore
	Body string
	ID   int
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
	db.AutoMigrate(&Document{})

	//password := "0654857340"
	metadata := json.RawMessage(`{"is_archived": 1}`)
	sampleDoc := Document{
		Body:     "This is a test document",
		Metadata: postgres.Jsonb{metadata},
		//Secrets:  postgres.Hstore{"password": &password},
	}

	//insert sampleDoc into the database
	db.Create(&sampleDoc)

	//retrieve the fields again to confirm if they were inserted correctly
	resultDoc := Document{}
	db.Where("id = ?", sampleDoc.ID).First(&resultDoc)

	metadataIsEqual := reflect.DeepEqual(resultDoc.Metadata, sampleDoc.Metadata)
	//secretsIsEqual := reflect.DeepEqual(resultDoc.Secrets, sampleDoc.Secrets)

	// this should print "true"
	//fmt.Println("Inserted fields are as expected:", metadataIsEqual && secretsIsEqual)
	fmt.Println("Inserted fields are as expected:", metadataIsEqual)
}
