package converter

import (
	model "Tiktok/biz/model/user"
	kitex "Tiktok/kitex_gen/user"
)

func ToBizUserInfo(src *kitex.UserInfo) *model.UserInfo {
	if src == nil {
		return nil
	}
	return &model.UserInfo{
		ID:        src.ID,
		Username:  src.Username,
		AvatarURL: src.AvatarURL,
		CreatedAt: src.CreatedAt,
		UpdatedAt: src.UpdatedAt,
		DeletedAt: src.DeletedAt,
	}
}
