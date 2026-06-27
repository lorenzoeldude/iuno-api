package services

import (
	"context"

	"iuno-api/db"
)

func DeleteUser(userID int) error {

	ctx := context.Background()

	_, err := db.Pool.Exec(
		ctx,
		`DELETE FROM users WHERE id = $1`,
		userID,
	)

	return err
}