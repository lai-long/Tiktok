package handler

import (
	react "Tiktok/kitex_gen/react"
	"Tiktok/pkg/consts"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockCommentService struct {
	mock.Mock
}

func (m *mockCommentService) CommentPublish(ctx context.Context, targetId, userId, content, targetType string) (int32, error) {
	args := m.Called(ctx, targetId, userId, content, targetType)
	return args.Get(0).(int32), args.Error(1)
}

func (m *mockCommentService) CommentList(ctx context.Context, targetId string, pageSize int64, pageNum int64) (int32, []*react.CommentInfo, error) {
	args := m.Called(ctx, targetId, pageSize, pageNum)
	return args.Get(0).(int32), args.Get(1).([]*react.CommentInfo), args.Error(2)
}

func (m *mockCommentService) CommentDelete(ctx context.Context, commentId string, targetId string, userId string, targetType string) (int32, error) {
	args := m.Called(ctx, commentId, targetId, userId, targetType)
	return args.Get(0).(int32), args.Error(1)
}

type mockLikeService struct {
	mock.Mock
}

func (m *mockLikeService) LikeAction(ctx context.Context, userId string, targetId string, action string, targetType string) (int32, error) {
	args := m.Called(ctx, userId, targetId, action, targetType)
	return args.Get(0).(int32), args.Error(1)
}

func (m *mockLikeService) LikeList(ctx context.Context, userId string, pageNum int64, pageSize int64) (int32, []string, int64, error) {
	args := m.Called(ctx, userId, pageNum, pageSize)
	return args.Get(0).(int32), args.Get(1).([]string), args.Get(2).(int64), args.Error(3)
}

func TestCommentServiceImpl_CommentPublish(t *testing.T) {
	tests := []struct {
		name     string
		mock     func(*mockCommentService)
		wantCode int32
		wantErr  bool
	}{
		{
			name: "success",
			mock: func(m *mockCommentService) {
				m.On("CommentPublish", mock.Anything, "target", "uid", "hello", "video").Return(consts.Success, nil)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name: "db_error",
			mock: func(m *mockCommentService) {
				m.On("CommentPublish", mock.Anything, "target", "uid", "hello", "video").Return(consts.ReactDBInsertError, errors.New("db down"))
			},
			wantCode: consts.ReactDBInsertError,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := new(mockCommentService)
			tt.mock(m)
			h := NewCommentService(m)
			resp, err := h.CommentPublish(context.Background(), &react.CommentPublishReq{
				TargetAt: "target", UserID: "uid", Content: "hello", TargetType: "video",
			})
			assert.Equal(t, tt.wantCode, resp.Code)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			m.AssertExpectations(t)
		})
	}
}

func TestCommentServiceImpl_CommentList(t *testing.T) {
	tests := []struct {
		name     string
		mock     func(*mockCommentService)
		wantCode int32
		wantErr  bool
	}{
		{
			name: "success",
			mock: func(m *mockCommentService) {
				m.On("CommentList", mock.Anything, "target", int64(10), int64(1)).Return(consts.Success, []*react.CommentInfo{{CommentId: "cid"}}, nil)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name: "db_error",
			mock: func(m *mockCommentService) {
				m.On("CommentList", mock.Anything, "target", int64(10), int64(1)).Return(consts.ReactDBSelectError, []*react.CommentInfo{}, errors.New("db down"))
			},
			wantCode: consts.ReactDBSelectError,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := new(mockCommentService)
			tt.mock(m)
			h := NewCommentService(m)
			resp, err := h.CommentList(context.Background(), &react.CommentListReq{TargetAt: "target", PageSize: 10, PageNum: 1})
			assert.Equal(t, tt.wantCode, resp.Code)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			m.AssertExpectations(t)
		})
	}
}

func TestCommentServiceImpl_CommentDelete(t *testing.T) {
	tests := []struct {
		name     string
		mock     func(*mockCommentService)
		wantCode int32
		wantErr  bool
	}{
		{
			name: "success",
			mock: func(m *mockCommentService) {
				m.On("CommentDelete", mock.Anything, "cid", "target", "uid", "video").Return(consts.Success, nil)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name: "db_error",
			mock: func(m *mockCommentService) {
				m.On("CommentDelete", mock.Anything, "cid", "target", "uid", "video").Return(consts.ReactDBDeleteError, errors.New("db down"))
			},
			wantCode: consts.ReactDBDeleteError,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := new(mockCommentService)
			tt.mock(m)
			h := NewCommentService(m)
			resp, err := h.CommentDelete(context.Background(), &react.CommentDeleteReq{
				CommentId: "cid", TargetAt: "target", UserID: "uid", TargetType: "video",
			})
			assert.Equal(t, tt.wantCode, resp.Code)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			m.AssertExpectations(t)
		})
	}
}

func TestLikeServiceImpl_LikeAction(t *testing.T) {
	tests := []struct {
		name     string
		mock     func(*mockLikeService)
		wantCode int32
		wantErr  bool
	}{
		{
			name: "success",
			mock: func(m *mockLikeService) {
				m.On("LikeAction", mock.Anything, "uid", "target", "1", "video").Return(consts.Success, nil)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name: "db_error",
			mock: func(m *mockLikeService) {
				m.On("LikeAction", mock.Anything, "uid", "target", "1", "video").Return(consts.ReactDBInsertError, errors.New("db down"))
			},
			wantCode: consts.ReactDBInsertError,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := new(mockLikeService)
			tt.mock(m)
			h := NewLikeService(m)
			resp, err := h.LikeAction(context.Background(), &react.LikeActionReq{UserID: "uid", TargetAt: "target", ActionType: "1", TargetType: "video"})
			assert.Equal(t, tt.wantCode, resp.Code)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			m.AssertExpectations(t)
		})
	}
}

func TestLikeServiceImpl_LikeList(t *testing.T) {
	tests := []struct {
		name     string
		mock     func(*mockLikeService)
		wantCode int32
		wantErr  bool
	}{
		{
			name: "success",
			mock: func(m *mockLikeService) {
				m.On("LikeList", mock.Anything, "uid", int64(1), int64(10)).Return(consts.Success, []string{"vid"}, int64(1), nil)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name: "db_error",
			mock: func(m *mockLikeService) {
				m.On("LikeList", mock.Anything, "uid", int64(1), int64(10)).Return(consts.ReactDBSelectError, []string{}, int64(0), errors.New("db down"))
			},
			wantCode: consts.ReactDBSelectError,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := new(mockLikeService)
			tt.mock(m)
			h := NewLikeService(m)
			resp, err := h.LikeList(context.Background(), &react.LikeListReq{UserId: "uid", PageNum: 1, PageSize: 10})
			assert.Equal(t, tt.wantCode, resp.Code)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			m.AssertExpectations(t)
		})
	}
}
