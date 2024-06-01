package controllers

import (
	"fmt"
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
		var expense *models.Expense

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

func SettleUpExpense() gin.HandlerFunc {
	return func(c *gin.Context) {
		var settleReq *models.SettleAmt

		if err := c.BindJSON(&settleReq); err != nil {
			log.Println("error while binding expense json")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := service.SettleUpExpense(settleReq)
		if err != nil {
			log.Println("error while settling up the expense")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		log.Println("expense settled successfully")
		c.JSON(http.StatusOK, gin.H{"msg": fmt.Sprintf("Expense between %v and %v of amount %v settled", settleReq.Payee_id, settleReq.Payer_id, settleReq.Amount)})
	}
}

func GetUnsettledAmountByUser() gin.HandlerFunc {
	return func(c *gin.Context) {

		user_id := c.Param("user_id")

		statements, err := service.GetUnsettledAmountByUser(user_id)
		if err != nil {
			log.Printf("Error while gettings unsettled statement:%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		log.Printf("statements retrived successfully")
		c.JSON(http.StatusOK, gin.H{"msg": "statements retrived successfully", "statements": statements})
	}
}
