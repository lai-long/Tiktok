package video

import (
	"Tiktok/internal/video/service"
	video "Tiktok/kitex_gen/video"
	"Tiktok/pkg/consts"
	"context"

	"github.com/alibaba/sentinel-golang/api"
)

// VideoServiceImpl implements the last service interface defined in the IDL.
type VideoServiceImpl struct {
	videoService *service.VideoRepo
}

func NewVideoServiceImpl(videoRepo *service.VideoRepo) *VideoServiceImpl {
	return &VideoServiceImpl{videoService: videoRepo}
}

// VideoPublish implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) VideoPublish(ctx context.Context, req *video.VideoPublishReq) (resp *video.VideoPublishResp, err error) {
	// TODO: Your code here...
	resp = &video.VideoPublishResp{}
	entry, blockErr := api.Entry("/video/publish")
	if blockErr != nil {
		resp.Code = consts.SentinelBlock
		return resp, nil
	}
	defer entry.Exit()

	code, err := s.videoService.VideoPublish(ctx, req.Title, req.Description, req.VideoURL, req.UserID)
	resp.Code = code
	return resp, nil
}

// VideoList implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) VideoList(ctx context.Context, req *video.VideoListReq) (resp *video.VideoListResp, err error) {
	// TODO: Your code here...
	resp = &video.VideoListResp{}
	entry, blockErr := api.Entry("/video/list")
	if blockErr != nil {
		resp.Code = consts.SentinelBlock
		return resp, nil
	}
	defer entry.Exit()

	code, info, err := s.videoService.VideoList(req.UserId, req.PageSize, req.PageNum)
	resp = &video.VideoListResp{
		Code: code,
		Data: &video.VideoData{
			Items: info,
			Total: int64(len(info)),
		},
	}
	return resp, nil
}

// VideoSearch implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) VideoSearch(ctx context.Context, req *video.VideoSearchReq) (resp *video.VideoSearchResp, err error) {
	// TODO: Your code here...
	resp = &video.VideoSearchResp{}
	entry, blockErr := api.Entry("/video/search")
	if blockErr != nil {
		resp.Code = consts.SentinelBlock
		return resp, nil
	}
	defer entry.Exit()
	code, infos, err := s.videoService.VideoSearch(req.KeyWord, req.PageNum, req.PageSize)
	resp = &video.VideoSearchResp{
		Code: code,
		Data: &video.VideoData{
			Items: infos,
			Total: int64(len(infos)),
		},
	}
	return resp, nil
}

// VideoPopular implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) VideoPopular(ctx context.Context, req *video.VideoHotReq) (resp *video.VideoHotResp, err error) {
	// TODO: Your code here...
	resp = &video.VideoHotResp{}
	entry, blockErr := api.Entry("/video/popular")
	if blockErr != nil {
		resp.Code = consts.SentinelBlock
		return resp, nil
	}
	defer entry.Exit()
	code, infos, err := s.videoService.VideoPopular(ctx, req.PageNum, req.PageSize)
	resp = &video.VideoHotResp{
		Code: code,
		Data: &video.VideoData{
			Items: infos,
			Total: int64(len(infos)),
		},
	}
	return resp, nil
}

// VideoStream implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) VideoStream(ctx context.Context, req *video.VideoStreamReq) (resp *video.VideoStreamResp, err error) {
	// TODO: Your code here...
	resp = &video.VideoStreamResp{}
	entry, blockErr := api.Entry("/video/stream")
	if blockErr != nil {
		resp.Code = consts.SentinelBlock
		return resp, nil
	}
	defer entry.Exit()
	code, infos, err := s.videoService.VideoStream()
	resp = &video.VideoStreamResp{
		Code: code,
		Data: &video.VideoData{
			Items: infos,
			Total: int64(len(infos)),
		},
	}
	return resp, nil
}
