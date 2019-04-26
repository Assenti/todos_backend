package routes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func greeting(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from Todos REST API")
}

// HandleRequests
func HandleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/api", greeting).Methods("GET")

	// Todos API
	router.HandleFunc("/api/todos", getTodos).Methods("GET")
	router.HandleFunc("/api/todo/{id}", getTodo).Methods("GET")
	router.HandleFunc("/api/todos", createTodo).Methods("POST")
	router.HandleFunc("/api/todos/{id}", updateTodo).Methods("PUT")
	router.HandleFunc("/api/todos/{id}", deleteTodo).Methods("DELETE")

	// Users API
	router.HandleFunc("/api/users", createUser).Methods("POST")

	fmt.Println("Server started...")
	log.Fatal(http.ListenAndServe(":3000", router))
}
