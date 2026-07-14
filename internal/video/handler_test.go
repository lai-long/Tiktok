package video

import (
	video "Tiktok/kitex_gen/video"
	"Tiktok/pkg/consts"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockVideoService struct {
	mock.Mock
}

func (m *mockVideoService) VideoPublish(
	ctx context.Context,
	title string,
	description string,
	url string,
	coverURL string,
	userID string,
) (int32, error) {
	args := m.Called(ctx, title, description, url, coverURL, userID)
	return args.Get(0).(int32), args.Error(1)
}

func (m *mockVideoService) VideoList(ctx context.Context, userId string, pageSize int64, pageNum int64) (int32, []*video.VideoInfo, error) {
	args := m.Called(ctx, userId, pageSize, pageNum)
	return args.Get(0).(int32), args.Get(1).([]*video.VideoInfo), args.Error(2)
}

func (m *mockVideoService) VideoSearch(ctx context.Context, keyword string, pageNum int64, pageSize int64) (int32, []*video.VideoInfo, error) {
	args := m.Called(ctx, keyword, pageNum, pageSize)
	return args.Get(0).(int32), args.Get(1).([]*video.VideoInfo), args.Error(2)
}

func (m *mockVideoService) VideoPopular(ctx context.Context, pageNum int64, pageSize int64) (int32, []*video.VideoInfo, error) {
	args := m.Called(ctx, pageNum, pageSize)
	return args.Get(0).(int32), args.Get(1).([]*video.VideoInfo), args.Error(2)
}

func (m *mockVideoService) VideoStream(ctx context.Context) (int32, []*video.VideoInfo, error) {
	args := m.Called(ctx)
	return args.Get(0).(int32), args.Get(1).([]*video.VideoInfo), args.Error(2)
}

func (m *mockVideoService) BatchGetVideo(ctx context.Context, ids []string) (int32, []*video.VideoInfo, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).(int32), args.Get(1).([]*video.VideoInfo), args.Error(2)
}

func TestVideoServiceImpl_VideoPublish(t *testing.T) {
	tests := []struct {
		name     string
		mock     func(*mockVideoService)
		wantCode int32
		wantErr  bool
	}{
		{
			name: "success",
			mock: func(m *mockVideoService) {
				m.On("VideoPublish", mock.Anything, "title", "desc", "url", "cover", "uid").Return(consts.Success, nil)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name: "db_error",
			mock: func(m *mockVideoService) {
				m.On("VideoPublish", mock.Anything, "title", "desc", "url", "cover", "uid").Return(consts.VideoDBInsertError, errors.New("db down"))
			},
			wantCode: consts.VideoDBInsertError,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := new(mockVideoService)
			tt.mock(m)
			h := NewVideoServiceImpl(m)
			resp, err := h.VideoPublish(context.Background(), &video.VideoPublishReq{
				Title: "title", Description: "desc", VideoURL: "url", CoverURL: "cover", UserID: "uid",
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

func TestVideoServiceImpl_VideoList(t *testing.T) {
	tests := []struct {
		name     string
		mock     func(*mockVideoService)
		wantCode int32
		wantErr  bool
	}{
		{
			name: "success",
			mock: func(m *mockVideoService) {
				m.On("VideoList", mock.Anything, "uid", int64(10), int64(1)).Return(consts.Success, []*video.VideoInfo{{ID: "vid"}}, nil)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name: "db_error",
			mock: func(m *mockVideoService) {
				m.On("VideoList", mock.Anything, "uid", int64(10), int64(1)).Return(consts.VideoDBSelectError, []*video.VideoInfo{}, errors.New("db down"))
			},
			wantCode: consts.VideoDBSelectError,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := new(mockVideoService)
			tt.mock(m)
			h := NewVideoServiceImpl(m)
			resp, err := h.VideoList(context.Background(), &video.VideoListReq{UserId: "uid", PageSize: 10, PageNum: 1})
			assert.Equal(t, tt.wantCode, resp.Code)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, resp.Data.Items, 1)
			}
			m.AssertExpectations(t)
		})
	}
}

func TestVideoServiceImpl_VideoSearch(t *testing.T) {
	tests := []struct {
		name     string
		mock     func(*mockVideoService)
		wantCode int32
		wantErr  bool
	}{
		{
			name: "success",
			mock: func(m *mockVideoService) {
				m.On("VideoSearch", mock.Anything, "kw", int64(1), int64(10)).Return(consts.Success, []*video.VideoInfo{{ID: "vid"}}, nil)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name: "db_error",
			mock: func(m *mockVideoService) {
				m.On("VideoSearch", mock.Anything, "kw", int64(1), int64(10)).Return(consts.VideoDBSelectError, []*video.VideoInfo{}, errors.New("db down"))
			},
			wantCode: consts.VideoDBSelectError,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := new(mockVideoService)
			tt.mock(m)
			h := NewVideoServiceImpl(m)
			resp, err := h.VideoSearch(context.Background(), &video.VideoSearchReq{KeyWord: "kw", PageNum: 1, PageSize: 10})
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

func TestVideoServiceImpl_VideoPopular(t *testing.T) {
	tests := []struct {
		name     string
		mock     func(*mockVideoService)
		wantCode int32
		wantErr  bool
	}{
		{
			name: "success",
			mock: func(m *mockVideoService) {
				m.On("VideoPopular", mock.Anything, int64(1), int64(10)).Return(consts.Success, []*video.VideoInfo{{ID: "vid"}}, nil)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name: "db_error",
			mock: func(m *mockVideoService) {
				m.On("VideoPopular", mock.Anything, int64(1), int64(10)).Return(consts.VideoDBSelectError, []*video.VideoInfo{}, errors.New("db down"))
			},
			wantCode: consts.VideoDBSelectError,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := new(mockVideoService)
			tt.mock(m)
			h := NewVideoServiceImpl(m)
			resp, err := h.VideoPopular(context.Background(), &video.VideoHotReq{PageNum: 1, PageSize: 10})
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

func TestVideoServiceImpl_VideoStream(t *testing.T) {
	tests := []struct {
		name     string
		mock     func(*mockVideoService)
		wantCode int32
		wantErr  bool
	}{
		{
			name: "success",
			mock: func(m *mockVideoService) {
				m.On("VideoStream", mock.Anything).Return(consts.Success, []*video.VideoInfo{{ID: "vid"}}, nil)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name: "db_error",
			mock: func(m *mockVideoService) {
				m.On("VideoStream", mock.Anything).Return(consts.VideoDBSelectError, []*video.VideoInfo{}, errors.New("db down"))
			},
			wantCode: consts.VideoDBSelectError,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := new(mockVideoService)
			tt.mock(m)
			h := NewVideoServiceImpl(m)
			resp, err := h.VideoStream(context.Background(), &video.VideoStreamReq{})
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

func TestVideoServiceImpl_BatchGetVideo(t *testing.T) {
	tests := []struct {
		name     string
		mock     func(*mockVideoService)
		wantCode int32
		wantErr  bool
	}{
		{
			name: "success",
			mock: func(m *mockVideoService) {
				m.On("BatchGetVideo", mock.Anything, []string{"vid"}).Return(consts.Success, []*video.VideoInfo{{ID: "vid"}}, nil)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name: "db_error",
			mock: func(m *mockVideoService) {
				m.On("BatchGetVideo", mock.Anything, []string{"vid"}).Return(consts.VideoDBSelectError, []*video.VideoInfo{}, errors.New("db down"))
			},
			wantCode: consts.VideoDBSelectError,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := new(mockVideoService)
			tt.mock(m)
			h := NewVideoServiceImpl(m)
			resp, err := h.BatchGetVideo(context.Background(), &video.BatchGetVideoReq{Ids: []string{"vid"}})
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
