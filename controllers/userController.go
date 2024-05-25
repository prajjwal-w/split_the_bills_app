package controllers

import (
	"log"
	"myJwtAuth/helpers"
	"myJwtAuth/models"
	"myJwtAuth/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// sign up handler
func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User

		//binding input json to the user struct
		if err := c.BindJSON(&user); err != nil {
			log.Println("Error while binding the input JSON")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//before registering user we are searching user with the same email in the db
		row, err := service.GetUserByEmail(user.Email)
		if err != nil {
			log.Println("Registering New user")
			//c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		}
		if row != nil {
			log.Println("User already present with the email")
			c.JSON(http.StatusConflict, gin.H{"msg": "User already exists with email"})
			return
		}
		//if user not present the we continue registering
		//hashing password
		password, err := helpers.HashPassword(user.Password)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		user.Password = password

		user.Created_at = time.Now()
		user.Updated_at = time.Now()

		//generating token and refresh token
		token, refreshToken, err := helpers.GenerateTokens(user.Email, user.User_Type)
		if err != nil {
			log.Println("Error while generating token", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "msg": "err while generating token"})
			return
		}
		user.Token = token
		user.Refresh_Token = refreshToken

		//inserting user into the db
		insertUser, insertErr := service.InsertUser(user)
		if insertErr != nil {
			log.Println("Error: User not created!...errror while inserting.")
			c.JSON(http.StatusInternalServerError, gin.H{"error": insertErr.Error(), "msg": "User not created"})
			return
		}
		log.Printf("User registeration sucessful with userId: %d", insertUser)
		c.JSON(http.StatusOK, gin.H{"msg": "User registerd successfully", "user_Id": insertUser})

	}
}

// login handler
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var login models.LoginUser

		if err := c.BindJSON(&login); err != nil {
			log.Println("Error while binding the input JSON")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		//search user with the email id
		user, err := service.GetUserByEmail(login.Email)
		if err != nil {
			log.Println("email is incorrect or maybe not registerd")
			c.JSON(http.StatusNotFound, gin.H{"msg": "user not found or incorrect email", "error": err.Error()})
			return
		}
		//log.Println(user)

		//check the password
		passwordIsValid, _ := helpers.VerifyPassword(login.Password, user.Password)
		if !passwordIsValid {
			log.Println("password incorrect")
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "incorrect password"})
			return
		}
		_, err = helpers.ValidateToken(user.Token)
		if err != nil {
			//generate new tokens
			token, refershToken, err := helpers.GenerateTokens(user.Email, user.User_Type)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"Error": err})
				return
			}

			//update new token to the db which is then used for Auth
			status, err := helpers.UpdateToken(token, refershToken, user.User_Id)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err})
				return
			} else if status {
				log.Println("Token updated Successfully!!")
			}

			//retrive user from db with the updated tokens
			userWithNewToken, err := service.GetUserByEmail(login.Email)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err})
				return
			}
			c.JSON(http.StatusOK, userWithNewToken)
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

// get users handler
func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := helpers.CheckUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Access Restricted to this resource"})
			return
		}

		users, err := service.GetAllUsers()
		if err != nil {
			log.Println("Error while reteriving all users data")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		UserCount := len(users)
		log.Printf("List of users retrived with user count: %d", UserCount)
		c.JSON(http.StatusOK, gin.H{"Total User Count": UserCount, "users": users})
	}
}

func GetUserbyId() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")
		if err := helpers.CheckUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Access Restricted to this resource"})
			return
		}
		user, err := service.GetUserById(userId)
		if err != nil {
			log.Println("Error while reteriving users data")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		log.Printf("user with Id: %s retrevied ", userId)
		c.JSON(http.StatusOK, user)
	}
}

func DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")
		if err := helpers.CheckUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Access Restricted to this resource"})
			return
		}

		status, err := service.DeleteUserById(userId)
		if err != nil || !status {
			log.Printf("Error while deleting %q", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		log.Printf("User %s Deleted Successfully", userId)
		c.JSON(http.StatusOK, gin.H{"msg": "user deleted successfully"})

	}

}

func ResetPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		var resetPass *models.ResetPass
		if err := c.BindJSON(&resetPass); err != nil {
			log.Println("error while binding the input")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		password, err := helpers.HashPassword(resetPass.NewPassword)
		if err != nil {
			log.Println("error while password hashing")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		_, err = service.ResetPassword(resetPass.UserId, password)
		if err != nil {
			log.Println("error while update password")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"msg": "Password reset sucessfully"})

	}
}
