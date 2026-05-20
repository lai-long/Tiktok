package service

import (
	"Tiktok/pkg/consts"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMfaDb struct {
	mock.Mock
}

func (m *MockMfaDb) SaveMfaSecret(mfa string, userId string) error {
	args := m.Called(mfa, userId)
	return args.Error(0)
}

func (m *MockMfaDb) GetMfaSecret(userId string) (string, error) {
	args := m.Called(userId)
	return args.String(0), args.Error(1)
}

func (m *MockMfaDb) MfaBindUpdate(userId string) error {
	args := m.Called(userId)
	return args.Error(0)
}

func (m *MockMfaDb) CheckMfaBind(userId string) (int, error) {
	args := m.Called(userId)
	return args.Int(0), args.Error(1)
}

func TestGenerateMfa(t *testing.T) {
	tests := []struct {
		name      string
		username  string
		userId    string
		mockSetup func(*MockMfaDb)
		wantCode  int32
		wantErr   bool
	}{
		{
			name:     "Success_generate",
			username: "testuser",
			userId:   "userID",
			mockSetup: func(m *MockMfaDb) {
				m.On("SaveMfaSecret", mock.Anything, "userID").Return(nil)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name:     "Fail_save_secret",
			username: "testuser",
			userId:   "userID",
			mockSetup: func(m *MockMfaDb) {
				m.On("SaveMfaSecret", mock.Anything, "userID").Return(errors.New("save error"))
			},
			wantCode: consts.UserDBUpdateError,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDb := new(MockMfaDb)
			tt.mockSetup(mockDb)
			mfaRepo := NewMfaRepo(mockDb)
			_, _, code, err := mfaRepo.GenerateMfa(tt.username, tt.userId)
			assert.Equal(t, tt.wantCode, code)
			assert.Equal(t, tt.wantErr, err != nil)
			mockDb.AssertExpectations(t)
		})
	}
}

func TestMfaBindByCode(t *testing.T) {
	tests := []struct {
		name      string
		code      string
		userId    string
		mockSetup func(*MockMfaDb)
		wantCode  int32
		wantErr   bool
	}{
		{
			name:   "Fail_get_secret_error",
			code:   "123456",
			userId: "userID",
			mockSetup: func(m *MockMfaDb) {
				m.On("GetMfaSecret", "userID").Return("", errors.New("db error"))
			},
			wantCode: consts.MfaDBSelectError,
			wantErr:  true,
		},
		{
			name:   "Fail_invalid_code",
			code:   "000000",
			userId: "userID",
			mockSetup: func(m *MockMfaDb) {
				m.On("GetMfaSecret", "userID").Return("VALID_SECRET_HERE", nil)
			},
			wantCode: consts.MfaCodeFalse,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDb := new(MockMfaDb)
			tt.mockSetup(mockDb)
			mfaRepo := NewMfaRepo(mockDb)
			code, err := mfaRepo.MfaBindByCode(tt.code, tt.userId)
			assert.Equal(t, tt.wantCode, code)
			assert.Equal(t, tt.wantErr, err != nil)
			mockDb.AssertExpectations(t)
		})
	}
}

func TestMfaBindBySecret(t *testing.T) {
	tests := []struct {
		name      string
		secret    string
		userId    string
		mockSetup func(*MockMfaDb)
		wantCode  int32
		wantErr   bool
	}{
		{
			name:   "Fail_get_secret_error",
			secret: "some_secret",
			userId: "userID",
			mockSetup: func(m *MockMfaDb) {
				m.On("GetMfaSecret", "userID").Return("", errors.New("db error"))
			},
			wantCode: consts.MfaDBSelectError,
			wantErr:  true,
		},
		{
			name:   "Fail_secret_mismatch",
			secret: "wrong_secret",
			userId: "userID",
			mockSetup: func(m *MockMfaDb) {
				m.On("GetMfaSecret", "userID").Return("correct_secret", nil)
			},
			wantCode: consts.MfaCodeFalse,
			wantErr:  false,
		},
		{
			name:   "Fail_update_error",
			secret: "correct_secret",
			userId: "userID",
			mockSetup: func(m *MockMfaDb) {
				m.On("GetMfaSecret", "userID").Return("correct_secret", nil)
				m.On("MfaBindUpdate", "userID").Return(errors.New("update error"))
			},
			wantCode: consts.UserDBUpdateError,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDb := new(MockMfaDb)
			tt.mockSetup(mockDb)
			mfaRepo := NewMfaRepo(mockDb)
			code, err := mfaRepo.MfaBindBySecret(tt.secret, tt.userId)
			assert.Equal(t, tt.wantCode, code)
			assert.Equal(t, tt.wantErr, err != nil)
			mockDb.AssertExpectations(t)
		})
	}
}
