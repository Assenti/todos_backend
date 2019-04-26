package main

import (
	"fmt"

	"github.com/Assenti/restapi/routes"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	fmt.Println("Main invoked")
	// todos = append(todos, Todo{ID: "1", Value: "Finish this service", Important: true, Completed: false, Owner: &User{Firstname: "Asset", Lastname: "Sultanov"}})
	// todos = append(todos, Todo{ID: "2", Value: "Upgrade this service", Important: true, Completed: false, Owner: &User{Firstname: "Asset", Lastname: "Sultanov"}})
	routes.HandleRequests()
}
