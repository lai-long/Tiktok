package rpc

import (
	"context"

	"Tiktok/internal/config"
	"Tiktok/kitex_gen/video"
	"Tiktok/kitex_gen/video/videoservice"
	"Tiktok/pkg/consts"
	"Tiktok/pkg/logger"
	"Tiktok/pkg/utils"

	model "Tiktok/biz/model/video"

	"github.com/cloudwego/kitex/client"
	etcd "github.com/kitex-contrib/registry-etcd"
	"go.uber.org/zap"
)

func InitVideoRpc() {
	r, err := etcd.NewEtcdResolver([]string{config.Cfg.EtcdAddr})
	if err != nil {
		logger.Fatal("etcd resolver error", zap.Error(err))
	}
	cli, err := videoservice.NewClient("videoService", client.WithResolver(r))
	if err != nil {
		logger.Fatal("video service client error", zap.Error(err))
	}
	videoClient = cli
}

func VideoPublish(ctx context.Context, req *video.VideoPublishReq) (int32, error) {
	resp, err := videoClient.VideoPublish(ctx, req)
	if err != nil || resp == nil {
		return consts.VideoReqValidError, err
	}
	return resp.Code, nil
}
func VideoList(ctx context.Context, req *video.VideoListReq) (int32, *model.VideoData, error) {
	resp, err := videoClient.VideoList(ctx, req)
	if err != nil || resp.Data == nil {
		return consts.VideoDBSelectError, nil, err
	}
	return resp.Code, utils.ToVideoData(resp.Data), nil
}
func VideoSearch(ctx context.Context, req *video.VideoSearchReq) (int32, *model.VideoData, error) {
	resp, err := videoClient.VideoSearch(ctx, req)
	if err != nil || resp.Data == nil {
		return consts.VideoDBSelectError, nil, err
	}
	return resp.Code, utils.ToVideoData(resp.Data), nil
}
func VideoPopular(ctx context.Context, req *video.VideoHotReq) (int32, *model.VideoData, error) {
	resp, err := videoClient.VideoPopular(ctx, req)
	if err != nil || resp.Data == nil {
		return consts.VideoDBSelectError, nil, err
	}
	return resp.Code, utils.ToVideoData(resp.Data), nil
}
func VideoStream(ctx context.Context, req *video.VideoStreamReq) (int32, *model.VideoData, error) {
	resp, err := videoClient.VideoStream(ctx, req)
	if err != nil || resp.Data == nil {
		return consts.VideoDBSelectError, nil, err
	}
	return resp.Code, utils.ToVideoData(resp.Data), nil
}

func BatchGetVideo(ctx context.Context, ids []string) (int32, []*model.VideoInfo, error) {
	resp, err := videoClient.BatchGetVideo(ctx, &video.BatchGetVideoReq{Ids: ids})
	if err != nil || resp == nil {
		return consts.VideoDBSelectError, nil, err
	}
	return resp.Code, utils.ToBizVideoInfoList(resp.Data.Items), nil
}
