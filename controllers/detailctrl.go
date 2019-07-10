package controllers

import (
	"github.com/Assenti/restapi/db"
	"github.com/Assenti/restapi/models"
	"github.com/kataras/iris"
)

// CreateDetail method
func CreateDetail(ctx iris.Context) {
	var detail models.TodoDetails

	err := ctx.ReadJSON(&detail)

	if err != nil || detail.Content == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"msg": "Content must be provided."})
		return
	}

	db := db.Connect()
	defer db.Close()

	db.Create(&models.TodoDetails{Content: detail.Content, TodoID: detail.TodoID})

	var newDetail models.TodoDetails
	db.Where("todo_id = ?", detail.TodoID).Last(&newDetail)
	ctx.JSON(iris.Map{"detail": newDetail})
}

// UpdateDetail method
func UpdateDetail(ctx iris.Context) {
	var detail models.TodoDetails

	err := ctx.ReadJSON(&detail)

	if err != nil || detail.Content == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"msg": "Value must be provided."})
		return
	}

	db := db.Connect()
	defer db.Close()

	db.Model(&detail).Where("id = ?", detail.ID).Update("content", detail.Content)

	var updated models.TodoDetails
	db.Where("id = ?", detail.ID).Last(&updated)
	ctx.JSON(iris.Map{"detail": updated})
}

// DeleteDetail method
func DeleteDetail(ctx iris.Context) {
	id := ctx.URLParam("id")

	db := db.Connect()
	defer db.Close()

	db.Where("id = ?", id).Delete(&models.TodoDetails{})
	ctx.JSON(iris.Map{"msg": "Detail successfully deleted."})
}

// GetDetails method
func GetDetails(ctx iris.Context) {
	var details []models.TodoDetails

	id := ctx.URLParam("todoid")

	db := db.Connect()
	defer db.Close()

	db.Where("todo_id = ?", id).Last(&details)

	ctx.JSON(iris.Map{"details": details})
}
