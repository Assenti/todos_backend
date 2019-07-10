package main

import (
	"os"

	"github.com/Assenti/restapi/api"
	"github.com/Assenti/restapi/models"

	// _ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/go-sql-driver/mysql"
	"github.com/kataras/iris"
)

func main() {
	var port string
	envPort := os.Getenv("PORT")
	if envPort != "" {
		port = envPort
	} else {
		port = "3000"
	}

	models.InitDb()

	app := iris.Default()

	api.InitRoutes(app)

	app.Run(iris.Addr(":" + port))
}
