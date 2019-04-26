package routes

const mysqlDbURI = "PgQXfyC4AD:CV3B9cSf2k@tcp(remotemysql.com:3306)/PgQXfyC4AD"
const port = "3000"

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
