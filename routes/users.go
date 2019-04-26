package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"
)

func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db, err := gorm.Open("mysql", mysqlDbURI)

	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to database")
	}
	defer db.Close()

	var users []User
	db.Find(&users)
	json.NewEncoder(w).Encode(users)

}

func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db, err := gorm.Open("mysql", mysqlDbURI)

	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to database")
	}
	defer db.Close()

	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)

	db.Create(&User{
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		Password:  user.Password})
	fmt.Fprintf(w, "User successfully created")
}
