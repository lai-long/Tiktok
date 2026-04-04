package service

import (
	"Tiktok/biz/entity"

	"Tiktok/biz/model/video"
	"Tiktok/pkg/consts"
	"log"
	"strconv"
)

type LikeCommentDatabase interface {
	CommentLikeCountUp(commentId string) error
	CommentLikeCountDown(commentId string) error
}
type LikeVideoDatabase interface {
	VideoLikeCountUp(videoId string) error
	VideoLikeCountDown(videoId string) error
	LikeVideoIds(userId string, pageNum int, pageSize int) (error, []string)
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

func (s *LikeService) LikeAction(userId string, targetId string, action string, targetType string) (int, string) {
	//target type 1视频 2评论
	if targetType == "1" {
		if action == "1" {
			err := s.likeDb.LikeCreate(userId, targetId, targetType)
			if err != nil {
				return consts.CodeDBCreateError, "LikeAction LikeCreate error"
			}
			err = s.videoDb.VideoLikeCountUp(targetId)
			if err != nil {
				return consts.CodeDBUpdateError, "LikeAction LikeCountUp error"
			}
			return consts.CodeSuccess, "LikeAction success"
		}
		if action == "2" {
			err := s.likeDb.LikeDelete(userId, targetId, targetType)
			if err != nil {
				log.Println(err)
				return consts.CodeDBDeleteError, "VideoLikeAction LikeDelete error"
			}
			err = s.videoDb.VideoLikeCountDown(targetId)
			if err != nil {
				log.Println(err)
				return consts.CodeDBUpdateError, "VideoLikeAction LikeCountDown error"
			}
			return consts.CodeSuccess, "VideoLikeAction success"
		}
		return consts.CodeLikeError, "VideoLikeAction action num error"
	} else if targetType == "2" {
		if action == "1" {
			err := s.likeDb.LikeCreate(userId, targetId, targetType)
			if err != nil {
				log.Println("CommentLikeAction LikeCreate error", err)
				return consts.CodeDBCreateError, "CommentLikeAction CommentCreate error"
			}
			err = s.commentDb.CommentLikeCountUp(targetId)
			if err != nil {
				log.Println("CommentLikeAction CommentCountUp error", err)
				return consts.CodeDBUpdateError, "CommentLikeAction CommentCountUp error"
			}
		}
		if action == "2" {
			err := s.likeDb.LikeDelete(userId, targetId, targetType)
			if err != nil {
				log.Println("CommentLikeAction CommentDelete error", err)
				return consts.CodeDBDeleteError, "CommentLikeAction CommentDelete error"
			}
			err = s.commentDb.CommentLikeCountDown(targetId)
			if err != nil {
				log.Println("CommentLikeAction CommentCountDown error", err)
				return consts.CodeDBUpdateError, "CommentLikeAction CommentCountDown error"
			}
		}
		return consts.CodeSuccess, "CommentLikeAction action num error"
	}
	return consts.CodeLikeError, "CommentLikeAction target type num error"
}

func (s *LikeService) LikeList(userId string, pageNum string, pageSize string) (int, string, []*video.VideoInfo, bool) {
	pageNumInt, err := strconv.Atoi(pageNum)
	if err != nil {
		log.Printf("LikeList pageNum strconv error : %v", err)
		return consts.CodeError, "LikeList pageNum strconv error", nil, false
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		log.Printf("LikeList pageSize strconv error : %v", err)
		return consts.CodeError, "LikeList pageSize strconv error", nil, false
	}
	err, videoId := s.videoDb.LikeVideoIds(userId, pageNumInt, pageSizeInt)
	if err != nil {
		log.Printf("LikeList err : %v", err)
		return consts.CodeDBSelectError, "LikeList db.LikeVideoIds error", nil, false
	}
	ok, videos := s.videoDb.LikeVideos(videoId)
	if !ok {
		return consts.CodeDBSelectError, "LikeList db.LikeVideos error", nil, false
	}
	var videoInfos []*video.VideoInfo
	for _, v := range videos {
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
	return consts.CodeSuccess, "LikeList success", videoInfos, true
}
