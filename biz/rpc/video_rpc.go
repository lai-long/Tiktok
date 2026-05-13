package rpc

import (
	"Tiktok/kitex_gen/video"
	"Tiktok/kitex_gen/video/videoservice"
	"context"
	"log"

	model "Tiktok/biz/model/video"

	"github.com/cloudwego/kitex/client"
	etcd "github.com/kitex-contrib/registry-etcd"
)

func InitVideoRpc() {
	r, err := etcd.NewEtcdResolver([]string{"http://127.0.0.1:2379"})
	if err != nil {
		log.Println(err)
		return
	}
	cli, err := videoservice.NewClient("videoService", client.WithResolver(r))
	if err != nil {
		log.Println(err)
		return
	}
	videoClient = cli
}

func VideoPublish(ctx context.Context, req *video.VideoPublishReq) (int32, error) {
	resp, err := videoClient.VideoPublish(ctx, req)
	return resp.Code, err
}
func VideoList(ctx context.Context, req *video.VideoListReq) (int32, *model.VideoData, error) {
	resp, err := videoClient.VideoList(ctx, req)
	infos := make([]*model.VideoInfo, resp.Data.Total)
	for i := 0; i < len(infos); i++ {
		infos[i] = &model.VideoInfo{
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
		}
	}
	data := &model.VideoData{
		Items: infos,
		Total: int64(len(infos)),
	}
	return resp.Code, data, err
}
func VideoSearch(ctx context.Context, req *video.VideoSearchReq) (int32, *model.VideoData, error) {
	resp, err := videoClient.VideoSearch(ctx, req)
	infos := make([]*model.VideoInfo, resp.Data.Total)
	for i := 0; i < len(infos); i++ {
		infos[i] = &model.VideoInfo{
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
		}
	}
	data := &model.VideoData{
		Items: infos,
		Total: int64(len(infos)),
	}
	return resp.Code, data, err
}
func VideoPopular(ctx context.Context, req *video.VideoHotReq) (int32, *model.VideoData, error) {
	resp, err := videoClient.VideoPopular(ctx, req)
	infos := make([]*model.VideoInfo, resp.Data.Total)
	for i := 0; i < len(infos); i++ {
		infos[i] = &model.VideoInfo{
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
		}
	}
	data := &model.VideoData{
		Items: infos,
		Total: int64(len(infos)),
	}
	return resp.Code, data, err
}
func VideoStream(ctx context.Context, req *video.VideoStreamReq) (int32, *model.VideoData, error) {
	resp, err := videoClient.VideoStream(ctx, req)
	infos := make([]*model.VideoInfo, resp.Data.Total)
	for i := 0; i < len(infos); i++ {
		infos[i] = &model.VideoInfo{
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
		}
	}
	data := &model.VideoData{
		Items: infos,
		Total: int64(len(infos)),
	}
	return resp.Code, data, err
}
