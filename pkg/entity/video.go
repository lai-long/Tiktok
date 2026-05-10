package entity

import (
	"Tiktok/biz/model/video"
	"database/sql"
)

type VideoEntity struct {
	ID           string       `db:"id"`
	UserID       string       `db:"user_id"`
	Title        string       `db:"title"`
	Description  string       `db:"description"`
	CommentCount int          `db:"comment_count"`
	CoverURL     string       `db:"cover_url"`
	CreatedAt    string       `db:"created_at"`
	DeletedAt    sql.NullTime `db:"deleted_at"`
	LikeCount    int          `db:"like_count"`
	UpdatedAt    string       `db:"updated_at"`
	VideoURL     string       `db:"video_url"`
	VisitCount   int          `db:"visit_count"`
}

func (u *VideoEntity) ToVideoInfo() *video.VideoInfo {
	return &video.VideoInfo{
		ID:           u.ID,
		UserID:       u.UserID,
		Title:        u.Title,
		Description:  u.Description,
		CommentCount: int64(u.CommentCount),
		CoverURL:     u.CoverURL,
		CreatedAt:    u.CreatedAt,
		LikeCount:    int64(u.LikeCount),
		UpdatedAt:    u.UpdatedAt,
		VideoURL:     u.VideoURL,
		VisitCount:   int64(u.VisitCount),
	}
}
