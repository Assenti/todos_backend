package routes

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

const mysqlDbURI = "PgQXfyC4AD:CV3B9cSf2k@tcp(remotemysql.com:3306)/PgQXfyC4AD?parseTime=true"
const port = "3001"

// User model
type User struct {
	ID        uint      `gorm:"primary_key"`
	CreatedAt time.Time `gorm:"default: CURRENT_TIMESTAMP"`
	Firstname string    `gorm:"type:varchar(100)"`
	Lastname  string    `gorm:"type:varchar(100)"`
	Email     string    `gorm:"unique_index"`
	Password  string    `gorm:"type:varchar(100)"`
}

// Todo Model
type Todo struct {
	ID        uint      `gorm:"primary_key"`
	CreatedAt time.Time `gorm:"default: CURRENT_TIMESTAMP"`
	UpdatedAt time.Time
	Value     string
	Important int8 `gorm:"default: 0"`
	Completed int8 `gorm:"default: 0"`
	UserID    uint
}

// InitMigration function
func InitMigration() {
	db, err = gorm.Open("mysql", mysqlDbURI)
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to database")
	}
	defer db.Close()

	db.AutoMigrate(&Todo{}, &User{})
}
