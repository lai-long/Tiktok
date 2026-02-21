package dto

type Comment struct {
	UserId    string `json:"user_id"`
	VideoId   string `json:"video_id"`
	CommentId string `json:"comment_id"`
	Content   string `json:"content"`
	LikeCount int64  `json:"like_count"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
}
