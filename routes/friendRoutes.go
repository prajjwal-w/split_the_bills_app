package routes

import (
	"myJwtAuth/controllers"

	"github.com/gin-gonic/gin"
)

// friend routes
func FriendRoutes(routes *gin.Engine) {
	routes.POST("/friends/addFriend", controllers.Addfriend())
	routes.GET("/friends/:user_id", controllers.GetAllFriendShipByUserId())
}
