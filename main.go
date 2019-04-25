package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

const mysqlDbURI = "PgQXfyC4AD:CV3B9cSf2k@tcp(remotemysql.com:3306)/PgQXfyC4AD"

// Todo model
type Todo struct {
	ID        string `json:"id"`
	Value     string `json:"value"`
	Important bool   `json:"important"`
	Completed bool   `json:"completed"`
	OwnerID   int    `json:"owner"`
	CreatedAt string `json:"createdAt"`
}

// User model
type User struct {
	ID        string `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

// Mock Todos
var todos []Todo

func getTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db, err := sql.Open("mysql", mysqlDbURI)

	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	fetch, err := db.Query("select * from todos")
	if err != nil {
		panic(err.Error())
	}

	var todo Todo

	for fetch.Next() {
		err = fetch.Scan(&todo.Value, &todo.ID, &todo.Completed, &todo.Important, &todo.OwnerID, &todo.CreatedAt)
		if err != nil {
			panic(err.Error())
		}
		todos = append(todos,
			Todo{ID: todo.ID,
				Value:     todo.Value,
				Completed: todo.Completed,
				Important: todo.Completed,
				OwnerID:   todo.OwnerID,
				CreatedAt: todo.CreatedAt})
	}

	json.NewEncoder(w).Encode(todos)
}

func getTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	db, err := sql.Open("mysql", mysqlDbURI)

	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	fetch, err := db.Query("select * from todos join users on todos.OwnerID = users.ID where todos.ID = " + params["id"])
	if err != nil {
		panic(err.Error())
	}

	var todo Todo

	for fetch.Next() {
		err = fetch.Scan(&todo.Value, &todo.ID, &todo.Completed, &todo.Important, &todo.OwnerID, &todo.CreatedAt)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println(todo)
		json.NewEncoder(w).Encode(todo)
	}
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db, err := sql.Open("mysql", "PgQXfyC4AD:CV3B9cSf2k@tcp(remotemysql.com:3306)/PgQXfyC4AD")

	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	var todo Todo
	_ = json.NewDecoder(r.Body).Decode(&todo)

	insert, err := db.Query("insert into todos (Value, OwnerID) values ('" + todo.Value + "', 1)")

	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()

	fetch, err := db.Query("select * from todos where Value = '" + todo.Value + "'")
	if err != nil {
		panic(err.Error())
	}

	for fetch.Next() {
		err = fetch.Scan(&todo.Value, &todo.ID, &todo.Completed, &todo.Important, &todo.OwnerID, &todo.CreatedAt)
		if err != nil {
			panic(err.Error())
		}
		json.NewEncoder(w).Encode(todo)
	}

	// Random ID generating
	// todo.ID = strconv.Itoa(rand.Intn(1000000))
	// todos = append(todos, todo)
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for index, item := range todos {
		if item.ID == params["id"] {
			todos = append(todos[:index], todos[index+1:]...)
			var todo Todo
			_ = json.NewDecoder(r.Body).Decode(&todo)
			todo.ID = params["id"]
			todos = append(todos, todo)
			json.NewEncoder(w).Encode(todo)
		}
	}
	json.NewEncoder(w).Encode(todos)
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for index, item := range todos {
		if item.ID == params["id"] {
			todos = append(todos[:index], todos[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(todos)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db, err := sql.Open("mysql", mysqlDbURI)

	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)

	insert, err := db.Query(
		"insert into users (Firstname, Lastname, Email, Password) values ('" + user.Firstname + "', '" + user.Lastname + "', '" + user.Email + "', '" + user.Password + "')")

	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()
	fmt.Fprint(w, "User successfully created")
}

func main() {
	r := mux.NewRouter()

	// todos = append(todos, Todo{ID: "1", Value: "Finish this service", Important: true, Completed: false, Owner: &User{Firstname: "Asset", Lastname: "Sultanov"}})
	// todos = append(todos, Todo{ID: "2", Value: "Upgrade this service", Important: true, Completed: false, Owner: &User{Firstname: "Asset", Lastname: "Sultanov"}})

	// Todos API
	r.HandleFunc("/api/todos", getTodos).Methods("GET")
	r.HandleFunc("/api/todo/{id}", getTodo).Methods("GET")
	r.HandleFunc("/api/todos", createTodo).Methods("POST")
	r.HandleFunc("/api/todos/{id}", updateTodo).Methods("PUT")
	r.HandleFunc("/api/todos/{id}", deleteTodo).Methods("DELETE")

	// Users API
	r.HandleFunc("/api/users", createUser).Methods("POST")

	fmt.Println("Server started...")
	log.Fatal(http.ListenAndServe(":3000", r))
}
