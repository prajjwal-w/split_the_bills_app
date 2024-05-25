package controllers

import (
	"fmt"
	"log"
	"myJwtAuth/helpers"
	"myJwtAuth/models"
	"myJwtAuth/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Add friend handler
func Addfriend() gin.HandlerFunc {
	return func(c *gin.Context) {
		var frd models.AddFriend

		if err := c.BindJSON(&frd); err != nil {
			log.Println("Error while input JSON binding")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var user_id int64
		var err error
		if helpers.IsEmail(frd.EmailOrPhone) || helpers.IsPhone(frd.EmailOrPhone) {
			if helpers.IsEmail(frd.EmailOrPhone) {
				user_id, err = service.GetUserByEmailLimitedData(frd.EmailOrPhone)
			} else {
				user_id, err = service.GetUserByPhoneLimitedData(frd.EmailOrPhone)
			}

			if err != nil {
				log.Printf("Error while checking the frienduser :%v", err.Error())
				c.JSON(http.StatusOK, gin.H{"error": err.Error(), "msg": "Sending Invitation link to the user"})
				return
			}
		}

		frdship_id, err := service.AddFriend(frd.UserId, user_id)
		if err != nil {
			log.Printf("Error while creating frdship: %v", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "msg": "maybe frd user is not registerd!!"})
			return
		}

		log.Printf("FriendShip created between %d and %d users and createdby %d user", frd.UserId, user_id, frd.UserId)
		c.JSON(http.StatusCreated, gin.H{"friendship_id": frdship_id, "msg": fmt.Sprintf("FriendShip created between %d and %d users and createdby %d user", frd.UserId, user_id, frd.UserId)})
	}

}

func GetAllFriendShipByUserId() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_Id := c.Param("user_id")

		friendList, err := service.GetAllFriends(user_Id)
		if err != nil {
			log.Printf("error while retriving friends details: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		friendCount := len(friendList)
		log.Printf("Retrived friends count: %d", friendCount)
		c.JSON(http.StatusOK, gin.H{"friendsCount": friendCount, "friends": friendList})

	}
}
