package helpers

import (
	"fmt"
	"log"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
)

func PasswordResetMail(userEmail string, token string) error {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error while loading env file: %v", err)
	}

	from := os.Getenv("MAIL_FROM")
	password := os.Getenv("MAIL_PASSWORD")
	app_host := os.Getenv("APP_HOST")

	//SMTP server configuration
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")

	//SUbject and message body
	subject := "Split the Bills: Password Reset Link"
	message := fmt.Sprintf("Click on the Below Link to reset password:\n" + app_host + "users/password-reset?token=" + token + "&email=" + userEmail + "\n")

	//Create messgae headers
	msg := "From:" + from + "\r\n"
	msg += "To:" + userEmail + "\r\n"
	msg += "Subject: " + subject + "\r\n\r\n"
	msg += message

	//Authentication
	auth := smtp.PlainAuth("", from, password, host)

	//Connect and send mail
	err = smtp.SendMail(fmt.Sprintf("%s:%s", host, port), auth, from, []string{userEmail}, []byte(msg))
	if err != nil {
		log.Printf("error while sending mail %v", err)
		return err
	}

	log.Println("Email sent Successfully")
	return nil
}
