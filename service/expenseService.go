package service

import (
	"context"
	"fmt"
	"strings"

	"myJwtAuth/models"
	"time"
)

func AddExpense(exp *models.AddExpense) (*models.AddExpense, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//need to add data into one or more table to ensure data is getting insert correctly in all table
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error: %v", err.Error())
	}

	query := `INSERT INTO expenses(amount, description, paid_by, group_id, created_at) VALUES ($1, $2, $3, $4, $5) returning expense_id`

	err = tx.QueryRowContext(ctx, query,
		exp.Amount,
		exp.Description,
		exp.PaidById,
		exp.GroupId,
		exp.Created_at,
	).Scan(&exp.ExpenseId)

	//if any error occer while inserting expense details we will rollback the transaction
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error while creating expense: %v", err.Error())
	}
	/////////////////////////////////////////////////////////////////////////////////////////////////////
	// //split the expenses
	// amt := SplitEqual(exp.SplitWith, exp.Amount)

	// splitQuery := `INSERT INTO expense_splits(expense_id, user_id, split_amt, status)
	//                SELECT $1, unnest($2::int[]), $3`

	// _, err = tx.ExecContext(ctx, splitQuery, exp.ExpenseId, pq.Array(exp.SplitWith), amt)
	// if err != nil {
	// 	tx.Rollback()
	// 	return nil, fmt.Errorf("error while inserting split expenses: %v", err.Error())
	// }
	/////////////////////////////////////////////////////////////////////////////////////////////////////
	splitquery := `INSERT INTO expense_splits (expense_id, user_id, split_amt, status) VALUES `
	values := []interface{}{}
	placeholders := []string{}

	for i, split := range exp.SplitWith {
		index := i * 4
		placeholders = append(placeholders, fmt.Sprintf("($%d,$%d,$%d,$%d)", index+1, index+2, index+3, index+4))
		values = append(values, exp.ExpenseId, split.UserId, split.AmountSplit, split.Status)
	}

	splitquery += strings.Join(placeholders, ", ")
	_, err = tx.ExecContext(ctx, splitquery, values...)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error while commiting transaction: %v", err.Error())
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("error while commiting transaction: %v", err.Error())
	}

	return exp, nil

}

func SplitEqual(users []models.SplitUsers, amt float64) float64 {
	userCount := len(users)

	splitAmt := amt / float64(userCount)

	return splitAmt
}

// func UpdateExpense(splitExp *models.SplitUsers) (*models.SplitUsers, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	query := `UPDATE`
// }
