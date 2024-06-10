package routes

import (
	"myJwtAuth/controllers"

	"github.com/gin-gonic/gin"
)

// func expense routes
func ExpenseRoutes(routes *gin.Engine) {
	routes.POST("/expenses", controllers.AddExpenses())
	routes.PUT("/expenses", controllers.SettleUpExpense())
	routes.GET("/expenses/:user_id", controllers.GetUnsettledAmountByUser())
}
