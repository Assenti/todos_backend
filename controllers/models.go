package controllers

import "time"

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
