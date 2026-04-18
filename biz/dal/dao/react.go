package dao

import (
	"Tiktok/biz/entity"
	"fmt"
	"log"

	"github.com/pkg/errors"
)

func (m *MySQLdb) VideoLikeCountUp(videoId string) error {
	sql := `UPDATE videos SET like_count=like_count + 1 WHERE id = ?`
	_, err := m.db.Exec(sql, videoId)
	return errors.Wrap(err, "dao VideoLikeCountUp")
}

func (m *MySQLdb) CommentLikeCountUp(commentId string) error {
	sql := `UPDATE comments SET like_count=like_count + 1 WHERE comment_id = ?`
	_, err := m.db.Exec(sql, commentId)
	return errors.Wrap(err, "dao CommentLikeCountUp")
}

func (m *MySQLdb) LikeCreate(userId string, targetId string, targetType string) error {
	sql := `INSERT INTO likes (target_id, user_id,target_type) VALUES (?, ?,?)`
	_, err := m.db.Exec(sql, targetId, userId, targetType)
	return errors.Wrap(err, "dao LikeCreate")
}

func (m *MySQLdb) VideoLikeCountDown(videoId string) error {
	sql := `UPDATE videos SET like_count=like_count - 1 WHERE id = ?`
	_, err := m.db.Exec(sql, videoId)
	return errors.Wrap(err, "dao VideoLikeCountDown")
}

func (m *MySQLdb) CommentLikeCountDown(commentId string) error {
	sql := `UPDATE comments SET like_count=like_count - 1 WHERE comment_id = ?`
	_, err := m.db.Exec(sql, commentId)
	return errors.Wrap(err, "dao CommentLikeCountDown")
}

func (m *MySQLdb) LikeDelete(userId, targetID, targetType string) error {
	sql := `DELETE FROM likes WHERE user_id=? AND target_id = ? AND target_type = ? LIMIT 1`
	result, err := m.db.Exec(sql, userId, targetID, targetType)
	if err != nil {
		return errors.Wrap(err, "dao LikeDelete")
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("no like found to delete")
	}
	return nil
}

func (m *MySQLdb) LikeVideoIds(userId string, pageNum int64, pageSize int64) ([]string, error) {
	sql := `SELECT target_id FROM likes WHERE  user_id = ? AND target_type = 1  ORDER BY created_at DESC LIMIT ? OFFSET ?`
	var videoId []string
	err := m.db.Select(&videoId, sql, userId, pageSize, pageNum*pageSize)
	return videoId, errors.Wrap(err, "dao Like video list")
}

func (m *MySQLdb) LikeVideos(videoId []string) (bool, []entity.VideoEntity) {
	videos := make([]entity.VideoEntity, len(videoId))
	var GetVideoErrors = 0
	for i := range videos {
		var err error
		videos[i], err = m.GetVideoByVideoId(videoId[i])
		if err != nil {
			log.Println("GetVideoByVideoId:", err)
			GetVideoErrors++
		}
	}
	if GetVideoErrors == 0 {
		return true, videos
	}
	return false, videos
}

func (m *MySQLdb) CreateComment(commentId string, videoId string, userId string, content string, targetType string) error {
	sql := `INSERT INTO comments (comment_id,target_id, user_id,content,target_type) VALUES (?, ?,?,?,?)`
	_, err := m.db.Exec(sql, commentId, videoId, userId, content, targetType)
	return err
}

func (m *MySQLdb) GetComments(videoId string, pageNum int64, pageSize int64) ([]entity.CommentEntity, error) {
	sql := `SELECT * FROM comments WHERE target_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`
	var comments []entity.CommentEntity
	err := m.db.Select(&comments, sql, videoId, pageSize, pageNum*pageSize)
	return comments, err
}

func (m *MySQLdb) CommentDelete(commentId string) error {
	sql := `UPDATE comments SET deleted_at = NOW() WHERE comment_id = ? AND deleted_at IS NULL`
	_, err := m.db.Exec(sql, commentId)
	return err
}

func (m *MySQLdb) GetCommentById(commentId string) (entity.CommentEntity, error) {
	sql := `SELECT * FROM comments WHERE comment_id = ? AND deleted_at IS NULL`
	var comment entity.CommentEntity
	err := m.db.Get(&comment, sql, commentId)
	return comment, err
}

func (m *MySQLdb) VideoCommentCountUp(videoID string) error {
	sql := `UPDATE videos SET comment_count = comment_count + 1 WHERE id = ?`
	_, err := m.db.Exec(sql, videoID)
	return err
}

func (m *MySQLdb) CommentCommentCountUp(commentID string) error {
	sql := `UPDATE comments SET comment_count = comment_count + 1 WHERE comment_id = ?`
	_, err := m.db.Exec(sql, commentID)
	return err
}

func (m *MySQLdb) VideoCommentCountDown(videoID string) error {
	sql := `UPDATE videos SET comment_count = comment_count - 1 WHERE id = ?`
	_, err := m.db.Exec(sql, videoID)
	return err
}

func (m *MySQLdb) CommentCommentCountDown(commentID string) error {
	sql := `UPDATE comments SET comment_count = comment_count - 1 WHERE comment_id = ?`
	_, err := m.db.Exec(sql, commentID)
	return err
}
