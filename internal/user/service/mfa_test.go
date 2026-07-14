package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"Tiktok/pkg/consts"

	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const testMfaSecret = "JBSWY3DPEHPK3PXP"

func TestGenerateMfa(t *testing.T) {
	tests := []struct {
		name     string
		username string
		userId   string
		setMock  func(*MockUser)
		wantErr  bool
		wantCode int32
	}{
		{
			name:     "Success_generate_mfa",
			username: "testuser",
			userId:   "userID",
			setMock: func(m *MockUser) {
				m.On("SaveMfaSecret", mock.Anything, mock.MatchedBy(func(secret string) bool { return secret != "" }), "userID").Return(nil)
			},
			wantErr:  false,
			wantCode: consts.Success,
		},
		{
			name:     "Fail_save_secret",
			username: "testuser",
			userId:   "userID",
			setMock: func(m *MockUser) {
				m.On("SaveMfaSecret", mock.Anything, mock.MatchedBy(func(secret string) bool { return secret != "" }), "userID").Return(errors.New("db error"))
			},
			wantErr:  true,
			wantCode: consts.UserDBUpdateError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUser := new(MockUser)
			tt.setMock(mockUser)
			user := NewUserRepo(mockUser, mockUser)
			url, secret, code, err := user.GenerateMfa(context.Background(), tt.username, tt.userId)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.wantCode, code)
			if tt.wantCode == consts.Success {
				assert.NotEmpty(t, secret)
				assert.Contains(t, url, "otpauth://totp/")
				genCode, genErr := totp.GenerateCode(secret, time.Now())
				assert.NoError(t, genErr)
				assert.True(t, totp.Validate(genCode, secret))
			}
			mockUser.AssertExpectations(t)
		})
	}
}

func TestMfaBindByCode(t *testing.T) {
	validCode, err := totp.GenerateCode(testMfaSecret, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name     string
		code     string
		userId   string
		setMock  func(*MockUser)
		wantErr  bool
		wantCode int32
	}{
		{
			name:   "Success_bind_by_code",
			code:   validCode,
			userId: "userID",
			setMock: func(m *MockUser) {
				m.On("GetMfaSecret", mock.Anything, "userID").Return(testMfaSecret, nil)
				m.On("MfaBindUpdate", mock.Anything, "userID").Return(nil)
			},
			wantErr:  false,
			wantCode: consts.Success,
		},
		{
			name:   "Fail_invalid_code",
			code:   "000000",
			userId: "userID",
			setMock: func(m *MockUser) {
				m.On("GetMfaSecret", mock.Anything, "userID").Return(testMfaSecret, nil)
			},
			wantErr:  false,
			wantCode: consts.MfaCodeFalse,
		},
		{
			name:   "Fail_get_secret",
			code:   validCode,
			userId: "userID",
			setMock: func(m *MockUser) {
				m.On("GetMfaSecret", mock.Anything, "userID").Return("", errors.New("db error"))
			},
			wantErr:  true,
			wantCode: consts.MfaDBSelectError,
		},
		{
			name:   "Fail_bind_update",
			code:   validCode,
			userId: "userID",
			setMock: func(m *MockUser) {
				m.On("GetMfaSecret", mock.Anything, "userID").Return(testMfaSecret, nil)
				m.On("MfaBindUpdate", mock.Anything, "userID").Return(errors.New("db error"))
			},
			wantErr:  true,
			wantCode: consts.UserDBUpdateError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUser := new(MockUser)
			tt.setMock(mockUser)
			user := NewUserRepo(mockUser, mockUser)
			code, err := user.MfaBindByCode(context.Background(), tt.code, tt.userId)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.wantCode, code)
			mockUser.AssertExpectations(t)
		})
	}
}

func TestMfaBindBySecret(t *testing.T) {
	tests := []struct {
		name     string
		secret   string
		userId   string
		setMock  func(*MockUser)
		wantErr  bool
		wantCode int32
	}{
		{
			name:   "Success_bind_by_secret",
			secret: testMfaSecret,
			userId: "userID",
			setMock: func(m *MockUser) {
				m.On("GetMfaSecret", mock.Anything, "userID").Return(testMfaSecret, nil)
				m.On("MfaBindUpdate", mock.Anything, "userID").Return(nil)
			},
			wantErr:  false,
			wantCode: consts.Success,
		},
		{
			name:   "Fail_secret_mismatch",
			secret: "MISMATCH",
			userId: "userID",
			setMock: func(m *MockUser) {
				m.On("GetMfaSecret", mock.Anything, "userID").Return(testMfaSecret, nil)
			},
			wantErr:  false,
			wantCode: consts.MfaCodeFalse,
		},
		{
			name:   "Fail_get_secret",
			secret: testMfaSecret,
			userId: "userID",
			setMock: func(m *MockUser) {
				m.On("GetMfaSecret", mock.Anything, "userID").Return("", errors.New("db error"))
			},
			wantErr:  true,
			wantCode: consts.MfaDBSelectError,
		},
		{
			name:   "Fail_bind_update",
			secret: testMfaSecret,
			userId: "userID",
			setMock: func(m *MockUser) {
				m.On("GetMfaSecret", mock.Anything, "userID").Return(testMfaSecret, nil)
				m.On("MfaBindUpdate", mock.Anything, "userID").Return(errors.New("db error"))
			},
			wantErr:  true,
			wantCode: consts.UserDBUpdateError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUser := new(MockUser)
			tt.setMock(mockUser)
			user := NewUserRepo(mockUser, mockUser)
			code, err := user.MfaBindBySecret(context.Background(), tt.secret, tt.userId)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.wantCode, code)
			mockUser.AssertExpectations(t)
		})
	}
}

func TestMfaConfirm(t *testing.T) {
	validCode, err := totp.GenerateCode(testMfaSecret, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name     string
		mfaCode  string
		userID   string
		setMock  func(*MockUser)
		wantErr  bool
		wantCode int32
	}{
		{
			name:    "Success_not_bound",
			mfaCode: "",
			userID:  "userID",
			setMock: func(m *MockUser) {
				m.On("CheckMfaBind", mock.Anything, "userID").Return(0, nil)
			},
			wantErr:  false,
			wantCode: consts.Success,
		},
		{
			name:    "Fail_bound_empty_code",
			mfaCode: "",
			userID:  "userID",
			setMock: func(m *MockUser) {
				m.On("CheckMfaBind", mock.Anything, "userID").Return(1, nil)
			},
			wantErr:  false,
			wantCode: consts.MfaReqValidError,
		},
		{
			name:    "Fail_bound_invalid_code",
			mfaCode: "000000",
			userID:  "userID",
			setMock: func(m *MockUser) {
				m.On("CheckMfaBind", mock.Anything, "userID").Return(1, nil)
				m.On("GetMfaSecret", mock.Anything, "userID").Return(testMfaSecret, nil)
			},
			wantErr:  false,
			wantCode: consts.MfaCodeFalse,
		},
		{
			name:    "Success_bound_valid_code",
			mfaCode: validCode,
			userID:  "userID",
			setMock: func(m *MockUser) {
				m.On("CheckMfaBind", mock.Anything, "userID").Return(1, nil)
				m.On("GetMfaSecret", mock.Anything, "userID").Return(testMfaSecret, nil)
			},
			wantErr:  false,
			wantCode: consts.Success,
		},
		{
			name:    "Fail_check_bind",
			mfaCode: "",
			userID:  "userID",
			setMock: func(m *MockUser) {
				m.On("CheckMfaBind", mock.Anything, "userID").Return(0, errors.New("db error"))
			},
			wantErr:  true,
			wantCode: consts.MfaDBSelectError,
		},
		{
			name:    "Fail_bound_get_secret",
			mfaCode: "000000",
			userID:  "userID",
			setMock: func(m *MockUser) {
				m.On("CheckMfaBind", mock.Anything, "userID").Return(1, nil)
				m.On("GetMfaSecret", mock.Anything, "userID").Return("", errors.New("db error"))
			},
			wantErr:  true,
			wantCode: consts.MfaDBSelectError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUser := new(MockUser)
			tt.setMock(mockUser)
			user := NewUserRepo(mockUser, mockUser)
			code, err := user.MfaConfirm(context.Background(), tt.mfaCode, tt.userID)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.wantCode, code)
			mockUser.AssertExpectations(t)
		})
	}
}
