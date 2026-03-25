package db

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
	sql := `SELECT * FROM users WHERE id IN(SELECT friend_id FROM friends WHERE user_id = ?) LIMIT ? OFFSET ?`
	var users []entity.UserEntity
	err := m.db.Select(&users, sql, userId, pageSize, offset)
	if err != nil {
		log.Println("Get FriendList err", err)
		return nil, false
	}
	return users, true
}

func (m *MySQLdb) CreateFriend(userId string, toUserId string) bool {
	sql := `INSERT INTO friends (user_id, friend_id) VALUES (?,?)`
	_, err := m.db.Exec(sql, userId, toUserId)
	if err != nil {
		log.Println("db CreateFriend err", err)
		return false
	}
	return true
}
func (m *MySQLdb) DeleteFriend(userId string, toUserId string) bool {
	sql := `DELETE FROM friends WHERE user_id = ? AND friend_id = ?`
	_, err := m.db.Exec(sql, userId, toUserId)
	if err != nil {
		log.Println("db DeleteFriend err", err)
	}
	return true
}
