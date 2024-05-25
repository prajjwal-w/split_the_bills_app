package service

import (
	"context"
	"fmt"
	"log"
	"myJwtAuth/models"
	"time"
)

// create group
func CreateGroup(group models.Group) (int64, []*models.GroupUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var grpId int64

	query := `INSERT INTO GROUPS(group_name,description,created_by,created_at) VALUES ($1, $2, $3, $4) returning group_id`

	err := db.QueryRowContext(ctx, query,
		group.GroupName,
		group.Description,
		group.CreatedBy,
		group.Created_at,
	).Scan(&grpId)

	if err != nil {
		return 0, nil, fmt.Errorf("error while creating group: %v", err)
	}
	//Adding the admin user to the group_user table
	err = AddUsersToGroup(grpId, group.CreatedBy)
	if err != nil {
		return grpId, nil, fmt.Errorf("error while adding user to group_users: %v", err)
	}

	//Retriviewing the group users details
	var grpUsers []*models.GroupUser
	grpUsers, err = GetGroupUsers(grpId)
	if err != nil {
		return grpId, nil, fmt.Errorf("error while retriviewing group_users: %v", err)
	}

	return grpId, grpUsers, nil
}

// add users to the group
func AddUsersToGroup(group_id, user_id int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := `INSERT INTO group_users(group_id,user_id) VALUES ($1, $2)`

	_, err := db.ExecContext(ctx, query,
		group_id,
		user_id,
	)
	if err != nil {
		return fmt.Errorf("error while adding member in group: %v", err)
	}

	return nil
}

// get the users from the group
func GetGroupUsers(group_id int64) ([]*models.GroupUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := `SELECT u.user_id, u.first_name, u.last_name, u.email, u.phone FROM group_users g
	          JOIN users u ON g.user_id = u.user_id WHERE g.group_id =$1`

	rows, err := db.QueryContext(ctx, query, group_id)
	if err != nil {
		return nil, fmt.Errorf("error while retriving group user data: %v", err)
	}
	defer rows.Close()

	var groupUsers []*models.GroupUser
	for rows.Next() {
		var grpUser models.GroupUser
		err := rows.Scan(
			&grpUser.UserId,
			&grpUser.First_Name,
			&grpUser.Last_Name,
			&grpUser.Email,
			&grpUser.Phone,
		)
		if err != nil {
			log.Printf("error while scanning the group users data: %v", err)
			return nil, fmt.Errorf("error while scanning the group users data: %v", err)
		}

		groupUsers = append(groupUsers, &grpUser)

	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows of group users data: %v", err)
	}

	return groupUsers, nil
}

func GetAllGroupsbyUser(user_id int64) ([]*models.Group, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//retrival of expenses in the below query is remaining make note
	query := `SELECT
	            g.group_id, g.group_name, g.description, g.created_by , u.user_id, u.first_name, u.last_name, u.email, u.phone , e.expense_id, e.amount,
				e.description, e.paid_by, e.created_at 
			  FROM
			    groups g
			  JOIN
			    group_users gu ON g.group_id = gu.group_id
			  JOIN
			    users u ON gu.user_id = u.user_id
			  JOIN
			    expenses e ON g.group_id = e.group_id
			  WHERE
			    g.group_id IN (
					SELECT group_id FROM group_users
					WHERE user_id = $1
				)
			  ORDER BY
			    g.group_id, u.user_id;
			`
	rows, err := db.QueryContext(ctx, query, user_id)
	if err != nil {
		return nil, fmt.Errorf("error while retrivewing group details: %v", err)
	}
	defer rows.Close()

	//creating map to map the group and users of that user
	groups := make(map[int64]*models.Group)

	for rows.Next() {
		var group_id, created_by int64
		var group_name, description string
		grp_member := &models.GroupUser{}
		grp_exp := &models.AddExpense{}

		err := rows.Scan(
			&group_id,
			&group_name,
			&description,
			&created_by,
			&grp_member.UserId,
			&grp_member.First_Name,
			&grp_member.Last_Name,
			&grp_member.Email,
			&grp_member.Phone,
			&grp_exp.ExpenseId,
			&grp_exp.Amount,
			&grp_exp.Description,
			&grp_exp.PaidById,
			&grp_exp.Created_at,
		)
		if err != nil {
			return nil, fmt.Errorf("error while scanning the group data: %v", err)
		}

		if group, ok := groups[group_id]; ok {
			group.Members = append(group.Members, grp_member)
			group.Expenses = append(group.Expenses, grp_exp)
			groups[group_id] = group
		} else {
			groups[group_id] = &models.Group{
				GroupId:     group_id,
				GroupName:   group_name,
				Description: description,
				CreatedBy:   created_by,
				Members:     []*models.GroupUser{grp_member},
				Expenses:    []*models.AddExpense{grp_exp},
			}
		}

	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows of groups data: %v", err)
	}
	var grpList []*models.Group
	for _, groupDetails := range groups {
		grpList = append(grpList, groupDetails)
	}

	return grpList, nil

}
