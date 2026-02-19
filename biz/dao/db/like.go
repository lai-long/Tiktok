package db

func LikeUp(id string) error {
	sql := `UPDATE videos SET like_count=like_count + 1 WHERE id = ?`
	_, err := db.Exec(sql, id)
	return err
}
func LikeDown(id string) error {
	sql := `UPDATE videos SET like_count=like_count - 1 WHERE id = ?`
	_, err := db.Exec(sql, id)
	return err
}
func CreateComment(videoId string, userId string, content string) error {
	sql := `INSERT INTO comments (video_id, user_id,content) VALUES (?, ?,?)`
	_, err := db.Exec(sql, videoId, userId, content)
	return err
}

func GetComments(videoId string) error {
	sql := `SELECT * FROM comments WHERE video_id = ?`
	_, err := db.Exec(sql, videoId)
	return err
}
