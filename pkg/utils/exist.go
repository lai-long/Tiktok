package utils

import (
	"Tiktok/biz/dao/db"
	"database/sql"
)

func IsUsernameExists(username string) (bool, error) {
	_, err := db.GetUserByUsername(username)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
