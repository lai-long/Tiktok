package handler

import (
	user "Tiktok/kitex_gen/user"
	"Tiktok/pkg/consts"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUserService struct {
	mock.Mock
}

func (m *mockUserService) GenerateMfa(ctx context.Context, username string, userId string) (string, string, int32, error) {
	args := m.Called(ctx, username, userId)
	return args.String(0), args.String(1), args.Get(2).(int32), args.Error(3)
}

func (m *mockUserService) MfaBindBySecret(ctx context.Context, secret string, userId string) (int32, error) {
	args := m.Called(ctx, secret, userId)
	return args.Get(0).(int32), args.Error(1)
}

func (m *mockUserService) MfaBindByCode(ctx context.Context, code string, userId string) (int32, error) {
	args := m.Called(ctx, code, userId)
	return args.Get(0).(int32), args.Error(1)
}

func (m *mockUserService) MfaConfirm(ctx context.Context, mfaCode string, userID string) (int32, error) {
	args := m.Called(ctx, mfaCode, userID)
	return args.Get(0).(int32), args.Error(1)
}

func (m *mockUserService) Register(ctx context.Context, userName string, password string) (int32, error) {
	args := m.Called(ctx, userName, password)
	return args.Get(0).(int32), args.Error(1)
}

func (m *mockUserService) Login(ctx context.Context, userName, password, mfaCode string) (int32, *user.UserInfo, string, string, error) {
	args := m.Called(ctx, userName, password, mfaCode)
	return args.Get(0).(int32), args.Get(1).(*user.UserInfo), args.String(2), args.String(3), args.Error(4)
}

func (m *mockUserService) UserInfo(ctx context.Context, userId string) (*user.UserInfo, int32, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).(*user.UserInfo), args.Get(1).(int32), args.Error(2)
}

func (m *mockUserService) UserAvatar(ctx context.Context, url string, userID string) (int32, *user.UserInfo, error) {
	args := m.Called(ctx, url, userID)
	return args.Get(0).(int32), args.Get(1).(*user.UserInfo), args.Error(2)
}

func (m *mockUserService) RefreshToken(ctx context.Context, refreshToken string) (int32, string, string, error) {
	args := m.Called(ctx, refreshToken)
	return args.Get(0).(int32), args.String(1), args.String(2), args.Error(3)
}

func TestUserServiceImpl_MfaQrcode(t *testing.T) {
	tests := []struct {
		name     string
		mock     func(*mockUserService)
		wantCode int32
		wantErr  bool
	}{
		{
			name: "success",
			mock: func(m *mockUserService) {
				m.On("GenerateMfa", mock.Anything, "alice", "uid").Return("qrcode", "secret", consts.Success, nil)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name: "db_error",
			mock: func(m *mockUserService) {
				m.On("GenerateMfa", mock.Anything, "alice", "uid").Return("", "", consts.UserDBUpdateError, errors.New("db down"))
			},
			wantCode: consts.UserDBUpdateError,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := new(mockUserService)
			tt.mock(m)
			h := NewUserService(m)
			resp, err := h.MfaQrcode(context.Background(), &user.MfaQrcodeReq{UserName: "alice", UserID: "uid"})
			assert.Equal(t, tt.wantCode, resp.Code)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "secret", resp.Data.Secret)
				assert.Equal(t, "qrcode", resp.Data.Qrcode)
			}
			m.AssertExpectations(t)
		})
	}
}

func TestUserServiceImpl_MfaBind(t *testing.T) {
	tests := []struct {
		name     string
		req      *user.MfaBindReq
		mock     func(*mockUserService)
		wantCode int32
		wantErr  bool
	}{
		{
			name: "success_secret",
			req:  &user.MfaBindReq{Type: "secret", Secret: "secret", UserID: "uid"},
			mock: func(m *mockUserService) {
				m.On("MfaBindBySecret", mock.Anything, "secret", "uid").Return(consts.Success, nil)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name: "success_qrcode",
			req:  &user.MfaBindReq{Type: "qrcode", MfaCode: "123456", UserID: "uid"},
			mock: func(m *mockUserService) {
				m.On("MfaBindByCode", mock.Anything, "123456", "uid").Return(consts.Success, nil)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name: "secret_error",
			req:  &user.MfaBindReq{Type: "secret", Secret: "secret", UserID: "uid"},
			mock: func(m *mockUserService) {
				m.On("MfaBindBySecret", mock.Anything, "secret", "uid").Return(consts.UserDBUpdateError, errors.New("db down"))
			},
			wantCode: consts.UserDBUpdateError,
			wantErr:  true,
		},
		{
			name:     "invalid_type",
			req:      &user.MfaBindReq{Type: "invalid", UserID: "uid"},
			wantCode: 1,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := new(mockUserService)
			if tt.mock != nil {
				tt.mock(m)
			}
			h := NewUserService(m)
			resp, err := h.MfaBind(context.Background(), tt.req)
			assert.Equal(t, tt.wantCode, resp.Code)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			m.AssertExpectations(t)
		})
	}
}

func TestUserServiceImpl_UserRegister(t *testing.T) {
	tests := []struct {
		name     string
		mock     func(*mockUserService)
		wantCode int32
		wantErr  bool
	}{
		{
			name: "success",
			mock: func(m *mockUserService) {
				m.On("Register", mock.Anything, "alice", "secret").Return(consts.Success, nil)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name: "db_error",
			mock: func(m *mockUserService) {
				m.On("Register", mock.Anything, "alice", "secret").Return(consts.UserDBInsertError, errors.New("db down"))
			},
			wantCode: consts.UserDBInsertError,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := new(mockUserService)
			tt.mock(m)
			h := NewUserService(m)
			resp, err := h.UserRegister(context.Background(), &user.RegisterReq{UserName: "alice", Password: "secret"})
			assert.Equal(t, tt.wantCode, resp.Code)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			m.AssertExpectations(t)
		})
	}
}

func TestUserServiceImpl_UserLogin(t *testing.T) {
	tests := []struct {
		name     string
		mock     func(*mockUserService)
		wantCode int32
		wantErr  bool
	}{
		{
			name: "success",
			mock: func(m *mockUserService) {
				m.On("Login", mock.Anything, "alice", "secret", "123456").Return(consts.Success, &user.UserInfo{ID: "uid"}, "re", "ac", nil)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name: "db_error",
			mock: func(m *mockUserService) {
				m.On("Login", mock.Anything, "alice", "secret", "123456").Return(consts.UserDBSelectError, &user.UserInfo{}, "", "", errors.New("db down"))
			},
			wantCode: consts.UserDBSelectError,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := new(mockUserService)
			tt.mock(m)
			h := NewUserService(m)
			resp, err := h.UserLogin(context.Background(), &user.LoginReq{UserName: "alice", Password: "secret", Code: "123456"})
			assert.Equal(t, tt.wantCode, resp.Code)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "uid", resp.Data.ID)
			}
			m.AssertExpectations(t)
		})
	}
}

func TestUserServiceImpl_UserInfo(t *testing.T) {
	tests := []struct {
		name     string
		mock     func(*mockUserService)
		wantCode int32
		wantErr  bool
	}{
		{
			name: "success",
			mock: func(m *mockUserService) {
				m.On("UserInfo", mock.Anything, "uid").Return(&user.UserInfo{ID: "uid"}, consts.Success, nil)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name: "redis_error",
			mock: func(m *mockUserService) {
				m.On("UserInfo", mock.Anything, "uid").Return(&user.UserInfo{}, consts.UserRedisSetError, errors.New("redis down"))
			},
			wantCode: consts.UserRedisSetError,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := new(mockUserService)
			tt.mock(m)
			h := NewUserService(m)
			resp, err := h.UserInfo(context.Background(), &user.UserInfoReq{UserId: "uid"})
			assert.Equal(t, tt.wantCode, resp.Code)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			m.AssertExpectations(t)
		})
	}
}

func TestUserServiceImpl_UserAvatar(t *testing.T) {
	tests := []struct {
		name     string
		mock     func(*mockUserService)
		wantCode int32
		wantErr  bool
	}{
		{
			name: "success",
			mock: func(m *mockUserService) {
				m.On("UserAvatar", mock.Anything, "url", "uid").Return(consts.Success, &user.UserInfo{ID: "uid"}, nil)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name: "db_error",
			mock: func(m *mockUserService) {
				m.On("UserAvatar", mock.Anything, "url", "uid").Return(consts.UserDBUpdateError, &user.UserInfo{}, errors.New("db down"))
			},
			wantCode: consts.UserDBUpdateError,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := new(mockUserService)
			tt.mock(m)
			h := NewUserService(m)
			resp, err := h.UserAvatar(context.Background(), &user.UserAvatarReq{AvatarURL: "url", UserID: "uid"})
			assert.Equal(t, tt.wantCode, resp.Code)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			m.AssertExpectations(t)
		})
	}
}

func TestUserServiceImpl_RefreshToken(t *testing.T) {
	tests := []struct {
		name     string
		mock     func(*mockUserService)
		wantCode int32
		wantErr  bool
	}{
		{
			name: "success",
			mock: func(m *mockUserService) {
				m.On("RefreshToken", mock.Anything, "token").Return(consts.Success, "re", "ac", nil)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name: "redis_error",
			mock: func(m *mockUserService) {
				m.On("RefreshToken", mock.Anything, "token").Return(consts.UserRedisGetError, "", "", errors.New("redis down"))
			},
			wantCode: consts.UserRedisGetError,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := new(mockUserService)
			tt.mock(m)
			h := NewUserService(m)
			resp, err := h.RefreshToken(context.Background(), &user.RefreshTokenReq{RefreshToken: "token"})
			assert.Equal(t, tt.wantCode, resp.Code)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			m.AssertExpectations(t)
		})
	}
}
