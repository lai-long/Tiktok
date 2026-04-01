package dao

import (
	"Tiktok/biz/model/entity"
	"fmt"
	"log"
)

func (m *MySQLdb) VideoLikeCountUp(videoId string) error {
	sql := `UPDATE videos SET like_count=like_count + 1 WHERE id = ?`
	_, err := m.db.Exec(sql, videoId)
	return err
}
func (m *MySQLdb) CommentLikeCountUp(commentId string) error {
	sql := `UPDATE comments SET like_count=like_count + 1 WHERE comment_id = ?`
	_, err := m.db.Exec(sql, commentId)
	return err
}
func (m *MySQLdb) LikeCreate(userId string, targetId string, targetType string) error {
	sql := `INSERT INTO likes (target_id, user_id,target_type) VALUES (?, ?)`
	_, err := m.db.Exec(sql, targetId, userId, targetType)
	return err
}
func (m *MySQLdb) VideoLikeCountDown(videoId string) error {
	sql := `UPDATE videos SET like_count=like_count - 1 WHERE id = ?`
	_, err := m.db.Exec(sql, videoId)
	return err
}
func (m *MySQLdb) CommentLikeCountDown(commentId string) error {
	sql := `UPDATE comments SET like_count=like_count - 1 WHERE comment_id = ?`
	_, err := m.db.Exec(sql, commentId)
	return err
}
func (m *MySQLdb) LikeDelete(userId, targetID, targetType string) error {
	sql := `DELETE FROM likes WHERE user_id=? AND target_id = ? AND target_type = ? LIMIT 1`
	result, err := m.db.Exec(sql, userId, targetID, targetType)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("no like found to delete")
	}
	return nil
}

func (m *MySQLdb) LikeVideoIds(userId string, pageNum int, pageSize int) (error, []string) {
	sql := `SELECT target_id FROM likes WHERE  user_id = ? AND target_type = 1  ORDER BY created_at DESC LIMIT ? OFFSET ?`
	var videoId []string
	offset := pageNum * pageSize
	err := m.db.Select(&videoId, sql, userId, pageSize, offset)
	return err, videoId
}

func (m *MySQLdb) LikeVideos(videoId []string) (bool, []entity.VideoEntity) {
	videos := make([]entity.VideoEntity, len(videoId))
	var GetVideoErrors = 0
	for i, _ := range videos {
		var err error
		videos[i], err = m.GetVideoByVideoId(videoId[i])
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

func (m *MySQLdb) CreateComment(commentId string, videoId string, userId string, content string) error {
	sql := `INSERT INTO comments (comment_id,video_id, user_id,content) VALUES (?, ?,?,?)`
	_, err := m.db.Exec(sql, commentId, videoId, userId, content)
	return err
}

func (m *MySQLdb) GetComments(videoId string, pageNum int, pageSize int) (error, []entity.CommentEntity) {
	sql := `SELECT * FROM comments WHERE video_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`
	offset := pageNum * pageSize
	var comments []entity.CommentEntity
	err := m.db.Select(&comments, sql, videoId, pageSize, offset)
	return err, comments
}

func (m *MySQLdb) CommentDelete(commentId string) error {
	sql := `DELETE FROM comments WHERE  comment_id = ?`
	_, err := m.db.Exec(sql, commentId)
	return err
}

func (m *MySQLdb) GetCommentById(commentId string) (entity.CommentEntity, error) {
	sql := `SELECT * FROM comments WHERE comment_id = ?`
	var comment entity.CommentEntity
	err := m.db.Get(&comment, sql, commentId)
	return comment, err
}

func (m *MySQLdb) CommentCountUp(videoId string) error {
	sql := `UPDATE videos SET comment_count = comment_count + 1 WHERE id = ?`
	_, err := m.db.Exec(sql, videoId)
	return err
}

func (m *MySQLdb) CommentCountDown(videoId string) error {
	sql := `UPDATE videos SET comment_count = comment_count - 1 WHERE id = ?`
	_, err := m.db.Exec(sql, videoId)
	return err
}
