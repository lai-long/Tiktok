package service

import (
	"Tiktok/kitex_gen/video"
	"Tiktok/pkg/consts"
	"Tiktok/pkg/entity"
	"context"
	"log"
	"time"

	"github.com/pkg/errors"
)

type LikeRedis interface {
	VideoLikeSAdd(ctx context.Context, userId string, videoId string) error
	VideoDislikeSRem(ctx context.Context, userId string, videoId string) error
	VideoLikeGet(ctx context.Context, userId string) ([]string, error)
}

type LikeCommentDatabase interface {
	CommentLikeCountUp(commentId string) error
	CommentLikeCountDown(commentId string) error
}
type LikeVideoDatabase interface {
	VideoLikeCountUp(videoId string) error
	VideoLikeCountDown(videoId string) error
	LikeVideoIds(userId string, pageNum int64, pageSize int64) ([]string, error)
	LikeVideos(videoId []string) (bool, []entity.VideoEntity)
}
type LikeDatabase interface {
	LikeCreate(userId string, targetId string, targetType string) error
	LikeDelete(userId, targetId, targetType string) error
}
type LikeRepo struct {
	videoDb   LikeVideoDatabase
	commentDb LikeCommentDatabase
	likeDb    LikeDatabase
	likeRedis LikeRedis
}

func NewLikeRepo(videoDb LikeVideoDatabase, commentDb LikeCommentDatabase, likeDb LikeDatabase, likeRedis LikeRedis) *LikeRepo {
	return &LikeRepo{
		videoDb:   videoDb,
		commentDb: commentDb,
		likeDb:    likeDb,
		likeRedis: likeRedis,
	}
}

func (s *LikeRepo) LikeVideo(ctx context.Context, userId string, targetId string, targetType string) (int32, error) {
	err := s.likeDb.LikeCreate(userId, targetId, targetType)
	if err != nil {
		return consts.ReactDBInsertError, errors.Wrap(err, "->LikeAction LikeCreate error")
	}
	err = s.videoDb.VideoLikeCountUp(targetId)
	if err != nil {
		log.Println("LikeAction VideoLikeCount Up error:", err)
	}
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := s.likeRedis.VideoLikeSAdd(ctx, userId, targetId)
		if err != nil {
			log.Println("LikeAction VideoLikeSAdd error:", err)
		}
	}()
	return consts.Success, nil
}

func (s *LikeRepo) DislikeVideo(ctx context.Context, userId string, targetId string, targetType string) (int32, error) {
	err := s.likeDb.LikeDelete(userId, targetId, targetType)
	if err != nil {
		return consts.ReactDBDeleteError, errors.Wrap(err, "->LikeAction LikeDelete error")
	}
	err = s.videoDb.VideoLikeCountDown(targetId)
	if err != nil {
		return consts.ReactDBUpdateError, errors.Wrap(err, "->LikeAction VideoLikeCount down error")
	}
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = s.likeRedis.VideoDislikeSRem(ctx, userId, targetId)
		if err != nil {
			log.Println("LikeAction VideoDislikeSRem error:", err)
		}
	}()
	return consts.Success, nil
}

func (s *LikeRepo) LikeComment(userId string, targetId string, targetType string) (int32, error) {
	err := s.likeDb.LikeCreate(userId, targetId, targetType)
	if err != nil {
		return consts.ReactDBInsertError, errors.Wrap(err, "->LikeAction LikeCreate error")
	}
	err = s.commentDb.CommentLikeCountUp(targetId)
	if err != nil {
		return consts.ReactDBUpdateError, errors.Wrap(err, "->LikeAction CommentLikeCount up error")
	}
	return consts.Success, nil
}

func (s *LikeRepo) DislikeComment(userId string, targetId string, targetType string) (int32, error) {
	err := s.likeDb.LikeDelete(userId, targetId, targetType)
	if err != nil {
		return consts.ReactDBDeleteError, errors.Wrap(err, "->LikeAction LikeDelete error")
	}
	err = s.commentDb.CommentLikeCountDown(targetId)
	if err != nil {
		return consts.ReactDBUpdateError, errors.Wrap(err, "->LikeAction CommentLikeCount down error")
	}
	return consts.Success, nil
}

func (s *LikeRepo) LikeAction(ctx context.Context, userId string, targetId string, action string, targetType string) (int32, error) {
	switch targetType {
	case consts.TargetVideo:
		switch action {
		case consts.ActionLike:
			code, err := s.LikeVideo(ctx, userId, targetId, targetType)
			return code, err
		case consts.ActionDislike:
			code, err := s.DislikeVideo(ctx, userId, targetId, targetType)
			return code, err
		default:
			return consts.ReactReqValueError, errors.Errorf("invalid action type: %s", action)
		}
	case consts.TargetComment:
		switch action {
		case consts.ActionLike:
			code, err := s.LikeComment(userId, targetId, targetType)
			return code, err
		case consts.ActionDislike:
			code, err := s.DislikeComment(userId, targetId, targetType)
			return code, err
		default:
			return consts.ReactReqValueError, errors.New("->LikeAction action type error")
		}
	default:
		return consts.ReactReqValueError, errors.New("->LikeAction targetType is not valid")
	}
}

func (s *LikeRepo) LikeList(ctx context.Context, userId string, pageNum int64, pageSize int64) (int32, []*video.VideoInfo, error) {
	results, err := s.likeRedis.VideoLikeGet(ctx, userId)
	if err == nil && len(results) > 0 {
		start := pageNum * pageSize
		end := pageSize + start
		if start >= int64(len(results)) {
			return consts.ReactReqValueError, nil, errors.New("pageNum out of range")
		}
		if end > int64(len(results)) {
			end = int64(len(results))
		}
		result := results[start:end]
		ok, videos := s.videoDb.LikeVideos(result)
		if !ok {
			return consts.ReactDBSelectError, nil, errors.New("->LikeList LikeVideos err")
		}

		var videoInfos []*video.VideoInfo
		for _, v := range videos {
			videoInfos = append(videoInfos, v.ToVideoInfo())
		}
		return consts.Success, videoInfos, nil
	}
	videoId, err := s.videoDb.LikeVideoIds(userId, pageNum, pageSize)
	if err != nil {
		return consts.ReactDBSelectError, nil, errors.Wrap(err, "->LikeList select LikeVideo error")
	}
	ok, videos := s.videoDb.LikeVideos(videoId)
	if !ok {
		return consts.ReactDBSelectError, nil, errors.New("->LikeList LikeVideos err")
	}
	var videoInfos []*video.VideoInfo
	for _, v := range videos {
		videoInfos = append(videoInfos, v.ToVideoInfo())
	}
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		for _, id := range videoId {
			err = s.likeRedis.VideoLikeSAdd(ctx, userId, id)
			if err != nil {
				log.Println("LikeAction LikeSAdd error:", err)
			}
		}
	}()
	return consts.Success, videoInfos, nil
}
