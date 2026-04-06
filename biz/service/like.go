package service

import (
	"Tiktok/biz/entity"

	"Tiktok/biz/model/video"
	"Tiktok/pkg/consts"

	"github.com/pkg/errors"
)

type LikeCommentDatabase interface {
	CommentLikeCountUp(commentId string) error
	CommentLikeCountDown(commentId string) error
}
type LikeVideoDatabase interface {
	VideoLikeCountUp(videoId string) error
	VideoLikeCountDown(videoId string) error
	LikeVideoIds(userId string, pageNum int64, pageSize int64) (error, []string)
	LikeVideos(videoId []string) (bool, []entity.VideoEntity)
}
type LikeDatabase interface {
	LikeCreate(userId string, targetId string, targetType string) error
	LikeDelete(userId, targetId, targetType string) error
}
type LikeService struct {
	videoDb   LikeVideoDatabase
	commentDb LikeCommentDatabase
	likeDb    LikeDatabase
}

func NewLikeVideoService(videoDb LikeVideoDatabase, commentDb LikeCommentDatabase, likeDb LikeDatabase) *LikeService {
	return &LikeService{
		videoDb:   videoDb,
		commentDb: commentDb,
		likeDb:    likeDb,
	}
}

func (s *LikeService) LikeAction(userId string, targetId string, action string, targetType string) (int32, error) {
	//target type 1视频 2评论
	if targetType == "1" {
		if action == "1" {
			err := s.likeDb.LikeCreate(userId, targetId, targetType)
			if err != nil {
				return consts.ReactDBInsertError, errors.Wrap(err, "->LikeAction LikeCreate error")
			}
			err = s.videoDb.VideoLikeCountUp(targetId)
			if err != nil {
				return consts.ReactDBUpdateError, errors.Wrap(err, "->LikeAction VideoLikeCount up error")
			}
			return consts.Success, nil
		}
		if action == "2" {
			err := s.likeDb.LikeDelete(userId, targetId, targetType)
			if err != nil {
				return consts.ReactDBDeleteError, errors.Wrap(err, "->LikeAction LikeDelete error")
			}
			err = s.videoDb.VideoLikeCountDown(targetId)
			if err != nil {
				return consts.ReactDBUpdateError, errors.Wrap(err, "->LikeAction VideoLikeCount down error")
			}
			return consts.Success, nil
		}
		return consts.ReactReqValueError, nil
	} else if targetType == "2" {
		if action == "1" {
			err := s.likeDb.LikeCreate(userId, targetId, targetType)
			if err != nil {
				return consts.ReactDBInsertError, errors.Wrap(err, "->LikeAction LikeCreate error")
			}
			err = s.commentDb.CommentLikeCountUp(targetId)
			if err != nil {
				return consts.ReactDBUpdateError, errors.Wrap(err, "->LikeAction CommentLikeCount up error")
			}
		}
		if action == "2" {
			err := s.likeDb.LikeDelete(userId, targetId, targetType)
			if err != nil {
				return consts.ReactDBDeleteError, errors.Wrap(err, "->LikeAction LikeDelete error")
			}
			err = s.commentDb.CommentLikeCountDown(targetId)
			if err != nil {
				return consts.ReactDBUpdateError, errors.Wrap(err, "->LikeAction CommentLikeCount down error")
			}
		}
		return consts.Success, nil
	}
	return consts.ReactReqValueError, nil
}

func (s *LikeService) LikeList(userId string, pageNum int64, pageSize int64) (int32, error, []*video.VideoInfo) {
	err, videoId := s.videoDb.LikeVideoIds(userId, pageNum, pageSize)
	if err != nil {
		return consts.ReactDBSelectError, errors.Wrap(err, "->LikeList select LikeVideo error"), nil
	}
	ok, videos := s.videoDb.LikeVideos(videoId)
	if !ok {
		return consts.ReactDBSelectError, errors.New("->LikeList LikeVideos err"), nil
	}
	var videoInfos []*video.VideoInfo
	for _, v := range videos {
		videoInfos = append(videoInfos, v.ToVideoInfo())
	}
	return consts.Success, nil, videoInfos
}
