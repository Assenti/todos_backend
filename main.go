package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Todo model
type Todo struct {
	ID        string `json:"id"`
	Value     string `json:"value"`
	Important bool   `json:"important"`
	Completed bool   `json:"completed"`
	Owner     *User  `json:"owner"`
}

// User model
type User struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	// email     string `json: "email"`
	// password  string `json: "password"`
}

// Mock Todos
var todos []Todo

func getTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func getTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for _, item := range todos {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Todo{})
}

func main() {
	r := mux.NewRouter()

	todos = append(todos, Todo{ID: "1", Value: "Finish this service", Important: true, Completed: false, Owner: &User{Firstname: "Asset", Lastname: "Sultanov"}})
	todos = append(todos, Todo{ID: "2", Value: "Upgrade this service", Important: true, Completed: false, Owner: &User{Firstname: "Asset", Lastname: "Sultanov"}})

	r.HandleFunc("/api/todos", getTodos).Methods("GET")
	r.HandleFunc("/api/todo/{id}", getTodo).Methods("GET")
	// r.HandleFunc("/api/todos", createTodos).Methods("POST")
	// r.HandleFunc("/api/todos/{id}", updateTodo).Methods("PUT")
	// r.HandleFunc("/api/todos/{id}", deleteTodo).Methods("DELETE")

	fmt.Println("Server started...")
	log.Fatal(http.ListenAndServe(":3000", r))
}
