package entity

type VideoEntity struct {
	ID           string `db:"id"`
	UserID       string `db:"user_id"`
	Title        string `db:"title"`
	Description  string `db:"description"`
	CommentCount int64  `db:"comment_count"`
	CoverURL     string `db:"cover_url"`
	CreatedAt    string `db:"created_at"`
	DeletedAt    string `db:"deleted_at"`
	LikeCount    int64  `db:"like_count"`
	UpdatedAt    string `db:"updated_at"`
	VideoURL     string `db:"video_url"`
	VisitCount   int64  `db:"visit_count"`
}
