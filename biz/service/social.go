package service

import (
	"Tiktok/biz/entity"
	"Tiktok/biz/model/user"
	"Tiktok/pkg/consts"
	"log"
)

type SocialDatabase interface {
	CreateFollowing(userId string, toUserId string) error
	DeleteFollowing(userId string, toUserId string) error
	FollowingList(userId string, pageNum int64, pageSize int64) ([]entity.UserEntity, error)
	FollowerList(userId string, pageNum int64, pageSize int64) ([]entity.UserEntity, error)
	FriendList(userId string, pageNum int64, pageSize int64) ([]entity.UserEntity, bool)
}
type SocialService struct {
	social SocialDatabase
	user   UserDatabase
}

func NewSocialService(social SocialDatabase, user UserDatabase) *SocialService {
	return &SocialService{
		social: social,
		user:   user,
	}
}
func (s *SocialService) RelationAction(toUserId string, actionType string, userId string) (int, string) {
	if actionType == "0" {
		err := s.social.CreateFollowing(userId, toUserId)
		if err != nil {
			log.Println(err)
			return consts.CodeDBCreateError, "RelationAction CreateFollowing error"
		}
		return consts.CodeSuccess, "RelationAction follow success"
	}
	if actionType == "1" {
		err := s.social.DeleteFollowing(userId, toUserId)
		if err != nil {
			return consts.CodeDBDeleteError, "RelationAction DeleteFollowing error"
		}
		return consts.CodeSuccess, "RelationAction delete follow success"
	}
	return consts.CodeRelationError, "RelationAction actionType error"
}

func (s *SocialService) FollowingList(userId string, pageNum int64, pageSize int64) (int, string, []*user.UserInfo, bool) {
	followings, err := s.social.FollowingList(userId, pageNum, pageSize)
	if err != nil {
		log.Println(err)
		return consts.CodeDBSelectError, "FollowingList  db.FollowingIdList error", nil, false
	}
	userInfos := []*user.UserInfo{}
	for i := 0; i < len(followings); i++ {
		userInfos = append(userInfos, followings[i].ToUserInfo())
	}
	return consts.CodeSuccess, "FollowingList success", userInfos, true
}

func (s *SocialService) FollowerList(userId string, pageNum int64, pageSize int64) (int, string, []*user.UserInfo, bool) {
	followers, err := s.social.FollowerList(userId, pageNum, pageSize)
	if err != nil {
		log.Printf("FollowerList db.FollowerIdList error: %v", err)
		return consts.CodeDBSelectError, "FollowerList db.FollowerIdList error", nil, false
	}
	userInfos := []*user.UserInfo{}
	for i := 0; i < len(followers); i++ {
		userInfos = append(userInfos, followers[i].ToUserInfo())
	}
	return consts.CodeSuccess, "FollowerList success", userInfos, true
}

func (s *SocialService) FriendList(userId string, pageNum int64, pageSize int64) (int, string, []*user.UserInfo, bool) {
	entityFriend, ok := s.social.FriendList(userId, pageNum, pageSize)
	if !ok {
		return consts.CodeDBSelectError, "FriendList db.FriendList error", nil, false
	}
	userInfos := []*user.UserInfo{}
	for i, _ := range entityFriend {
		userInfos = append(userInfos, entityFriend[i].ToUserInfo())
	}
	return consts.CodeSuccess, "FollowerList success", userInfos, true
}
