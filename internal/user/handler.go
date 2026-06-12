package handler

import (
	"Tiktok/internal/user/service"
	user "Tiktok/kitex_gen/user"
	"Tiktok/pkg/logger"
	"Tiktok/pkg/utils"
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

// MfaQrcode implements the UserServiceImpl interface.
func (s *UserServiceImpl) MfaQrcode(ctx context.Context, req *user.MfaQrcodeReq) (resp *user.MfaQrcodeResp, err error) {
	qrCode, secret, code, _ := s.service.GenerateMfa(ctx, req.UserName, req.UserID)
	resp = &user.MfaQrcodeResp{
		Code: code,
		Data: &user.MfaData{
			Secret: secret,
			Qrcode: qrCode,
		},
	}
	return resp, nil
}

// MfaBind implements the UserServiceImpl interface.
func (s *UserServiceImpl) MfaBind(ctx context.Context, req *user.MfaBindReq) (resp *user.MfaBindResp, err error) {
	var code int32
	switch req.Type {
	case "secret":
		code, _ = s.service.MfaBindBySecret(ctx, req.Secret, req.UserID)
	case "qrcode":
		code, _ = s.service.MfaBindByCode(ctx, req.MfaCode, req.UserID)
	default:
		code = 1
	}
	resp = &user.MfaBindResp{
		Code: code,
	}
	return resp, nil
}

// MfaConfirm implements the UserServiceImpl interface.
func (s *UserServiceImpl) MfaConfirm(ctx context.Context, req *user.MfaConfirmReq) (resp *user.MfaConfirmResp, err error) {
	resp = &user.MfaConfirmResp{}
	code, err := s.service.MfaConfirm(ctx, req.QrCode, req.UserID)
	if err != nil {
		resp.Code = code
		return resp, err
	}
	resp.Code = code
	return resp, nil
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
		logger.Error("UserInfo failed", zap.Error(err), logger.WithUserID(req.UserId), logger.WithTraceID(utils.GetTraceID(ctx)))
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
		logger.Error("UserAvatar failed", zap.Error(err), logger.WithUserID(req.UserID), logger.WithTraceID(utils.GetTraceID(ctx)))
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
		logger.Error("RefreshToken failed", zap.Error(err), logger.WithTraceID(utils.GetTraceID(ctx)))
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
