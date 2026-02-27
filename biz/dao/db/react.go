package db

import (
	"Tiktok/biz/model/entity"
	"fmt"
)

func LikeCountUp(video_id string) error {
	sql := `UPDATE videos SET like_count=like_count + 1 WHERE id = ?`
	_, err := db.Exec(sql, video_id)
	return err
}

func LikeCreate(user_id string, video_id string) error {
	sql := `INSERT INTO likes (video_id, user_id) VALUES (?, ?)`
	_, err := db.Exec(sql, video_id, user_id)
	return err
}
func LikeCountDown(video_id string) error {
	sql := `UPDATE videos SET like_count=like_count - 1 WHERE id = ?`
	_, err := db.Exec(sql, video_id)
	return err
}
func LikeDelete(user_id string, video_id string) error {
	sql := `DELETE FROM likes WHERE video_id = ? AND user_id = ?`
	result, err := db.Exec(sql, video_id, user_id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("no like found to delete")
	}
	return nil
}
func LikeVideoIds(user_id string, pageNum int, pageSize int) (error, []string) {
	sql := `SELECT video_id FROM likes WHERE user_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`
	var video_id []string
	offset := pageNum * pageSize
	err := db.Select(&video_id, sql, user_id, pageSize, offset)
	return err, video_id
}
func LikeVideos(videoId []string) (bool, []entity.VideoEntity) {
	videos := make([]entity.VideoEntity, len(videoId))
	var GetVideoErrors = 0
	for i, _ := range videos {
		var err error
		videos[i], err = GetVideoByVideoId(videoId[i])
		if err != nil {
			GetVideoErrors = GetVideoErrors + 1
		}
	}
	if GetVideoErrors == 0 {
		return true, videos
	} else {
		return false, videos
	}
}
func CreateComment(commentId string, videoId string, userId string, content string) error {
	sql := `INSERT INTO comments (comment_id,video_id, user_id,content) VALUES (?, ?,?,?)`
	_, err := db.Exec(sql, commentId, videoId, userId, content)
	return err
}

func GetComments(videoId string, pageNum int, pageSize int) (error, []entity.CommentEntity) {
	sql := `SELECT * FROM comments WHERE video_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`
	offset := pageNum * pageSize
	var comments []entity.CommentEntity
	err := db.Select(&comments, sql, videoId, pageSize, offset)
	return err, comments
}
func CommentDelete(videoId string, commentId string) error {
	sql := `DELETE FROM comments WHERE video_id = ? AND comment_id = ?`
	_, err := db.Exec(sql, videoId, commentId)
	return err
}
func GetCommentById(commentId string) (entity.CommentEntity, error) {
	sql := `SELECT * FROM comments WHERE comment_id = ?`
	var comment entity.CommentEntity
	err := db.Get(&comment, sql, commentId)
	return comment, err
}
func CommentCountUp(videoId string) error {
	sql := `UPDATE videos SET comment_count = comment_count + 1 WHERE id = ?`
	_, err := db.Exec(sql, videoId)
	return err
}
func CommentCountDown(videoId string) error {
	sql := `UPDATE videos SET comment_count = comment_count - 1 WHERE id = ?`
	_, err := db.Exec(sql, videoId)
	return err
}
