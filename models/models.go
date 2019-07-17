package models

import (
	"time"

	"github.com/Assenti/restapi/db"
)

// User model
type User struct {
	ID           uint64    `gorm:"primary_key" json:"id"`
	CreatedAt    time.Time `gorm:"default: CURRENT_TIMESTAMP" json:"createdAt"`
	Firstname    string    `gorm:"type:varchar(100)" json:"firstname"`
	Lastname     string    `gorm:"type:varchar(100)" json:"lastname"`
	Email        string    `gorm:"unique_index" json:"email"`
	Password     string    `gorm:"type:varchar(100)" json:"password"`
	LastLoggedOn time.Time `gorm:"type: TIMESTAMP" json:"lastLoggedOn"`
}

// Group model
type Group struct {
	ID        uint64    `gorm:"primary_key" json:"id"`
	Name      string    `gorm:"type:varchar(100)" json:"name"`
	CreatedAt time.Time `gorm:"default: CURRENT_TIMESTAMP" json:"createdAt"`
	UserID    uint64    `gorm:"unique_index" json:"userId"`
}

// GroupParticipants model
type GroupParticipants struct {
	ID      uint64 `gorm:"primary_key" json:"id"`
	UserID  uint64 `json:"userId"`
	GroupID uint64 `json:"groupId"`
}

// JoinedGroupParticipants model
type JoinedGroupParticipants struct {
	ID        uint64 `gorm:"primary_key" json:"id"`
	UserID    uint64 `json:"userId"`
	GroupID   uint64 `json:"groupId"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

// UserSubmit model
type UserSubmit struct {
	Email    string `gorm:"unique_index" json:"email"`
	Password string `gorm:"type:varchar(100)" json:"password"`
	Remember bool   `json:"remember"`
}

// UserInfo model to pass to frontend
type UserInfo struct {
	ID           uint64    `gorm:"primary_key" json:"id"`
	CreatedAt    time.Time `gorm:"default: CURRENT_TIMESTAMP" json:"createdAt"`
	Firstname    string    `gorm:"type:varchar(100)" json:"firstname"`
	Lastname     string    `gorm:"type:varchar(100)" json:"lastname"`
	Email        string    `gorm:"unique_index" json:"email"`
	Token        string    `json:"token"`
	LastLoggedOn time.Time `gorm:"type: TIMESTAMP" json:"lastLoggedOn"`
}

// Todo model
type Todo struct {
	ID           uint64    `gorm:"primary_key" json:"id"`
	CreatedAt    time.Time `gorm:"default: CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	CompleteDate string    `gorm:"nullable" json:"completeDate"`
	Status       string    `json:"status"`
	Value        string    `json:"value"`
	Important    int8      `gorm:"default: 0" json:"important"`
	Completed    int8      `gorm:"default: 0" json:"completed"`
	UserID       uint      `json:"userId"`
	Performer    uint64    `json:"performer"`
	GroupID      uint      `json:"groupId"`
}

// JoinedTodo model
type JoinedTodo struct {
	ID           uint64    `gorm:"primary_key" json:"id"`
	CreatedAt    time.Time `gorm:"default: CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	CompleteDate string    `gorm:"nullable" json:"completeDate"`
	Status       string    `json:"status"`
	Value        string    `json:"value"`
	Important    int8      `gorm:"default: 0" json:"important"`
	Completed    int8      `gorm:"default: 0" json:"completed"`
	UserID       uint      `json:"userId"`
	Performer    uint64    `json:"performer"`
	GroupID      uint      `json:"groupId"`
	Firstname    string    `json:"firstname"`
	Lastname     string    `json:"lastname"`
}

// TodoDetails model
type TodoDetails struct {
	ID      uint64 `gorm:"primary_key" json:"id"`
	TodoID  uint64 `json:"todoId"`
	Content string `gorm:"type:text" json:"content"`
}

// Performer model
type Performer struct {
	ID        uint64 `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
}

// InitDb function
func InitDb() {
	db := db.Connect()

	db.AutoMigrate(
		&Todo{},
		&TodoDetails{},
		&User{},
		&Group{},
		&GroupParticipants{},
	)
}
