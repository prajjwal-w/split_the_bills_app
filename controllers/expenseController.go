package controllers

import (
	"log"
	"myJwtAuth/models"
	"myJwtAuth/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// add expense handler
func AddExpenses() gin.HandlerFunc {
	return func(c *gin.Context) {
		var expense *models.AddExpense

		if err := c.BindJSON(&expense); err != nil {
			log.Println("error while binding expense json")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		expense.Created_at = time.Now()
		exp, err := service.AddExpense(expense)
		if err != nil {
			log.Printf("error while adding expense: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		log.Printf("expenses added: %v", exp)
		c.JSON(http.StatusOK, gin.H{"msg": "expense added successfully", "response": exp})

	}
}

// func UpdateSplitExpense() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var updateExp *models.SplitUsers

// 		if err := c.BindJSON(&updateExp); err != nil {
// 			log.Printf("error while updateExp json binding: %v", err)
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}

// 		exp, err := service.UpdateExpense(updateExp)
// 	}
// }
