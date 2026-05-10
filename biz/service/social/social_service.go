package social

import (
	userModel "Tiktok/biz/model/user"

	user "Tiktok/internal/user/service"
	"Tiktok/pkg/consts"
	"Tiktok/pkg/entity"

	"github.com/pkg/errors"
)

type SocialDatabase interface {
	CreateFollowing(userId string, toUserId string) error
	DeleteFollowing(userId string, toUserId string) error
	FollowingList(userId string, pageNum int64, pageSize int64) ([]entity.UserEntity, error)
	FollowerList(userId string, pageNum int64, pageSize int64) ([]entity.UserEntity, error)
	FriendList(userId string, pageNum int64, pageSize int64) ([]entity.UserEntity, bool)
}
type SocialRepo struct {
	social SocialDatabase
	user   user.UserDatabase
}

func NewSocialRepo(social SocialDatabase, userDb user.UserDatabase) *SocialRepo {
	return &SocialRepo{
		social: social,
		user:   userDb,
	}
}
func (s *SocialRepo) RelationAction(toUserId string, actionType string, userId string) (int32, error) {
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

func (s *SocialRepo) FollowingList(userId string, pageNum int64, pageSize int64) (int32, []*userModel.UserInfo, error) {
	followings, err := s.social.FollowingList(userId, pageNum, pageSize)
	if err != nil {
		return consts.SocialDBSelectError, nil, errors.Wrap(err, "->Following List Get Following List err")
	}
	userInfos := []*userModel.UserInfo{}
	for i := 0; i < len(followings); i++ {
		user2 := &userModel.UserInfo{
			ID:        followings[i].ToUserInfo().ID,
			Username:  followings[i].ToUserInfo().Username,
			AvatarURL: followings[i].ToUserInfo().AvatarURL,
			CreatedAt: followings[i].ToUserInfo().CreatedAt,
			UpdatedAt: followings[i].ToUserInfo().UpdatedAt,
			DeletedAt: followings[i].ToUserInfo().DeletedAt,
		}
		userInfos = append(userInfos, user2)
	}
	return consts.Success, userInfos, nil
}

func (s *SocialRepo) FollowerList(userId string, pageNum int64, pageSize int64) (int32, []*userModel.UserInfo, error) {
	followers, err := s.social.FollowerList(userId, pageNum, pageSize)
	if err != nil {
		return consts.SocialDBSelectError, nil, errors.Wrap(err, "->FollowerList Get List err")
	}
	userInfos := []*userModel.UserInfo{}
	for i := 0; i < len(followers); i++ {
		user2 := &userModel.UserInfo{
			ID:        followers[i].ToUserInfo().ID,
			Username:  followers[i].ToUserInfo().Username,
			AvatarURL: followers[i].ToUserInfo().AvatarURL,
			CreatedAt: followers[i].ToUserInfo().CreatedAt,
			UpdatedAt: followers[i].ToUserInfo().UpdatedAt,
			DeletedAt: followers[i].ToUserInfo().DeletedAt,
		}
		userInfos = append(userInfos, user2)
	}
	return consts.Success, userInfos, nil
}

func (s *SocialRepo) FriendList(userId string, pageNum int64, pageSize int64) (int32, []*userModel.UserInfo, error) {
	entityFriend, ok := s.social.FriendList(userId, pageNum, pageSize)
	if !ok {
		return consts.SocialDBSelectError, nil, errors.New("->FriendList Get List err")
	}
	userInfos := []*userModel.UserInfo{}
	for i := range entityFriend {
		user2 := &userModel.UserInfo{
			ID:        entityFriend[i].ToUserInfo().ID,
			Username:  entityFriend[i].ToUserInfo().Username,
			AvatarURL: entityFriend[i].ToUserInfo().AvatarURL,
			CreatedAt: entityFriend[i].ToUserInfo().CreatedAt,
			UpdatedAt: entityFriend[i].ToUserInfo().UpdatedAt,
			DeletedAt: entityFriend[i].ToUserInfo().DeletedAt,
		}
		userInfos = append(userInfos, user2)
	}
	return consts.Success, userInfos, nil
}
