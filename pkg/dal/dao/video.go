package dao

import (
	"Tiktok/pkg/entity"
	"context"

	"github.com/jmoiron/sqlx"
)

func (m *MySQLdb) CreatVideo(ctx context.Context, entity entity.VideoEntity) error {
	sql := `INSERT INTO videos (title,description,id,user_id,video_url,cover_url,visit_count) VALUES(?,?,?,?,?,?,?)`
	_, err := m.db.ExecContext(ctx, sql, entity.Title, entity.Description, entity.ID, entity.UserID, entity.VideoURL,
		entity.CoverURL, entity.VisitCount)
	return err
}

func (m *MySQLdb) GetVideoByUserID(ctx context.Context, userID string, pageSize int64, pageNum int64) ([]entity.VideoEntity, error) {
	var video []entity.VideoEntity
	sql := `select * from videos where user_id = ? ORDER BY id DESC LIMIT ? OFFSET ?`
	offset := pageNum * pageSize
	err := m.db.SelectContext(ctx, &video, sql, userID, pageSize, offset)
	return video, err
}

func (m *MySQLdb) GetVideoByKeyWord(ctx context.Context, keyword string, pageNum int64, pageSize int64) ([]entity.VideoEntity, error) {
	var video []entity.VideoEntity
	keywords := "%" + keyword + "%"
	sql := `select * from videos where title like ? or description like ? ORDER BY id DESC LIMIT ? OFFSET ? `
	offset := pageNum * pageSize
	err := m.db.SelectContext(ctx, &video, sql, keywords, keywords, pageSize, offset)
	return video, err
}

func (m *MySQLdb) GetVideoByVideoId(ctx context.Context, videoID string) (entity.VideoEntity, error) {
	var video entity.VideoEntity
	sql := `select * from videos where id= ?`
	err := m.db.GetContext(ctx, &video, sql, videoID)
	return video, err
}

func (m *MySQLdb) GetVideoStream(ctx context.Context) ([]entity.VideoEntity, error) {
	var video []entity.VideoEntity
	sql := `SELECT * FROM videos ORDER BY created_at DESC LIMIT 10`
	err := m.db.SelectContext(ctx, &video, sql)
	if err != nil {
		return video, err
	}
	return video, nil
}

func (m *MySQLdb) GetVideoByIds(ctx context.Context, ids []string) ([]entity.VideoEntity, error) {
	var videos []entity.VideoEntity
	query, args, err := sqlx.In("SELECT * FROM videos WHERE id IN (?)", ids)
	if err != nil {
		return nil, err
	}
	err = m.db.SelectContext(ctx, &videos, query, args...)
	return videos, err
}
