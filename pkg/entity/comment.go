package entity

import (
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
