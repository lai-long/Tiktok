package handler

import (
	"Tiktok/internal/user/service"
	user "Tiktok/kitex_gen/user"
	"Tiktok/pkg/logger"
	"context"
	"go.uber.org/zap"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct {
	service *service.UserRepo
}

func NewUserService(service *service.UserRepo) *UserServiceImpl {
	return &UserServiceImpl{service: service}
}

// UserRegister implements the UserServiceImpl interface.
func (s *UserServiceImpl) UserRegister(ctx context.Context, req *user.RegisterReq) (resp *user.RegisterResp, err error) {
	resp = &user.RegisterResp{}
	code, _ := s.service.Register(ctx, req.UserName, req.Password)
	resp.Code = code
	return resp, nil
}

// UserLogin implements the UserServiceImpl interface.
func (s *UserServiceImpl) UserLogin(ctx context.Context, req *user.LoginReq) (resp *user.LoginResp, err error) {
	resp = &user.LoginResp{}
	code, userInfo, reToken, acToken, err := s.service.Login(ctx, req.UserName, req.Password, req.Code)
	if err != nil {
		resp.Code = code
	}
	resp.Code = code
	resp.Data = userInfo
	resp.RefreshToken = reToken
	resp.AccessToken = acToken
	return resp, nil
}

// UserInfo implements the UserServiceImpl interface.
func (s *UserServiceImpl) UserInfo(ctx context.Context, req *user.UserInfoReq) (resp *user.UserInfoResp, err error) {
	resp = &user.UserInfoResp{}
	userInfo, code, err := s.service.UserInfo(ctx, req.UserId)
	if err != nil {
		logger.Error("UserInfo failed", zap.Error(err), zap.String("user_id", req.UserId))
		resp.Code = code
		resp.Data = userInfo
		return resp, nil
	}
	resp.Code = code
	resp.Data = userInfo
	return resp, nil
}

// UserAvatar implements the UserServiceImpl interface.
func (s *UserServiceImpl) UserAvatar(ctx context.Context, req *user.UserAvatarReq) (resp *user.UserAvatarResp, err error) {
	resp = &user.UserAvatarResp{}
	code, userInfo, err := s.service.UserAvatar(ctx, req.AvatarURL, req.UserID)
	if err != nil {
		logger.Error("UserAvatar failed", zap.Error(err), zap.String("user_id", req.UserID))
		resp.Code = code
		resp.Data = userInfo
		return resp, nil
	}
	resp.Code = code
	resp.Data = userInfo
	return resp, nil
}

// RefreshToken implements the UserServiceImpl interface.
func (s *UserServiceImpl) RefreshToken(ctx context.Context, req *user.RefreshTokenReq) (resp *user.RefreshTokenResp, err error) {
	resp = &user.RefreshTokenResp{}
	code, reToken, acToken, err := s.service.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		logger.Error("RefreshToken failed", zap.Error(err))
		resp.Code = code
		resp.RefreshToken = reToken
		resp.AccessToken = acToken
		return resp, nil
	}
	resp.Code = code
	resp.RefreshToken = reToken
	resp.AccessToken = acToken
	return resp, nil
}
