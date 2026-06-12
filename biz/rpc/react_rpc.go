package rpc

import (
	"Tiktok/internal/config"
	"Tiktok/kitex_gen/react"
	"Tiktok/kitex_gen/react/commentservice"
	"Tiktok/kitex_gen/react/likeservice"
	"Tiktok/pkg/consts"
	"Tiktok/pkg/logger"
	"Tiktok/pkg/utils"
	"context"

	model "Tiktok/biz/model/react"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/transmeta"
	etcd "github.com/kitex-contrib/registry-etcd"
	"go.uber.org/zap"
)

func InitReactRpc() {
	r, err := etcd.NewEtcdResolver([]string{config.GetCfg().EtcdAddr})
	if err != nil {
		logger.Fatal("etcd resolver error", zap.Error(err))
	}
	commentCli, err := commentservice.NewClient("commentService", client.WithResolver(r), client.WithMetaHandler(transmeta.MetainfoClientHandler))
	if err != nil {
		logger.Fatal("comment service client error", zap.Error(err))
	}
	likeCli, err := likeservice.NewClient("likeService", client.WithResolver(r), client.WithMetaHandler(transmeta.MetainfoClientHandler))
	if err != nil {
		logger.Fatal("like service client error", zap.Error(err))
	}
	commentClient = commentCli
	likeClient = likeCli
}

func CommentPublish(ctx context.Context, req *react.CommentPublishReq) (int32, error) {
	defer utils.TrackRPC(ctx, "comment", "CommentPublish")()
	resp, err := commentClient.CommentPublish(ctx, req)
	if err != nil || resp == nil {
		return consts.ReactError, err
	}
	return resp.Code, nil
}

func CommentList(ctx context.Context, req *react.CommentListReq) (int32, *model.CommentData, error) {
	defer utils.TrackRPC(ctx, "comment", "CommentList")()
	resp, err := commentClient.CommentList(ctx, req)
	if err != nil || resp.Data == nil {
		return consts.ReactDBSelectError, nil, err
	}
	infos := make([]*model.CommentInfo, 0, resp.Data.Total)
	for i := 0; i < len(resp.Data.Items); i++ {
		infos = append(infos, &model.CommentInfo{
			UserId:       resp.Data.Items[i].UserId,
			TargetId:     resp.Data.Items[i].TargetId,
			CommentId:    resp.Data.Items[i].CommentId,
			Content:      resp.Data.Items[i].Content,
			LikeCount:    resp.Data.Items[i].LikeCount,
			CreatedAt:    resp.Data.Items[i].CreatedAt,
			UpdatedAt:    resp.Data.Items[i].UpdatedAt,
			DeletedAt:    resp.Data.Items[i].DeletedAt,
			TargetType:   resp.Data.Items[i].TargetType,
			CommentCount: resp.Data.Items[i].CommentCount,
		})
	}
	data := &model.CommentData{
		Items: infos,
		Total: int64(len(infos)),
	}
	return resp.Code, data, err
}

func CommentDelete(ctx context.Context, req *react.CommentDeleteReq) (int32, error) {
	defer utils.TrackRPC(ctx, "comment", "CommentDelete")()
	resp, err := commentClient.CommentDelete(ctx, req)
	if err != nil || resp == nil {
		return consts.ReactError, err
	}
	return resp.Code, nil
}

func LikeAction(ctx context.Context, req *react.LikeActionReq) (int32, error) {
	defer utils.TrackRPC(ctx, "like", "LikeAction")()
	resp, err := likeClient.LikeAction(ctx, req)
	if err != nil || resp == nil {
		return consts.ReactError, err
	}
	return resp.Code, nil
}

func LikeList(ctx context.Context, req *react.LikeListReq) (int32, *model.LikeVideoData, error) {
	defer utils.TrackRPC(ctx, "like", "LikeList")()
	resp, err := likeClient.LikeList(ctx, req)
	if err != nil || resp.Data == nil {
		return consts.ReactDBSelectError, nil, err
	}
	if len(resp.Data.VideoIds) == 0 {
		return resp.Code, &model.LikeVideoData{Total: 0}, nil
	}
	_, videoInfos, err := BatchGetVideo(ctx, resp.Data.VideoIds)
	if err != nil {
		return consts.VideoDBSelectError, nil, err
	}
	data := &model.LikeVideoData{
		Items: videoInfos,
		Total: int64(len(videoInfos)),
	}
	return resp.Code, data, nil
}
