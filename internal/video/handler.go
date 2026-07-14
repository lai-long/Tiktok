package video

import (
	video "Tiktok/kitex_gen/video"
	"context"
)

// videoService 是 VideoServiceImpl 依赖的最小接口，便于测试时注入 mock。
type videoService interface {
	VideoPublish(ctx context.Context, title string, description string, url string, coverURL string, userID string) (int32, error)
	VideoList(ctx context.Context, userId string, pageSize int64, pageNum int64) (int32, []*video.VideoInfo, error)
	VideoSearch(ctx context.Context, keyword string, pageNum int64, pageSize int64) (int32, []*video.VideoInfo, error)
	VideoPopular(ctx context.Context, pageNum int64, pageSize int64) (int32, []*video.VideoInfo, error)
	VideoStream(ctx context.Context) (int32, []*video.VideoInfo, error)
	BatchGetVideo(ctx context.Context, ids []string) (int32, []*video.VideoInfo, error)
}

type VideoServiceImpl struct {
	videoService videoService
}

func NewVideoServiceImpl(videoRepo videoService) *VideoServiceImpl {
	return &VideoServiceImpl{videoService: videoRepo}
}

func (s *VideoServiceImpl) VideoPublish(ctx context.Context, req *video.VideoPublishReq) (resp *video.VideoPublishResp, err error) {
	resp = &video.VideoPublishResp{}
	code, err := s.videoService.VideoPublish(ctx, req.Title, req.Description, req.VideoURL, req.CoverURL, req.UserID)
	if err != nil {
		resp.Code = code
		return resp, err
	}
	resp.Code = code
	return resp, nil
}

func (s *VideoServiceImpl) VideoList(ctx context.Context, req *video.VideoListReq) (resp *video.VideoListResp, err error) {
	code, info, err := s.videoService.VideoList(ctx, req.UserId, req.PageSize, req.PageNum)
	if err != nil {
		resp = &video.VideoListResp{
			Code: code,
			Data: &video.VideoData{
				Items: info,
				Total: int64(len(info)),
			},
		}
		return resp, err
	}
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
	code, infos, err := s.videoService.VideoSearch(ctx, req.KeyWord, req.PageNum, req.PageSize)
	if err != nil {
		resp = &video.VideoSearchResp{
			Code: code,
			Data: &video.VideoData{
				Items: infos,
				Total: int64(len(infos)),
			},
		}
		return resp, err
	}
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
	code, infos, err := s.videoService.VideoPopular(ctx, req.PageNum, req.PageSize)
	if err != nil {
		resp = &video.VideoHotResp{
			Code: code,
			Data: &video.VideoData{
				Items: infos,
				Total: int64(len(infos)),
			},
		}
		return resp, err
	}
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
	code, infos, err := s.videoService.VideoStream(ctx)
	if err != nil {
		resp = &video.VideoStreamResp{
			Code: code,
			Data: &video.VideoData{
				Items: infos,
				Total: int64(len(infos)),
			},
		}
		return resp, err
	}
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
	code, infos, err := s.videoService.BatchGetVideo(ctx, req.Ids)
	if err != nil {
		resp = &video.BatchGetVideoResp{
			Code: code,
			Data: &video.VideoData{
				Items: infos,
				Total: int64(len(infos)),
			},
		}
		return resp, err
	}
	resp = &video.BatchGetVideoResp{
		Code: code,
		Data: &video.VideoData{
			Items: infos,
			Total: int64(len(infos)),
		},
	}
	return resp, nil
}
