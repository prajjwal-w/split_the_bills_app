package routes

import (
	"myJwtAuth/controllers"

	"github.com/gin-gonic/gin"
)

func ExpenseRoutes(routes *gin.Engine) {
	routes.POST("/expenses", controllers.AddExpenses())
}
