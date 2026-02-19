package db

import (
	"Tiktok/biz/model/entity"
)

func CreateUser(user entity.UserEntity) error {
	sql := `INSERT INTO users (username,  password, id) VALUES (?, ? ,?)`
	_, err := db.Exec(sql, user.Username, user.Password, user.Id)
	return err
}

func GetUserByUsername(username string) (entity.UserEntity, error) {
	var user entity.UserEntity
	sql := `SELECT * FROM users WHERE username = ?`
	err := db.Get(&user, sql, username)
	return user, err
}

func GetUserByUserId(userId string) (entity.UserEntity, error) {
	var user entity.UserEntity
	sql := `SELECT * FROM users WHERE id = ?`
	err := db.Get(&user, sql, userId)
	return user, err
}

func UpdateUserAvatar(url string, userId interface{}) error {
	sql := `UPDATE users SET avatar_url=? WHERE id=?`
	_, err := db.Exec(sql, url, userId)
	return err
}
