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
	CreateFollower(userId string, toUserId string) error
	DeleteFollowing(userId string, toUserId string) error
	DeleteFollower(userId string, toUserId string) error
	FollowingIdList(userId string, pageNum int, pageSize int) ([]string, error)
	FollowerIdList(userId string, pageNum int, pageSize int) ([]string, error)
	FriendIdList(userId string, pageNum, pageSize int) (followingIds []string, followerIds []string, err1 error, err2 error)
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
		err = s.social.CreateFollower(userId, toUserId)
		if err != nil {
			return consts.CodeDBCreateError, "RelationAction CreateFollower error"
		}
		return consts.CodeSuccess, "RelationAction follow success"
	}
	if actionType == "1" {
		err := s.social.DeleteFollowing(userId, toUserId)
		if err != nil {
			return consts.CodeDBDeleteError, "RelationAction DeleteFollowing error"
		}
		err = s.social.DeleteFollower(userId, toUserId)
		if err != nil {
			return consts.CodeDBDeleteError, "RelationAction DeleteFollower error"
		}
		return consts.CodeSuccess, "RelationAction delete follow success"
	}
	return consts.CodeRelationError, "RelationAction actionType error"
}
func (s *SocialService) FollowingList(userId string, pageNum string, pageSize string) (int, string, []dto.User, bool) {
	pageNumInt, err := strconv.Atoi(pageNum)
	if err != nil {
		return consts.CodeError, "FollowingList PageNum strconv error", []dto.User{}, false
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		return consts.CodeError, "FollowingList PageSize strconv error", []dto.User{}, false
	}
	followingIds, err := s.social.FollowingIdList(userId, pageNumInt, pageSizeInt)
	if err != nil {
		log.Println(err)
		return consts.CodeDBSelectError, "FollowingList  db.FollowingIdList error", []dto.User{}, false
	}
	followingUsers := make([]entity.UserEntity, len(followingIds))
	for i := 0; i < len(followingIds); i++ {
		followingUsers[i], err = s.user.GetUserByUserId(followingIds[i])
		if err != nil {
			return consts.CodeDBSelectError, "FollowingList  db.GetUserByUserId error", []dto.User{}, false
		}
	}
	dtoFollowings := make([]dto.User, len(followingUsers))
	for i := 0; i < len(followingIds); i++ {
		dtoFollowings[i].ID = followingIds[i]
		dtoFollowings[i].Username = followingUsers[i].Username
		dtoFollowings[i].AvatarURL = followingUsers[i].Avatar_url
	}
	return consts.CodeSuccess, "FollowingList success", dtoFollowings, true
}
func (s *SocialService) FollowerList(userId string, pageNum string, pageSize string) (int, string, []dto.User, bool) {
	pageNumInt, err := strconv.Atoi(pageNum)
	if err != nil {
		log.Printf("FollowerList PageNum strconv error: %v", err)
		return consts.CodeError, "FollowerList PageNum strconv error", []dto.User{}, false
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		log.Printf("FollowerList PageSize strconv error: %v", err)
		return consts.CodeError, "FollowerList PageSize strconv error", []dto.User{}, false
	}
	followerIds, err := s.social.FollowerIdList(userId, pageNumInt, pageSizeInt)
	if err != nil {
		log.Printf("FollowerList db.FollowerIdList error: %v", err)
		return consts.CodeDBSelectError, "FollowerList db.FollowerIdList error", []dto.User{}, false
	}
	followerUsers := make([]entity.UserEntity, len(followerIds))
	for i := 0; i < len(followerIds); i++ {
		followerUsers[i], err = s.user.GetUserByUserId(followerIds[i])
		if err != nil {
			log.Printf("FollowerList db.GetUserByUserId error: %v", err)
			return consts.CodeDBSelectError, "FollowerList db.GetUserByUserId error", []dto.User{}, false
		}
	}
	dtoFollowers := make([]dto.User, len(followerUsers))
	for i := 0; i < len(followerIds); i++ {
		dtoFollowers[i].ID = followerIds[i]
		dtoFollowers[i].Username = followerUsers[i].Username
		dtoFollowers[i].AvatarURL = followerUsers[i].Avatar_url
	}
	return consts.CodeSuccess, "FollowerList success", dtoFollowers, true
}
func (s *SocialService) FriendList(userId string, pageNum string, pageSize string) (int, string, []dto.User, bool) {
	pageNumInt, err := strconv.Atoi(pageNum)
	if err != nil {
		log.Printf("FriendList PageNum strconv error: %v", err)
		return consts.CodeError, "FollowerList PageNum strconv error", []dto.User{}, false
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		log.Printf("FriendList PageSize strconv error: %v", err)
		return consts.CodeError, "FollowerList PageSize strconv error", []dto.User{}, false
	}
	followings, followers, err1, err2 := s.social.FriendIdList(userId, pageNumInt, pageSizeInt)
	if err1 != nil || err2 != nil {
		log.Printf("FriendList db.FriendIdList error: %v and %v", err1, err2)
		return consts.CodeDBSelectError, "FollowerList db.FriendIdList error", []dto.User{}, false
	}
	dtoFollowers := make([]dto.User, len(followings)+len(followers))
	for i := 0; i < len(followings); i++ {
		follow, err := s.user.GetUserByUserId(followings[i])
		if err != nil {
			log.Printf("FriendList db.GetUserByUserId error: %v", err)
			return consts.CodeDBSelectError, "FollowerList followings db.GetUserByUserId error", []dto.User{}, false
		}
		dtoFollowers[i].ID = follow.Id
		dtoFollowers[i].Username = follow.Username
		dtoFollowers[i].AvatarURL = follow.Avatar_url
	}
	for i := 0; i < len(followers); i++ {
		follow, err := s.user.GetUserByUserId(followers[i])
		if err != nil {
			log.Printf("FriendList db.GetUserByUserId error: %v", err)
			return consts.CodeDBSelectError, "FollowerList followers db.GetUserByUserId error", []dto.User{}, false
		}
		dtoFollowers[i+len(followings)].ID = follow.Id
		dtoFollowers[i+len(followings)].Username = follow.Username
		dtoFollowers[i].AvatarURL = follow.Avatar_url
	}
	return consts.CodeSuccess, "FollowerList success", dtoFollowers, true
}
