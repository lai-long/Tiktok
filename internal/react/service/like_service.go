package service

import (
	"Tiktok/pkg/consts"
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
	CommentLikeCountUp(ctx context.Context, commentId string) error
	CommentLikeCountDown(ctx context.Context, commentId string) error
}
type LikeVideoDatabase interface {
	VideoLikeCountUp(ctx context.Context, videoId string) error
	VideoLikeCountDown(ctx context.Context, videoId string) error
	LikeVideoIds(ctx context.Context, userId string, pageNum int64, pageSize int64) ([]string, error)
}
type LikeDatabase interface {
	LikeCreate(ctx context.Context, userId string, targetId string, targetType string) error
	LikeDelete(ctx context.Context, userId string, targetId string, targetType string) error
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
	err := s.likeDb.LikeCreate(ctx, userId, targetId, targetType)
	if err != nil {
		return consts.ReactDBInsertError, errors.Wrap(err, "->LikeAction LikeCreate error")
	}
	err = s.videoDb.VideoLikeCountUp(ctx, targetId)
	if err != nil {
		log.Println("LikeAction VideoLikeCount Up error:", err)
	}
	go func() {
		asyncCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := s.likeRedis.VideoLikeSAdd(asyncCtx, userId, targetId)
		if err != nil {
			log.Println("LikeAction VideoLikeSAdd error:", err)
		}
	}()
	return consts.Success, nil
}

func (s *LikeRepo) DislikeVideo(ctx context.Context, userId string, targetId string, targetType string) (int32, error) {
	err := s.likeDb.LikeDelete(ctx, userId, targetId, targetType)
	if err != nil {
		return consts.ReactDBDeleteError, errors.Wrap(err, "->LikeAction LikeDelete error")
	}
	err = s.videoDb.VideoLikeCountDown(ctx, targetId)
	if err != nil {
		return consts.ReactDBUpdateError, errors.Wrap(err, "->LikeAction VideoLikeCount down error")
	}
	go func() {
		asyncCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = s.likeRedis.VideoDislikeSRem(asyncCtx, userId, targetId)
		if err != nil {
			log.Println("LikeAction VideoDislikeSRem error:", err)
		}
	}()
	return consts.Success, nil
}

func (s *LikeRepo) LikeComment(ctx context.Context, userId string, targetId string, targetType string) (int32, error) {
	err := s.likeDb.LikeCreate(ctx, userId, targetId, targetType)
	if err != nil {
		return consts.ReactDBInsertError, errors.Wrap(err, "->LikeAction LikeCreate error")
	}
	err = s.commentDb.CommentLikeCountUp(ctx, targetId)
	if err != nil {
		return consts.ReactDBUpdateError, errors.Wrap(err, "->LikeAction CommentLikeCount up error")
	}
	return consts.Success, nil
}

func (s *LikeRepo) DislikeComment(ctx context.Context, userId string, targetId string, targetType string) (int32, error) {
	err := s.likeDb.LikeDelete(ctx, userId, targetId, targetType)
	if err != nil {
		return consts.ReactDBDeleteError, errors.Wrap(err, "->LikeAction LikeDelete error")
	}
	err = s.commentDb.CommentLikeCountDown(ctx, targetId)
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
			code, err := s.LikeComment(ctx, userId, targetId, targetType)
			return code, err
		case consts.ActionDislike:
			code, err := s.DislikeComment(ctx, userId, targetId, targetType)
			return code, err
		default:
			return consts.ReactReqValueError, errors.New("->LikeAction action type error")
		}
	default:
		return consts.ReactReqValueError, errors.New("->LikeAction targetType is not valid")
	}
}

func (s *LikeRepo) LikeList(ctx context.Context, userId string, pageNum int64, pageSize int64) (int32, []string, int64, error) {
	results, err := s.likeRedis.VideoLikeGet(ctx, userId)
	if err == nil && len(results) > 0 {
		start := pageNum * pageSize
		end := pageSize + start
		if start >= int64(len(results)) {
			return consts.ReactReqValueError, nil, 0, errors.New("pageNum out of range")
		}
		if end > int64(len(results)) {
			end = int64(len(results))
		}
		result := results[start:end]
		return consts.Success, result, int64(len(result)), nil
	}
	videoIds, err := s.videoDb.LikeVideoIds(ctx, userId, pageNum, pageSize)
	if err != nil {
		return consts.ReactDBSelectError, nil, 0, errors.Wrap(err, "->LikeList select LikeVideo error")
	}
	go func() {
		asyncCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		for _, id := range videoIds {
			err = s.likeRedis.VideoLikeSAdd(asyncCtx, userId, id)
			if err != nil {
				log.Println("LikeAction LikeSAdd error:", err)
			}
		}
	}()
	return consts.Success, videoIds, int64(len(videoIds)), nil
}
