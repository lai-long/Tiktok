package service

import (
	"Tiktok/pkg/config"
	"Tiktok/pkg/consts"
	"Tiktok/pkg/entity"
	"Tiktok/pkg/utils"
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUser struct {
	mock.Mock
}

func (m *MockUser) UserTokenSet(ctx context.Context, refreshToken string, userId string) error {
	args := m.Called(ctx, refreshToken, userId)
	return args.Error(0)
}
func (m *MockUser) UserGetByRefreshToken(ctx context.Context, refreshToken string) (userId string, err error) {
	args := m.Called(ctx, refreshToken)
	return args.String(0), args.Error(1)
}
func (m *MockUser) UserTokenDelete(ctx context.Context, refreshToken string) error {
	args := m.Called(ctx, refreshToken)
	return args.Error(0)
}
func (m *MockUser) CreateUser(ctx context.Context, user entity.UserEntity) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}
func (m *MockUser) GetUserByUsername(ctx context.Context, username string) (entity.UserEntity, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(entity.UserEntity), args.Error(1)
}
func (m *MockUser) GetUserByUserId(ctx context.Context, userId string) (entity.UserEntity, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).(entity.UserEntity), args.Error(1)
}
func (m *MockUser) UpdateUserAvatar(ctx context.Context, url string, userId interface{}) error {
	args := m.Called(ctx, url, userId)
	return args.Error(0)
}
func (m *MockUser) CheckMfaBind(ctx context.Context, userId string) (int, error) {
	args := m.Called(ctx, userId)
	return args.Int(0), args.Error(1)
}
func (m *MockUser) GetMfaSecret(ctx context.Context, userId string) (string, error) {
	args := m.Called(ctx, userId)
	return args.String(0), args.Error(1)
}
func (m *MockUser) GetCachedUserInfo(ctx context.Context, userId string) (*entity.UserEntity, error) {
	args := m.Called(ctx, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.UserEntity), args.Error(1)
}
func (m *MockUser) SetCachedUserInfo(ctx context.Context, userId string, info *entity.UserEntity) error {
	args := m.Called(ctx, userId, info)
	return args.Error(0)
}
func (m *MockUser) DelCachedUserInfo(ctx context.Context, userId string) error {
	args := m.Called(ctx, userId)
	return args.Error(0)
}

func TestRegister(t *testing.T) {
	tests := []struct {
		name     string
		userName string
		password string
		setMock  func(*MockUser)
		wantErr  bool
		wantCode int32
	}{
		{
			name:     "Success_register",
			userName: "username",
			password: "password",
			setMock: func(m *MockUser) {
				m.On("GetUserByUsername", mock.Anything, "username").Return(entity.UserEntity{}, sql.ErrNoRows)
				m.On("CreateUser", mock.Anything, mock.Anything).Return(nil)
			},
			wantErr:  false,
			wantCode: consts.Success,
		},
		{
			name:     "Fail_register_exist",
			userName: "username",
			password: "password",
			setMock: func(m *MockUser) {
				m.On("GetUserByUsername", mock.Anything, "username").Return(entity.UserEntity{Username: "username"}, nil)
			},
			wantErr:  false,
			wantCode: consts.UserNameExists,
		},
		{
			name:     "Fail_register_err",
			userName: "username",
			password: "password",
			setMock: func(m *MockUser) {
				m.On("GetUserByUsername", mock.Anything, "username").Return(entity.UserEntity{}, errors.New("some error"))
			},
			wantErr:  true,
			wantCode: consts.UserDBSelectError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUser := new(MockUser)
			tt.setMock(mockUser)
			user := NewUserRepo(mockUser, mockUser)
			code, err := user.Register(context.Background(), tt.userName, tt.password)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.wantCode, code)
		})
	}
}

func TestLogin(t *testing.T) {
	_, err := config.Load([]string{"../../../pkg/config"})
	if err != nil {
		t.Fatal(err)
	}
	hashPassword, _ := utils.HashPassword("password")
	tests := []struct {
		name     string
		userName string
		password string
		mfaCode  string
		ctx      context.Context
		setMock  func(*MockUser)
		wantErr  bool
		wantCode int32
	}{
		{
			name:     "Success_login_normal",
			userName: "username",
			password: "password",
			ctx:      context.Background(),
			setMock: func(m *MockUser) {
				m.On("GetUserByUsername", mock.Anything, "username").Return(entity.UserEntity{ID: "ID", Username: "username", Password: hashPassword}, nil)
				m.On("CheckMfaBind", mock.Anything, "ID").Return(0, nil)
				m.On("UserTokenSet", mock.Anything, mock.Anything, "ID").Return(nil)
			},
			wantErr:  false,
			wantCode: consts.Success,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUser := new(MockUser)
			tt.setMock(mockUser)
			user := NewUserRepo(mockUser, mockUser)
			code, _, _, _, err := user.Login(tt.ctx, tt.userName, tt.password, tt.mfaCode)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.wantCode, code)
			mockUser.AssertExpectations(t)
		})
	}
}

func TestUserInfo(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		userId   string
		setMock  func(*MockUser)
		wantCode int32
		wantErr  bool
	}{
		{
			name:   "Success_cache_hit",
			ctx:    context.Background(),
			userId: "userID",
			setMock: func(m *MockUser) {
				m.On("GetCachedUserInfo", mock.Anything, "userID").Return(&entity.UserEntity{ID: "userID", Username: "cachedUser"}, nil)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name:   "Success_cache_miss_db_hit",
			ctx:    context.Background(),
			userId: "userID",
			setMock: func(m *MockUser) {
				m.On("GetCachedUserInfo", mock.Anything, "userID").Return(nil, errors.New("cache miss"))
				m.On("GetUserByUserId", mock.Anything, "userID").Return(entity.UserEntity{ID: "userID", Username: "dbUser"}, nil)
				m.On("SetCachedUserInfo", mock.Anything, "userID", mock.Anything).Return(nil)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name:   "Fail_db_error",
			ctx:    context.Background(),
			userId: "userID",
			setMock: func(m *MockUser) {
				m.On("GetCachedUserInfo", mock.Anything, "userID").Return(nil, errors.New("cache miss"))
				m.On("GetUserByUserId", mock.Anything, "userID").Return(entity.UserEntity{}, errors.New("db error"))
			},
			wantCode: consts.UserDBSelectError,
			wantErr:  true,
		},
		{
			name:   "Fail_cache_set_error",
			ctx:    context.Background(),
			userId: "userID",
			setMock: func(m *MockUser) {
				m.On("GetCachedUserInfo", mock.Anything, "userID").Return(nil, errors.New("cache miss"))
				m.On("GetUserByUserId", mock.Anything, "userID").Return(entity.UserEntity{ID: "userID", Username: "dbUser"}, nil)
				m.On("SetCachedUserInfo", mock.Anything, "userID", mock.Anything).Return(errors.New("cache set error"))
			},
			wantCode: consts.UserDBSelectError,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUser := new(MockUser)
			tt.setMock(mockUser)
			user := NewUserRepo(mockUser, mockUser)
			_, code, err := user.UserInfo(tt.ctx, tt.userId)
			assert.Equal(t, tt.wantCode, code)
			assert.Equal(t, tt.wantErr, err != nil)
			mockUser.AssertExpectations(t)
		})
	}
}

func TestUserAvatar(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		userID   string
		setMock  func(*MockUser)
		wantCode int32
		wantErr  bool
	}{
		{
			name:   "Success_avatar_update",
			url:    "http://example.com/avatar.jpg",
			userID: "userID",
			setMock: func(m *MockUser) {
				m.On("UpdateUserAvatar", mock.Anything, "http://example.com/avatar.jpg", "userID").Return(nil)
				m.On("DelCachedUserInfo", mock.Anything, "userID").Return(nil)
				userEntity := entity.UserEntity{ID: "userID", Username: "user", Avatar_url: "http://example.com/avatar.jpg"}
				m.On("GetUserByUserId", mock.Anything, "userID").Return(userEntity, nil)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name:   "Fail_update_avatar_error",
			url:    "http://example.com/avatar.jpg",
			userID: "userID",
			setMock: func(m *MockUser) {
				m.On("UpdateUserAvatar", mock.Anything, "http://example.com/avatar.jpg", "userID").Return(errors.New("update error"))
			},
			wantCode: consts.UserDBUpdateError,
			wantErr:  true,
		},
		{
			name:   "Fail_get_user_after_update",
			url:    "http://example.com/avatar.jpg",
			userID: "userID",
			setMock: func(m *MockUser) {
				m.On("UpdateUserAvatar", mock.Anything, "http://example.com/avatar.jpg", "userID").Return(nil)
				m.On("DelCachedUserInfo", mock.Anything, "userID").Return(nil)
				m.On("GetUserByUserId", mock.Anything, "userID").Return(entity.UserEntity{}, errors.New("get error"))
			},
			wantCode: consts.UserDBSelectError,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUser := new(MockUser)
			tt.setMock(mockUser)
			user := NewUserRepo(mockUser, mockUser)
			code, _, err := user.UserAvatar(context.Background(), tt.url, tt.userID)
			assert.Equal(t, tt.wantCode, code)
			assert.Equal(t, tt.wantErr, err != nil)
			mockUser.AssertExpectations(t)
		})
	}
}

func TestRefreshToken(t *testing.T) {
	_, err := config.Load([]string{"../../../pkg/config"})
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name         string
		ctx          context.Context
		refreshToken string
		setMock      func(*MockUser)
		wantCode     int32
		wantErr      bool
	}{
		{
			name:         "Success_refresh_token",
			ctx:          context.Background(),
			refreshToken: "valid_refresh_token",
			setMock: func(m *MockUser) {
				m.On("UserGetByRefreshToken", mock.Anything, "valid_refresh_token").Return("userID", nil)
				m.On("GetUserByUserId", mock.Anything, "userID").Return(entity.UserEntity{ID: "userID", Username: "user"}, nil)
				m.On("UserTokenDelete", mock.Anything, "valid_refresh_token").Return(nil)
				m.On("UserTokenSet", mock.Anything, mock.Anything, "userID").Return(nil)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name:         "Fail_redis_get_error",
			ctx:          context.Background(),
			refreshToken: "invalid_token",
			setMock: func(m *MockUser) {
				m.On("UserGetByRefreshToken", mock.Anything, "invalid_token").Return("", errors.New("redis get error"))
			},
			wantCode: consts.UserRedisGetError,
			wantErr:  true,
		},
		{
			name:         "Fail_db_error",
			ctx:          context.Background(),
			refreshToken: "valid_token",
			setMock: func(m *MockUser) {
				m.On("UserGetByRefreshToken", mock.Anything, "valid_token").Return("userID", nil)
				m.On("GetUserByUserId", mock.Anything, "userID").Return(entity.UserEntity{}, errors.New("db error"))
			},
			wantCode: consts.UserDBSelectError,
			wantErr:  true,
		},
		{
			name:         "Fail_token_delete_error",
			ctx:          context.Background(),
			refreshToken: "valid_token",
			setMock: func(m *MockUser) {
				m.On("UserGetByRefreshToken", mock.Anything, "valid_token").Return("userID", nil)
				m.On("GetUserByUserId", mock.Anything, "userID").Return(entity.UserEntity{ID: "userID", Username: "user"}, nil)
				m.On("UserTokenDelete", mock.Anything, "valid_token").Return(errors.New("delete error"))
			},
			wantCode: consts.UserRedisDelError,
			wantErr:  true,
		},
		{
			name:         "Fail_redis_set_error",
			ctx:          context.Background(),
			refreshToken: "valid_token",
			setMock: func(m *MockUser) {
				m.On("UserGetByRefreshToken", mock.Anything, "valid_token").Return("userID", nil)
				m.On("GetUserByUserId", mock.Anything, "userID").Return(entity.UserEntity{ID: "userID", Username: "user"}, nil)
				m.On("UserTokenDelete", mock.Anything, "valid_token").Return(nil)
				m.On("UserTokenSet", mock.Anything, mock.Anything, "userID").Return(errors.New("set error"))
			},
			wantCode: consts.UserRedisSetError,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUser := new(MockUser)
			tt.setMock(mockUser)
			user := NewUserRepo(mockUser, mockUser)
			code, _, _, err := user.RefreshToken(tt.ctx, tt.refreshToken)
			assert.Equal(t, tt.wantCode, code)
			assert.Equal(t, tt.wantErr, err != nil)
			mockUser.AssertExpectations(t)
		})
	}
}
