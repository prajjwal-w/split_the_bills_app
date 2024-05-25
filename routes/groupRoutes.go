package routes

import (
	"myJwtAuth/controllers"
	"myJwtAuth/middleware"

	"github.com/gin-gonic/gin"
)

func GroupRoutes(routes *gin.Engine) {
	routes.Use(middleware.Authenticate())
	routes.POST("/groups", controllers.CreateGroup())
	routes.POST("/groups/adduser", controllers.AddUsersInAGroup())
	routes.GET("/groups/:user_id", controllers.GetAllGroups())
	//routes.DELETE("/users/:user_id", controllers.DeleteUser())

}
