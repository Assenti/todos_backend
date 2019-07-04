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

// UserSubmit model
type UserSubmit struct {
	Email    string `gorm:"unique_index"`
	Password string `gorm:"type:varchar(100)"`
	Remember bool
}

// UserInfo model to pass to frontend
type UserInfo struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `gorm:"default: CURRENT_TIMESTAMP" json:"createdAt"`
	Firstname string    `gorm:"type:varchar(100)" json:"firstname"`
	Lastname  string    `gorm:"type:varchar(100)" json:"lastname"`
	Email     string    `gorm:"unique_index" json:"email"`
	Token     string    `json:"token"`
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
