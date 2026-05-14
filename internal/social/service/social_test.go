package service

import (
	"Tiktok/pkg/consts"
	"Tiktok/pkg/entity"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSocialDb struct {
	mock.Mock
}

func (m *MockSocialDb) CreateFollowing(userId string, toUserId string) error {
	args := m.Called(userId, toUserId)
	return args.Error(0)
}

func (m *MockSocialDb) DeleteFollowing(userId string, toUserId string) error {
	args := m.Called(userId, toUserId)
	return args.Error(0)
}

func (m *MockSocialDb) FollowingList(userId string, pageNum int64, pageSize int64) ([]entity.UserEntity, error) {
	args := m.Called(userId, pageNum, pageSize)
	return args.Get(0).([]entity.UserEntity), args.Error(1)
}

func (m *MockSocialDb) FollowerList(userId string, pageNum int64, pageSize int64) ([]entity.UserEntity, error) {
	args := m.Called(userId, pageNum, pageSize)
	return args.Get(0).([]entity.UserEntity), args.Error(1)
}

func (m *MockSocialDb) FriendList(userId string, pageNum int64, pageSize int64) ([]entity.UserEntity, bool) {
	args := m.Called(userId, pageNum, pageSize)
	return args.Get(0).([]entity.UserEntity), args.Bool(1)
}

func TestRelationAction(t *testing.T) {
	tests := []struct {
		name       string
		toUserId   string
		actionType string
		userId     string
		mockSetup  func(*MockSocialDb)
		wantCode   int32
		wantErr    bool
	}{
		{
			name:       "Success_follow",
			toUserId:   "toUserID",
			actionType: "0",
			userId:     "userID",
			mockSetup: func(m *MockSocialDb) {
				m.On("CreateFollowing", "userID", "toUserID").Return(nil)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name:       "Fail_follow_db_error",
			toUserId:   "toUserID",
			actionType: "0",
			userId:     "userID",
			mockSetup: func(m *MockSocialDb) {
				m.On("CreateFollowing", "userID", "toUserID").Return(errors.New("db error"))
			},
			wantCode: consts.SocialDBInsertError,
			wantErr:  true,
		},
		{
			name:       "Success_unfollow",
			toUserId:   "toUserID",
			actionType: "1",
			userId:     "userID",
			mockSetup: func(m *MockSocialDb) {
				m.On("DeleteFollowing", "userID", "toUserID").Return(nil)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name:       "Fail_unfollow_db_error",
			toUserId:   "toUserID",
			actionType: "1",
			userId:     "userID",
			mockSetup: func(m *MockSocialDb) {
				m.On("DeleteFollowing", "userID", "toUserID").Return(errors.New("db error"))
			},
			wantCode: consts.SocialDBDeleteError,
			wantErr:  true,
		},
		{
			name:       "Fail_invalid_action",
			toUserId:   "toUserID",
			actionType: "2",
			userId:     "userID",
			mockSetup: func(m *MockSocialDb) {
			},
			wantCode: consts.SocialReqValueError,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDb := new(MockSocialDb)
			tt.mockSetup(mockDb)
			socialRepo := NewSocialRepo(mockDb)
			code, err := socialRepo.RelationAction(context.Background(), tt.toUserId, tt.actionType, tt.userId)
			assert.Equal(t, tt.wantCode, code)
			assert.Equal(t, tt.wantErr, err != nil)
			mockDb.AssertExpectations(t)
		})
	}
}

func TestFollowingList(t *testing.T) {
	tests := []struct {
		name      string
		userId    string
		pageNum   int64
		pageSize  int64
		mockSetup func(*MockSocialDb)
		wantCode  int32
		wantErr   bool
	}{
		{
			name:     "Success_following_list",
			userId:   "userID",
			pageNum:  1,
			pageSize: 10,
			mockSetup: func(m *MockSocialDb) {
				m.On("FollowingList", "userID", int64(1), int64(10)).Return([]entity.UserEntity{
					{ID: "1", Username: "user1", Created_at: time.Now(), Updated_at: time.Now()},
					{ID: "2", Username: "user2", Created_at: time.Now(), Updated_at: time.Now()},
				}, nil)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name:     "Fail_db_error",
			userId:   "userID",
			pageNum:  1,
			pageSize: 10,
			mockSetup: func(m *MockSocialDb) {
				m.On("FollowingList", "userID", int64(1), int64(10)).Return([]entity.UserEntity{}, errors.New("db error"))
			},
			wantCode: consts.SocialDBSelectError,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDb := new(MockSocialDb)
			tt.mockSetup(mockDb)
			socialRepo := NewSocialRepo(mockDb)
			code, _, err := socialRepo.FollowingList(tt.userId, tt.pageNum, tt.pageSize)
			assert.Equal(t, tt.wantCode, code)
			assert.Equal(t, tt.wantErr, err != nil)
			mockDb.AssertExpectations(t)
		})
	}
}

func TestFollowerList(t *testing.T) {
	tests := []struct {
		name      string
		userId    string
		pageNum   int64
		pageSize  int64
		mockSetup func(*MockSocialDb)
		wantCode  int32
		wantErr   bool
	}{
		{
			name:     "Success_follower_list",
			userId:   "userID",
			pageNum:  1,
			pageSize: 10,
			mockSetup: func(m *MockSocialDb) {
				m.On("FollowerList", "userID", int64(1), int64(10)).Return([]entity.UserEntity{
					{ID: "1", Username: "follower1", Created_at: time.Now(), Updated_at: time.Now()},
				}, nil)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name:     "Fail_db_error",
			userId:   "userID",
			pageNum:  1,
			pageSize: 10,
			mockSetup: func(m *MockSocialDb) {
				m.On("FollowerList", "userID", int64(1), int64(10)).Return([]entity.UserEntity{}, errors.New("db error"))
			},
			wantCode: consts.SocialDBSelectError,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDb := new(MockSocialDb)
			tt.mockSetup(mockDb)
			socialRepo := NewSocialRepo(mockDb)
			code, _, err := socialRepo.FollowerList(tt.userId, tt.pageNum, tt.pageSize)
			assert.Equal(t, tt.wantCode, code)
			assert.Equal(t, tt.wantErr, err != nil)
			mockDb.AssertExpectations(t)
		})
	}
}

func TestFriendList(t *testing.T) {
	tests := []struct {
		name      string
		userId    string
		pageNum   int64
		pageSize  int64
		mockSetup func(*MockSocialDb)
		wantCode  int32
		wantErr   bool
	}{
		{
			name:     "Success_friend_list",
			userId:   "userID",
			pageNum:  1,
			pageSize: 10,
			mockSetup: func(m *MockSocialDb) {
				m.On("FriendList", "userID", int64(1), int64(10)).Return([]entity.UserEntity{
					{ID: "1", Username: "friend1", Created_at: time.Now(), Updated_at: time.Now()},
				}, true)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name:     "Fail_friend_list_not_ok",
			userId:   "userID",
			pageNum:  1,
			pageSize: 10,
			mockSetup: func(m *MockSocialDb) {
				m.On("FriendList", "userID", int64(1), int64(10)).Return([]entity.UserEntity{}, false)
			},
			wantCode: consts.SocialDBSelectError,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDb := new(MockSocialDb)
			tt.mockSetup(mockDb)
			socialRepo := NewSocialRepo(mockDb)
			code, _, err := socialRepo.FriendList(tt.userId, tt.pageNum, tt.pageSize)
			assert.Equal(t, tt.wantCode, code)
			assert.Equal(t, tt.wantErr, err != nil)
			mockDb.AssertExpectations(t)
		})
	}
}
