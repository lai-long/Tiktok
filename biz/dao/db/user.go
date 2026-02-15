package db

import (
	"Tiktok/biz/model/entity"
)

func CreateUser(user entity.UserEntity) error {
	query := `INSERT INTO users (username,  password, id) VALUES (?, ? ,?)`
	_, err := db.Exec(query, user.Username, user.Password, user.Id)
	return err
}
func GetUserByUsername(username string) (entity.UserEntity, error) {
	var user entity.UserEntity
	query := `SELECT * FROM users WHERE username = ?`
	err := db.Get(&user, query, username)
	return user, err
}
func GetUserByUserId(userId string) (entity.UserEntity, error) {
	var user entity.UserEntity
	query := `SELECT * FROM users WHERE id = ?`
	err := db.Get(&user, query, userId)
	return user, err
}
