package routes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func greeting(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from Todos REST API")
}

// HandleRequests function
func HandleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"http://localhost:8080"})
	router.HandleFunc("/api", greeting).Methods("GET")

	// Todos API
	router.HandleFunc("/api/todos/{userid}", getTodos).Methods("GET")
	router.HandleFunc("/api/todo/{id}", getTodo).Methods("GET")
	router.HandleFunc("/api/todos", createTodo).Methods("POST")
	router.HandleFunc("/api/todos", updateTodo).Methods("PUT")
	router.HandleFunc("/api/todos/{id}", deleteTodo).Methods("DELETE")

	// Users API
	router.HandleFunc("/api/login", login).Methods("POST")
	router.HandleFunc("/api/users", getUsers).Methods("GET")
	router.HandleFunc("/api/users", createUser).Methods("POST")

	fmt.Println("Server started...")
	log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(headers, methods, origins)(router)))
}
