package db

import (
	"Tiktok/biz/model/entity"
)

func LikeCountUp(user_id string, video_id string) error {
	sql := `UPDATE videos SET like_count=like_count + 1 WHERE id = ?`
	_, err := db.Exec(sql, video_id)
	return err
}
func LikeCreate(user_id string, video_id string) error {
	sql := `INSERT INTO likes (video_id, user_id) VALUES (?, ?)`
	_, err := db.Exec(sql, video_id, user_id)
	return err
}
func LikeCountDown(user_id string, video_id string) error {
	sql := `UPDATE videos SET like_count=like_count - 1 WHERE id = ?`
	_, err := db.Exec(sql, video_id)
	return err
}
func LikeDelete(user_id string, video_id string) error {
	sql := `DELETE FROM likes WHERE video_id = ? AND user_id = ?`
	_, err := db.Exec(sql, video_id, user_id)
	return err
}
func LikeVideoIds(user_id string, pageNum int, pageSize int) (error, []string) {
	sql := `SELECT video_id FROM likes WHERE user_id = ? ORDER BY video_id DESC LIMIT ? OFFSET ?`
	var video_id []string
	offset := (pageNum - 1) * pageSize
	err := db.Select(&video_id, sql, user_id, offset)
	return err, video_id
}
func LikeVideos(videoId []string) (bool, []entity.VideoEntity) {
	videos := make([]entity.VideoEntity, len(videoId))
	var GetVideoErrors = 0
	for i, _ := range videos {
		var err error
		videos[i], err = GetVideoByVideoId(videoId[i])
		if err != nil {
			GetVideoErrors = GetVideoErrors + 1
		}
	}
	if GetVideoErrors == 0 {
		return true, videos
	} else {
		return false, videos
	}
}
