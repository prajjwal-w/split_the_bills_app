package helpers

import (
	"fmt"
	"log"
	"myJwtAuth/service"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

var SECRET_KEY string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	SECRET_KEY = os.Getenv("SECRET_KEY")
}

type SignedDetails struct {
	Email    string
	UserType string
	jwt.StandardClaims
}

func GenerateTokens(email string, userType string) (string, string, error) {
	claims := &SignedDetails{
		Email:    email,
		UserType: userType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {

		return "", "", err
	}
	refresh_token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))
	if err != nil {

		return "", "", err
	}

	return token, refresh_token, nil

}

func UpdateToken(token string, refreshToken string, userId int64) (bool, error) {

	status, err := service.UpdateTokenInDB(token, refreshToken, userId)
	return status, err

}

func ValidateToken(signedToken string) (claims *SignedDetails, err error) {

	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("while parsing jwt token %s", err)

	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		err = fmt.Errorf("invalid Token")
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = fmt.Errorf("token Expired")
		return
	}

	return claims, err
}
