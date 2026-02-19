package db

import "Tiktok/biz/model/entity"

func CreatVideo(entity entity.VideoEntity) error {
	sql := `INSERT INTO videos (title ,description,id,user_id,video_url,cover_url) VALUES(?,?,?,?,?,?)`
	_, err := db.Exec(sql, entity.Title, entity.Description, entity.ID, entity.UserID, entity.VideoURL)
	return err
}
func GetVideoByUserID(userId string, pageSize string, pageNum string) (entity.VideoEntity, error) {
	var video entity.VideoEntity
	sql := `select * from videos where user_id = ? ORDER BY id DESC LIMIT ? OFFSET ?`
	err := db.Get(&video, sql, userId, pageSize, pageNum)
	return video, err
}
func GetVideoByVideoTitleOrDescription(title string, description string) (entity.VideoEntity, error) {
	var video entity.VideoEntity
	sql := `select * from videos where title = ? and description = ?`
	err := db.Get(&video, sql, title, description)
	return video, err
}
