package controllers

import (
	"fmt"
	"log"
	"myJwtAuth/models"
	"myJwtAuth/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// create group
func CreateGroup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var group models.Group

		if err := c.BindJSON(&group); err != nil {
			log.Println("Error while input JSON binding")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		group.Created_at = time.Now()
		grpId, grpUsers, err := service.CreateGroup(group)

		if err != nil {
			log.Printf("error while creating group: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "msg": "error while creating group"})
			return
		}
		group.GroupId = grpId
		group.Members = grpUsers

		log.Printf("Group %v Created Successfully by user %v , groupId: %v ", group.GroupName, group.CreatedBy, grpId)
		c.JSON(http.StatusCreated, gin.H{"msg": "Group created Successfully", "group_details": group})

	}
}

// retrive the groups in which the user is present
func GetAllGroups() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Param("user_id")
		useridInt, _ := strconv.ParseInt(user_id, 10, 64)
		groups, err := service.GetAllGroupsbyUser(useridInt)
		if err != nil {
			log.Printf("error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		log.Printf("Group details successfully retrivewed of user %v", user_id)
		c.JSON(http.StatusOK, gin.H{"group_details": groups})

	}
}

// add a user in group
func AddUsersInAGroup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.AddUserGrp

		if err := c.BindJSON(&user); err != nil {
			log.Printf("error while binding JSON: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		err := service.AddUsersToGroup(user.GroupId, user.UserId)

		if err != nil {
			log.Printf("error while adding user to group: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while adding user to group"})
			return
		}

		log.Printf("User %v added successfully to group %v", user.UserId, user.GroupId)
		c.JSON(http.StatusOK, gin.H{"msg": fmt.Sprintf("User %v added successfully to group %v", user.UserId, user.GroupId)})
	}
}
