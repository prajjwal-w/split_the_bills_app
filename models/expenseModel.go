package models

import "time"

//Expense
type Expense struct {
	ExpenseId   int64        `json:"expense_id"`
	Amount      float64      `json:"amount"`
	Description string       `json:"description"`
	PaidById    int64        `json:"paid_by_id"`
	PaidBYName  string       `json:"paid_by_name"`
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

//expenseStatement
type Statement struct {
	PayeeId        int64   `json:"payee_id"`
	PayeeFirstName string  `json:"payee_firstname"`
	PayeeLastName  string  `json:"payee_lastname"`
	UnsettledAmt   float64 `json:"unsettled_amount"`
}

//expense Settle
type SettleAmt struct {
	Payer_id int64   `json:"payer_id"`
	Payee_id int64   `json:"payee_id"`
	Amount   float64 `json:"amt"`
}
