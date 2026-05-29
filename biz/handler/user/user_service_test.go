package user

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"Tiktok/pkg/consts"

	Rpc "Tiktok/biz/rpc"
	user2 "Tiktok/kitex_gen/user"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/cloudwego/kitex/client/callopt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// jsonBody 辅助函数：构造 JSON body 和 Content-Type Header
func jsonBody(v any) (*ut.Body, []ut.Header) {
	b, _ := json.Marshal(v)
	return &ut.Body{Body: bytes.NewReader(b), Len: len(b)},
		[]ut.Header{{Key: "Content-Type", Value: "application/json"}}
}

// assertResponseCode 验证响应体中的 base.code 字段
func assertResponseCode(t *testing.T, c *app.RequestContext, wantCode int32) {
	var result map[string]interface{}
	assert.NoError(t, json.Unmarshal(c.Response.Body(), &result))
	base := result["base"].(map[string]interface{})
	assert.Equal(t, wantCode, int32(base["code"].(float64)))
}

type MockUserClient struct {
	mock.Mock
}

func (m *MockUserClient) UserRegister(ctx context.Context, req *user2.RegisterReq, callOptions ...callopt.Option) (*user2.RegisterResp, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user2.RegisterResp), args.Error(1)
}

func (m *MockUserClient) UserLogin(ctx context.Context, req *user2.LoginReq, callOptions ...callopt.Option) (*user2.LoginResp, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user2.LoginResp), args.Error(1)
}

func (m *MockUserClient) UserInfo(ctx context.Context, req *user2.UserInfoReq, callOptions ...callopt.Option) (*user2.UserInfoResp, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user2.UserInfoResp), args.Error(1)
}

func (m *MockUserClient) UserAvatar(ctx context.Context, req *user2.UserAvatarReq, callOptions ...callopt.Option) (*user2.UserAvatarResp, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user2.UserAvatarResp), args.Error(1)
}

func (m *MockUserClient) RefreshToken(ctx context.Context, req *user2.RefreshTokenReq, callOptions ...callopt.Option) (*user2.RefreshTokenResp, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user2.RefreshTokenResp), args.Error(1)
}

func TestUserRegister(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(*MockUserClient)
		wantCode  int32
		wantErr   bool
	}{
		{
			name: "Success_register",
			mockSetup: func(m *MockUserClient) {
				m.On("UserRegister", mock.Anything, mock.MatchedBy(func(req *user2.RegisterReq) bool {
					return req.UserName != "" && req.Password != ""
				})).Return(&user2.RegisterResp{Code: consts.Success}, nil)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name: "Fail_register_error",
			mockSetup: func(m *MockUserClient) {
				m.On("UserRegister", mock.Anything, mock.MatchedBy(func(req *user2.RegisterReq) bool {
					return req.UserName != "" && req.Password != ""
				})).Return(&user2.RegisterResp{Code: consts.UserDBSelectError}, nil)
			},
			wantCode: consts.UserDBSelectError,
			wantErr:  true,
		},
		{
			name: "Fail_register_rpc_error",
			mockSetup: func(m *MockUserClient) {
				m.On("UserRegister", mock.Anything, mock.MatchedBy(func(req *user2.RegisterReq) bool {
					return req.UserName != "" && req.Password != ""
				})).Return(nil, assert.AnError)
			},
			wantCode: consts.UserReqValidError,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockUserClient)
			Rpc.SetUserClient(mockClient)
			tt.mockSetup(mockClient)

			body, header := jsonBody(map[string]any{"username": "test", "password": "123456"})
			c := ut.CreateUtRequestContext("POST", "/user/register", body, header...)

			UserRegister(context.Background(), c)

			assert.Equal(t, 200, c.Response.StatusCode())
			assertResponseCode(t, c, tt.wantCode)
			mockClient.AssertExpectations(t)
		})
	}
}

func TestUserLogin(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(*MockUserClient)
		wantCode  int32
		wantErr   bool
	}{
		{
			name: "Success_login",
			mockSetup: func(m *MockUserClient) {
				m.On("UserLogin", mock.Anything, mock.MatchedBy(func(req *user2.LoginReq) bool {
					return req.UserName != "" && req.Password != ""
				})).Return(&user2.LoginResp{
					Code:         consts.Success,
					Data:         &user2.UserInfo{ID: "123", Username: "testuser"},
					RefreshToken: "refresh_token",
					AccessToken:  "access_token",
				}, nil)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name: "Fail_login_error",
			mockSetup: func(m *MockUserClient) {
				m.On("UserLogin", mock.Anything, mock.MatchedBy(func(req *user2.LoginReq) bool {
					return req.UserName != "" && req.Password != ""
				})).Return(&user2.LoginResp{Code: consts.UserPasswordError}, nil)
			},
			wantCode: consts.UserPasswordError,
			wantErr:  true,
		},
		{
			name: "Fail_login_rpc_error",
			mockSetup: func(m *MockUserClient) {
				m.On("UserLogin", mock.Anything, mock.MatchedBy(func(req *user2.LoginReq) bool {
					return req.UserName != "" && req.Password != ""
				})).Return(nil, assert.AnError)
			},
			wantCode: consts.UserReqValidError,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockUserClient)
			Rpc.SetUserClient(mockClient)
			tt.mockSetup(mockClient)

			body, header := jsonBody(map[string]any{"username": "test", "password": "123"})
			c := ut.CreateUtRequestContext("POST", "/user/login", body, header...)

			UserLogin(context.Background(), c)

			assert.Equal(t, 200, c.Response.StatusCode())
			assertResponseCode(t, c, tt.wantCode)
			mockClient.AssertExpectations(t)
		})
	}
}

func TestUserInfo(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(*MockUserClient)
		wantCode  int32
		wantErr   bool
	}{
		{
			name: "Success_userInfo",
			mockSetup: func(m *MockUserClient) {
				m.On("UserInfo", mock.Anything, mock.MatchedBy(func(req *user2.UserInfoReq) bool {
					return req.UserId != ""
				})).Return(&user2.UserInfoResp{
					Code: consts.Success,
					Data: &user2.UserInfo{ID: "123", Username: "testuser"},
				}, nil)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name: "Fail_userInfo_error",
			mockSetup: func(m *MockUserClient) {
				m.On("UserInfo", mock.Anything, mock.MatchedBy(func(req *user2.UserInfoReq) bool {
					return req.UserId != ""
				})).Return(&user2.UserInfoResp{Code: consts.UserDBSelectError}, nil)
			},
			wantCode: consts.UserDBSelectError,
			wantErr:  true,
		},
		{
			name: "Fail_userInfo_rpc_error",
			mockSetup: func(m *MockUserClient) {
				m.On("UserInfo", mock.Anything, mock.MatchedBy(func(req *user2.UserInfoReq) bool {
					return req.UserId != ""
				})).Return(nil, assert.AnError)
			},
			wantCode: consts.UserReqValidError,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockUserClient)
			Rpc.SetUserClient(mockClient)
			tt.mockSetup(mockClient)

			c := ut.CreateUtRequestContext("GET", "/user/info?user_id=123", nil)
			c.Set("userid", "123")

			UserInfo(context.Background(), c)

			assert.Equal(t, 200, c.Response.StatusCode())
			assertResponseCode(t, c, tt.wantCode)
			mockClient.AssertExpectations(t)
		})
	}
}

func TestRefreshToken(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(*MockUserClient)
		wantCode  int32
		wantErr   bool
	}{
		{
			name: "Success_refreshToken",
			mockSetup: func(m *MockUserClient) {
				m.On("RefreshToken", mock.Anything, mock.MatchedBy(func(req *user2.RefreshTokenReq) bool {
					return req.RefreshToken != ""
				})).Return(&user2.RefreshTokenResp{
					Code:         consts.Success,
					RefreshToken: "new_refresh_token",
					AccessToken:  "new_access_token",
				}, nil)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name: "Fail_refreshToken_error",
			mockSetup: func(m *MockUserClient) {
				m.On("RefreshToken", mock.Anything, mock.MatchedBy(func(req *user2.RefreshTokenReq) bool {
					return req.RefreshToken != ""
				})).Return(&user2.RefreshTokenResp{Code: consts.UserRedisGetError}, nil)
			},
			wantCode: consts.UserRedisGetError,
			wantErr:  true,
		},
		{
			name: "Fail_refreshToken_rpc_error",
			mockSetup: func(m *MockUserClient) {
				m.On("RefreshToken", mock.Anything, mock.MatchedBy(func(req *user2.RefreshTokenReq) bool {
					return req.RefreshToken != ""
				})).Return(nil, assert.AnError)
			},
			wantCode: consts.UserReqValidError,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockUserClient)
			Rpc.SetUserClient(mockClient)
			tt.mockSetup(mockClient)

			body, header := jsonBody(map[string]any{"refresh_token": "test_token"})
			c := ut.CreateUtRequestContext("POST", "/user/refresh", body, header...)

			RefreshToken(context.Background(), c)

			assert.Equal(t, 200, c.Response.StatusCode())
			assertResponseCode(t, c, tt.wantCode)
			mockClient.AssertExpectations(t)
		})
	}
}
