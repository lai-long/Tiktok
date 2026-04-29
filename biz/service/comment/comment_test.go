package comment

import (
	"Tiktok/biz/entity"
	"Tiktok/pkg/consts"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCommentRepo struct {
	mock.Mock
}

func (m *MockCommentRepo) GetComments(videoId string, pageNum int64, pageSize int64) ([]entity.CommentEntity, error) {
	args := m.Called(videoId, pageNum, pageSize)
	return args.Get(0).([]entity.CommentEntity), args.Error(1)
}
func (m *MockCommentRepo) CommentDelete(commentId string) error {
	args := m.Called(commentId)
	return args.Error(0)
}
func (m *MockCommentRepo) GetCommentById(commentId string) (entity.CommentEntity, error) {
	args := m.Called(commentId)
	return args.Get(0).(entity.CommentEntity), args.Error(1)
}
func (m *MockCommentRepo) VideoCommentCountUp(videoId string) error {
	args := m.Called(videoId)
	return args.Error(0)
}
func (m *MockCommentRepo) CommentCommentCountUp(commentId string) error {
	args := m.Called(commentId)
	return args.Error(0)
}
func (m *MockCommentRepo) VideoCommentCountDown(videoId string) error {
	args := m.Called(videoId)
	return args.Error(0)
}
func (m *MockCommentRepo) CommentCommentCountDown(commentId string) error {
	args := m.Called(commentId)
	return args.Error(0)
}
func (m *MockCommentRepo) CreateComment(commentId string, targetId string, userId string, content string, targetType string) error {
	args := m.Called(commentId, targetId, userId, content, targetType)
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
				m.On("GetComments", "123", int64(2), int64(10)).Return([]entity.CommentEntity{{}}, nil)
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
				m.On("GetComments", "123", int64(2), int64(10)).Return([]entity.CommentEntity{}, errors.New("fail"))
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
			code, comments, err := svc.CommentList(tt.targetId, tt.pageSize, tt.pageNum)
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
				m.On("CreateComment", mock.Anything, "123", "1212", "testing", "1").Return(nil)
				m.On("VideoCommentCountUp", "123").Return(nil)
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
				m.On("CreateComment", mock.Anything, "123", "1212", "testing", "1").Return(errors.New("fail"))
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
				m.On("CreateComment", mock.Anything, "123", "1212", "testing", "1").Return(nil)
				m.On("VideoCommentCountUp", "123").Return(errors.New("fail"))
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
				m.On("CreateComment", mock.Anything, "123", "1212", "testing", "2").Return(nil)
				m.On("CommentCommentCountUp", "123").Return(nil)
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
				m.On("CreateComment", mock.Anything, "123", "1212", "testing", "2").Return(nil)
				m.On("CommentCommentCountUp", "123").Return(errors.New("fail"))
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
			code, err := svc.CommentPublish(tt.targetId, tt.userId, tt.content, tt.targetType)
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
				m.On("GetCommentById", "123").Return(entity.CommentEntity{UserID: "1212"}, nil)
				m.On("CommentDelete", "123").Return(nil)
				m.On("VideoCommentCountDown", "1234").Return(nil)
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
				m.On("GetCommentById", "123").Return(entity.CommentEntity{UserID: "1212"}, nil)
				m.On("CommentDelete", "123").Return(errors.New("fail"))
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
				m.On("GetCommentById", "123").Return(entity.CommentEntity{UserID: "1212"}, nil)
				m.On("CommentDelete", "123").Return(nil)
				m.On("VideoCommentCountDown", "1234").Return(errors.New("fail"))
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
			code, err := svc.CommentDelete(tt.commentId, tt.targetId, tt.userId, tt.targetType)
			assert.Equal(t, tt.wantCode, code)
			assert.Equal(t, tt.wantErr, err != nil)
			mockComment.AssertExpectations(t)
		})
	}
}
