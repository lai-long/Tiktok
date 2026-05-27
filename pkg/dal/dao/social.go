package dao

import (
	"Tiktok/pkg/entity"
	"context"
)

func (m *MySQLdb) CreateFollowing(ctx context.Context, userID string, toUserID string) error {
	sql := `INSERT INTO following (user_id, following_id) VALUES (?,?)`
	_, err := m.db.ExecContext(ctx, sql, userID, toUserID)
	return err
}
func (m *MySQLdb) DeleteFollowing(ctx context.Context, userID string, toUserID string) error {
	sql := `UPDATE following SET deleted_at = NOW() WHERE user_id = ? AND following_id = ? AND deleted_at IS NULL`
	_, err := m.db.ExecContext(ctx, sql, userID, toUserID)
	return err
}
func (m *MySQLdb) FollowingList(ctx context.Context, userID string, pageNum int64, pageSize int64) ([]entity.UserEntity, error) {
	sql := `SELECT * FROM users WHERE id IN (SELECT following_id FROM following WHERE user_id = ? AND deleted_at
	is null) LIMIT ? OFFSET ?`
	var users []entity.UserEntity
	offset := pageNum * pageSize
	err := m.db.SelectContext(ctx, &users, sql, userID, pageSize, offset)
	return users, err
}
func (m *MySQLdb) FollowerList(ctx context.Context, userID string, pageNum int64, pageSize int64) ([]entity.UserEntity, error) {
	sql := `SELECT * FROM users WHERE id IN (SELECT user_id FROM following WHERE following_id = ? and
        	users.deleted_at is null) LIMIT ? OFFSET ?`
	offset := pageNum * pageSize
	var users []entity.UserEntity
	err := m.db.SelectContext(ctx, &users, sql, userID, pageSize, offset)
	return users, err
}
func (m *MySQLdb) FriendList(ctx context.Context, userID string, pageNum int64, pageSize int64) ([]entity.UserEntity, bool) {
	offset := pageNum * pageSize
	sql := `SELECT * FROM users WHERE id IN(SELECT following_id FROM following WHERE user_id = ? and deleted_at is null)
            AND id IN (SELECT user_id from following where following_id=? and deleted_at is null) LIMIT ? OFFSET ?`
	var users []entity.UserEntity
	err := m.db.SelectContext(ctx, &users, sql, userID, userID, pageSize, offset)
	if err != nil {
		return nil, false
	}
	return users, true
}
