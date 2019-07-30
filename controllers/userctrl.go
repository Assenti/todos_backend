package controllers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Assenti/restapi/db"
	"github.com/Assenti/restapi/models"
	"github.com/kataras/iris"
	"github.com/sethvargo/go-password/password"
	"gopkg.in/dgrijalva/jwt-go.v3"
	"gopkg.in/gomail.v2"
)

const mailHost = "smtp.gmail.com"
const mailPort = 587
const mailSender = "Personal Planner <testyfy7@gmail.com>"
const mailUser = "testyfy7@gmail.com"
const mailPassword = "qwgqebnounjwnmfo"
const jwtSecret = "ju$tTe$t1t"
const sessionTime = 30

// CreateUser method (New user registration)
func CreateUser(ctx iris.Context) {
	var user models.User

	err := ctx.ReadJSON(&user)
	if err != nil || (user.Firstname == "" && user.Lastname == "" && user.Email == "" && user.Password == "") {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"msg": "Firstname, Lastname, Email and Password must be provided."})
		return
	}

	db := db.Connect()
	defer db.Close()

	hash, _ := HashPassword(user.Password)

	errors := db.Create(&models.User{
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		Password:  hash}).GetErrors()

	if len(errors) > 0 {
		var existErr string
		for _, err := range errors {
			existErr = err.Error()
		}
		if strings.Contains(existErr, "Duplicate entry") {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.JSON(iris.Map{"msg": "Such Email is registered yet"})
			return
		}
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"msg": existErr})
		return
	}

	// Configure the email message
	htmlBody := fmt.Sprintf(`
		<h1>Personal Planner</h1>
		<p style="font-size: 16px">Hello, %s %s, you are successfully registered on Personal Planner web app.</p>
		<hr>
		<p style="font-size: 11px">Do not reply to this message, it was generated automatically.</p>`,
		user.Firstname, user.Lastname)
	m := gomail.NewMessage()
	m.SetHeader("From", mailSender)
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
	var user models.User

	err := ctx.ReadJSON(&user)
	if err != nil || (user.Firstname == "" && user.Lastname == "") {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"msg": "Firstname, Lastname must be provided."})
		return
	}

	db := db.Connect()
	defer db.Close()

	db.Model(&user).Where("id = ?", user.ID).Updates(models.User{Firstname: user.Firstname, Lastname: user.Lastname})

	var updatedUser models.User
	db.Where("id = ?", user.ID).Last(&updatedUser)
	ctx.JSON(iris.Map{"user": updatedUser})
}

// GetUsersList method
func GetUsersList(ctx iris.Context) {
	var users []models.Performer
	db := db.Connect()
	defer db.Close()

	db.Table("users").Scan(&users)
	ctx.JSON(iris.Map{"users": users})
}

// Login method
func Login(ctx iris.Context) {
	var user models.UserSubmit

	err := ctx.ReadJSON(&user)
	if err != nil || (user.Email == "" && user.Password == "") {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"message": "Email and Password must be provided."})
		return
	}

	db := db.Connect()
	defer db.Close()

	var foundedUser models.User
	db.Where(&models.User{Email: user.Email}).First(&foundedUser)

	passMatch := IsPasswordMatch(user.Password, foundedUser.Password)

	if foundedUser.Email == user.Email && passMatch {
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

		temp := time.Now()
		db.Model(&foundedUser).Where("email = ?", foundedUser.Email).Update("last_logged_on", temp)

		var userInfo models.UserInfo
		userInfo.CreatedAt = foundedUser.CreatedAt
		userInfo.Email = foundedUser.Email
		userInfo.Firstname = foundedUser.Firstname
		userInfo.ID = foundedUser.ID
		userInfo.Lastname = foundedUser.Lastname
		userInfo.Token = tokenString
		userInfo.LastLoggedOn = temp

		userInfoInString, err := json.Marshal(&userInfo)
		if err != nil {
			panic(err)
		}
		encodedUserInfo := base64.StdEncoding.EncodeToString(userInfoInString)

		ctx.JSON(iris.Map{"user": encodedUserInfo})
	} else {
		// In case of wrong password send alerting email to account owner
		m := gomail.NewMessage()
		m.SetHeader("From", mailSender)
		m.SetHeader("To", foundedUser.Email)
		m.SetHeader("Subject", "Notification from Personal Planner app")
		m.SetBody("text/html", `
			<h1>Personal Planner</h1>
			<p style="font-size: 16px">Warning! Someone try to sign in to app using your account.</p>
			<hr>
			<p style="font-size: 11px">Do not reply to this message, it was generated automatically.</p>`)

		d := gomail.NewDialer(mailHost, mailPort, mailUser, mailPassword)

		// Send the email
		if err := d.DialAndSend(m); err != nil {
			panic(err)
		}

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

	db := db.Connect()
	defer db.Close()

	var foundedUser models.User
	db.Where(&models.User{Email: email}).First(&foundedUser)

	if foundedUser.Email == "" {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(iris.Map{"message": "There is no user with same Email."})
		return
	}

	res, err := password.Generate(12, 2, 2, true, false)
	if err != nil {
		log.Fatal(err)
	}

	db.Model(&foundedUser).Where(&models.User{ID: foundedUser.ID}).Updates(models.User{Password: res})

	// Configure the email message
	htmlBody := fmt.Sprintf(`
		<h1>Personal Planner</h1>
		<p style="font-size: 16px">You sent password restoring request. Generated new password: %s. Change it from your cabinet in app.</p>
		<hr>
		<p style="font-size: 11px">Do not reply to this message, it was generated automatically.</p>`, res)
	m := gomail.NewMessage()
	m.SetHeader("From", mailSender)
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
	var user models.User

	err := ctx.ReadJSON(&user)
	if err != nil || user.Password == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"msg": "Password must be provided."})
		return
	}

	db := db.Connect()
	defer db.Close()

	var foundedUser models.User
	db.Where(&models.User{ID: user.ID}).First(&foundedUser)

	db.Model(&foundedUser).Where(&models.User{ID: foundedUser.ID}).Updates(models.User{Password: user.Password})

	// Configure the email message
	htmlBody := fmt.Sprintf(`
		<h1>Personal Planner</h1>
		<p style="font-size: 16px">Your password successfully changed.</p>
		<hr>
		<p style="font-size: 11px">Do not reply to this message, it was generated automatically.</p>`)
	m := gomail.NewMessage()
	m.SetHeader("From", mailSender)
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
	var user models.User

	err := ctx.ReadJSON(&user)
	if err != nil && (user.Password == "") {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"message": "Password must be provided."})
		return
	}

	db := db.Connect()
	defer db.Close()

	var foundedUser models.User
	db.Where(&models.User{ID: user.ID}).First(&foundedUser)

	passMatch := IsPasswordMatch(user.Password, foundedUser.Password)

	if !passMatch {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(iris.Map{"message": "Invalid current Password."})
		return
	}
	ctx.StatusCode(iris.StatusOK)
}

// SendInvitation method
func SendInvitation(ctx iris.Context) {
	email := ctx.URLParam("email")
	inviter := ctx.URLParam("inviter")

	if email == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"msg": "Email must be provided"})
		return
	}

	appLink := "https://planner-2.herokuapp.com"

	// Configure the email message
	htmlBody := fmt.Sprintf(`
		<h1>Personal Planner</h1>
		<p style="font-size: 16px">Hey there, %s invited you to the Personal Planner web application. Try it now, follow this link %s</p>
		<hr>
		<p style="font-size: 11px">Do not reply to this message, it was generated automatically.</p>`,
		inviter, appLink)
	m := gomail.NewMessage()
	m.SetHeader("From", mailSender)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Invitation to app")
	m.SetBody("text/html", htmlBody)

	d := gomail.NewDialer(mailHost, mailPort, mailUser, mailPassword)

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
