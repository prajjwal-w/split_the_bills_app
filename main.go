package main

import (
	"log"
	"myJwtAuth/routes"

	"github.com/gin-gonic/gin"
)

const port = "8080"

// server
func main() {
	router := gin.New()
	router.Use(gin.Logger())

	routes.AuthRoutes(router)
	routes.FriendRoutes(router)
	routes.UserRoutes(router)
	routes.GroupRoutes(router)
	routes.ExpenseRoutes(router)

	err := router.Run(":" + port)
	if err != nil {
		log.Fatal("Error while starting server: ", err)
	}
}
