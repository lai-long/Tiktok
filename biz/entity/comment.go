package entity

import (
	"Tiktok/biz/model/react"
	"database/sql"
)

// CommentEntity is used to react with database
type CommentEntity struct {
	UserID       string       `db:"user_id"`
	TargetID     string       `db:"target_id"`
	CommentID    string       `db:"comment_id"`
	Content      string       `db:"content"`
	LikeCount    int64        `db:"like_count"`
	CommentCount int64        `db:"comment_count"`
	CreatedAt    string       `db:"created_at"`
	UpdatedAt    string       `db:"updated_at"`
	DeletedAt    sql.NullTime `db:"deleted_at"`
	TargetType   string       `db:"target_type"`
}

// ToCommentInfo is used to exchange CommentEntity to commentInfo
func (c *CommentEntity) ToCommentInfo() *react.CommentInfo {
	return &react.CommentInfo{
		UserId:       c.UserID,
		TargetId:     c.TargetID,
		CommentId:    c.CommentID,
		Content:      c.Content,
		LikeCount:    c.LikeCount,
		CreatedAt:    c.CreatedAt,
		UpdatedAt:    c.UpdatedAt,
		TargetType:   c.TargetType,
		CommentCount: c.CommentCount,
	}
}
