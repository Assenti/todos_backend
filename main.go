package main

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
)

// Todo model
type Todo struct {
	ID        uint      `gorm:"primary_key"`
	CreatedAt time.Time `gorm:"default: CURRENT_TIMESTAMP"`
	UpdatedAt time.Time
	Value     string `json:"value"`
	Important int8   `gorm:"default: 0"`
	Completed int8   `gorm:"default: 0"`
	UserID    uint
}

const mysqlDbURI = "PgQXfyC4AD:CV3B9cSf2k@tcp(remotemysql.com:3306)/PgQXfyC4AD?parseTime=true"

var db *gorm.DB
var err error

func main() {
	app := iris.Default()
	app.Logger().SetLevel("debug")
	app.Use(recover.New())
	app.Use(logger.New())

	app.Get("/ping", func(ctx iris.Context) {
		ctx.JSON(iris.Map{
			"message": "pong",
		})
	})

	app.Get("/todo", func(ctx iris.Context) {
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
	})

	app.Get("/todos", func(ctx iris.Context) {
		var todos []Todo

		db, err = gorm.Open("mysql", mysqlDbURI)
		if err != nil {
			fmt.Println(err.Error())
			panic("Failed to connect to database")
		}
		defer db.Close()

		db.Find(&todos)
		ctx.JSON(iris.Map{"todos": todos})
	})

	app.Post("/todos", func(ctx iris.Context) {
		var todo Todo

		err := ctx.ReadJSON(&todo)
		println(todo.Value)
		if err != nil || todo.Value == "" {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.JSON(iris.Map{"message": "Post body must be a JSON object with at least a value!"})
			return
		}

		db, dbErr := gorm.Open("mysql", mysqlDbURI)

		if dbErr != nil {
			fmt.Println(dbErr.Error())
			panic("Failed to connect to database")
		}
		defer db.Close()

		db.Create(&Todo{Value: todo.Value, UserID: todo.UserID})
		ctx.JSON(iris.Map{"message": "Success", "todo": todo})
	})

	app.Run(iris.Addr(":3000"))
}
