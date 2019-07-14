package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB
var err error

func InitialMigration() {
	db, err = gorm.Open("mysql", "root:password@/tunaiku_test?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to database")
	}
	defer db.Close()

	// Migrate table from model
	db.AutoMigrate(&Loan{})
}

func main() {
	fmt.Println("Go API Tunaiku Test")

	InitialMigration()

	handleRequest()
}
