package video

import (
	"Tiktok/internal/video/service"
	video "Tiktok/kitex_gen/video"
	"context"
)

type VideoServiceImpl struct {
	videoService *service.VideoRepo
}

func NewVideoServiceImpl(videoRepo *service.VideoRepo) *VideoServiceImpl {
	return &VideoServiceImpl{videoService: videoRepo}
}

func (s *VideoServiceImpl) VideoPublish(ctx context.Context, req *video.VideoPublishReq) (resp *video.VideoPublishResp, err error) {
	resp = &video.VideoPublishResp{}
	code, _ := s.videoService.VideoPublish(ctx, req.Title, req.Description, req.VideoURL, req.CoverURL, req.UserID)
	resp.Code = code
	return resp, nil
}

func (s *VideoServiceImpl) VideoList(ctx context.Context, req *video.VideoListReq) (resp *video.VideoListResp, err error) {
	code, info, _ := s.videoService.VideoList(req.UserId, req.PageSize, req.PageNum)
	resp = &video.VideoListResp{
		Code: code,
		Data: &video.VideoData{
			Items: info,
			Total: int64(len(info)),
		},
	}
	return resp, nil
}

func (s *VideoServiceImpl) VideoSearch(ctx context.Context, req *video.VideoSearchReq) (resp *video.VideoSearchResp, err error) {
	code, infos, _ := s.videoService.VideoSearch(req.KeyWord, req.PageNum, req.PageSize)
	resp = &video.VideoSearchResp{
		Code: code,
		Data: &video.VideoData{
			Items: infos,
			Total: int64(len(infos)),
		},
	}
	return resp, nil
}

func (s *VideoServiceImpl) VideoPopular(ctx context.Context, req *video.VideoHotReq) (resp *video.VideoHotResp, err error) {
	code, infos, _ := s.videoService.VideoPopular(ctx, req.PageNum, req.PageSize)
	resp = &video.VideoHotResp{
		Code: code,
		Data: &video.VideoData{
			Items: infos,
			Total: int64(len(infos)),
		},
	}
	return resp, nil
}

func (s *VideoServiceImpl) VideoStream(ctx context.Context, req *video.VideoStreamReq) (resp *video.VideoStreamResp, err error) {
	code, infos, _ := s.videoService.VideoStream()
	resp = &video.VideoStreamResp{
		Code: code,
		Data: &video.VideoData{
			Items: infos,
			Total: int64(len(infos)),
		},
	}
	return resp, nil
}

func (s *VideoServiceImpl) BatchGetVideo(ctx context.Context, req *video.BatchGetVideoReq) (resp *video.BatchGetVideoResp, err error) {
	code, infos, _ := s.videoService.BatchGetVideo(req.Ids)
	resp = &video.BatchGetVideoResp{
		Code: code,
		Data: &video.VideoData{
			Items: infos,
			Total: int64(len(infos)),
		},
	}
	return resp, nil
}
