package dao

import (
	entity2 "Tiktok/pkg/entity"
	"context"
	"fmt"

	"github.com/pkg/errors"
)

func (m *MySQLdb) VideoLikeCountUp(ctx context.Context, videoId string) error {
	sql := `UPDATE videos SET like_count=like_count + 1 WHERE id = ?`
	_, err := m.db.ExecContext(ctx, sql, videoId)
	return errors.Wrap(err, "dao VideoLikeCountUp")
}

func (m *MySQLdb) CommentLikeCountUp(ctx context.Context, commentId string) error {
	sql := `UPDATE comments SET like_count=like_count + 1 WHERE comment_id = ?`
	_, err := m.db.ExecContext(ctx, sql, commentId)
	return errors.Wrap(err, "dao CommentLikeCountUp")
}

func (m *MySQLdb) LikeCreate(ctx context.Context, userId string, targetId string, targetType string) error {
	sql := `INSERT INTO likes (target_id, user_id,target_type) VALUES (?, ?,?)`
	_, err := m.db.ExecContext(ctx, sql, targetId, userId, targetType)
	return errors.Wrap(err, "dao LikeCreate")
}

func (m *MySQLdb) VideoLikeCountDown(ctx context.Context, videoId string) error {
	sql := `UPDATE videos SET like_count=like_count - 1 WHERE id = ?`
	_, err := m.db.ExecContext(ctx, sql, videoId)
	return errors.Wrap(err, "dao VideoLikeCountDown")
}

func (m *MySQLdb) CommentLikeCountDown(ctx context.Context, commentId string) error {
	sql := `UPDATE comments SET like_count=like_count - 1 WHERE comment_id = ?`
	_, err := m.db.ExecContext(ctx, sql, commentId)
	return errors.Wrap(err, "dao CommentLikeCountDown")
}

func (m *MySQLdb) LikeDelete(ctx context.Context, userId, targetID, targetType string) error {
	sql := `DELETE FROM likes WHERE user_id=? AND target_id = ? AND target_type = ? LIMIT 1`
	result, err := m.db.ExecContext(ctx, sql, userId, targetID, targetType)
	if err != nil {
		return errors.Wrap(err, "dao LikeDelete")
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("no like found to delete")
	}
	return nil
}

func (m *MySQLdb) LikeVideoIds(ctx context.Context, userId string, pageNum int64, pageSize int64) ([]string, error) {
	sql := `SELECT target_id FROM likes WHERE  user_id = ? AND target_type = 1  ORDER BY created_at DESC LIMIT ? OFFSET ?`
	var videoId []string
	err := m.db.SelectContext(ctx, &videoId, sql, userId, pageSize, pageNum*pageSize)
	return videoId, errors.Wrap(err, "dao Like video list")
}

func (m *MySQLdb) CreateComment(ctx context.Context, commentId string, videoId string, userId string, content string, targetType string) error {
	sql := `INSERT INTO comments (comment_id,target_id, user_id,content,target_type) VALUES (?, ?,?,?,?)`
	_, err := m.db.ExecContext(ctx, sql, commentId, videoId, userId, content, targetType)
	return err
}

func (m *MySQLdb) GetComments(ctx context.Context, targetId string, pageNum int64, pageSize int64) ([]entity2.CommentEntity, error) {
	sql := `SELECT * FROM comments WHERE target_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`
	var comments []entity2.CommentEntity
	err := m.db.SelectContext(ctx, &comments, sql, targetId, pageSize, pageNum*pageSize)
	return comments, err
}

func (m *MySQLdb) CommentDelete(ctx context.Context, commentId string) error {
	sql := `UPDATE comments SET deleted_at = NOW() WHERE comment_id = ? AND deleted_at IS NULL`
	_, err := m.db.ExecContext(ctx, sql, commentId)
	return err
}

func (m *MySQLdb) GetCommentById(ctx context.Context, commentId string) (entity2.CommentEntity, error) {
	sql := `SELECT * FROM comments WHERE comment_id = ? AND deleted_at IS NULL`
	var comment entity2.CommentEntity
	err := m.db.GetContext(ctx, &comment, sql, commentId)
	return comment, err
}

func (m *MySQLdb) VideoCommentCountUp(ctx context.Context, videoID string) error {
	sql := `UPDATE videos SET comment_count = comment_count + 1 WHERE id = ?`
	_, err := m.db.ExecContext(ctx, sql, videoID)
	return err
}

func (m *MySQLdb) CommentCommentCountUp(ctx context.Context, commentID string) error {
	sql := `UPDATE comments SET comment_count = comment_count + 1 WHERE comment_id = ?`
	_, err := m.db.ExecContext(ctx, sql, commentID)
	return err
}

func (m *MySQLdb) VideoCommentCountDown(ctx context.Context, videoID string) error {
	sql := `UPDATE videos SET comment_count = comment_count - 1 WHERE id = ?`
	_, err := m.db.ExecContext(ctx, sql, videoID)
	return err
}

func (m *MySQLdb) CommentCommentCountDown(ctx context.Context, commentID string) error {
	sql := `UPDATE comments SET comment_count = comment_count - 1 WHERE comment_id = ?`
	_, err := m.db.ExecContext(ctx, sql, commentID)
	return err
}
