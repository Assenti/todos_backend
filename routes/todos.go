package routes

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

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
	var todos []Todo

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
	var todos []Todo

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
	var todos []Todo

	for index, item := range todos {
		if item.ID == params["id"] {
			todos = append(todos[:index], todos[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(todos)
}
