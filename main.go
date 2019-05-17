package main

import (
	"github.com/Assenti/restapi/controllers"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
)

func main() {
	app := iris.Default()
	app.Logger().SetLevel("debug")
	app.Use(recover.New())
	app.Use(logger.New())

	app.Get("/", func(ctx iris.Context) {
		ctx.HTML("Server successfully started")
	})

	// Todos API
	app.Get("/api/todo", controllers.GetSingleTodo)
	app.Get("/api/usertodos", controllers.GetUserTodos)
	app.Get("/api/todocompletion", controllers.ToggleTodoCompletion)
	app.Get("/api/todoimportance", controllers.ToggleTodoImportance)
	app.Get("/api/todos", controllers.GetTodos)
	app.Post("/api/todos", controllers.CreateTodo)
	app.Put("/api/todos", controllers.UpdateTodo)

	// User API
	app.Post("/api/users", controllers.CreateUser)
	app.Put("/api/users", controllers.UpdateUser)
	app.Get("/api/users", controllers.GetUsersList)
	app.Post("/api/login", controllers.Login)

	app.Run(iris.Addr(":3000"))
}
