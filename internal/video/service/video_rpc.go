package service

import (
	"Tiktok/kitex_gen/video"

	"Tiktok/pkg/consts"
	"Tiktok/pkg/entity"
	"Tiktok/pkg/logger"
	"Tiktok/pkg/utils"
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type VideoRedis interface {
	VideoHotSet(ctx context.Context, key string, member interface{}, score float64) error
	VideoHotGet(ctx context.Context, key string, pageNum int64, pageSize int64) ([]redis.Z, error)
	VideoHotIncrBy(ctx context.Context, key string, videoID string, delta float64) error
	VideoInfoSet(ctx context.Context, VideoID string, video *entity.VideoEntity) error
	VideoInfoGet(ctx context.Context, VideoID string) (*entity.VideoEntity, error)
	VideoInfoDelete(ctx context.Context, VideoID string) error
}
type VideoDatabase interface {
	CreateVideo(ctx context.Context, entity entity.VideoEntity) error
	GetVideoByUserID(ctx context.Context, userId string, pageSize int64, pageNum int64) ([]entity.VideoEntity, error)
	GetVideoByKeyWord(ctx context.Context, keyword string, pageNum int64, pageSize int64) ([]entity.VideoEntity, error)
	GetVideoByVideoId(ctx context.Context, videoId string) (entity.VideoEntity, error)
	GetVideoStream(ctx context.Context) ([]entity.VideoEntity, error)
	GetVideoByIds(ctx context.Context, ids []string) ([]entity.VideoEntity, error)
	VideoVisitCountUp(ctx context.Context, videoID string) error
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
	defer utils.TrackTime(ctx, "VideoPublish")()
	var videoEntity entity.VideoEntity
	videoEntity.Title = title
	videoEntity.Description = description
	videoEntity.VideoURL = url
	videoEntity.CoverURL = coverURL
	videoEntity.UserID = userID
	videoEntity.VisitCount = 0
	videoEntity.ID = utils.IDGenerate()
	err := s.VideoRedis.VideoHotSet(ctx, "videoHot", videoEntity.ID, 0)
	if err != nil {
		return consts.VideoRedisSetError, errors.Wrap(err, "->VideoPublish redis hot set err")
	}
	err = s.videoDb.CreateVideo(ctx, videoEntity)
	if err != nil {
		return consts.VideoDBInsertError, errors.Wrap(err, "->VideoPublish create video err")
	}
	go func() {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		_ = s.VideoRedis.VideoInfoSet(ctx, videoEntity.ID, &videoEntity)
	}()
	return consts.Success, nil
}

func (s *VideoRepo) VideoList(ctx context.Context, userId string, pageSize int64, pageNum int64) (int32, []*video.VideoInfo, error) {
	defer utils.TrackTime(ctx, "VideoList")()
	videoList, err := s.videoDb.GetVideoByUserID(ctx, userId, pageSize, pageNum)
	if err != nil {
		return consts.VideoDBSelectError, nil, errors.Wrap(err, "->VideoList GetVideo err")
	}
	return consts.Success, toVideoInfoList(videoList), nil
}

func (s *VideoRepo) VideoSearch(ctx context.Context, keyword string, pageNum int64, pageSize int64) (int32, []*video.VideoInfo, error) {
	defer utils.TrackTime(ctx, "VideoSearch")()
	videoEntity, err := s.videoDb.GetVideoByKeyWord(ctx, keyword, pageNum, pageSize)
	if err != nil {
		return consts.VideoDBSelectError, nil, errors.Wrap(err, "->VideoSearch GetVideo Error")
	}
	return consts.Success, toVideoInfoList(videoEntity), nil
}

func (s *VideoRepo) VideoPopular(ctx context.Context, pageNum int64, pageSize int64) (int32, []*video.VideoInfo, error) {
	defer utils.TrackTime(ctx, "VideoPopular")()
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
			videoEntity[i], err = s.videoDb.GetVideoByVideoId(ctx, z[i].Member.(string))
			if err != nil {
				return consts.VideoDBSelectError, nil, errors.Wrap(err, "->video popular select video")
			}
			_ = s.VideoRedis.VideoInfoSet(ctx, z[i].Member.(string), &videoEntity[i])
		}
	}
	return consts.Success, toVideoInfoList(videoEntity), nil
}

func (s *VideoRepo) VideoStream(ctx context.Context) (int32, []*video.VideoInfo, error) {
	defer utils.TrackTime(ctx, "VideoStream")()
	videoEntity, err := s.videoDb.GetVideoStream(ctx)
	if err != nil {
		return consts.VideoDBSelectError, nil, errors.Wrap(err, "->video stream select video error")
	}
	for _, v := range videoEntity {
		if _, err := s.VideoVisitCountUp(ctx, v.ID); err != nil {
			logger.Error("VideoStream visit count error", zap.Error(err), logger.WithTraceID(utils.GetTraceID(ctx)))
		}
	}
	return consts.Success, toVideoInfoList(videoEntity), nil
}

func (s *VideoRepo) BatchGetVideo(ctx context.Context, ids []string) (int32, []*video.VideoInfo, error) {
	defer utils.TrackTime(ctx, "BatchGetVideo")()
	videos, err := s.videoDb.GetVideoByIds(ctx, ids)
	if err != nil {
		return consts.VideoDBSelectError, nil, errors.Wrap(err, "->BatchGetVideo GetVideoByIds err")
	}
	return consts.Success, toVideoInfoList(videos), nil
}

func (s *VideoRepo) VideoVisitCountUp(ctx context.Context, videoID string) (int32, error) {
	defer utils.TrackTime(ctx, "VideoVisitCountUp")()
	err := s.videoDb.VideoVisitCountUp(ctx, videoID)
	if err != nil {
		return consts.VideoDBUpdateError, errors.Wrap(err, "->VideoVisitCountUp update visit count err")
	}
	go func() {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		if err := s.VideoRedis.VideoHotIncrBy(ctx, "videoHot", videoID, consts.HotScoreWeightView); err != nil {
			logger.Error("VideoVisitCountUp hot incr error", zap.Error(err), logger.WithTraceID(utils.GetTraceID(ctx)))
		}
		if err := s.VideoRedis.VideoInfoDelete(ctx, videoID); err != nil {
			logger.Error("VideoVisitCountUp cache delete error", zap.Error(err), logger.WithTraceID(utils.GetTraceID(ctx)))
		}
	}()
	return consts.Success, nil
}
