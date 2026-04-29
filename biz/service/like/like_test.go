package like

import (
	"Tiktok/biz/entity"

	"Tiktok/pkg/consts"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func (m *MockLike) CommentLikeCountUp(commentId string) error {
	args := m.Called(commentId)
	return args.Error(0)
}
func (m *MockLike) CommentLikeCountDown(commentId string) error {
	args := m.Called(commentId)
	return args.Error(0)
}
func (m *MockLike) VideoLikeCountUp(videoId string) error {
	args := m.Called(videoId)
	return args.Error(0)
}
func (m *MockLike) VideoLikeCountDown(videoId string) error {
	args := m.Called(videoId)
	return args.Error(0)
}
func (m *MockLike) LikeVideoIds(userId string, pageNum int64, pageSize int64) ([]string, error) {
	args := m.Called(userId, pageNum, pageSize)
	return args.Get(0).([]string), args.Error(1)
}
func (m *MockLike) LikeVideos(videoId []string) (bool, []entity.VideoEntity) {
	args := m.Called(videoId)
	return args.Bool(0), args.Get(1).([]entity.VideoEntity)
}
func (m *MockLike) LikeCreate(userId string, targetID string, targetType string) error {
	args := m.Called(userId, targetID, targetType)
	return args.Error(0)
}
func (m *MockLike) LikeDelete(userId, targetId string, targetType string) error {
	args := m.Called(userId, targetId, targetType)
	return args.Error(0)
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
				m.On("VideoLikeCountUp", "targetID").Return(nil)
				m.On("LikeCreate", "userID", "targetID", "1").Return(nil)
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
				m.On("CommentLikeCountUp", "targetID").Return(nil)
				m.On("LikeCreate", "userID", "targetID", "2").Return(nil)
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
				m.On("LikeCreate", "userID", "targetID", "1").Return(errors.New("fail"))
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
				m.On("LikeCreate", "userID", "targetID", "2").Return(errors.New("fail"))
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
				m.On("CommentLikeCountDown", "targetID").Return(nil)
				m.On("LikeDelete", "userID", "targetID", "2").Return(nil)
			},
			wantErr:  false,
			wantCode: consts.Success,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLike := new(MockLike)
			tt.mockSetUp(mockLike)
			like := NewLikeRepo(mockLike, mockLike, mockLike)
			code, err := like.LikeAction(tt.userId, tt.targetId, tt.action, tt.targetType)
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
				m.On("LikeVideoIds", "userID", int64(1), int64(10)).Return([]string{"1", "2"}, nil)
				m.On("LikeVideos", []string{"1", "2"}).Return(true, []entity.VideoEntity{})
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
				m.On("LikeVideoIds", "userID", int64(1), int64(10)).Return([]string{"1", "2"}, errors.New("fail"))
			},
			wantErr:  true,
			wantCode: consts.ReactDBSelectError,
		},
		{
			name:     "Fail_likeList2",
			userId:   "userID",
			pageNum:  1,
			pageSize: 10,
			mockSetUp: func(m *MockLike) {
				m.On("LikeVideoIds", "userID", int64(1), int64(10)).Return([]string{"1", "2"}, nil)
				m.On("LikeVideos", []string{"1", "2"}).Return(false, []entity.VideoEntity{})
			},
			wantErr:  true,
			wantCode: consts.ReactDBSelectError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLike := new(MockLike)
			tt.mockSetUp(mockLike)
			like := NewLikeRepo(mockLike, mockLike, mockLike)
			code, _, err := like.LikeList(tt.userId, tt.pageNum, tt.pageSize)
			assert.Equal(t, tt.wantCode, code)
			assert.Equal(t, tt.wantErr, err != nil)
			mockLike.AssertExpectations(t)
		})
	}
}
