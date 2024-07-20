package controllers

import (
	"log"
	"myJwtAuth/helpers"
	"myJwtAuth/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// get users handler
func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := helpers.CheckUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Access Restricted to this resource"})
			return
		}

		users, err := service.GetAllUsers()
		if err != nil {
			log.Println("Error while reteriving all users data")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		UserCount := len(users)
		log.Printf("List of users retrived with user count: %d", UserCount)
		c.JSON(http.StatusOK, gin.H{"Total User Count": UserCount, "users": users})
	}
}

func GetUserbyId() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")
		if err := helpers.CheckUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Access Restricted to this resource"})
			return
		}
		user, err := service.GetUserById(userId)
		if err != nil {
			log.Println("Error while reteriving users data")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		log.Printf("user with Id: %s retrevied ", userId)
		c.JSON(http.StatusOK, user)
	}
}

func DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")
		if err := helpers.CheckUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Access Restricted to this resource"})
			return
		}

		status, err := service.DeleteUserById(userId)
		if err != nil || !status {
			log.Printf("Error while deleting %q", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		log.Printf("User %s Deleted Successfully", userId)
		c.JSON(http.StatusOK, gin.H{"msg": "user deleted successfully"})

	}

}
