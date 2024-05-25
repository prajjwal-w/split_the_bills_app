package helpers

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

// hash the password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return string(bytes), nil
}

func VerifyPassword(loginPass string, foundPass string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(foundPass), []byte(loginPass))
	if err != nil {
		return false, err
	}

	return true, nil
}
