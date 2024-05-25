package service

import (
	"context"
	"fmt"
	"log"
	"myJwtAuth/models"
	"strconv"

	"time"
)

func AddFriend(user_id int64, frdsUserId int64) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var frdship_id int64
	query := `INSERT INTO friendships(user1_id, user2_id) VALUES($1,$2) RETURNING friendship_id`

	err := db.QueryRowContext(ctx, query,
		user_id,
		frdsUserId).Scan(&frdship_id)

	if err != nil {
		return 0, fmt.Errorf("error inserting user: %v", err)
	}

	return frdship_id, nil

}

func GetAllFriends(user_id string) ([]*models.Friend, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	userId, err := strconv.ParseInt(user_id, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %v", err)
	}
	query := `SELECT f.friendship_id,u.user_id, u.first_name, u.last_name, u.email, u.phone FROM friendships f
	          JOIN users u ON f.user1_id = u.user_id OR f.user2_id = u.user_id
			  WHERE (f.user1_id = $1 OR f.user2_id = $1)
			  AND u.user_id != $1`
	rows, err := db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, fmt.Errorf("error while retriving the user's friends data: %v", err)
	}
	defer rows.Close()

	var friends []*models.Friend
	for rows.Next() {
		var frd models.Friend
		err := rows.Scan(
			&frd.FriendShipID,
			&frd.UserId,
			&frd.FirstName,
			&frd.LastName,
			&frd.Email,
			&frd.Phone,
		)

		if err != nil {
			log.Printf("error while scanning the friends: %v", err)
			return nil, fmt.Errorf("error while scanning the friends: %v", err)
		}
		friends = append(friends, &frd)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return friends, nil
}
