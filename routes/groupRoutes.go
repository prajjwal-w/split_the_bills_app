package routes

import (
	"myJwtAuth/controllers"
	"myJwtAuth/middleware"

	"github.com/gin-gonic/gin"
)

// group routes
func GroupRoutes(routes *gin.Engine) {
	routes.Use(middleware.Authenticate())
	routes.POST("/groups", controllers.CreateGroup())
	routes.POST("/groups/adduser", controllers.AddUsersInAGroup())
	routes.GET("/groups/:user_id", controllers.GetAllGroups())
	routes.DELETE("/groups/deleteUser", controllers.RemoveUserFromGroup())

}
