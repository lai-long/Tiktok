package db

import "Tiktok/biz/model/entity"

func CreateComment(commentId string, videoId string, userId string, content string) error {
	sql := `INSERT INTO comments (comment_id,video_id, user_id,content) VALUES (?, ?,?,?)`
	_, err := db.Exec(sql, commentId, videoId, userId, content)
	return err
}

func GetComments(videoId string, pageNum int, pageSize int) (error, []entity.CommentEntity) {
	sql := `SELECT * FROM comments WHERE video_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`
	offset := (pageNum - 1) * pageSize
	var comments []entity.CommentEntity
	err := db.Select(&comments, sql, videoId, pageSize, offset)
	return err, comments
}
func CommentDelete(videoId string, commentId string) error {
	sql := `DELETE FROM comments WHERE video_id = ? AND comment_id = ?`
	_, err := db.Exec(sql, videoId, commentId)
	return err
}
