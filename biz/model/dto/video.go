package dto

// 视频
type Video struct {
	ID           string `json:"id" `
	UserID       string `json:"user_id" form:"user_id"`
	Title        string `form:"title" json:"title" binding:"required"`
	Description  string `form:"description" json:"description" binding:"required"`
	CommentCount int64  `json:"comment_count"`
	CoverURL     string `json:"cover_url"`
	CreatedAt    string `json:"created_at"`
	DeletedAt    string `json:"deleted_at"`
	LikeCount    int64  `json:"like_count"`
	UpdatedAt    string `json:"updated_at"`
	VideoURL     string `json:"video_url"`
	VisitCount   int64  `json:"visit_count"`
}
