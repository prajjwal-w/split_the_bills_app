package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"myJwtAuth/database"
	"myJwtAuth/models"
	"strconv"
	"time"
)

var db = database.DatabaseConn()

func GetUserByEmail(userEmail string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	query := `SELECT * FROM USERS WHERE email = $1`
	var user models.User
	err := db.QueryRowContext(ctx, query, userEmail).Scan(
		&user.User_Id,
		&user.First_Name,
		&user.Last_Name,
		&user.Password,
		&user.Email,
		&user.Phone,
		&user.Token,
		&user.User_Type,
		&user.Refresh_Token,
		&user.Created_at,
		&user.Updated_at,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("error retrieving user details: %v", err)
	}
	return &user, nil
}

func GetUserByPhone(phone string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	query := `SELECT * FROM USERS WHERE phone = $1`
	var user models.User
	err := db.QueryRowContext(ctx, query, phone).Scan(
		&user.User_Id,
		&user.First_Name,
		&user.Last_Name,
		&user.Password,
		&user.Email,
		&user.Phone,
		&user.Token,
		&user.User_Type,
		&user.Refresh_Token,
		&user.Created_at,
		&user.Updated_at,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("error retrieving user details: %v", err)
	}
	return &user, nil
}

func InsertUser(user models.User) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var userId int64
	query := `INSERT INTO USERS (first_name, last_name, password, phone, email, 
		      token, user_type, refresh_token, created_at, updated_at) values ($1, $2, $3, $4, $5, $6,$7,$8,$9,$10) 
			  returning user_id`

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("error preparing statement: %v", err)
	}
	defer stmt.Close()

	err = db.QueryRowContext(ctx, query,
		user.First_Name,
		user.Last_Name,
		user.Password,
		user.Phone,
		user.Email,
		user.Token,
		user.User_Type,
		user.Refresh_Token,
		user.Created_at,
		user.Updated_at,
	).Scan(&userId)

	if err != nil {
		return 0, fmt.Errorf("error inserting user: %v", err)
	}

	return userId, nil
}

func UpdateTokenInDB(token string, refreshToken string, userId int64) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	time := time.Now()
	query := "UPDATE USERS SET token = $1, refresh_token = $2, updated_at = $3 WHERE user_id = $4"

	_, err := db.ExecContext(ctx, query, token, refreshToken, time, userId)
	if err != nil {
		return false, fmt.Errorf("error while updating token in db")
	}

	return true, nil

}

func GetAllUsers() ([]*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := "SELECT * FROM USERS"

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.User_Id,
			&user.First_Name,
			&user.Last_Name,
			&user.Password,
			&user.Phone,
			&user.Email,
			&user.Token,
			&user.User_Type,
			&user.Refresh_Token,
			&user.Created_at,
			&user.Updated_at,
		)
		if err != nil {
			log.Println("Error while scanning the users data")
			return nil, fmt.Errorf("error scanning user data: %v", err)
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}
	return users, nil
}

func GetUserById(userId string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	userIdInt, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %v", err)
	}
	query := `SELECT * FROM USERS WHERE user_id = $1`
	var user models.User
	err = db.QueryRowContext(ctx, query, userIdInt).Scan(
		&user.User_Id,
		&user.First_Name,
		&user.Last_Name,
		&user.Password,
		&user.Phone,
		&user.Email,
		&user.Token,
		&user.User_Type,
		&user.Refresh_Token,
		&user.Created_at,
		&user.Updated_at,
	)

	if err != nil {
		return nil, fmt.Errorf("error retrieving user data: %v", err)
	}

	return &user, nil
}

func DeleteUserById(userId string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	userIdInt, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		return false, fmt.Errorf("invalid user ID format: %v", err)
	}

	query := "DELETE FROM USERS WHERE user_id = $1"

	_, err = db.ExecContext(ctx, query, userIdInt)
	if err != nil {
		return false, fmt.Errorf("error deleting user: %v", err)
	}

	return true, nil
}

func ResetPassword(userId int64, newPass string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := "UPDATE USERS SET password=$1 WHERE user_id=$2"
	_, err := db.ExecContext(ctx, query, newPass, userId)
	if err != nil {
		return false, err
	}

	return true, nil

}

func GetUserByEmailLimitedData(email string) (int64, error) {
	user, err := GetUserByEmail(email)
	if err != nil {
		return 0, err
	}

	return user.User_Id, nil
}

func GetUserByPhoneLimitedData(phone string) (int64, error) {
	user, err := GetUserByPhone(phone)
	if err != nil {
		return 0, err
	}

	return user.User_Id, nil
}
