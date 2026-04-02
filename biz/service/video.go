package service

import (
	"Tiktok/biz/model/dto"
	"Tiktok/biz/model/entity"
	"Tiktok/pkg/consts"
	"Tiktok/pkg/utils"
	"context"
	"io"
	"log"
	"math/rand"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type VideoRedis interface {
	VideoHotSet(ctx context.Context, key string, member interface{}, score float64) error
	VideoHotGet(ctx context.Context, key string, pageNum int, pageSize int) ([]redis.Z, error)
}
type VideoDatabase interface {
	CreatVideo(entity entity.VideoEntity) error
	GetVideoByUserID(userId string, pageSize int, pageNum int) ([]entity.VideoEntity, error)
	GetVideoByKeyWord(keyword string, pageNum int, pageSize int) ([]entity.VideoEntity, error)
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

func (s *VideoService) VideoPublish(video dto.Video, data *multipart.FileHeader, ctx context.Context) (int, string) {
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
	videoEntity.Title = video.Title
	videoEntity.Description = video.Description
	videoEntity.VideoURL = "/home/lai-long/Tiktok/a/" + filename + filepath.Ext(data.Filename)
	videoEntity.UserID = video.UserID
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

func (s *VideoService) VideoList(userId string, pageSize string, pageNum string) (int, string, []dto.Video, bool) {
	pageNumInt := 0
	pageSizeInt := 10
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		log.Printf("strconv.Atoi error: %v", err)
		return consts.CodeError, "VideoList pageSize strconv error", []dto.Video{}, false
	}
	pageNumInt, err = strconv.Atoi(pageNum)
	if err != nil {
		log.Printf("strconv.Atoi error: %v", err)
		return consts.CodeError, "VideoList pageNum error", []dto.Video{}, false
	}
	videoList, err := s.videoDb.GetVideoByUserID(userId, pageSizeInt, pageNumInt)
	if err != nil {
		log.Printf("GetVideoByUserID err: %v", err)
		return consts.CodeDBSelectError, "VideoList GetVideoByUserID error", []dto.Video{}, false
	}
	videoDTOs := make([]dto.Video, len(videoList))
	for i := 0; i < len(videoList); i++ {
		videoDTOs[i] = dto.Video{
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
			VisitCount:   int64(videoList[i].VisitCount),
		}
	}
	return consts.CodeSuccess, "success", videoDTOs, true
}

func (s *VideoService) VideoSearch(keyword string, pageNum string, pageSize string) (int, string, []dto.Video, bool) {
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		log.Printf("strconv.Atoi error: %v", err)
		return consts.CodeError, "VideoSearch pageSize strconv error", []dto.Video{}, false
	}
	pageNumInt, err := strconv.Atoi(pageNum)
	if err != nil {
		log.Printf("VideoSearch pageNum  strconv error: %v", err)
		return consts.CodeError, "VideoSearch pageNum error", []dto.Video{}, false
	}
	video, err := s.videoDb.GetVideoByKeyWord(keyword, pageNumInt, pageSizeInt)
	if err != nil {
		log.Printf("db.GetVideoByKeyWord err: %v", err)
		return consts.CodeVideoError, "GetVideoByVideoTitleOrDescription error", []dto.Video{}, false
	}
	videoDTOs := make([]dto.Video, len(video))
	for i := 0; i < len(video); i++ {
		videoDTOs[i].ID = video[i].ID
		videoDTOs[i].Title = video[i].Title
		videoDTOs[i].Description = video[i].Description
		videoDTOs[i].VideoURL = video[i].VideoURL
		videoDTOs[i].CreatedAt = video[i].CreatedAt
		videoDTOs[i].LikeCount = int64(video[i].LikeCount)
		videoDTOs[i].UpdatedAt = video[i].UpdatedAt
		videoDTOs[i].VideoURL = video[i].VideoURL
		videoDTOs[i].CoverURL = video[i].CoverURL
		videoDTOs[i].CommentCount = int64(video[i].CommentCount)
		videoDTOs[i].CreatedAt = video[i].CreatedAt
	}
	return consts.CodeSuccess, "success", videoDTOs, true
}

func (s *VideoService) VideoPopular(ctx context.Context, pageNum string, pageSize string) (int, string, []dto.Video, bool) {
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		log.Printf("strconv.Atoi error: %v", err)
		return consts.CodeVideoError, "VideoPopular pageSize strconv error", []dto.Video{}, false
	}
	pageNumInt, err := strconv.Atoi(pageNum)
	if err != nil {
		log.Printf("strconv.Atoi error: %v", err)
		return consts.CodeVideoError, "VideoPopular pageNum strconv error", []dto.Video{}, false
	}
	z, err := s.VideoRedis.VideoHotGet(ctx, "videoHot", pageNumInt, pageSizeInt)
	if err != nil {
		log.Printf("re.VideoHotGet err: %v", err)
		return consts.CodeRedisError, "VideoPopular re.VideoHotGet err", []dto.Video{}, false
	}
	videoEntity := make([]entity.VideoEntity, len(z))
	for i, _ := range z {
		videoEntity[i], err = s.videoDb.GetVideoByVideoId(z[i].Member.(string))
		if err != nil {
			log.Printf("GetVideoByVideoId %v", err)
			return consts.CodeDBSelectError, "VideoPopular db.GetVideoByVideoId err", []dto.Video{}, false
		}
	}
	videoDTOs := make([]dto.Video, len(z))
	for i := 0; i < len(z); i++ {
		videoDTOs[i].ID = videoEntity[i].ID
		videoDTOs[i].Title = videoEntity[i].Title
		videoDTOs[i].Description = videoEntity[i].Description
		videoDTOs[i].VideoURL = videoEntity[i].VideoURL
		videoDTOs[i].CreatedAt = videoEntity[i].CreatedAt
		videoDTOs[i].VisitCount = int64(videoEntity[i].VisitCount)
	}
	return consts.CodeSuccess, "success", videoDTOs, true
}

func (s *VideoService) VideoStream() (int, string, []dto.Video) {
	videoEntity, err := s.videoDb.GetVideoStream()
	if err != nil {
		log.Printf("videoDb.GetVideoStream err: %v", err)
		return consts.CodeDBSelectError, "videoDb.GetVideoStream err", nil
	}
	video := make([]dto.Video, len(videoEntity))
	for i, v := range videoEntity {
		video[i].CreatedAt = v.CreatedAt
		video[i].UpdatedAt = v.UpdatedAt
		video[i].VideoURL = v.VideoURL
		video[i].CoverURL = v.CoverURL
		video[i].Title = v.Title
		video[i].Description = v.Description
		video[i].LikeCount = int64(v.LikeCount)
		video[i].CommentCount = int64(v.CommentCount)
	}
	return consts.CodeSuccess, "success", video
}
