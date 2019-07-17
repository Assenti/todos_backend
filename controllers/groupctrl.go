package controllers

import (
	"fmt"
	"strconv"
	"strings"

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
	var participants []models.JoinedGroupParticipants
	groupID := ctx.URLParam("groupid")

	db := db.Connect()
	defer db.Close()

	db.Table("group_participants").Raw(`SELECT 
		group_participants.id, 
		group_participants.user_id, 
		group_participants.group_id, 
		users.firstname, 
		users.lastname
		FROM group_participants
		LEFT JOIN users ON users.id = group_participants.user_id
		WHERE group_id = ?`, groupID).Scan(&participants)

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

// GetGroupsWhereParticipate method
func GetGroupsWhereParticipate(ctx iris.Context) {
	var groups []models.Group
	userID := ctx.URLParam("userid")

	db := db.Connect()
	defer db.Close()

	var groupParticipants []models.GroupParticipants
	db.Where("user_id = ?", userID).Find(&groupParticipants)

	var groupIDs []uint64

	for _, p := range groupParticipants {
		groupIDs = append(groupIDs, p.GroupID)
	}

	var IDsInString []string

	for _, ID := range groupIDs {
		IDsInString = append(IDsInString, strconv.FormatUint(ID, 10))
	}

	uniqueIDsInString := Unique(IDsInString)
	stringifiedIDs := strings.Join(uniqueIDsInString, ",")

	query := fmt.Sprintf("SELECT * FROM PgQXfyC4AD.groups WHERE id IN (%s)", stringifiedIDs)

	db.Table("groups").Raw(query).Scan(&groups)

	ctx.JSON(iris.Map{"groups": groups})
}

// GetGroupsTodos method
func GetGroupsTodos(ctx iris.Context) {
	var todos []models.JoinedTodo

	id := ctx.URLParam("groupId")

	db := db.Connect()
	defer db.Close()

	db.Table("users").Raw(`SELECT 
				todos.id, 
				todos.value, 
				todos.created_at, 
				todos.updated_at, 
				todos.important, 
				todos.completed, 
				todos.user_id, 
				todos.complete_date, 
				todos.status, 
				todos.group_id, 
				todos.performer, 
				users.firstname, 
				users.lastname
				FROM todos
				LEFT JOIN users ON todos.performer = users.id
				where group_id = ?`, id).Scan(&todos)

	ctx.JSON(iris.Map{"todos": todos})
}
