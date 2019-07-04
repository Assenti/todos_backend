package main

import (
	"os"

	"github.com/Assenti/restapi/controllers"
	_ "github.com/go-sql-driver/mysql"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
)

func main() {
	var port string
	var envPort string
	envPort = os.Getenv("PORT")
	if envPort != "" {
		port = envPort
	} else {
		port = "3000"
	}

	app := iris.Default()
	app.Logger().SetLevel("debug")
	app.Use(recover.New())
	app.Use(logger.New())

	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080", "https://planner-2.herokuapp.com"},
		AllowedMethods:   []string{"HEAD", "GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	// Rendering static files
	// app.StaticWeb("/", "./dist")

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
		api.Post("/sendViaEmail", controllers.SendTodosListViaEmail)

		// User API
		api.Post("/users", controllers.CreateUser)
		api.Put("/users", controllers.UpdateUser)
		api.Get("/users", controllers.GetUsersList)
		api.Post("/login", controllers.Login)
		api.Get("/restorepassword", controllers.RestorePassword)
		api.Post("/changepassword", controllers.ChangePassword)
		api.Post("/checkpassword", controllers.CheckPassword)
	}

	app.Run(iris.Addr(":" + port))
}
