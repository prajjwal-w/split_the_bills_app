package models

import "time"

//Expense
type AddExpense struct {
	ExpenseId   int64        `json:"expense_id"`
	Amount      float64      `json:"amount"`
	Description string       `json:"description"`
	PaidById    int64        `json:"paid_by_id"`
	GroupId     int64        `json:"group_id"`
	SplitWith   []SplitUsers `json:"splitwith"`
	Created_at  time.Time    `json:"created_at"`
}

//Splitusers
type SplitUsers struct {
	ExpneseId   int64   `json:"expense_id"`
	UserId      int64   `json:"user_id"`
	AmountSplit float64 `json:"amt_split"`
	Status      string  `json:"status"`
}
