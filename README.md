# REST API for Personal Planner ver. 2.0 web app

> Used Stack: Go lang, IRIS web framework, GORM, MySQL DB

## Build app
```
$ go build
```

## Run app
```
$ go run main.go
```

## Methods List

> GET /todo - controllers.GetSingleTodo

> GET /usertodos - controllers.GetUserTodos

> GET /todocompletion - controllers.ToggleTodoCompletion

> GET /todoimportance - controllers.ToggleTodoImportance

> GET /todos - controllers.GetTodos

> POST /todos - controllers.CreateTodo

> PUT /todos - controllers.UpdateTodo

> DELETE /todos - controllers.DeleteTodo

> POST /sendViaEmail - controllers.SendTodosListViaEmail

> POST /users - controllers.CreateUser

> PUT /users - controllers.UpdateUser

> GET /users - controllers.GetUsersList

> POST /login - controllers.Login

> GET /restorepassword - controllers.RestorePassword

> POST /changepassword - controllers.ChangePassword

> POST /checkpassword - controllers.CheckPassword 
