package service

import (
	"Tiktok/biz/model/dto"
	"Tiktok/biz/model/entity"
	"Tiktok/pkg/consts"
	"log"
	"strconv"
)

type SocialDatabase interface {
	CreateFollowing(userId string, toUserId string) error
	DeleteFollowing(userId string, toUserId string) error
	FollowingList(userId string, pageNum int, pageSize int) ([]entity.UserEntity, error)
	FollowerList(userId string, pageNum int, pageSize int) ([]entity.UserEntity, error)
	FriendList(userId string, pageNum int, pageSize int) ([]entity.UserEntity, bool)
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
func (s *SocialService) FollowingList(userId string, pageNum string, pageSize string) (int, string, []dto.User, bool) {
	pageNumInt := 0
	pageSizeInt := 10
	pageNumInt, err := strconv.Atoi(pageNum)
	if err != nil {
		return consts.CodeError, "FollowingList PageNum strconv error", []dto.User{}, false
	}
	pageSizeInt, err = strconv.Atoi(pageSize)
	if err != nil {
		return consts.CodeError, "FollowingList PageSize strconv error", []dto.User{}, false
	}
	followings, err := s.social.FollowingList(userId, pageNumInt, pageSizeInt)
	if err != nil {
		log.Println(err)
		return consts.CodeDBSelectError, "FollowingList  db.FollowingIdList error", []dto.User{}, false
	}
	dtoFollowings := make([]dto.User, len(followings))
	for i := 0; i < len(dtoFollowings); i++ {
		dtoFollowings[i].ID = followings[i].Id
		dtoFollowings[i].Username = followings[i].Username
		dtoFollowings[i].AvatarURL = followings[i].Avatar_url
	}
	return consts.CodeSuccess, "FollowingList success", dtoFollowings, true
}
func (s *SocialService) FollowerList(userId string, pageNum string, pageSize string) (int, string, []dto.User, bool) {
	pageNumInt := 0
	pageSizeInt := 10
	pageNumInt, err := strconv.Atoi(pageNum)
	if err != nil {
		log.Printf("FollowerList PageNum strconv error: %v", err)
		return consts.CodeError, "FollowerList PageNum strconv error", []dto.User{}, false
	}
	pageSizeInt, err = strconv.Atoi(pageSize)
	if err != nil {
		log.Printf("FollowerList PageSize strconv error: %v", err)
		return consts.CodeError, "FollowerList PageSize strconv error", []dto.User{}, false
	}
	followers, err := s.social.FollowerList(userId, pageNumInt, pageSizeInt)
	if err != nil {
		log.Printf("FollowerList db.FollowerIdList error: %v", err)
		return consts.CodeDBSelectError, "FollowerList db.FollowerIdList error", []dto.User{}, false
	}

	dtoFollowers := make([]dto.User, len(followers))
	for i := 0; i < len(followers); i++ {
		dtoFollowers[i].ID = followers[i].Id
		dtoFollowers[i].Username = followers[i].Username
		dtoFollowers[i].AvatarURL = followers[i].Avatar_url
	}
	return consts.CodeSuccess, "FollowerList success", dtoFollowers, true
}
func (s *SocialService) FriendList(userId string, pageNum string, pageSize string) (int, string, []dto.User, bool) {
	pageNumInt := 0
	pageSizeInt := 10
	pageNumInt, err := strconv.Atoi(pageNum)
	if err != nil {
		log.Printf("FriendList PageNum strconv error: %v", err)
		return consts.CodeError, "FollowerList PageNum strconv error", []dto.User{}, false
	}
	pageSizeInt, err = strconv.Atoi(pageSize)
	if err != nil {
		log.Printf("FriendList PageSize strconv error: %v", err)
		return consts.CodeError, "FollowerList PageSize strconv error", []dto.User{}, false
	}
	entityFriend, ok := s.social.FriendList(userId, pageNumInt, pageSizeInt)
	if !ok {
		return consts.CodeDBSelectError, "FriendList db.FriendList error", []dto.User{}, false
	}
	friends := make([]dto.User, len(entityFriend))
	for i, _ := range entityFriend {
		friends[i].Username = entityFriend[i].Username
		friends[i].AvatarURL = entityFriend[i].Avatar_url
		friends[i].ID = entityFriend[i].Id
	}
	return consts.CodeSuccess, "FollowerList success", friends, true
}
