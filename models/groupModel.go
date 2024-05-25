package models

import "time"

type Group struct {
	GroupId     int64         `json:"group_id"`
	GroupName   string        `json:"group_name"`
	Description string        `json:"description"`
	Members     []*GroupUser  `json:"group_members"`
	Expenses    []*AddExpense `json:"group_expenses"`
	CreatedBy   int64         `json:"created_by"`
	Created_at  time.Time
}

type AddUserGrp struct {
	GroupId int64 `json:"group_id"`
	UserId  int64 `json:"user_id"`
}

type GroupUser struct {
	UserId     int64  `json:"user_id"`
	First_Name string `json:"first_name"`
	Last_Name  string `json:"last_name"`
	Email      string `json:"email"`
	Phone      string `json:"Phone"`
}