package routes

import (
	"myJwtAuth/controllers"
	"myJwtAuth/middleware"

	"github.com/gin-gonic/gin"
)

// user routes
func UserRoutes(routes *gin.Engine) {
	routes.PUT("/users/", controllers.ResetPassword())
	routes.Use(middleware.Authenticate())
	routes.GET("/users", controllers.GetUsers())
	routes.GET("/users/:user_id", controllers.GetUserbyId())
	routes.DELETE("/users/:user_id", controllers.DeleteUser())

}
