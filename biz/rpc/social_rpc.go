package rpc

import (
	user "Tiktok/biz/model/user"
	"Tiktok/kitex_gen/social"
	"Tiktok/kitex_gen/social/socialservice"
	"Tiktok/pkg/consts"
	"context"
	"log"

	"github.com/cloudwego/kitex/client"
	etcd "github.com/kitex-contrib/registry-etcd"
)

func InitSocialRpc() {
	r, err := etcd.NewEtcdResolver([]string{"127.0.0.1:2379"})
	if err != nil {
		log.Fatal(err)
	}
	cli, err := socialservice.NewClient("socialService", client.WithResolver(r))
	if err != nil {
		log.Fatal(err)
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
	items := make([]*user.UserInfo, len(resp.Data.Items))
	for i, item := range resp.Data.Items {
		items[i] = &user.UserInfo{
			ID:        item.ID,
			Username:  item.Username,
			AvatarURL: item.AvatarURL,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
			DeletedAt: item.DeletedAt,
		}
	}
	data := &SocialData{
		Items: items,
		Total: resp.Data.Total,
	}
	return resp.Code, data, nil
}

func FollowerListRpc(ctx context.Context, req *social.FollowerListReq) (int32, *SocialData, error) {
	resp, err := socialClient.FollowerList(ctx, req)
	if err != nil || resp.Data == nil {
		return consts.SocialDBSelectError, nil, err
	}
	items := make([]*user.UserInfo, len(resp.Data.Items))
	for i, item := range resp.Data.Items {
		items[i] = &user.UserInfo{
			ID:        item.ID,
			Username:  item.Username,
			AvatarURL: item.AvatarURL,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
			DeletedAt: item.DeletedAt,
		}
	}
	data := &SocialData{
		Items: items,
		Total: resp.Data.Total,
	}
	return resp.Code, data, nil
}

func FriendListRpc(ctx context.Context, req *social.FriendListReq) (int32, *SocialData, error) {
	resp, err := socialClient.FriendList(ctx, req)
	if err != nil || resp.Data == nil {
		return consts.SocialDBSelectError, nil, err
	}
	items := make([]*user.UserInfo, len(resp.Data.Items))
	for i, item := range resp.Data.Items {
		items[i] = &user.UserInfo{
			ID:        item.ID,
			Username:  item.Username,
			AvatarURL: item.AvatarURL,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
			DeletedAt: item.DeletedAt,
		}
	}
	data := &SocialData{
		Items: items,
		Total: resp.Data.Total,
	}
	return resp.Code, data, nil
}

type SocialData struct {
	Items []*user.UserInfo
	Total int64
}
