package comment

import (
	"Tiktok/biz/entity"
	"Tiktok/pkg/consts"
	"database/sql"
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
func (m *MockCommentRepo) CreateComment(commentId string, videoId string, userId string, content string, targetType string) error {
	args := m.Called(commentId, videoId, userId, content, targetType)
	return args.Error(0)
}

func TestCommentList_Success(t *testing.T) {
	mockComment := new(MockCommentRepo)
	mockComment.On("GetComments", "123", int64(2), int64(10)).Return([]entity.CommentEntity{
		{
			UserID: "1212", TargetID: "3344", CommentID: "5566", Content: "testing",
			LikeCount: 0, CommentCount: 0, CreatedAt: "", UpdatedAt: "",
			DeletedAt: sql.NullTime{}, TargetType: "1",
		},
	}, nil)
	svc := NewCommentService(mockComment)
	code, comments, err := svc.CommentList("123", int64(10), int64(2))
	assert.NoError(t, err)
	assert.Equal(t, consts.Success, code)
	assert.Equal(t, 1, len(comments))
	mockComment.AssertExpectations(t)
}

func TestCommentList_Fail(t *testing.T) {
	mockComment := new(MockCommentRepo)
	mockComment.On("GetComments", "123", int64(2), int64(10)).Return([]entity.CommentEntity{}, errors.New("fail"))
	svc := NewCommentService(mockComment)
	code, comments, err := svc.CommentList("123", int64(10), int64(2))
	assert.Error(t, err)
	assert.Equal(t, 0, len(comments))
	assert.Equal(t, consts.ReactDBSelectError, code)
	mockComment.AssertExpectations(t)
}

func TestCommentPublish_Success(t *testing.T) {
	mockComment := new(MockCommentRepo)
	mockComment.On("CreateComment", mock.Anything, "123", "123", "123", "1").Return(nil)
	mockComment.On("VideoCommentCountUp", "123").Return(nil)
	svc := NewCommentService(mockComment)
	code, err := svc.CommentPublish("123", "123", "123", "1")
	assert.NoError(t, err)
	assert.Equal(t, consts.Success, code)
	mockComment.AssertExpectations(t)
}
