package controllers

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/kataras/iris"
	"gopkg.in/gomail.v2"
)

const mysqlDbURI = "PgQXfyC4AD:CV3B9cSf2k@tcp(remotemysql.com:3306)/PgQXfyC4AD?parseTime=true"

var db *gorm.DB
var err error

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

	db.Where("user_id = ?", id).Order("created_at desc").Find(&todos)

	// const shortDate = "2019-01-01"
	// start := todos[len(todos)-1].CreatedAt
	// end := todos[0].CreatedAt
	// compareStart := todos[len(todos)-1].CreatedAt
	// var dates []string

	// fmt.Println(start)
	// fmt.Println(end)

	// for i := 0; i < 1000000; i++ {
	// 	fmt.Println(i)
	// 	if IsDateInPeriod(compareStart, end, start) {
	// 		dates = append(dates, start.String()[0:10])
	// 		start = start.AddDate(0, 0, 1)
	// 		i++
	// 	} else {
	// 		break
	// 	}
	// }

	// dates = Unique(dates)

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
		return
	}

	db, dbErr := gorm.Open("mysql", mysqlDbURI)

	if dbErr != nil {
		fmt.Println(dbErr.Error())
		panic("Failed to connect to database")
	}
	defer db.Close()

	db.Create(&Todo{Value: todo.Value, UserID: todo.UserID})

	var newTodo Todo
	db.Where("user_id = ?", todo.UserID).Last(&newTodo)
	ctx.JSON(iris.Map{"todo": newTodo})
}

// UpdateTodo method
func UpdateTodo(ctx iris.Context) {
	var todo Todo

	err := ctx.ReadJSON(&todo)

	if err != nil || todo.Value == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"msg": "Todo must be provided."})
		return
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

// DeleteTodo method
func DeleteTodo(ctx iris.Context) {
	id := ctx.URLParam("id")

	db, dbErr := gorm.Open("mysql", mysqlDbURI)

	if dbErr != nil {
		fmt.Println(dbErr.Error())
		panic("Failed to connect to database")
	}
	defer db.Close()

	db.Where("id = ?", id).Delete(&Todo{})
	ctx.JSON(iris.Map{"msg": "Todo successfully deleted."})
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

// SendTodosListViaEmail method
func SendTodosListViaEmail(ctx iris.Context) {
	var todos []Todo
	email := ctx.URLParam("email")

	err := ctx.ReadJSON(&todos)

	if err != nil || len(todos) == 0 {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"msg": "Todos list must be provided."})
		return
	}

	fmt.Println(len(todos))

	var todosHTML string
	htmlBody := "<h4>Your todos list:</h4>"
	footer := "<hr><p style='font-size: 11px'>Message was generated automatically. Please don't reply.</p>"

	for i, todo := range todos {
		index := i + 1
		todosHTML += fmt.Sprintf("<div>%d) %s;</div>", index, todo.Value)
	}

	fmt.Println(todosHTML)

	htmlBody = htmlBody + todosHTML + footer

	m := gomail.NewMessage()
	m.SetHeader("From", mailUser)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Todos List")
	m.SetBody("text/html", htmlBody)

	d := gomail.NewDialer(mailHost, mailPort, mailUser, mailPassword)

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}

	ctx.StatusCode(iris.StatusOK)
}
