package controllers

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/kataras/iris"
)

// User model
type User struct {
	ID        uint      `gorm:"primary_key"`
	CreatedAt time.Time `gorm:"default: CURRENT_TIMESTAMP"`
	Firstname string    `gorm:"type:varchar(100)"`
	Lastname  string    `gorm:"type:varchar(100)"`
	Email     string    `gorm:"unique_index"`
	Password  string    `gorm:"type:varchar(100)"`
}

// CreateUser method
func CreateUser(ctx iris.Context) {
	var user User

	err := ctx.ReadJSON(&user)
	if err != nil || (user.Firstname == "" && user.Lastname == "" && user.Email == "" && user.Password == "") {
		ctx.StatusCode(400)
		ctx.JSON(iris.Map{"msg": "Firstname, Lastname, Email and Password must be provided."})
		return
	}

	db, dbErr := gorm.Open("mysql", mysqlDbURI)

	if dbErr != nil {
		fmt.Println(dbErr.Error())
		panic("Failed to connect to database")
	}
	defer db.Close()

	db.Create(&User{Firstname: user.Firstname, Lastname: user.Lastname, Email: user.Email, Password: user.Password})
	ctx.JSON(iris.Map{"msg": "New user successfully created"})
}

// UpdateUser method
func UpdateUser(ctx iris.Context) {
	var user User

	err := ctx.ReadJSON(&user)
	if err != nil || (user.Firstname == "" && user.Lastname == "") {
		ctx.StatusCode(400)
		ctx.JSON(iris.Map{"msg": "Firstname, Lastname must be provided."})
		return
	}

	db, dbErr := gorm.Open("mysql", mysqlDbURI)

	if dbErr != nil {
		fmt.Println(dbErr.Error())
		panic("Failed to connect to database")
	}
	defer db.Close()

	db.Model(&user).Where("id = ?", user.ID).Updates(User{Firstname: user.Firstname, Lastname: user.Lastname})

	var updatedUser User
	db.Where("id = ?", user.ID).Last(&updatedUser)
	ctx.JSON(iris.Map{"user": updatedUser})
}

// GetUsersList method
func GetUsersList(ctx iris.Context) {
	var users []User
	db, dbErr := gorm.Open("mysql", mysqlDbURI)

	if dbErr != nil {
		fmt.Println(dbErr.Error())
		panic("Failed to connect to database")
	}
	defer db.Close()

	db.Find(&users)
	ctx.JSON(iris.Map{"users": users})
}

// Login method
func Login(ctx iris.Context) {
	var user User

	err := ctx.ReadJSON(&user)
	if err != nil || (user.Email == "" && user.Password == "") {
		ctx.StatusCode(400)
		ctx.JSON(iris.Map{"message": "Email and Password must be provided."})
		return
	}

	db, dbErr := gorm.Open("mysql", mysqlDbURI)

	if dbErr != nil {
		fmt.Println(dbErr.Error())
		panic("Failed to connect to database")
	}
	defer db.Close()

	var foundedUser User
	db.Where(&User{Email: user.Email, Password: user.Password}).First(&foundedUser)

	if foundedUser.Email == user.Email && foundedUser.Password == user.Password {
		ctx.JSON(iris.Map{"user": &foundedUser})
	} else {
		ctx.StatusCode(400)
		ctx.JSON(iris.Map{"msg": "Invalid Email or Password."})
	}
}
