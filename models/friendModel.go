package models

type AddFriend struct {
	UserId       int64  `json:"userId" validate:"required"`
	EmailOrPhone string `json:"emailOrPhone" validate:"required"`
}

type Friend struct {
	FriendShipID int64  `json:"friendshipId"`
	UserId       int64  `json:"user_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
}
