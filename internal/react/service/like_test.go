package service

import (
	"Tiktok/pkg/consts"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func (m *MockLike) CommentLikeCountUp(ctx context.Context, commentId string) error {
	args := m.Called(ctx, commentId)
	return args.Error(0)
}
func (m *MockLike) CommentLikeCountDown(ctx context.Context, commentId string) error {
	args := m.Called(ctx, commentId)
	return args.Error(0)
}
func (m *MockLike) VideoLikeCountUp(ctx context.Context, videoId string) error {
	args := m.Called(ctx, videoId)
	return args.Error(0)
}
func (m *MockLike) VideoLikeCountDown(ctx context.Context, videoId string) error {
	args := m.Called(ctx, videoId)
	return args.Error(0)
}
func (m *MockLike) LikeVideoIds(ctx context.Context, userId string, pageNum int64, pageSize int64) ([]string, error) {
	args := m.Called(ctx, userId, pageNum, pageSize)
	return args.Get(0).([]string), args.Error(1)
}
func (m *MockLike) LikeCreate(ctx context.Context, userId string, targetID string, targetType string) error {
	args := m.Called(ctx, userId, targetID, targetType)
	return args.Error(0)
}
func (m *MockLike) LikeDelete(ctx context.Context, userId, targetId string, targetType string) error {
	args := m.Called(ctx, userId, targetId, targetType)
	return args.Error(0)
}
func (m *MockLike) VideoLikeSAdd(ctx context.Context, userId string, videoId string) error {
	return nil
}
func (m *MockLike) VideoDislikeSRem(ctx context.Context, userId string, videoId string) error {
	return nil
}
func (m *MockLike) VideoLikeGet(ctx context.Context, userId string) ([]string, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).([]string), args.Error(1)
}

type MockLike struct {
	mock.Mock
}

func TestLikeAction(t *testing.T) {
	tests := []struct {
		name       string
		userId     string
		targetId   string
		action     string
		targetType string
		mockSetUp  func(like *MockLike)
		wantErr    bool
		wantCode   int32
	}{
		{
			name:       "Success_likeVideo",
			userId:     "userID",
			targetId:   "targetID",
			action:     "1",
			targetType: "1",
			mockSetUp: func(m *MockLike) {
				m.On("VideoLikeCountUp", mock.Anything, "targetID").Return(nil)
				m.On("LikeCreate", mock.Anything, "userID", "targetID", "1").Return(nil)
			},
			wantErr:  false,
			wantCode: consts.Success,
		},
		{
			name:       "Success_likeComment",
			userId:     "userID",
			targetId:   "targetID",
			action:     "1",
			targetType: "2",
			mockSetUp: func(m *MockLike) {
				m.On("CommentLikeCountUp", mock.Anything, "targetID").Return(nil)
				m.On("LikeCreate", mock.Anything, "userID", "targetID", "2").Return(nil)
			},
			wantErr:  false,
			wantCode: consts.Success,
		},
		{
			name:       "Fail_likeVideo1",
			userId:     "userID",
			targetId:   "targetID",
			action:     "1",
			targetType: "1",
			mockSetUp: func(m *MockLike) {
				m.On("LikeCreate", mock.Anything, "userID", "targetID", "1").Return(errors.New("fail"))
			},
			wantErr:  true,
			wantCode: consts.ReactDBInsertError,
		},
		{
			name:       "Fail_likeComment1",
			userId:     "userID",
			targetId:   "targetID",
			action:     "1",
			targetType: "2",
			mockSetUp: func(m *MockLike) {
				m.On("LikeCreate", mock.Anything, "userID", "targetID", "2").Return(errors.New("fail"))
			},
			wantCode: consts.ReactDBInsertError,
			wantErr:  true,
		},
		{
			name:       "Success_dislikeComment",
			userId:     "userID",
			targetId:   "targetID",
			action:     "2",
			targetType: "2",
			mockSetUp: func(m *MockLike) {
				m.On("CommentLikeCountDown", mock.Anything, "targetID").Return(nil)
				m.On("LikeDelete", mock.Anything, "userID", "targetID", "2").Return(nil)
			},
			wantErr:  false,
			wantCode: consts.Success,
		},
		{
			name:       "Fail_dislikeVideo1",
			userId:     "userID",
			targetId:   "targetID",
			action:     "2",
			targetType: "1",
			mockSetUp: func(m *MockLike) {
				m.On("LikeDelete", mock.Anything, "userID", "targetID", "1").Return(errors.New("fail"))
			},
			wantErr:  true,
			wantCode: consts.ReactDBDeleteError,
		},
		{
			name:       "Fail_dislikeVideo2",
			userId:     "userID",
			targetId:   "targetID",
			action:     "2",
			targetType: "1",
			mockSetUp: func(m *MockLike) {
				m.On("LikeDelete", mock.Anything, "userID", "targetID", "1").Return(nil)
				m.On("VideoLikeCountDown", mock.Anything, "targetID").Return(errors.New("fail"))
			},
			wantErr:  true,
			wantCode: consts.ReactDBUpdateError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLike := new(MockLike)
			tt.mockSetUp(mockLike)
			like := NewLikeRepo(mockLike, mockLike, mockLike, mockLike)
			code, err := like.LikeAction(context.Background(), tt.userId, tt.targetId, tt.action, tt.targetType)
			assert.Equal(t, tt.wantCode, code)
			assert.Equal(t, tt.wantErr, err != nil)
			mockLike.AssertExpectations(t)
		})
	}
}

func TestLikeList(t *testing.T) {
	tests := []struct {
		name      string
		userId    string
		pageNum   int64
		pageSize  int64
		mockSetUp func(like *MockLike)
		wantErr   bool
		wantCode  int32
	}{
		{
			name:     "Success_likeList",
			userId:   "userID",
			pageNum:  1,
			pageSize: 10,
			mockSetUp: func(m *MockLike) {
				m.On("VideoLikeGet", mock.Anything, "userID").Return([]string{}, errors.New("cache miss"))
				m.On("LikeVideoIds", mock.Anything, "userID", int64(1), int64(10)).Return([]string{"1", "2"}, nil)
			},
			wantErr:  false,
			wantCode: consts.Success,
		},
		{
			name:     "Fail_likeList1",
			userId:   "userID",
			pageNum:  1,
			pageSize: 10,
			mockSetUp: func(m *MockLike) {
				m.On("VideoLikeGet", mock.Anything, "userID").Return([]string{}, errors.New("cache miss"))
				m.On("LikeVideoIds", mock.Anything, "userID", int64(1), int64(10)).Return([]string{"1", "2"}, errors.New("fail"))
			},
			wantErr:  true,
			wantCode: consts.ReactDBSelectError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLike := new(MockLike)
			tt.mockSetUp(mockLike)
			like := NewLikeRepo(mockLike, mockLike, mockLike, mockLike)
			code, _, _, err := like.LikeList(context.Background(), tt.userId, tt.pageNum, tt.pageSize)
			assert.Equal(t, tt.wantCode, code)
			assert.Equal(t, tt.wantErr, err != nil)
			mockLike.AssertExpectations(t)
		})
	}
}
