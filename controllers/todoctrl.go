package controllers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Assenti/restapi/db"
	"github.com/Assenti/restapi/models"
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris"
	"gopkg.in/gomail.v2"
)

// GetTodos method
func GetTodos(ctx iris.Context) {
	var todos []models.Todo

	db := db.Connect()
	defer db.Close()

	db.Find(&todos)
	ctx.JSON(iris.Map{"todos": todos})
}

// GetUserTodos method
func GetUserTodos(ctx iris.Context) {
	var todos []models.JoinedTodo

	id := ctx.URLParam("userid")

	db := db.Connect()
	defer db.Close()

	db.Table("users").Raw(`SELECT 
				todos.id, 
				todos.value, 
				todos.created_at, 
				todos.updated_at, 
				todos.important, 
				todos.completed, 
				todos.user_id, 
				todos.complete_date, 
				todos.status, 
				todos.group_id, 
				todos.performer, 
				users.firstname, 
				users.lastname
				FROM todos
				LEFT JOIN users ON todos.performer = users.id
				where user_id = ?`, id).Scan(&todos)

	ctx.JSON(iris.Map{"todos": todos})
}

// GetSingleTodo method
func GetSingleTodo(ctx iris.Context) {
	id := ctx.URLParam("id")
	var result models.Todo

	db := db.Connect()
	defer db.Close()

	db.Where("id = ?", id).Last(&result)
	ctx.JSON(iris.Map{"todo": result})
}

// CreateTodo method
func CreateTodo(ctx iris.Context) {
	var todo models.Todo

	err := ctx.ReadJSON(&todo)

	if err != nil || todo.Value == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"msg": "Value must be provided."})
		return
	}

	db := db.Connect()
	defer db.Close()

	db.Create(&models.Todo{Value: todo.Value, UserID: todo.UserID})

	var newTodo models.Todo
	db.Where("user_id = ?", todo.UserID).Last(&newTodo)
	ctx.JSON(iris.Map{"todo": newTodo})
}

// UpdateTodo method
func UpdateTodo(ctx iris.Context) {
	var todo models.Todo

	err := ctx.ReadJSON(&todo)

	if err != nil || todo.Value == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"msg": "Todo must be provided."})
		return
	}

	db := db.Connect()
	defer db.Close()

	db.Model(&todo).Where("id = ?", todo.ID).Update("value", todo.Value)

	var updatedTodo models.Todo
	db.Where("id = ?", todo.ID).Last(&updatedTodo)
	ctx.JSON(iris.Map{"todo": updatedTodo})
}

// DeleteTodo method
func DeleteTodo(ctx iris.Context) {
	id := ctx.URLParam("id")

	db := db.Connect()
	defer db.Close()

	db.Where("id = ?", id).Delete(&models.Todo{})
	ctx.JSON(iris.Map{"msg": "Todo successfully deleted."})
}

// ToggleTodoCompletion method
func ToggleTodoCompletion(ctx iris.Context) {
	id := ctx.URLParam("id")
	db := db.Connect()
	defer db.Close()

	var todo models.Todo
	db.Where("id = ?", id).Last(&todo)

	if todo.Completed == 1 {
		db.Model(&todo).Where("id = ?", id).Updates(map[string]interface{}{"completed": 0, "complete_date": gorm.Expr("NULL")})
	} else {
		temp := time.Now()
		current := temp.Format("2006-01-02T15:04:05")
		db.Model(&todo).Where("id = ?", id).Updates(map[string]interface{}{"completed": 1, "complete_date": current})
	}

	var updatedTodo models.Todo
	db.Where("id = ?", id).Last(&updatedTodo)

	ctx.JSON(iris.Map{"todo": &updatedTodo})
}

// ToggleTodoImportance method
func ToggleTodoImportance(ctx iris.Context) {
	id := ctx.URLParam("id")
	db := db.Connect()
	defer db.Close()

	var todo models.Todo
	db.Where("id = ?", id).Last(&todo)

	if todo.Important == 1 {
		db.Model(&todo).Where("id = ?", id).Update("important", 0)
	} else {
		db.Model(&todo).Where("id = ?", id).Update("important", 1)
	}

	var updatedTodo models.Todo
	db.Where("id = ?", id).Last(&updatedTodo)
	ctx.JSON(iris.Map{"todo": &updatedTodo})
}

// SendTodosListViaEmail method
func SendTodosListViaEmail(ctx iris.Context) {
	var todos []models.Todo
	email := ctx.URLParam("email")

	err := ctx.ReadJSON(&todos)

	if err != nil || len(todos) == 0 {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"msg": "Todos list must be provided."})
		return
	}

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

// SetTodoPerformer method
func SetTodoPerformer(ctx iris.Context) {
	todoID := ctx.URLParam("todoId")
	userID := ctx.URLParam("userId")

	db := db.Connect()
	defer db.Close()

	var todo models.Todo
	db.Where("id = ?", todoID).Last(&todo)

	intUserID, _ := strconv.ParseUint(userID, 10, 64)

	db.Model(&todo).Where("id = ?", todoID).Update("performer", intUserID)

	var updatedTodo models.Todo
	db.Where("id = ?", todoID).Last(&updatedTodo)
	ctx.JSON(iris.Map{"todo": &updatedTodo})
}
