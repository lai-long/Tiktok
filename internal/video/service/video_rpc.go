package service

import (
	"Tiktok/kitex_gen/video"

	"Tiktok/pkg/consts"
	"Tiktok/pkg/entity"
	"Tiktok/pkg/utils"
	"context"
	"math/rand"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type VideoRedis interface {
	VideoHotSet(ctx context.Context, key string, member interface{}, score float64) error
	VideoHotGet(ctx context.Context, key string, pageNum int64, pageSize int64) ([]redis.Z, error)
	VideoInfoSet(ctx context.Context, VideoID string, video *entity.VideoEntity) error
	VideoInfoGet(ctx context.Context, VideoID string) (*entity.VideoEntity, error)
}
type VideoDatabase interface {
	CreatVideo(entity entity.VideoEntity) error
	GetVideoByUserID(userId string, pageSize int64, pageNum int64) ([]entity.VideoEntity, error)
	GetVideoByKeyWord(keyword string, pageNum int64, pageSize int64) ([]entity.VideoEntity, error)
	GetVideoByVideoId(videoId string) (entity.VideoEntity, error)
	GetVideoStream() ([]entity.VideoEntity, error)
	GetVideoByIds(ids []string) ([]entity.VideoEntity, error)
}
type VideoRepo struct {
	videoDb    VideoDatabase
	VideoRedis VideoRedis
}

func NewVideoRepo(videoDb VideoDatabase, videoRedis VideoRedis) *VideoRepo {
	return &VideoRepo{videoDb: videoDb, VideoRedis: videoRedis}
}

func toVideoInfo(e entity.VideoEntity) *video.VideoInfo {
	return &video.VideoInfo{
		ID:           e.ID,
		UserID:       e.UserID,
		Title:        e.Title,
		Description:  e.Description,
		CommentCount: int64(e.CommentCount),
		CoverURL:     utils.SignQiNiuURL(e.CoverURL),
		CreatedAt:    e.CreatedAt,
		LikeCount:    int64(e.LikeCount),
		UpdatedAt:    e.UpdatedAt,
		VideoURL:     utils.SignQiNiuURL(e.VideoURL),
		VisitCount:   int64(e.VisitCount),
	}
}

func toVideoInfoList(videos []entity.VideoEntity) []*video.VideoInfo {
	result := make([]*video.VideoInfo, len(videos))
	for i, v := range videos {
		result[i] = toVideoInfo(v)
	}
	return result
}

func (s *VideoRepo) VideoPublish(ctx context.Context, title string, description string, url string, coverURL string, userID string) (int32, error) {
	var videoEntity entity.VideoEntity
	videoEntity.Title = title
	videoEntity.Description = description
	videoEntity.VideoURL = url
	videoEntity.CoverURL = coverURL
	videoEntity.UserID = userID
	videoEntity.VisitCount = rand.Intn(100)
	videoEntity.ID = utils.IDGenerate()
	err := s.VideoRedis.VideoHotSet(ctx, "videoHot", videoEntity.ID, float64(videoEntity.VisitCount))
	if err != nil {
		return consts.VideoRedisSetError, errors.Wrap(err, "->VideoPublish redis hot set err")
	}
	err = s.videoDb.CreatVideo(videoEntity)
	if err != nil {
		return consts.VideoDBInsertError, errors.Wrap(err, "->VideoPublish create video err")
	}
	go func() {
		_ = s.VideoRedis.VideoInfoSet(ctx, videoEntity.ID, &videoEntity)
	}()
	return consts.Success, nil
}

func (s *VideoRepo) VideoList(userId string, pageSize int64, pageNum int64) (int32, []*video.VideoInfo, error) {
	videoList, err := s.videoDb.GetVideoByUserID(userId, pageSize, pageNum)
	if err != nil {
		return consts.VideoDBSelectError, nil, errors.Wrap(err, "->VideoList GetVideo err")
	}
	return consts.Success, toVideoInfoList(videoList), nil
}

func (s *VideoRepo) VideoSearch(keyword string, pageNum int64, pageSize int64) (int32, []*video.VideoInfo, error) {
	videoEntity, err := s.videoDb.GetVideoByKeyWord(keyword, pageNum, pageSize)
	if err != nil {
		return consts.VideoDBSelectError, nil, errors.Wrap(err, "->VideoSearch GetVideo Error")
	}
	return consts.Success, toVideoInfoList(videoEntity), nil
}

func (s *VideoRepo) VideoPopular(ctx context.Context, pageNum int64, pageSize int64) (int32, []*video.VideoInfo, error) {
	z, err := s.VideoRedis.VideoHotGet(ctx, "videoHot", pageNum, pageSize)
	if err != nil {
		return consts.VideoRedisGetError, nil, errors.Wrap(err, "->VideoPopular GetVideoHot error")
	}
	videoEntity := make([]entity.VideoEntity, len(z))
	for i := range z {
		videoEntityTemp, err := s.VideoRedis.VideoInfoGet(ctx, z[i].Member.(string))
		if err == nil {
			videoEntity[i] = *videoEntityTemp
		} else {
			videoEntity[i], err = s.videoDb.GetVideoByVideoId(z[i].Member.(string))
			if err != nil {
				return consts.VideoDBSelectError, nil, errors.Wrap(err, "->video popular select video")
			}
			_ = s.VideoRedis.VideoInfoSet(ctx, z[i].Member.(string), &videoEntity[i])
		}
	}
	return consts.Success, toVideoInfoList(videoEntity), nil
}

func (s *VideoRepo) VideoStream() (int32, []*video.VideoInfo, error) {
	videoEntity, err := s.videoDb.GetVideoStream()
	if err != nil {
		return consts.VideoDBSelectError, nil, errors.Wrap(err, "->video stream select video error")
	}
	return consts.Success, toVideoInfoList(videoEntity), nil
}

func (s *VideoRepo) BatchGetVideo(ids []string) (int32, []*video.VideoInfo, error) {
	videos, err := s.videoDb.GetVideoByIds(ids)
	if err != nil {
		return consts.VideoDBSelectError, nil, errors.Wrap(err, "->BatchGetVideo GetVideoByIds err")
	}
	return consts.Success, toVideoInfoList(videos), nil
}
