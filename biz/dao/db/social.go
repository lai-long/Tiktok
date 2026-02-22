package db

func CreateFollowing(userId string, toUserId string) error {
	sql := `INSERT INTO relations (user_id, following_id) VALUES (?,?)`
	_, err := db.Exec(sql, userId, toUserId)
	return err
}
func CreateFollower(userId string, toUserId string) error {
	sql := `INSERT INTO relations (user_id, follower_id) VALUES (?,?)`
	_, err := db.Exec(sql, toUserId, userId)
	return err
}
func DeleteFollowing(userId string, toUserId string) error {
	sql := `DELETE FROM relations WHERE user_id = ? AND follower_id = ?`
	_, err := db.Exec(sql, userId, toUserId)
	return err
}
func DeleteFollower(userId string, toUserId string) error {
	sql := `DELETE FROM relations WHERE user_id = ? AND follower_id = ?`
	_, err := db.Exec(sql, toUserId, userId)
	return err
}
func FollowingIdList(userId string, pageNum int, pageSize int) ([]string, error) {
	sql := `SELECT following_id FROM relations WHERE user_id = ? LIMIT ? OFFSET ?`
	offset := (pageNum - 1) * pageSize
	var followingIds []string
	err := db.Select(&followingIds, sql, userId, pageSize, offset)
	return followingIds, err
}
func FollowerIdList(userId string, pageNum int, pageSize int) ([]string, error) {
	sql := `SELECT follower_id FROM relations WHERE user_id = ? LIMIT ? OFFSET ?`
	offset := (pageNum - 1) * pageSize
	var followerIds []string
	err := db.Select(&followerIds, sql, userId, pageSize, offset)
	return followerIds, err
}
