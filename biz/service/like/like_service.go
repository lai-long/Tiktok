package like

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
}

func NewLikeRepo(videoDb LikeVideoDatabase, commentDb LikeCommentDatabase, likeDb LikeDatabase) *LikeRepo {
	return &LikeRepo{
		videoDb:   videoDb,
		commentDb: commentDb,
		likeDb:    likeDb,
	}
}

func (s *LikeRepo) LikeAction(userId string, targetId string, action string, targetType string) (int32, error) {
	switch targetType {
	case "1":
		switch action {
		case "1":
			err := s.likeDb.LikeCreate(userId, targetId, targetType)
			if err != nil {
				return consts.ReactDBInsertError, errors.Wrap(err, "->LikeAction LikeCreate error")
			}
			err = s.videoDb.VideoLikeCountUp(targetId)
			if err != nil {
				return consts.ReactDBUpdateError, errors.Wrap(err, "->LikeAction VideoLikeCount up error")
			}
			return consts.Success, nil
		case "2":
			err := s.likeDb.LikeDelete(userId, targetId, targetType)
			if err != nil {
				return consts.ReactDBDeleteError, errors.Wrap(err, "->LikeAction LikeDelete error")
			}
			err = s.videoDb.VideoLikeCountDown(targetId)
			if err != nil {
				return consts.ReactDBUpdateError, errors.Wrap(err, "->LikeAction VideoLikeCount down error")
			}
			return consts.Success, nil
		default:
			return consts.ReactReqValueError, errors.Errorf("invalid action type: %s", action)
		}
	case "2":
		switch action {
		case "1":
			err := s.likeDb.LikeCreate(userId, targetId, targetType)
			if err != nil {
				return consts.ReactDBInsertError, errors.Wrap(err, "->LikeAction LikeCreate error")
			}
			err = s.commentDb.CommentLikeCountUp(targetId)
			if err != nil {
				return consts.ReactDBUpdateError, errors.Wrap(err, "->LikeAction CommentLikeCount up error")
			}
		case "2":
			err := s.likeDb.LikeDelete(userId, targetId, targetType)
			if err != nil {
				return consts.ReactDBDeleteError, errors.Wrap(err, "->LikeAction LikeDelete error")
			}
			err = s.commentDb.CommentLikeCountDown(targetId)
			if err != nil {
				return consts.ReactDBUpdateError, errors.Wrap(err, "->LikeAction CommentLikeCount down error")
			}
		default:
			return consts.ReactReqValueError, errors.New("->LikeAction action type error")
		}
	default:
		return consts.ReactReqValueError, errors.New("->LikeAction targetType is not valid")
	}
	return consts.ReactReqValueError, nil
}

func (s *LikeRepo) LikeList(userId string, pageNum int64, pageSize int64) (int32, []*video.VideoInfo, error) {
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
	return consts.Success, videoInfos, nil
}
