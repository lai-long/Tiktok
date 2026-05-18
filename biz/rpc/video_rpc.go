package rpc

import (
	"Tiktok/kitex_gen/video"
	"Tiktok/kitex_gen/video/videoservice"
	"Tiktok/pkg/consts"
	"Tiktok/pkg/logger"
	"context"

	model "Tiktok/biz/model/video"

	"github.com/cloudwego/kitex/client"
	etcd "github.com/kitex-contrib/registry-etcd"
	"go.uber.org/zap"
)

func InitVideoRpc() {
	r, err := etcd.NewEtcdResolver([]string{"127.0.0.1:2379"})
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
	infos := make([]*model.VideoInfo, len(resp.Data.Items))
	for i := 0; i < len(resp.Data.Items); i++ {
		infos = append(infos, &model.VideoInfo{
			ID:           resp.Data.Items[i].ID,
			UserID:       resp.Data.Items[i].UserID,
			Title:        resp.Data.Items[i].Title,
			Description:  resp.Data.Items[i].Description,
			CommentCount: resp.Data.Items[i].CommentCount,
			CoverURL:     resp.Data.Items[i].CoverURL,
			CreatedAt:    resp.Data.Items[i].CreatedAt,
			LikeCount:    resp.Data.Items[i].LikeCount,
			UpdatedAt:    resp.Data.Items[i].UpdatedAt,
			VideoURL:     resp.Data.Items[i].VideoURL,
			VisitCount:   resp.Data.Items[i].VisitCount,
		})
	}
	data := &model.VideoData{
		Items: infos,
		Total: int64(len(infos)),
	}
	return resp.Code, data, err
}
func VideoSearch(ctx context.Context, req *video.VideoSearchReq) (int32, *model.VideoData, error) {
	resp, err := videoClient.VideoSearch(ctx, req)
	if err != nil || resp.Data == nil {
		return consts.VideoDBSelectError, nil, err
	}
	infos := make([]*model.VideoInfo, 0, resp.Data.Total)
	for i := 0; i < len(resp.Data.Items); i++ {
		infos = append(infos, &model.VideoInfo{
			ID:           resp.Data.Items[i].ID,
			UserID:       resp.Data.Items[i].UserID,
			Title:        resp.Data.Items[i].Title,
			Description:  resp.Data.Items[i].Description,
			CommentCount: resp.Data.Items[i].CommentCount,
			CoverURL:     resp.Data.Items[i].CoverURL,
			CreatedAt:    resp.Data.Items[i].CreatedAt,
			LikeCount:    resp.Data.Items[i].LikeCount,
			UpdatedAt:    resp.Data.Items[i].UpdatedAt,
			VideoURL:     resp.Data.Items[i].VideoURL,
			VisitCount:   resp.Data.Items[i].VisitCount,
		})
	}
	data := &model.VideoData{
		Items: infos,
		Total: int64(len(infos)),
	}
	return resp.Code, data, err
}
func VideoPopular(ctx context.Context, req *video.VideoHotReq) (int32, *model.VideoData, error) {
	resp, err := videoClient.VideoPopular(ctx, req)
	if err != nil || resp.Data == nil {
		return consts.VideoDBSelectError, nil, err
	}
	infos := make([]*model.VideoInfo, 0, resp.Data.Total)
	for i := 0; i < len(resp.Data.Items); i++ {
		infos = append(infos, &model.VideoInfo{
			ID:           resp.Data.Items[i].ID,
			UserID:       resp.Data.Items[i].UserID,
			Title:        resp.Data.Items[i].Title,
			Description:  resp.Data.Items[i].Description,
			CommentCount: resp.Data.Items[i].CommentCount,
			CoverURL:     resp.Data.Items[i].CoverURL,
			CreatedAt:    resp.Data.Items[i].CreatedAt,
			LikeCount:    resp.Data.Items[i].LikeCount,
			UpdatedAt:    resp.Data.Items[i].UpdatedAt,
			VideoURL:     resp.Data.Items[i].VideoURL,
			VisitCount:   resp.Data.Items[i].VisitCount,
		})
	}
	data := &model.VideoData{
		Items: infos,
		Total: int64(len(infos)),
	}
	return resp.Code, data, err
}
func VideoStream(ctx context.Context, req *video.VideoStreamReq) (int32, *model.VideoData, error) {
	resp, err := videoClient.VideoStream(ctx, req)
	if err != nil || resp.Data == nil {
		return consts.VideoDBSelectError, nil, err
	}
	infos := make([]*model.VideoInfo, 0, resp.Data.Total)
	for i := 0; i < len(resp.Data.Items); i++ {
		infos = append(infos, &model.VideoInfo{
			ID:           resp.Data.Items[i].ID,
			UserID:       resp.Data.Items[i].UserID,
			Title:        resp.Data.Items[i].Title,
			Description:  resp.Data.Items[i].Description,
			CommentCount: resp.Data.Items[i].CommentCount,
			CoverURL:     resp.Data.Items[i].CoverURL,
			CreatedAt:    resp.Data.Items[i].CreatedAt,
			LikeCount:    resp.Data.Items[i].LikeCount,
			UpdatedAt:    resp.Data.Items[i].UpdatedAt,
			VideoURL:     resp.Data.Items[i].VideoURL,
			VisitCount:   resp.Data.Items[i].VisitCount,
		})
	}
	data := &model.VideoData{
		Items: infos,
		Total: int64(len(infos)),
	}
	return resp.Code, data, err
}
