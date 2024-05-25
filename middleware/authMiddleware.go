package middleware

import (
	"log"
	"myJwtAuth/helpers"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("Authorization")
		if clientToken == "" || !strings.HasPrefix(clientToken, "Bearer ") {
			log.Println("No Authorization headers provided")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No Authorization headers provided"})
			c.Abort()
			return
		}

		clientToken = strings.TrimPrefix(clientToken, "Bearer ")
		claims, err := helpers.ValidateToken(clientToken)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("user_type", claims.UserType)
		c.Next()
	}
}
