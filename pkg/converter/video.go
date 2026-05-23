package converter

import (
	model "Tiktok/biz/model/video"
	kitex "Tiktok/kitex_gen/video"
)

func ToBizVideoInfo(src *kitex.VideoInfo) *model.VideoInfo {
	if src == nil {
		return nil
	}
	return &model.VideoInfo{
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

func ToBizVideoInfoList(src []*kitex.VideoInfo) []*model.VideoInfo {
	if src == nil {
		return nil
	}
	dst := make([]*model.VideoInfo, len(src))
	for i, v := range src {
		dst[i] = ToBizVideoInfo(v)
	}
	return dst
}

func ToVideoData(src *kitex.VideoData) *model.VideoData {
	if src == nil {
		return nil
	}
	return &model.VideoData{
		Items: ToBizVideoInfoList(src.Items),
		Total: src.Total,
	}
}
