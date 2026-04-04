package service

import (
	"Tiktok/biz/entity"
	"Tiktok/biz/model/video"
	"Tiktok/pkg/consts"
	"Tiktok/pkg/utils"
	"context"
	"io"
	"log"
	"math/rand"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/redis/go-redis/v9"
)

type VideoRedis interface {
	VideoHotSet(ctx context.Context, key string, member interface{}, score float64) error
	VideoHotGet(ctx context.Context, key string, pageNum int64, pageSize int64) ([]redis.Z, error)
}
type VideoDatabase interface {
	CreatVideo(entity entity.VideoEntity) error
	GetVideoByUserID(userId string, pageSize int64, pageNum int64) ([]entity.VideoEntity, error)
	GetVideoByKeyWord(keyword string, pageNum int64, pageSize int64) ([]entity.VideoEntity, error)
	GetVideoByVideoId(videoId string) (entity.VideoEntity, error)
	GetVideoStream() ([]entity.VideoEntity, error)
}
type VideoService struct {
	videoDb    VideoDatabase
	VideoRedis VideoRedis
}

func NewVideoService(videoDb VideoDatabase, videoRedis VideoRedis) *VideoService {
	return &VideoService{videoDb: videoDb, VideoRedis: videoRedis}
}

func (s *VideoService) VideoPublish(videoInfo *video.VideoInfo, data *multipart.FileHeader, ctx context.Context) (int, string) {
	dataFile, err := data.Open()
	if err != nil {
		return consts.CodeIOError, "VideoPublish data.Open err"
	}
	defer dataFile.Close()
	filename := utils.IdGenerate()
	err = os.MkdirAll("/home/lai-long/Tiktok/a", os.ModePerm)
	if err != nil {
		log.Println(err)
		return consts.CodeIOError, "VideoPublish os.MkdirAll err"
	}
	file, err := os.Create("/home/lai-long/Tiktok/a/" + filename + filepath.Ext(data.Filename))
	if err != nil {
		log.Println(err)
		return consts.CodeIOError, "VideoPublish os.Create err"
	}
	defer file.Close()
	if _, err := io.Copy(file, dataFile); err != nil {
		log.Println(err)
		return consts.CodeIOError, "VideoPublish io.copy err"
	}
	var videoEntity entity.VideoEntity
	videoEntity.Title = videoInfo.Title
	videoEntity.Description = videoInfo.Description
	videoEntity.VideoURL = "/home/lai-long/Tiktok/a/" + filename + filepath.Ext(data.Filename)
	videoEntity.UserID = videoInfo.UserID
	videoEntity.ID = filename
	videoEntity.VisitCount = rand.Intn(100)
	err = s.VideoRedis.VideoHotSet(ctx, "videoHot", videoEntity.ID, float64(videoEntity.VisitCount))
	if err != nil {
		return consts.CodeRedisError, `VideoPublish re.VideoHotSet err`
	}
	err = s.videoDb.CreatVideo(videoEntity)
	if err != nil {
		log.Println("VideoPublish db.CreateVideo err: %v", err)
		return consts.CodeDBCreateError, "VideoPublish db.Create err"
	}
	return consts.CodeSuccess, "success"
}

func (s *VideoService) VideoList(userId string, pageSize int64, pageNum int64) (int, string, []*video.VideoInfo, bool) {
	videoList, err := s.videoDb.GetVideoByUserID(userId, pageSize, pageNum)
	if err != nil {
		log.Printf("GetVideoByUserID err: %v", err)
		return consts.CodeDBSelectError, "VideoList GetVideoByUserID error", nil, false
	}
	videoInfos := []*video.VideoInfo{}
	for i := 0; i < len(videoList); i++ {
		videoInfos = append(videoInfos, &video.VideoInfo{
			ID:           videoList[i].ID,
			UserID:       videoList[i].UserID,
			Title:        videoList[i].Title,
			Description:  videoList[i].Description,
			CommentCount: int64(videoList[i].CommentCount),
			CoverURL:     videoList[i].CoverURL,
			CreatedAt:    videoList[i].CreatedAt,
			LikeCount:    int64(videoList[i].LikeCount),
			UpdatedAt:    videoList[i].UpdatedAt,
			VideoURL:     videoList[i].VideoURL,
			VisitCount:   videoInfos[i].VisitCount,
		})
	}
	return consts.CodeSuccess, "success", videoInfos, true
}

func (s *VideoService) VideoSearch(keyword string, pageNum int64, pageSize int64) (int, string, []*video.VideoInfo, bool) {
	videoEntity, err := s.videoDb.GetVideoByKeyWord(keyword, pageNum, pageSize)
	if err != nil {
		log.Printf("db.GetVideoByKeyWord err: %v", err)
		return consts.CodeVideoError, "GetVideoByVideoTitleOrDescription error", nil, false
	}
	videoInfos := []*video.VideoInfo{}
	for i := 0; i < len(videoEntity); i++ {
		videoInfos = append(videoInfos, &video.VideoInfo{
			ID:           videoEntity[i].ID,
			UserID:       videoEntity[i].UserID,
			Title:        videoEntity[i].Title,
			Description:  videoEntity[i].Description,
			CommentCount: int64(videoEntity[i].CommentCount),
			CoverURL:     videoEntity[i].CoverURL,
			CreatedAt:    videoEntity[i].CreatedAt,
			LikeCount:    int64(videoEntity[i].LikeCount),
			UpdatedAt:    videoEntity[i].UpdatedAt,
			VideoURL:     videoEntity[i].VideoURL,
			VisitCount:   videoInfos[i].VisitCount,
		})
	}
	return consts.CodeSuccess, "success", videoInfos, true
}

func (s *VideoService) VideoPopular(ctx context.Context, pageNum int64, pageSize int64) (int, string, []*video.VideoInfo, bool) {
	z, err := s.VideoRedis.VideoHotGet(ctx, "videoHot", pageNum, pageSize)
	if err != nil {
		log.Printf("re.VideoHotGet err: %v", err)
		return consts.CodeRedisError, "VideoPopular re.VideoHotGet err", nil, false
	}
	videoEntity := make([]entity.VideoEntity, len(z))
	for i, _ := range z {
		videoEntity[i], err = s.videoDb.GetVideoByVideoId(z[i].Member.(string))
		if err != nil {
			log.Printf("GetVideoByVideoId %v", err)
			return consts.CodeDBSelectError, "VideoPopular db.GetVideoByVideoId err", nil, false
		}
	}
	var videoInfos []*video.VideoInfo
	for i := 0; i < len(z); i++ {
		videoInfos = append(videoInfos, &video.VideoInfo{
			ID:           videoEntity[i].ID,
			UserID:       videoEntity[i].UserID,
			Title:        videoEntity[i].Title,
			Description:  videoEntity[i].Description,
			CommentCount: int64(videoEntity[i].CommentCount),
			CoverURL:     videoEntity[i].CoverURL,
			CreatedAt:    videoEntity[i].CreatedAt,
			LikeCount:    int64(videoEntity[i].LikeCount),
			UpdatedAt:    videoEntity[i].UpdatedAt,
			VideoURL:     videoEntity[i].VideoURL,
			VisitCount:   videoInfos[i].VisitCount,
		})
	}
	return consts.CodeSuccess, "success", videoInfos, true
}

func (s *VideoService) VideoStream() (int, string, []*video.VideoInfo) {
	videoEntity, err := s.videoDb.GetVideoStream()
	if err != nil {
		log.Printf("videoDb.GetVideoStream err: %v", err)
		return consts.CodeDBSelectError, "videoDb.GetVideoStream err", nil
	}
	videoInfos := []*video.VideoInfo{}
	for _, v := range videoEntity {
		videoInfos = append(videoInfos, &video.VideoInfo{
			ID:           v.ID,
			UserID:       v.UserID,
			Title:        v.Title,
			Description:  v.Description,
			CommentCount: int64(v.CommentCount),
			CoverURL:     v.CoverURL,
			CreatedAt:    v.CreatedAt,
			LikeCount:    int64(v.LikeCount),
			UpdatedAt:    v.UpdatedAt,
			VideoURL:     v.VideoURL,
			VisitCount:   int64(v.VisitCount),
		})
	}
	return consts.CodeSuccess, "success", videoInfos
}
