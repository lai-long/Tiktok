package social

import (
	"Tiktok/biz/entity"
	userModel "Tiktok/biz/model/user"
	"Tiktok/biz/service/user"
	"Tiktok/pkg/consts"

	"github.com/pkg/errors"
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
	user   user.UserDatabase
}

func NewSocialService(social SocialDatabase, userDb user.UserDatabase) *SocialService {
	return &SocialService{
		social: social,
		user:   userDb,
	}
}
func (s *SocialService) RelationAction(toUserId string, actionType string, userId string) (int32, error) {
	if actionType == "0" {
		err := s.social.CreateFollowing(userId, toUserId)
		if err != nil {
			return consts.SocialDBInsertError, errors.Wrap(err, "->RelationAction CreateFollowing err")
		}
		return consts.Success, nil
	}
	if actionType == "1" {
		err := s.social.DeleteFollowing(userId, toUserId)
		if err != nil {
			return consts.SocialDBDeleteError, errors.Wrap(err, "->RelationACtion DeleteFollowing err")
		}
		return consts.Success, nil
	}
	return consts.SocialReqValueError, nil
}

func (s *SocialService) FollowingList(userId string, pageNum int64, pageSize int64) (int32, error, []*userModel.UserInfo) {
	followings, err := s.social.FollowingList(userId, pageNum, pageSize)
	if err != nil {
		return consts.SocialDBSelectError, errors.Wrap(err, "->Following List Get Following List err"), nil
	}
	userInfos := []*userModel.UserInfo{}
	for i := 0; i < len(followings); i++ {
		userInfos = append(userInfos, followings[i].ToUserInfo())
	}
	return consts.Success, nil, userInfos
}

func (s *SocialService) FollowerList(userId string, pageNum int64, pageSize int64) (int32, error, []*userModel.UserInfo) {
	followers, err := s.social.FollowerList(userId, pageNum, pageSize)
	if err != nil {
		return consts.SocialDBSelectError, errors.Wrap(err, "->FollowerList Get List err"), nil
	}
	userInfos := []*userModel.UserInfo{}
	for i := 0; i < len(followers); i++ {
		userInfos = append(userInfos, followers[i].ToUserInfo())
	}
	return consts.Success, nil, userInfos
}

func (s *SocialService) FriendList(userId string, pageNum int64, pageSize int64) (int32, error, []*userModel.UserInfo) {
	entityFriend, ok := s.social.FriendList(userId, pageNum, pageSize)
	if !ok {
		return consts.SocialDBSelectError, errors.New("->FriendList Get List err"), nil
	}
	userInfos := []*userModel.UserInfo{}
	for i, _ := range entityFriend {
		userInfos = append(userInfos, entityFriend[i].ToUserInfo())
	}
	return consts.Success, nil, userInfos
}
