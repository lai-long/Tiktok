package service

import (
	"Tiktok/pkg/consts"
	"Tiktok/pkg/entity"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCommentRepo struct {
	mock.Mock
}

func (m *MockCommentRepo) GetComments(ctx context.Context, videoId string, pageNum int64, pageSize int64) ([]entity.CommentEntity, error) {
	args := m.Called(ctx, videoId, pageNum, pageSize)
	return args.Get(0).([]entity.CommentEntity), args.Error(1)
}
func (m *MockCommentRepo) CommentDelete(ctx context.Context, commentId string) error {
	args := m.Called(ctx, commentId)
	return args.Error(0)
}
func (m *MockCommentRepo) GetCommentById(ctx context.Context, commentId string) (entity.CommentEntity, error) {
	args := m.Called(ctx, commentId)
	return args.Get(0).(entity.CommentEntity), args.Error(1)
}
func (m *MockCommentRepo) VideoCommentCountUp(ctx context.Context, videoId string) error {
	args := m.Called(ctx, videoId)
	return args.Error(0)
}
func (m *MockCommentRepo) CommentCommentCountUp(ctx context.Context, commentId string) error {
	args := m.Called(ctx, commentId)
	return args.Error(0)
}
func (m *MockCommentRepo) VideoCommentCountDown(ctx context.Context, videoId string) error {
	args := m.Called(ctx, videoId)
	return args.Error(0)
}
func (m *MockCommentRepo) CommentCommentCountDown(ctx context.Context, commentId string) error {
	args := m.Called(ctx, commentId)
	return args.Error(0)
}
func (m *MockCommentRepo) CreateComment(ctx context.Context, commentId string, targetId string, userId string, content string, targetType string) error {
	args := m.Called(ctx, commentId, targetId, userId, content, targetType)
	return args.Error(0)
}

func TestCommentList(t *testing.T) {
	tests := []struct {
		name      string
		targetId  string
		pageNum   int64
		pageSize  int64
		mockSetup func(*MockCommentRepo)
		wantCode  int32
		wantLen   int
		wantErr   bool
	}{
		{
			name:     "Success",
			targetId: "123",
			pageNum:  2,
			pageSize: 10,
			mockSetup: func(m *MockCommentRepo) {
				m.On("GetComments", mock.Anything, "123", int64(2), int64(10)).Return([]entity.CommentEntity{{}}, nil)
			},
			wantCode: consts.Success,
			wantLen:  1,
			wantErr:  false,
		},
		{
			name:     "Fail",
			targetId: "123",
			pageNum:  2,
			pageSize: 10,
			mockSetup: func(m *MockCommentRepo) {
				m.On("GetComments", mock.Anything, "123", int64(2), int64(10)).Return([]entity.CommentEntity{}, errors.New("fail"))
			},
			wantCode: consts.ReactDBSelectError,
			wantLen:  0,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockComment := new(MockCommentRepo)
			tt.mockSetup(mockComment)
			svc := NewCommentService(mockComment)
			code, comments, err := svc.CommentList(context.Background(), tt.targetId, tt.pageSize, tt.pageNum)
			assert.Equal(t, tt.wantCode, code)
			assert.Equal(t, tt.wantLen, len(comments))
			assert.Equal(t, tt.wantErr, err != nil)
			mockComment.AssertExpectations(t)
		})
	}
}

func TestCommentPublish(t *testing.T) {
	tests := []struct {
		name       string
		targetId   string
		userId     string
		content    string
		targetType string
		mockSetup  func(*MockCommentRepo)
		wantCode   int32
		wantErr    bool
	}{
		{
			name:       "Success_VideoComment",
			targetId:   "123",
			userId:     "1212",
			content:    "testing",
			targetType: "1",
			mockSetup: func(m *MockCommentRepo) {
				m.On("CreateComment", mock.Anything, mock.Anything, "123", "1212", "testing", "1").Return(nil)
				m.On("VideoCommentCountUp", mock.Anything, "123").Return(nil)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name:       "Fail_CommentInsert",
			targetId:   "123",
			userId:     "1212",
			content:    "testing",
			targetType: "1",
			mockSetup: func(m *MockCommentRepo) {
				m.On("CreateComment", mock.Anything, mock.Anything, "123", "1212", "testing", "1").Return(errors.New("fail"))
			},
			wantCode: consts.ReactDBInsertError,
			wantErr:  true,
		},
		{
			name:       "Fail_VideoCommentUpdate",
			targetId:   "123",
			userId:     "1212",
			content:    "testing",
			targetType: "1",
			mockSetup: func(m *MockCommentRepo) {
				m.On("CreateComment", mock.Anything, mock.Anything, "123", "1212", "testing", "1").Return(nil)
				m.On("VideoCommentCountUp", mock.Anything, "123").Return(errors.New("fail"))
			},
			wantCode: consts.ReactDBUpdateError,
			wantErr:  true,
		},
		{
			name:       "Success_CommentComment",
			targetId:   "123",
			userId:     "1212",
			content:    "testing",
			targetType: "2",
			mockSetup: func(m *MockCommentRepo) {
				m.On("CreateComment", mock.Anything, mock.Anything, "123", "1212", "testing", "2").Return(nil)
				m.On("CommentCommentCountUp", mock.Anything, "123").Return(nil)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name:       "Fail_CommentCommentUpdate",
			targetId:   "123",
			userId:     "1212",
			content:    "testing",
			targetType: "2",
			mockSetup: func(m *MockCommentRepo) {
				m.On("CreateComment", mock.Anything, mock.Anything, "123", "1212", "testing", "2").Return(nil)
				m.On("CommentCommentCountUp", mock.Anything, "123").Return(errors.New("fail"))
			},
			wantCode: consts.ReactDBUpdateError,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockComment := new(MockCommentRepo)
			tt.mockSetup(mockComment)
			svc := NewCommentService(mockComment)
			code, err := svc.CommentPublish(context.Background(), tt.targetId, tt.userId, tt.content, tt.targetType)
			assert.Equal(t, tt.wantCode, code)
			assert.Equal(t, tt.wantErr, err != nil)
			mockComment.AssertExpectations(t)
		})
	}
}

func TestCommentDelete(t *testing.T) {
	tests := []struct {
		name       string
		commentId  string
		targetId   string
		userId     string
		targetType string
		mockSetup  func(*MockCommentRepo)
		wantErr    bool
		wantCode   int32
	}{
		{
			name:       "Success_Video",
			commentId:  "123",
			targetId:   "1234",
			userId:     "1212",
			targetType: "1",
			mockSetup: func(m *MockCommentRepo) {
				m.On("GetCommentById", mock.Anything, "123").Return(entity.CommentEntity{UserID: "1212"}, nil)
				m.On("CommentDelete", mock.Anything, "123").Return(nil)
				m.On("VideoCommentCountDown", mock.Anything, "1234").Return(nil)
			},
			wantErr:  false,
			wantCode: consts.Success,
		},
		{
			name:       "Fail_VideoCommentDelete",
			commentId:  "123",
			targetId:   "1234",
			userId:     "1212",
			targetType: "1",
			mockSetup: func(m *MockCommentRepo) {
				m.On("GetCommentById", mock.Anything, "123").Return(entity.CommentEntity{UserID: "1212"}, nil)
				m.On("CommentDelete", mock.Anything, "123").Return(errors.New("fail"))
			},
			wantErr:  true,
			wantCode: consts.ReactDBDeleteError,
		},
		{
			name:       "Fail_VideoCommentCountDown",
			commentId:  "123",
			targetId:   "1234",
			userId:     "1212",
			targetType: "1",
			mockSetup: func(m *MockCommentRepo) {
				m.On("GetCommentById", mock.Anything, "123").Return(entity.CommentEntity{UserID: "1212"}, nil)
				m.On("CommentDelete", mock.Anything, "123").Return(nil)
				m.On("VideoCommentCountDown", mock.Anything, "1234").Return(errors.New("fail"))
			},
			wantErr:  true,
			wantCode: consts.ReactDBUpdateError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockComment := new(MockCommentRepo)
			tt.mockSetup(mockComment)
			svc := NewCommentService(mockComment)
			code, err := svc.CommentDelete(context.Background(), tt.commentId, tt.targetId, tt.userId, tt.targetType)
			assert.Equal(t, tt.wantCode, code)
			assert.Equal(t, tt.wantErr, err != nil)
			mockComment.AssertExpectations(t)
		})
	}
}
