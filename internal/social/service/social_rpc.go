package service

import (
	"Tiktok/kitex_gen/social"
	"Tiktok/pkg/consts"
	"Tiktok/pkg/entity"
	"context"

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
	socialDb SocialDatabase
}

func NewSocialRepo(socialDb SocialDatabase) *SocialRepo {
	return &SocialRepo{socialDb: socialDb}
}
func (s *SocialRepo) RelationAction(ctx context.Context, toUserId string, actionType string, userId string) (int32, error) {
	if actionType == consts.ActionFollow {
		err := s.socialDb.CreateFollowing(userId, toUserId)
		if err != nil {
			return consts.SocialDBInsertError, errors.Wrap(err, "->RelationAction CreateFollowing err")
		}
		return consts.Success, nil
	}
	if actionType == consts.ActionUnfollow {
		err := s.socialDb.DeleteFollowing(userId, toUserId)
		if err != nil {
			return consts.SocialDBDeleteError, errors.Wrap(err, "->RelationACtion DeleteFollowing err")
		}
		return consts.Success, nil
	}
	return consts.SocialReqValueError, nil
}

func (s *SocialRepo) FollowingList(userId string, pageNum int64, pageSize int64) (int32, []*social.UserInfo, error) {
	entities, err := s.socialDb.FollowingList(userId, pageNum, pageSize)
	if err != nil {
		return consts.SocialDBSelectError, nil, errors.Wrap(err, "->FollowingList Get Following List err")
	}
	return consts.Success, toUserInfoList(entities), nil
}

func (s *SocialRepo) FollowerList(userId string, pageNum int64, pageSize int64) (int32, []*social.UserInfo, error) {
	entities, err := s.socialDb.FollowerList(userId, pageNum, pageSize)
	if err != nil {
		return consts.SocialDBSelectError, nil, errors.Wrap(err, "->FollowerList Get List err")
	}
	return consts.Success, toUserInfoList(entities), nil
}

func (s *SocialRepo) FriendList(userId string, pageNum int64, pageSize int64) (int32, []*social.UserInfo, error) {
	entityFriend, ok := s.socialDb.FriendList(userId, pageNum, pageSize)
	if !ok {
		return consts.SocialDBSelectError, nil, errors.New("->FriendList Get List err")
	}
	return consts.Success, toUserInfoList(entityFriend), nil
}

func toUserInfoList(entities []entity.UserEntity) []*social.UserInfo {
	userInfos := make([]*social.UserInfo, 0, len(entities))
	for _, e := range entities {
		userInfos = append(userInfos, toSocialUserInfo(e))
	}
	return userInfos
}

func toSocialUserInfo(e entity.UserEntity) *social.UserInfo {
	return &social.UserInfo{
		ID:        e.ID,
		Username:  e.Username,
		AvatarURL: e.Avatar_url,
		CreatedAt: e.Created_at.String(),
		UpdatedAt: e.Updated_at.String(),
	}
}
