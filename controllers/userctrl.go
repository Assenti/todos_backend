package controllers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/kataras/iris"
	"github.com/sethvargo/go-password/password"
	"gopkg.in/dgrijalva/jwt-go.v3"
	"gopkg.in/gomail.v2"
)

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
		ctx.StatusCode(iris.StatusBadRequest)
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
	htmlBody := fmt.Sprintf("<h4>Hello, %s %s, you are successfully registered on Personal Planner web app.</h4>", user.Firstname, user.Lastname)
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
		ctx.StatusCode(iris.StatusBadRequest)
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
	var user UserSubmit

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
			Email string `json:"email"`
			jwt.StandardClaims
		}

		var standardClaims jwt.StandardClaims
		expirationTime := time.Now().Add(sessionTime * time.Minute)
		standardClaims.ExpiresAt = expirationTime.Unix()
		claims := &Claims{
			Email: user.Email,
		}

		if !user.Remember {
			claims.StandardClaims = standardClaims
		}

		// Declare the token with the algorithm used for signing, and the claims
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// Create the JWT string
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
		}

		var userInfo UserInfo
		userInfo.CreatedAt = foundedUser.CreatedAt
		userInfo.Email = foundedUser.Email
		userInfo.Firstname = foundedUser.Firstname
		userInfo.ID = foundedUser.ID
		userInfo.Lastname = foundedUser.Lastname
		userInfo.Token = tokenString

		userInfoInString, err := json.Marshal(&userInfo)
		if err != nil {
			panic(err)
		}
		encodedUserInfo := base64.StdEncoding.EncodeToString(userInfoInString)
		ctx.JSON(iris.Map{"user": encodedUserInfo})
	} else {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"msg": "Invalid Email or Password."})
	}
}

// RestorePassword method
func RestorePassword(ctx iris.Context) {
	email := ctx.URLParam("email")

	if email == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"message": "Email must be provided."})
		return
	}

	db, dbErr := gorm.Open("mysql", mysqlDbURI)

	if dbErr != nil {
		fmt.Println(dbErr.Error())
		panic("Failed to connect to database")
	}
	defer db.Close()

	var foundedUser User
	db.Where(&User{Email: email}).First(&foundedUser)

	if foundedUser.Email == "" {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(iris.Map{"message": "There is no user with same Email."})
		return
	}

	res, err := password.Generate(12, 2, 2, true, false)
	if err != nil {
		log.Fatal(err)
	}

	db.Model(&foundedUser).Where(&User{ID: foundedUser.ID}).Updates(User{Password: res})

	// Configure the email message
	htmlBody := fmt.Sprintf("<h4>You sent password restoring request. Generated new password: %s. Change it from your cabinet in app.</h4>", res)
	m := gomail.NewMessage()
	m.SetHeader("From", mailUser)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Password Restore")
	m.SetBody("text/html", htmlBody)

	d := gomail.NewDialer(mailHost, mailPort, mailUser, mailPassword)

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
	ctx.JSON(iris.Map{"message": "Password was sent to Email."})
}

// ChangePassword method
func ChangePassword(ctx iris.Context) {
	var user User

	err := ctx.ReadJSON(&user)
	if err != nil || user.Password == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"msg": "Password must be provided."})
		return
	}

	db, dbErr := gorm.Open("mysql", mysqlDbURI)

	if dbErr != nil {
		fmt.Println(dbErr.Error())
		panic("Failed to connect to database")
	}
	defer db.Close()

	var foundedUser User
	db.Where(&User{ID: user.ID}).First(&foundedUser)

	db.Model(&foundedUser).Where(&User{ID: foundedUser.ID}).Updates(User{Password: user.Password})

	// Configure the email message
	htmlBody := fmt.Sprintf("<p>Your password successfully changed.</p>")
	m := gomail.NewMessage()
	m.SetHeader("From", mailUser)
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", "Password Change")
	m.SetBody("text/html", htmlBody)

	d := gomail.NewDialer(mailHost, mailPort, mailUser, mailPassword)

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
	ctx.JSON(iris.Map{"message": "Password successfully changed."})
}

// CheckPassword method
func CheckPassword(ctx iris.Context) {
	var user User

	err := ctx.ReadJSON(&user)
	if err != nil && (user.Password == "") {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"message": "Password must be provided."})
		return
	}

	db, dbErr := gorm.Open("mysql", mysqlDbURI)

	if dbErr != nil {
		fmt.Println(dbErr.Error())
		panic("Failed to connect to database")
	}
	defer db.Close()

	var foundedUser User
	db.Where(&User{ID: user.ID}).First(&foundedUser)

	if foundedUser.Password != user.Password {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(iris.Map{"message": "Invalid current Password."})
		return
	}
	ctx.StatusCode(iris.StatusOK)
}
