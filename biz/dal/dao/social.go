package dao

import (
	"Tiktok/biz/entity"
	"log"
)

func (m *MySQLdb) CreateFollowing(userId string, toUserId string) error {
	sql := `INSERT INTO following (user_id, following_id) VALUES (?,?)`
	_, err := m.db.Exec(sql, userId, toUserId)
	return err
}
func (m *MySQLdb) DeleteFollowing(userId string, toUserId string) error {
	sql := `UPDATE following SET deleted_at = NOW() WHERE user_id = ? AND following_id = ? AND deleted_at IS NULL`
	_, err := m.db.Exec(sql, userId, toUserId)
	return err
}
func (m *MySQLdb) FollowingList(userId string, pageNum int64, pageSize int64) ([]entity.UserEntity, error) {
	sql := `SELECT * FROM users WHERE id IN (SELECT following_id FROM following WHERE user_id = ? AND deleted_at is null) LIMIT ? OFFSET ?`
	var users []entity.UserEntity
	offset := pageNum * pageSize
	err := m.db.Select(&users, sql, userId, pageSize, offset)
	return users, err
}
func (m *MySQLdb) FollowerList(userId string, pageNum int64, pageSize int64) ([]entity.UserEntity, error) {
	sql := `SELECT * FROM users WHERE id IN (SELECT user_id FROM following WHERE following_id = ? and users.deleted_at is null) LIMIT ? OFFSET ?`
	offset := pageNum * pageSize
	var users []entity.UserEntity
	err := m.db.Select(&users, sql, userId, pageSize, offset)
	return users, err
}
func (m *MySQLdb) FriendList(userId string, pageNum int64, pageSize int64) ([]entity.UserEntity, bool) {
	offset := pageNum * pageSize
	sql := `SELECT * FROM users WHERE id IN(SELECT following_id FROM following WHERE user_id = ? and deleted_at is null)
                      AND id IN (SELECT user_id from following where following_id=? and deleted_at is null) LIMIT ? OFFSET ?`
	var users []entity.UserEntity
	err := m.db.Select(&users, sql, userId, userId, pageSize, offset)
	if err != nil {
		log.Println("Get FriendList err", err)
		return nil, false
	}
	return users, true
}
