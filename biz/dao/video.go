package dao

import (
	"Tiktok/biz/entity"
)

func (m *MySQLdb) CreatVideo(entity entity.VideoEntity) error {
	sql := `INSERT INTO videos (title ,description,id,user_id,video_url,visit_count) VALUES(?,?,?,?,?,?)`
	_, err := m.db.Exec(sql, entity.Title, entity.Description, entity.ID, entity.UserID, entity.VideoURL, entity.VisitCount)
	return err
}

func (m *MySQLdb) GetVideoByUserID(userId string, pageSize int, pageNum int) ([]entity.VideoEntity, error) {
	var video []entity.VideoEntity
	sql := `select * from videos where user_id = ? ORDER BY id DESC LIMIT ? OFFSET ?`
	offset := pageNum * pageSize
	err := m.db.Select(&video, sql, userId, pageSize, offset)
	return video, err
}

func (m *MySQLdb) GetVideoByKeyWord(keyword string, pageNum int, pageSize int) ([]entity.VideoEntity, error) {
	var video []entity.VideoEntity
	keywords := "%" + keyword + "%"
	sql := `select * from videos where title like ? or description like ? ORDER BY id DESC LIMIT ? OFFSET ? `
	offset := pageNum * pageSize
	err := m.db.Select(&video, sql, keywords, keywords, pageSize, offset)
	return video, err
}

func (m *MySQLdb) GetVideoByVideoId(videoId string) (entity.VideoEntity, error) {
	var video entity.VideoEntity
	sql := `select * from videos where id= ?`
	err := m.db.Get(&video, sql, videoId)
	return video, err
}

func (m *MySQLdb) GetVideoStream() ([]entity.VideoEntity, error) {
	var video []entity.VideoEntity
	sql := `select * from videos  ORDER BY RAND() LIMIT 10`
	err := m.db.Select(&video, sql)
	if err != nil {
		return video, err
	}
	return video, nil
}
