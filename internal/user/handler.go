package handler

import (
	"Tiktok/internal/user/service"
	user "Tiktok/kitex_gen/user"
	"Tiktok/pkg/consts"
	"context"
	"log"

	"github.com/alibaba/sentinel-golang/api"
	"github.com/pkg/errors"
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
	// TODO: Your code here...
	entry, blockErr := api.Entry("/user/register")
	resp = &user.RegisterResp{}
	if blockErr != nil {
		resp.Code = consts.SentinelBlock
		return resp, blockErr
	}
	defer entry.Exit()
	code, err := s.service.Register(req.UserName, req.Password)
	resp.Code = code
	return resp, nil
}

// UserLogin implements the UserServiceImpl interface.
func (s *UserServiceImpl) UserLogin(ctx context.Context, req *user.LoginReq) (resp *user.LoginResp, err error) {
	// TODO: Your code here...
	entry, blockErr := api.Entry("/user/login")
	resp = &user.LoginResp{}
	if blockErr != nil {
		resp.Code = consts.SentinelBlock
		return resp, blockErr
	}
	defer entry.Exit()
	code, userInfo, reToken, acToken, err := s.service.Login(req.UserName, req.Password, req.Code, ctx)
	if err != nil {
		resp.Code = code
		return resp, errors.Wrap(err, "login")
	}
	resp = &user.LoginResp{
		Code:         code,
		Data:         userInfo,
		RefreshToken: reToken,
		AccessToken:  acToken,
	}
	return resp, nil
}

// UserInfo implements the UserServiceImpl interface.
func (s *UserServiceImpl) UserInfo(ctx context.Context, req *user.UserInfoReq) (resp *user.UserInfoResp, err error) {
	// TODO: Your code here...
	entry, blockErr := api.Entry("/user/info")
	resp = &user.UserInfoResp{}
	if blockErr != nil {
		resp.Code = consts.SentinelBlock
		return resp, nil
	}
	defer entry.Exit()
	userInfo, code, err := s.service.UserInfo(ctx, req.UserId)
	if err != nil {
		resp.Code = code
		return resp, errors.Wrap(err, "UserInfo")
	}
	resp = &user.UserInfoResp{
		Code: code,
		Data: userInfo,
	}
	return resp, nil
}

// UserAvatar implements the UserServiceImpl interface.
func (s *UserServiceImpl) UserAvatar(ctx context.Context, req *user.UserAvatarReq) (resp *user.UserAvatarResp, err error) {
	// TODO: Your code here...
	resp = &user.UserAvatarResp{}
	entry, blockErr := api.Entry("/user/avatar/upload")
	if blockErr != nil {
		resp.Code = consts.SentinelBlock
		return resp, nil
	}
	defer entry.Exit()

	code, userInfo, err := s.service.UserAvatar(req.AvatarURL, req.UserID)
	if err != nil {
		log.Println("userService.UserAvatar error:", err)
		resp.Code = code
		return resp, errors.Wrap(err, "UserAvatar")
	}
	resp = &user.UserAvatarResp{
		Code: code,
		Data: userInfo,
	}
	return resp, nil
}

// RefreshToken implements the UserServiceImpl interface.
func (s *UserServiceImpl) RefreshToken(ctx context.Context, req *user.RefreshTokenReq) (resp *user.RefreshTokenResp, err error) {
	// TODO: Your code here...
	resp = &user.RefreshTokenResp{}
	entry, blockErr := api.Entry("/user/refresh")
	if blockErr != nil {
		resp.Code = consts.SentinelBlock
		return resp, nil
	}
	defer entry.Exit()
	code, reToken, acToken, err := s.service.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		resp.Code = code
		return resp, errors.Wrap(err, "RefreshToken")
	}
	resp = &user.RefreshTokenResp{
		Code:         code,
		RefreshToken: reToken,
		AccessToken:  acToken,
	}
	return resp, nil
}
