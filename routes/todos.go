package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB
var err error

func getTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db, err = gorm.Open("mysql", mysqlDbURI)
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to database")
	}
	defer db.Close()

	params := mux.Vars(r)
	var todos []Todo
	db.Find(&todos, "user_id = ?", params["userid"])
	json.NewEncoder(w).Encode(todos)
}

func getTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db, err = gorm.Open("mysql", mysqlDbURI)
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to database")
	}
	defer db.Close()

	params := mux.Vars(r)
	var todo Todo
	db.Take(&todo, "id = ?", params["id"])
	json.NewEncoder(w).Encode(todo)
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db, err := gorm.Open("mysql", mysqlDbURI)

	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to database")
	}
	defer db.Close()

	var todo Todo
	_ = json.NewDecoder(r.Body).Decode(&todo)

	db.Create(&Todo{Value: todo.Value, UserID: todo.UserID})
	fmt.Fprintf(w, "Todo successfully created")

	// Random ID generating
	// todo.ID = strconv.Itoa(rand.Intn(1000000))
	// todos = append(todos, todo)
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db, err := gorm.Open("mysql", mysqlDbURI)

	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to database")
	}
	defer db.Close()

	var todo Todo
	var todoU Todo
	_ = json.NewDecoder(r.Body).Decode(&todoU)

	db.Model(&todo).Updates(Todo{Value: todoU.Value, Completed: todoU.Completed, Important: todoU.Important, UpdatedAt: time.Now()})
	fmt.Fprintf(w, "Todo successfully updated")
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var todo Todo
	db, err := gorm.Open("mysql", mysqlDbURI)

	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to database")
	}
	defer db.Close()

	db.Delete(&todo, params["id"])

	fmt.Fprintf(w, "Todo successfully deleted")
}
