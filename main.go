package main

import (
	"github.com/Assenti/restapi/controllers"
	_ "github.com/go-sql-driver/mysql"
	"github.com/iris-contrib/middleware/cors"
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

	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"HEAD", "GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	api := app.Party("/api", crs).AllowMethods(iris.MethodOptions)
	{
		// Todos API
		api.Get("/todo", controllers.GetSingleTodo)
		api.Get("/usertodos", controllers.GetUserTodos)
		api.Get("/todocompletion", controllers.ToggleTodoCompletion)
		api.Get("/todoimportance", controllers.ToggleTodoImportance)
		api.Get("/todos", controllers.GetTodos)
		api.Post("/todos", controllers.CreateTodo)
		api.Put("/todos", controllers.UpdateTodo)
		api.Delete("/todos", controllers.DeleteTodo)

		// User API
		api.Post("/users", controllers.CreateUser)
		api.Put("/users", controllers.UpdateUser)
		api.Get("/users", controllers.GetUsersList)
		api.Post("/login", controllers.Login)
	}

	app.Run(iris.Addr(":3000"))
}
