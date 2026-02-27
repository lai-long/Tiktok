package dto

type Comment struct {
	UserId    string `json:"user_id" form:"user_id"`
	VideoId   string `json:"video_id" form:"video_id"`
	CommentId string `json:"comment_id" form:"comment_id"`
	Content   string `json:"content" form:"content"`
	LikeCount int64  `json:"like_count" form:"like_count"`
	CreatedAt string `json:"created_at" form:"created_at"`
	UpdatedAt string `json:"updated_at" form:"updated_at"`
	DeletedAt string `json:"deleted_at" form:"deleted_at"`
}
