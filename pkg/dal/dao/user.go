package dao

import (
	"Tiktok/pkg/entity"
	"context"
)

func (m *MySQLdb) CreateUser(ctx context.Context, user entity.UserEntity) error {
	sql := `INSERT INTO users (username,  password, id) VALUES (?, ? ,?)`
	_, err := m.db.ExecContext(ctx, sql, user.Username, user.Password, user.ID)
	return err
}

func (m *MySQLdb) GetUserByUsername(ctx context.Context, username string) (entity.UserEntity, error) {
	var user entity.UserEntity
	sql := `SELECT * FROM users WHERE username = ?`
	err := m.db.GetContext(ctx, &user, sql, username)
	return user, err
}

func (m *MySQLdb) GetUserByUserId(ctx context.Context, userID string) (entity.UserEntity, error) {
	var user entity.UserEntity
	sql := `SELECT * FROM users WHERE id = ?`
	err := m.db.GetContext(ctx, &user, sql, userID)
	return user, err
}

func (m *MySQLdb) UpdateUserAvatar(ctx context.Context, url string, userID interface{}) error {
	sql := `UPDATE users SET avatar_url=? WHERE id=?`
	_, err := m.db.ExecContext(ctx, sql, url, userID)
	return err
}
