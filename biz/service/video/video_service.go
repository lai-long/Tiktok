package video

import (
	"Tiktok/biz/entity"
	"Tiktok/biz/model/video"
	"Tiktok/pkg/config"
	"Tiktok/pkg/consts"
	"Tiktok/pkg/utils"
	"context"
	"log"
	"math/rand"
	"mime/multipart"
	"path/filepath"

	"github.com/pkg/errors"
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

func (s *VideoService) VideoPublish(videoInfo *video.VideoInfo, data *multipart.FileHeader, ctx context.Context) (int32, error) {
	dataFile, err := data.Open()
	if err != nil {
		return consts.IOOsError, errors.Wrap(err, "->VideoPublish data.Open err")
	}
	defer func() {
		err := dataFile.Close()
		if err != nil {
			log.Println(errors.Wrap(err, "VideoPublish data close"))
		}
	}()
	filename := utils.IdGenerate()
	code, err := utils.SaveUploadFile(dataFile, config.Cfg.Path.VideoPath, filename+filepath.Ext(data.Filename))
	if err != nil {
		return code, errors.Wrap(err, " VideoPublish ")
	}
	var videoEntity entity.VideoEntity
	videoEntity.Title = videoInfo.Title
	videoEntity.Description = videoInfo.Description
	videoEntity.VideoURL = config.Cfg.Path.VideoPath + filename
	videoEntity.UserID = videoInfo.UserID
	videoEntity.ID = filename
	videoEntity.VisitCount = rand.Intn(100)
	err = s.VideoRedis.VideoHotSet(ctx, "videoHot", videoEntity.ID, float64(videoEntity.VisitCount))
	if err != nil {
		return consts.VideoRedisSetError, errors.Wrap(err, "->VideoPublish redis hot set err")
	}
	err = s.videoDb.CreatVideo(videoEntity)
	if err != nil {
		return consts.VideoDBInsertError, errors.Wrap(err, "->VideoPublish create video err")
	}
	return consts.Success, nil
}

func (s *VideoService) VideoList(userId string, pageSize int64, pageNum int64) (int32, error, []*video.VideoInfo) {
	videoList, err := s.videoDb.GetVideoByUserID(userId, pageSize, pageNum)
	if err != nil {
		return consts.VideoDBSelectError, errors.Wrap(err, "->VideoList GetVideo err"), nil
	}
	videoInfos := []*video.VideoInfo{}
	for i := 0; i < len(videoList); i++ {
		videoInfos = append(videoInfos, videoList[i].ToVideoInfo())
	}
	return consts.Success, nil, videoInfos
}

func (s *VideoService) VideoSearch(keyword string, pageNum int64, pageSize int64) (int32, error, []*video.VideoInfo) {
	videoEntity, err := s.videoDb.GetVideoByKeyWord(keyword, pageNum, pageSize)
	if err != nil {
		return consts.VideoDBSelectError, errors.Wrap(err, "->VideoSearch GetVideo Error"), nil
	}
	videoInfos := []*video.VideoInfo{}
	for i := 0; i < len(videoEntity); i++ {
		videoInfos = append(videoInfos, videoEntity[i].ToVideoInfo())
	}
	return consts.Success, nil, videoInfos
}

func (s *VideoService) VideoPopular(ctx context.Context, pageNum int64, pageSize int64) (int32, error, []*video.VideoInfo) {
	z, err := s.VideoRedis.VideoHotGet(ctx, "videoHot", pageNum, pageSize)
	if err != nil {
		return consts.VideoRedisGetError, errors.Wrap(err, "->VideoPopular GetVideoHot error"), nil
	}
	videoEntity := make([]entity.VideoEntity, len(z))
	for i, _ := range z {
		videoEntity[i], err = s.videoDb.GetVideoByVideoId(z[i].Member.(string))
		if err != nil {
			return consts.VideoDBSelectError, errors.Wrap(err, "->video popular select video"), nil
		}
	}
	var videoInfos []*video.VideoInfo
	for i := 0; i < len(z); i++ {
		videoInfos = append(videoInfos, videoEntity[i].ToVideoInfo())
	}
	return consts.Success, nil, videoInfos
}

func (s *VideoService) VideoStream() (int32, error, []*video.VideoInfo) {
	videoEntity, err := s.videoDb.GetVideoStream()
	if err != nil {
		return consts.VideoDBSelectError, errors.Wrap(err, "->video stream select video error"), nil
	}
	videoInfos := []*video.VideoInfo{}
	for i, _ := range videoEntity {
		videoInfos = append(videoInfos, videoEntity[i].ToVideoInfo())
	}
	return consts.Success, nil, videoInfos
}
