package models

import "time"

//user
type User struct {
	User_Id       int64     `json:"user_id"`
	First_Name    string    `json:"first_name" validate:"required"`
	Last_Name     string    `json:"last_name"  validate:"required"`
	Password      string    `json:"password" validate:"required"`
	Phone         string    `json:"phone" validate:"required"`
	Email         string    `json:"email" validate:"required"`
	Token         string    `json:"token"`
	User_Type     string    `json:"user_type"   validate:"required, eq=ADMIN||eq=USER"`
	Refresh_Token string    `json:"refresh_token"`
	Created_at    time.Time `json:"created_at"`
	Updated_at    time.Time `json:"updated_at"`
}

type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ResetPass struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password"`
	Token    string `json:"token"`
}
