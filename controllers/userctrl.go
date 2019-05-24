package controllers

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/kataras/iris"
	"gopkg.in/dgrijalva/jwt-go.v3"
	"gopkg.in/gomail.v2"
)

// User model
type User struct {
	ID        uint      `gorm:"primary_key"`
	CreatedAt time.Time `gorm:"default: CURRENT_TIMESTAMP"`
	Firstname string    `gorm:"type:varchar(100)"`
	Lastname  string    `gorm:"type:varchar(100)"`
	Email     string    `gorm:"unique_index"`
	Password  string    `gorm:"type:varchar(100)"`
}

const mailHost = "smtp.gmail.com"
const mailPort = 587
const mailUser = "testyfy7@gmail.com"
const mailPassword = "qwgqebnounjwnmfo"
const jwtSecret = "ju$tTe$t1t"
const sessionTime = 30

// CreateUser method
func CreateUser(ctx iris.Context) {
	var user User

	err := ctx.ReadJSON(&user)
	if err != nil || (user.Firstname == "" && user.Lastname == "" && user.Email == "" && user.Password == "") {
		ctx.StatusCode(400)
		ctx.JSON(iris.Map{"msg": "Firstname, Lastname, Email and Password must be provided."})
		return
	}

	db, dbErr := gorm.Open("mysql", mysqlDbURI)

	if dbErr != nil {
		fmt.Println(dbErr.Error())
		panic("Failed to connect to database")
	}
	defer db.Close()

	db.Create(&User{Firstname: user.Firstname, Lastname: user.Lastname, Email: user.Email, Password: user.Password})

	// Configure the email message
	htmlBody := fmt.Sprintf("<h3>Hello, %s %s, you are successfully registered on Personal Planner web app.</h3>", user.Firstname, user.Lastname)
	m := gomail.NewMessage()
	m.SetHeader("From", mailUser)
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", "Registration confirmation")
	m.SetBody("text/html", htmlBody)

	d := gomail.NewDialer(mailHost, mailPort, mailUser, mailPassword)

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}

	ctx.JSON(iris.Map{"msg": "New user successfully created"})
}

// UpdateUser method
func UpdateUser(ctx iris.Context) {
	var user User

	err := ctx.ReadJSON(&user)
	if err != nil || (user.Firstname == "" && user.Lastname == "") {
		ctx.StatusCode(400)
		ctx.JSON(iris.Map{"msg": "Firstname, Lastname must be provided."})
		return
	}

	db, dbErr := gorm.Open("mysql", mysqlDbURI)

	if dbErr != nil {
		fmt.Println(dbErr.Error())
		panic("Failed to connect to database")
	}
	defer db.Close()

	db.Model(&user).Where("id = ?", user.ID).Updates(User{Firstname: user.Firstname, Lastname: user.Lastname})

	var updatedUser User
	db.Where("id = ?", user.ID).Last(&updatedUser)
	ctx.JSON(iris.Map{"user": updatedUser})
}

// GetUsersList method
func GetUsersList(ctx iris.Context) {
	var users []User
	db, dbErr := gorm.Open("mysql", mysqlDbURI)

	if dbErr != nil {
		fmt.Println(dbErr.Error())
		panic("Failed to connect to database")
	}
	defer db.Close()

	db.Find(&users)
	ctx.JSON(iris.Map{"users": users})
}

// Login method
func Login(ctx iris.Context) {
	var user User

	err := ctx.ReadJSON(&user)
	if err != nil || (user.Email == "" && user.Password == "") {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"message": "Email and Password must be provided."})
		return
	}

	db, dbErr := gorm.Open("mysql", mysqlDbURI)

	if dbErr != nil {
		fmt.Println(dbErr.Error())
		panic("Failed to connect to database")
	}
	defer db.Close()

	var foundedUser User
	db.Where(&User{Email: user.Email, Password: user.Password}).First(&foundedUser)

	if foundedUser.Email == user.Email && foundedUser.Password == user.Password {
		// Create the JWT key used to create the signature
		var jwtKey = []byte(jwtSecret)

		// Create a struct that will be encoded to a JWT.
		// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
		type Claims struct {
			Firstname string `json:"firstname"`
			jwt.StandardClaims
		}

		expirationTime := time.Now().Add(sessionTime * time.Minute)
		// Create the JWT claims, which includes the firstname and expiry time
		claims := &Claims{
			Firstname: user.Firstname,
			StandardClaims: jwt.StandardClaims{
				// In JWT, the expiry time is expressed as unix milliseconds
				ExpiresAt: expirationTime.Unix(),
			},
		}

		// Declare the token with the algorithm used for signing, and the claims
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// Create the JWT string
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
		}

		ctx.JSON(iris.Map{"user": &foundedUser, "token": tokenString})
	} else {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"msg": "Invalid Email or Password."})
	}
}
