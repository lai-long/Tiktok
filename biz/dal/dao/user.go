package dao

import (
	"Tiktok/biz/entity"
	"log"
)

func (m *MySQLdb) CreateUser(user entity.UserEntity) error {
	sql := `INSERT INTO users (username,  password, id) VALUES (?, ? ,?)`
	_, err := m.db.Exec(sql, user.Username, user.Password, user.Id)
	return err
}

func (m *MySQLdb) GetUserByUsername(username string) (entity.UserEntity, error) {
	var user entity.UserEntity
	sql := `SELECT * FROM users WHERE username = ?`
	err := m.db.Get(&user, sql, username)
	return user, err
}

func (m *MySQLdb) GetUserByUserId(userId string) (entity.UserEntity, error) {
	var user entity.UserEntity
	sql := `SELECT * FROM users WHERE id = ?`
	err := m.db.Get(&user, sql, userId)
	if err != nil {
		log.Println(err)
	}
	return user, err
}

func (m *MySQLdb) UpdateUserAvatar(url string, userId interface{}) error {
	sql := `UPDATE users SET avatar_url=? WHERE id=?`
	_, err := m.db.Exec(sql, url, userId)
	return err
}
