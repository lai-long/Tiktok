package db

import "Tiktok/biz/model/entity"

func CreatVideo(entity entity.VideoEntity) error {
	sql := `INSERT INTO videos (title ,description,id,user_id,video_url,visit_count) VALUES(?,?,?,?,?,?)`
	_, err := db.Exec(sql, entity.Title, entity.Description, entity.ID, entity.UserID, entity.VideoURL, entity.VisitCount)
	return err
}
func GetVideoByUserID(userId string, pageSize int, pageNum int) ([]entity.VideoEntity, error) {
	var video []entity.VideoEntity
	sql := `select * from videos where user_id = ? ORDER BY id DESC LIMIT ? OFFSET ?`
	offset := pageNum * pageSize
	err := db.Select(&video, sql, userId, pageSize, offset)
	return video, err
}
func GetVideoByKeyWord(keyword string, pageNum int, pageSize int) ([]entity.VideoEntity, error) {
	var video []entity.VideoEntity
	keywords := "%" + keyword + "%"
	sql := `select * from videos where title like ? or description like ? ORDER BY id DESC LIMIT ? OFFSET ? `
	offset := pageNum * pageSize
	err := db.Select(&video, sql, keywords, keywords, pageSize, offset)
	return video, err
}
func GetVideoByVideoId(videoId string) (entity.VideoEntity, error) {
	var video entity.VideoEntity
	sql := `select * from videos where id= ?`
	err := db.Get(&video, sql, videoId)
	return video, err
}
