package db

import (
	"Tiktok/biz/model/entity"
	"fmt"
	"log"
)

func VideoLikeCountUp(videoId string) error {
	sql := `UPDATE videos SET like_count=like_count + 1 WHERE id = ?`
	_, err := db.Exec(sql, videoId)
	return err
}
func CommentLikeCountUp(commentId string) error {
	sql := `UPDATE comments SET like_count=like_count + 1 WHERE comment_id = ?`
	_, err := db.Exec(sql, commentId)
	return err
}
func VideoLikeCreate(userId string, videoId string) error {
	sql := `INSERT INTO likes (to_video_id, user_id) VALUES (?, ?)`
	_, err := db.Exec(sql, videoId, userId)
	return err
}
func CommentLikeCreate(userId string, commentId string) error {
	sql := `INSERT INTO likes (to_comment_id, user_id) VALUES (?, ?)`
	_, err := db.Exec(sql, commentId, userId)
	return err
}
func VideoLikeCountDown(videoId string) error {
	sql := `UPDATE videos SET like_count=like_count - 1 WHERE id = ?`
	_, err := db.Exec(sql, videoId)
	return err
}
func CommentLikeCountDown(commentId string) error {
	sql := `UPDATE comments SET like_count=like_count - 1 WHERE comment_id = ?`
	_, err := db.Exec(sql, commentId)
	return err
}
func VideoLikeDelete(userId string, videoId string) error {
	sql := `DELETE FROM likes WHERE to_video_id = ? AND user_id = ?`
	result, err := db.Exec(sql, videoId, userId)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("no like found to delete")
	}
	return nil
}
func CommentLikeDelete(userId string, commentId string) error {
	sql := `DELETE FROM likes WHERE to_comment_id = ? AND user_id = ?`
	result, err := db.Exec(sql, commentId, userId)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("no like found to delete")
	}
	return nil
}
func LikeVideoIds(userId string, pageNum int, pageSize int) (error, []string) {
	sql := `SELECT to_video_id FROM likes WHERE to_video_id IS NOT NULL AND user_id = ?  ORDER BY created_at DESC LIMIT ? OFFSET ?`
	var videoId []string
	offset := pageNum * pageSize
	err := db.Select(&videoId, sql, userId, pageSize, offset)
	return err, videoId
}
func LikeVideos(videoId []string) (bool, []entity.VideoEntity) {
	videos := make([]entity.VideoEntity, len(videoId))
	var GetVideoErrors = 0
	for i, _ := range videos {
		var err error
		videos[i], err = GetVideoByVideoId(videoId[i])
		if err != nil {
			log.Println("GetVideoByVideoId:", err)
			GetVideoErrors = GetVideoErrors + 1
		}
	}
	if GetVideoErrors == 0 {
		return true, videos
	}
	return false, videos
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
