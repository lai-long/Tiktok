package dao

import (
	"Tiktok/biz/model/entity"
	"log"
)

func (m *MySQLdb) CreateFollowing(userId string, toUserId string) error {
	sql := `INSERT INTO relations (user_id, following_id) VALUES (?,?)`
	_, err := m.db.Exec(sql, userId, toUserId)
	return err
}
func (m *MySQLdb) CreateFollower(userId string, toUserId string) error {
	sql := `INSERT INTO relations (user_id, follower_id) VALUES (?,?)`
	_, err := m.db.Exec(sql, toUserId, userId)
	return err
}
func (m *MySQLdb) DeleteFollowing(userId string, toUserId string) error {
	sql := `DELETE FROM relations WHERE user_id = ? AND follower_id = ?`
	_, err := m.db.Exec(sql, userId, toUserId)
	return err
}
func (m *MySQLdb) DeleteFollower(userId string, toUserId string) error {
	sql := `DELETE FROM relations WHERE user_id = ? AND follower_id = ?`
	_, err := m.db.Exec(sql, toUserId, userId)
	return err
}
func (m *MySQLdb) FollowingIdList(userId string, pageNum int, pageSize int) ([]string, error) {
	sql := `SELECT following_id FROM relations WHERE user_id = ? AND following_id IS NOT NULL LIMIT ? OFFSET ?`
	offset := pageNum * pageSize
	var followingIds []string
	err := m.db.Select(&followingIds, sql, userId, pageSize, offset)
	return followingIds, err
}
func (m *MySQLdb) FollowerIdList(userId string, pageNum int, pageSize int) ([]string, error) {
	sql := `SELECT follower_id FROM relations WHERE user_id = ? AND follower_id IS NOT NULL LIMIT ? OFFSET ?`
	offset := pageNum * pageSize
	var followerIds []string
	err := m.db.Select(&followerIds, sql, userId, pageSize, offset)
	return followerIds, err
}
func (m *MySQLdb) FriendList(userId string, pageNum int, pageSize int) ([]entity.UserEntity, bool) {
	offset := pageNum * pageSize
	sql := `SELECT * FROM users WHERE id IN(SELECT following_id FROM relations WHERE user_id = ?) 
                      AND id IN (SELECT user_id from relations where following_id=?) LIMIT ? OFFSET ?`
	var users []entity.UserEntity
	err := m.db.Select(&users, sql, userId, userId, pageSize, offset)
	if err != nil {
		log.Println("Get FriendList err", err)
		return nil, false
	}
	return users, true
}
