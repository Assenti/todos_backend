package controllers

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/kataras/iris"
)

const mysqlDbURI = "PgQXfyC4AD:CV3B9cSf2k@tcp(remotemysql.com:3306)/PgQXfyC4AD?parseTime=true"

var db *gorm.DB
var err error

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

// GetTodos method
func GetTodos(ctx iris.Context) {
	var todos []Todo

	db, err = gorm.Open("mysql", mysqlDbURI)
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to database")
	}
	defer db.Close()

	db.Find(&todos)
	ctx.JSON(iris.Map{"todos": todos})
}

// GetUserTodos method
func GetUserTodos(ctx iris.Context) {
	var todos []Todo

	id := ctx.URLParam("userid")

	db, err = gorm.Open("mysql", mysqlDbURI)
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to database")
	}
	defer db.Close()

	db.Where("user_id = ?", id).Find(&todos)
	ctx.JSON(iris.Map{"todos": todos})
}

// GetSingleTodo method
func GetSingleTodo(ctx iris.Context) {
	id := ctx.URLParam("id")
	var result Todo

	db, err = gorm.Open("mysql", mysqlDbURI)
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to database")
	}
	defer db.Close()

	db.Where("id = ?", id).Last(&result)
	ctx.JSON(iris.Map{"todo": result})
}

// CreateTodo method
func CreateTodo(ctx iris.Context) {
	var todo Todo

	err := ctx.ReadJSON(&todo)

	if err != nil || todo.Value == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"msg": "Value must be provided."})
	}

	db, dbErr := gorm.Open("mysql", mysqlDbURI)

	if dbErr != nil {
		fmt.Println(dbErr.Error())
		panic("Failed to connect to database")
	}
	defer db.Close()

	db.Create(&Todo{Value: todo.Value, UserID: todo.UserID})
	ctx.JSON(iris.Map{"msg": "Todo successfully created."})
}

// UpdateTodo method
func UpdateTodo(ctx iris.Context) {
	var todo Todo

	err := ctx.ReadJSON(&todo)

	if err != nil || todo.Value == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"msg": "Todo must be provided."})
	}

	db, dbErr := gorm.Open("mysql", mysqlDbURI)

	if dbErr != nil {
		fmt.Println(dbErr.Error())
		panic("Failed to connect to database")
	}
	defer db.Close()

	db.Model(&todo).Where("id = ?", todo.ID).Update("value", todo.Value)

	var updatedTodo Todo
	db.Where("id = ?", todo.ID).Last(&updatedTodo)
	ctx.JSON(iris.Map{"todo": updatedTodo})
}

// ToggleTodoCompletion method
func ToggleTodoCompletion(ctx iris.Context) {
	id := ctx.URLParam("id")
	db, dbErr := gorm.Open("mysql", mysqlDbURI)

	if dbErr != nil {
		fmt.Println(dbErr.Error())
		panic("Failed to connect to database")
	}
	defer db.Close()

	var todo Todo
	db.Where("id = ?", id).Last(&todo)

	if todo.Completed == 1 {
		db.Model(&todo).Where("id = ?", id).Update("completed", 0)
	} else {
		db.Model(&todo).Where("id = ?", id).Update("completed", 1)
	}

	var updatedTodo Todo
	db.Where("id = ?", id).Last(&updatedTodo)

	ctx.JSON(iris.Map{"todo": &updatedTodo})
}

// ToggleTodoImportance method
func ToggleTodoImportance(ctx iris.Context) {
	id := ctx.URLParam("id")
	db, dbErr := gorm.Open("mysql", mysqlDbURI)

	if dbErr != nil {
		fmt.Println(dbErr.Error())
		panic("Failed to connect to database")
	}
	defer db.Close()

	var todo Todo
	db.Where("id = ?", id).Last(&todo)

	if todo.Important == 1 {
		db.Model(&todo).Where("id = ?", id).Update("important", 0)
	} else {
		db.Model(&todo).Where("id = ?", id).Update("important", 1)
	}

	var updatedTodo Todo
	db.Where("id = ?", id).Last(&updatedTodo)
	ctx.JSON(iris.Map{"todo": &updatedTodo})
}
