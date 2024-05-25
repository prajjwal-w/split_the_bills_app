package routes

import (
	"myJwtAuth/controllers"

	"github.com/gin-gonic/gin"
)

func FriendRoutes(routes *gin.Engine) {
	routes.POST("/friends/addFriend", controllers.Addfriend())
	routes.GET("/friends/:user_id", controllers.GetAllFriendShipByUserId())
}
