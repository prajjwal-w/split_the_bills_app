package helpers

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// check the user type
func CheckUserType(c *gin.Context, role string) (err error) {
	userType := c.GetString("user_type")
	if userType != role {
		return fmt.Errorf("unauthorized access to this resource")

	}
	return nil
}
