package db

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

var db *gorm.DB
var err error

// Connect to db
func Connect() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		panic(err.Error())
	}

	host := os.Getenv("MS_HOST")
	port := os.Getenv("MS_PORT")
	user := os.Getenv("MS_USER")
	dbName := os.Getenv("MS_DB")
	pass := os.Getenv("MS_PASSWORD")

	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		user, pass, host, port, dbName)

	// configs := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
	// 	host, port, user, dbName, pass)

	db, err := gorm.Open("mysql", dbURI)

	if err != nil {
		panic(err.Error())
	}

	return db
}
