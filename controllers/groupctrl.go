package controllers

import (
	"strconv"

	"github.com/Assenti/restapi/db"
	"github.com/Assenti/restapi/models"
	"github.com/kataras/iris"
)

// CreateGroup method
func CreateGroup(ctx iris.Context) {
	var group models.Group

	err := ctx.ReadJSON(&group)
	if err != nil || (group.Name == "") {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"msg": "Group Name"})
		return
	}

	db := db.Connect()
	defer db.Close()

	db.Create(&models.Group{Name: group.Name, UserID: group.UserID})

	var created models.Group
	db.Where("user_id = ?", group.UserID).Last(&created)

	db.Create(&models.GroupParticipants{GroupID: created.ID, UserID: created.UserID})

	var groups []models.Group
	db.Where("user_id = ?", group.UserID).Find(&groups)
	ctx.JSON(iris.Map{"groups": groups})
}

// GetUserGroups method
func GetUserGroups(ctx iris.Context) {
	var groups []models.Group
	id := ctx.URLParam("userid")
	db := db.Connect()
	defer db.Close()
	db.Where("user_id = ?", id).Find(&groups)
	ctx.JSON(iris.Map{"groups": groups})
}

// ChangeGroupName method
func ChangeGroupName(ctx iris.Context) {
	var group models.Group

	err := ctx.ReadJSON(&group)

	if err != nil || group.Name == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"msg": "Name must be provided"})
		return
	}

	db := db.Connect()
	defer db.Close()

	db.Model(&group).Where("id = ?", group.ID).Update("name", group.Name)

	var changed models.Group
	db.Where("id = ?", group.ID).Last(&changed)
	ctx.JSON(iris.Map{"group": changed})
}

// DeleteGroup method
func DeleteGroup(ctx iris.Context) {
	id := ctx.URLParam("id")

	db := db.Connect()
	defer db.Close()

	db.Where("id = ?", id).Delete(&models.Group{})
	ctx.JSON(iris.Map{"msg": "Group successfully deleted."})
}

// GetGroupParticipants method
func GetGroupParticipants(ctx iris.Context) {
	var participants []models.GroupParticipants
	groupID := ctx.URLParam("groupid")

	db := db.Connect()
	defer db.Close()

	db.Where("group_id = ?", groupID).Find(&participants)
	ctx.JSON(iris.Map{"participants": participants})
}

// AddParticipant method
func AddParticipant(ctx iris.Context) {
	var participants []models.GroupParticipants

	groupID := ctx.URLParam("groupid")
	userID := ctx.URLParam("userid")

	intGroupID, _ := strconv.ParseUint(groupID, 10, 64)
	intUserID, _ := strconv.ParseUint(userID, 10, 64)

	db := db.Connect()
	defer db.Close()

	db.Create(&models.GroupParticipants{GroupID: intGroupID, UserID: intUserID})
	db.Where("group_id = ?", intGroupID).Find(&participants)
	ctx.JSON(iris.Map{"participants": participants})
}
