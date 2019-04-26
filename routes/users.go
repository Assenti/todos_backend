package routes

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

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
