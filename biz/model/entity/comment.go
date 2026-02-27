package entity

import "database/sql"

type CommentEntity struct {
	UserId    string       `db:"user_id"`
	VideoId   string       `db:"video_id"`
	CommentId string       `db:"comment_id"`
	Content   string       `db:"content"`
	LikeCount int64        `db:"like_count"`
	CreatedAt string       `db:"created_at"`
	UpdatedAt string       `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}
