package utils

import (
	modelUser "Tiktok/biz/model/user"
	modelVideo "Tiktok/biz/model/video"
	kitexUser "Tiktok/kitex_gen/user"
	kitexVideo "Tiktok/kitex_gen/video"
)

func ToBizVideoInfo(src *kitexVideo.VideoInfo) *modelVideo.VideoInfo {
	if src == nil {
		return nil
	}
	return &modelVideo.VideoInfo{
		ID:           src.ID,
		UserID:       src.UserID,
		Title:        src.Title,
		Description:  src.Description,
		CommentCount: src.CommentCount,
		CoverURL:     src.CoverURL,
		CreatedAt:    src.CreatedAt,
		LikeCount:    src.LikeCount,
		UpdatedAt:    src.UpdatedAt,
		VideoURL:     src.VideoURL,
		VisitCount:   src.VisitCount,
	}
}

func ToBizVideoInfoList(src []*kitexVideo.VideoInfo) []*modelVideo.VideoInfo {
	if src == nil {
		return nil
	}
	dst := make([]*modelVideo.VideoInfo, len(src))
	for i, v := range src {
		dst[i] = ToBizVideoInfo(v)
	}
	return dst
}

func ToVideoData(src *kitexVideo.VideoData) *modelVideo.VideoData {
	if src == nil {
		return nil
	}
	return &modelVideo.VideoData{
		Items: ToBizVideoInfoList(src.Items),
		Total: src.Total,
	}
}

func ToBizUserInfo(src *kitexUser.UserInfo) *modelUser.UserInfo {
	if src == nil {
		return nil
	}
	return &modelUser.UserInfo{
		ID:        src.ID,
		Username:  src.Username,
		AvatarURL: src.AvatarURL,
		CreatedAt: src.CreatedAt,
		UpdatedAt: src.UpdatedAt,
		DeletedAt: src.DeletedAt,
	}
}
