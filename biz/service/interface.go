package service

import (
	"Tiktok/biz/dao/db"
	"Tiktok/biz/model/entity"
)

type Database interface {
	CreateUser(user entity.UserEntity) error
	GetUserByUsername(username string) (entity.UserEntity, error)
	GetUserByUserId(userId string) (entity.UserEntity, error)
	UpdateUserAvatar(url string, userId interface{}) error
	CreatVideo(entity entity.VideoEntity) error
	GetVideoByUserID(userId string, pageSize int, pageNum int) ([]entity.VideoEntity, error)
	GetVideoByKeyWord(keyword string, pageNum int, pageSize int) ([]entity.VideoEntity, error)
	GetVideoByVideoId(videoId string) (entity.VideoEntity, error)
	CreateFollowing(userId string, toUserId string) error
	CreateFollower(userId string, toUserId string) error
	DeleteFollowing(userId string, toUserId string) error
	DeleteFollower(userId string, toUserId string) error
	FollowingIdList(userId string, pageNum int, pageSize int) ([]string, error)
	FollowerIdList(userId string, pageNum int, pageSize int) ([]string, error)
	FriendIdList(userId string, pageNum, pageSize int) (followingIds []string, followerIds []string, err1 error, err2 error)
	VideoLikeCountUp(videoId string) error
	CommentLikeCountUp(commentId string) error
	VideoLikeCreate(userId string, videoId string) error
	CommentLikeCreate(userId string, commentId string) error
	VideoLikeCountDown(videoId string) error
	CommentLikeCountDown(commentId string) error
	VideoLikeDelete(userId string, videoId string) error
	CommentLikeDelete(userId string, commentId string) error
	LikeVideoIds(userId string, pageNum int, pageSize int) (error, []string)
	LikeVideos(videoId []string) (bool, []entity.VideoEntity)
	CreateComment(commentId string, videoId string, userId string, content string) error
	GetComments(videoId string, pageNum int, pageSize int) (error, []entity.CommentEntity)
	CommentDelete(videoId string, commentId string) error
	GetCommentById(commentId string) (entity.CommentEntity, error)
	CommentCountUp(videoId string) error
	CommentCountDown(videoId string) error
	SaveMfaSecret(mfa string, userId string) error
	GetMfaSecret(userId string) (string, error)
	MfaBindUpdate(userId string) error
	CheckMfaBind(userId string) (error, int)
}
type UserDatabase interface {
	CreateUser(user entity.UserEntity) error
	GetUserByUsername(username string) (entity.UserEntity, error)
	GetUserByUserId(userId string) (entity.UserEntity, error)
	UpdateUserAvatar(url string, userId interface{}) error
}
type MfaDatabase interface {
	SaveMfaSecret(mfa string, userId string) error
	GetMfaSecret(userId string) (string, error)
	MfaBindUpdate(userId string) error
	CheckMfaBind(userId string) (error, int)
}
type Service struct {
	db Database
}

func NewService(db *db.MySQLdb) *Service {
	return &Service{db: db}
}

type UserService struct {
	userDb UserDatabase
	mfaDb  MfaDatabase
}

func NewUserService(userDb UserDatabase, mfaDb MfaDatabase) *UserService {
	return &UserService{userDb: userDb, mfaDb: mfaDb}
}
