package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"myJwtAuth/models"
	"time"
)

// add expense service
func AddExpense(exp *models.Expense) (*models.Expense, error) {
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

func GetUnsettledAmountByUser(user_id string) ([]*models.Statement, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	userIdInt, err := strconv.ParseInt(user_id, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %v", err)
	}
	query := `SELECT
	             e.paid_by AS payee_id,
				 u.first_name As payee_first_name,
				 u.last_name AS payee_last_name,
				 SUM(es.split_amt) AS total_unsettled_amount
				FROM
				   expense_splits es
				JOIN
				   expenses e ON es.expense_id = e.expense_id
				JOIN
				   users u ON e.paid_by = u.user_id
				WHERE
				   es.user_id = $1 AND es.status = 'unsettled'
				GROUP BY
				   e.paid_by, u.first_name, u.last_name`
	rows, err := db.QueryContext(ctx, query, userIdInt)
	if err != nil {
		return nil, fmt.Errorf("error while retrivewing the unsettled ")
	}
	defer rows.Close()

	var settlements []*models.Statement
	for rows.Next() {
		var settlement models.Statement

		err := rows.Scan(&settlement.PayeeId, &settlement.PayeeFirstName, &settlement.PayeeLastName, &settlement.UnsettledAmt)
		if err != nil {
			return nil, fmt.Errorf("error while scaning the retrived rows : %v", err)
		}
		settlements = append(settlements, &settlement)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while iterating rows of unsettled amt :%v", err)
	}

	return settlements, nil

}

func SettleUpExpense(s *models.SettleAmt) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	query := `SELECT expense_id,user_id, split_amt FROM expense_splits 
	          WHERE user_id = $1 AND expense_id IN 
			  (SELECT expense_id FROM expenses WHERE paid_by =$2) AND status = 'unsettled'
			  ORDER BY expense_id`

	rows, err := tx.QueryContext(ctx, query, s.Payer_id, s.Payee_id)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to query unsettled splits: %v", err)
	}

	defer rows.Close()

	var splits []struct {
		Expense_Id int64
		User_Id    int64
		SplitAmt   float64
	}

	for rows.Next() {
		var split struct {
			Expense_Id int64
			User_Id    int64
			SplitAmt   float64
		}

		if err := rows.Scan(&split.Expense_Id, &split.User_Id, &split.SplitAmt); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to scan unsettled splits: %v", err)
		}
		splits = append(splits, split)
	}

	if err := rows.Err(); err != nil {
		tx.Rollback()
		return fmt.Errorf("error iterating unsettled splits: %w", err)
	}

	for _, split := range splits {
		if s.Amount == 0 {
			break
		}

		if s.Amount >= split.SplitAmt {
			//fully settle this split
			query := `UPDATE expense_splits SET status = 'settled' WHERE expense_id = $1 AND user_id = $2`
			_, err := tx.ExecContext(ctx, query, split.Expense_Id, split.User_Id)

			if err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to update split: %v", err)
			}

			s.Amount -= split.SplitAmt
		} else {
			//partially settle this split
			query := `UPDATE expense_splits SET split_amt = $1 WHERE expense_id = $2 AND user_id = $3`
			_, err := tx.ExecContext(ctx, query, split.SplitAmt-s.Amount, split.Expense_Id, split.User_Id)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to update split: %w", err)
			}

			s.Amount = 0
		}
	}

	if s.Amount > 0 {
		tx.Rollback()
		return fmt.Errorf("amount to settle exceeds total unsettled amount")
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil

}
