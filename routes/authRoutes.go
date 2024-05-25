package routes

import (
	"myJwtAuth/controllers"

	"github.com/gin-gonic/gin"
)

// auth routes
func AuthRoutes(routes *gin.Engine) {
	routes.POST("/users/signup", controllers.SignUp())
	routes.POST("/users/login", controllers.Login())
}
