package helpers

import "regexp"

func IsEmail(email string) bool {
	//Regular expression pattern for email validation
	emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	//Compile the regular expression pattern
	re := regexp.MustCompile(emailPattern)

	//Check if the input match with the pattern
	return re.MatchString(email)
}

func IsPhone(phone string) bool {
	//Regular expression pattern for phone number validation
	phonePattern := `^\d{10,}$`

	//Compile the regular expression pattern
	re := regexp.MustCompile(phonePattern)

	//Check if the input match with the pattern
	return re.MatchString(phone)

}

// func StrToInt(v string) int64{

// }
