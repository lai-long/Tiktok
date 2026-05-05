package user

import (
	"Tiktok/biz/entity"
	"Tiktok/pkg/config"
	"Tiktok/pkg/consts"
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
func (m *MockUser) CreateUser(user entity.UserEntity) error {
	args := m.Called(user)
	return args.Error(0)
}
func (m *MockUser) GetUserByUsername(username string) (entity.UserEntity, error) {
	args := m.Called(username)
	return args.Get(0).(entity.UserEntity), args.Error(1)
}
func (m *MockUser) GetUserByUserId(userId string) (entity.UserEntity, error) {
	args := m.Called(userId)
	return args.Get(0).(entity.UserEntity), args.Error(1)
}
func (m *MockUser) UpdateUserAvatar(url string, userId interface{}) error {
	args := m.Called(url, userId)
	return args.Error(0)
}
func (m *MockUser) SaveMfaSecret(mfa string, userId string) error {
	args := m.Called(mfa, userId)
	return args.Error(0)
}
func (m *MockUser) GetMfaSecret(userId string) (string, error) {
	args := m.Called(userId)
	return args.String(0), args.Error(1)
}
func (m *MockUser) MfaBindUpdate(userId string) error {
	args := m.Called(userId)
	return args.Error(0)
}
func (m *MockUser) CheckMfaBind(userId string) (int, error) {
	args := m.Called(userId)
	return args.Int(0), args.Error(1)
}
func (m *MockUser) GetCachedUserInfo(ctx context.Context, userId string) (*entity.UserEntity, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).(*entity.UserEntity), args.Error(1)
}
func (m *MockUser) SetCachedUserInfo(ctx context.Context, userId string, info *entity.UserEntity) error {
	args := m.Called(ctx, userId, info)
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
				m.On("GetUserByUsername", "username").Return(entity.UserEntity{}, sql.ErrNoRows)
				m.On("CreateUser", mock.Anything).Return(nil)
			},
			wantErr:  false,
			wantCode: consts.Success,
		},
		{
			name:     "Fail_register_exist",
			userName: "username",
			password: "password",
			setMock: func(m *MockUser) {
				m.On("GetUserByUsername", "username").Return(entity.UserEntity{Username: "username"}, nil)
			},
			wantErr:  false,
			wantCode: consts.UserNameExists,
		},
		{
			name:     "Fail_register_err",
			userName: "username",
			password: "password",
			setMock: func(m *MockUser) {
				m.On("GetUserByUsername", "username").Return(entity.UserEntity{}, errors.New("some error"))
			},
			wantErr:  true,
			wantCode: consts.UserDBSelectError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUser := new(MockUser)
			tt.setMock(mockUser)
			user := NewUserRepo(mockUser, mockUser, mockUser)
			code, err := user.Register(tt.userName, tt.password)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.wantCode, code)
		})
	}
}

func TestLogin(t *testing.T) {
	_, err := config.Load([]string{"/home/lai-long/Tiktok/pkg/config"})
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
				m.On("GetUserByUsername", "username").Return(entity.UserEntity{ID: "ID", Username: "username", Password: hashPassword}, nil)
				m.On("CheckMfaBind", "ID").Return(0, nil)
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
			user := NewUserRepo(mockUser, mockUser, mockUser)
			code, _, _, _, err := user.Login(tt.userName, tt.password, tt.mfaCode, tt.ctx)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.wantCode, code)
			mockUser.AssertExpectations(t)
		})
	}
}
