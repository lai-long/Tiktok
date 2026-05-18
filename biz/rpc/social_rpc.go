package rpc

import (
	user "Tiktok/biz/model/user"
	"Tiktok/kitex_gen/social"
	"Tiktok/kitex_gen/social/socialservice"
	"Tiktok/pkg/consts"
	"Tiktok/pkg/logger"
	"context"

	"github.com/cloudwego/kitex/client"
	etcd "github.com/kitex-contrib/registry-etcd"
	"go.uber.org/zap"
)

func InitSocialRpc() {
	r, err := etcd.NewEtcdResolver([]string{"127.0.0.1:2379"})
	if err != nil {
		logger.Fatal("etcd resolver error", zap.Error(err))
	}
	cli, err := socialservice.NewClient("socialService", client.WithResolver(r))
	if err != nil {
		logger.Fatal("social service client error", zap.Error(err))
	}
	socialClient = cli
}

func RelationActionRpc(ctx context.Context, req *social.RelationActionReq) (int32, error) {
	resp, err := socialClient.RelationAction(ctx, req)
	if err != nil {
		return consts.SocialReqValueError, err
	}
	return resp.Code, nil
}

func FollowingListRpc(ctx context.Context, req *social.FollowingListReq) (int32, *SocialData, error) {
	resp, err := socialClient.FollowingList(ctx, req)
	if err != nil || resp.Data == nil {
		return consts.SocialDBSelectError, nil, err
	}
	return resp.Code, buildSocialData(resp.Data), nil
}

func FollowerListRpc(ctx context.Context, req *social.FollowerListReq) (int32, *SocialData, error) {
	resp, err := socialClient.FollowerList(ctx, req)
	if err != nil || resp.Data == nil {
		return consts.SocialDBSelectError, nil, err
	}
	return resp.Code, buildSocialData(resp.Data), nil
}

func FriendListRpc(ctx context.Context, req *social.FriendListReq) (int32, *SocialData, error) {
	resp, err := socialClient.FriendList(ctx, req)
	if err != nil || resp.Data == nil {
		return consts.SocialDBSelectError, nil, err
	}
	return resp.Code, buildSocialData(resp.Data), nil
}

func buildSocialData(data *social.SocialData) *SocialData {
	items := make([]*user.UserInfo, len(data.Items))
	for i, item := range data.Items {
		items[i] = &user.UserInfo{
			ID:        item.ID,
			Username:  item.Username,
			AvatarURL: item.AvatarURL,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
			DeletedAt: item.DeletedAt,
		}
	}
	return &SocialData{Items: items, Total: data.Total}
}

type SocialData struct {
	Items []*user.UserInfo
	Total int64
}
