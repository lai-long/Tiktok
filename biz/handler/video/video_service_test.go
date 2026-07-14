package video

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"Tiktok/pkg/consts"

	Rpc "Tiktok/biz/rpc"
	video2 "Tiktok/kitex_gen/video"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/cloudwego/kitex/client/callopt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// jsonBody 辅助函数
func jsonBody(v any) (*ut.Body, []ut.Header) {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return &ut.Body{Body: bytes.NewReader(b), Len: len(b)},
		[]ut.Header{{Key: "Content-Type", Value: "application/json"}}
}

// assertResponseCode 验证响应体中的 base.code 字段
func assertResponseCode(t *testing.T, c *app.RequestContext, wantCode int32) {
	var result map[string]interface{}
	assert.NoError(t, json.Unmarshal(c.Response.Body(), &result))
	base := result["base"].(map[string]interface{})
	assert.Equal(t, wantCode, int32(base["code"].(float64)))
}

type MockVideoClient struct {
	mock.Mock
}

func (m *MockVideoClient) VideoPublish(
	ctx context.Context, req *video2.VideoPublishReq, callOptions ...callopt.Option,
) (*video2.VideoPublishResp, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*video2.VideoPublishResp), args.Error(1)
}

func (m *MockVideoClient) VideoList(ctx context.Context, req *video2.VideoListReq, callOptions ...callopt.Option) (*video2.VideoListResp, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*video2.VideoListResp), args.Error(1)
}

func (m *MockVideoClient) VideoSearch(
	ctx context.Context, req *video2.VideoSearchReq, callOptions ...callopt.Option,
) (*video2.VideoSearchResp, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*video2.VideoSearchResp), args.Error(1)
}

func (m *MockVideoClient) VideoPopular(ctx context.Context, req *video2.VideoHotReq, callOptions ...callopt.Option) (*video2.VideoHotResp, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*video2.VideoHotResp), args.Error(1)
}

func (m *MockVideoClient) VideoStream(
	ctx context.Context, req *video2.VideoStreamReq, callOptions ...callopt.Option,
) (*video2.VideoStreamResp, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*video2.VideoStreamResp), args.Error(1)
}

func (m *MockVideoClient) BatchGetVideo(
	ctx context.Context, req *video2.BatchGetVideoReq, callOptions ...callopt.Option,
) (*video2.BatchGetVideoResp, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*video2.BatchGetVideoResp), args.Error(1)
}

func TestVideoList(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(*MockVideoClient)
		wantCode  int32
		wantErr   bool
	}{
		{
			name: "Success_list",
			mockSetup: func(m *MockVideoClient) {
				m.On("VideoList", mock.Anything, mock.MatchedBy(func(req *video2.VideoListReq) bool {
					return req.UserId != ""
				})).Return(&video2.VideoListResp{
					Code: consts.Success,
					Data: &video2.VideoData{
						Items: []*video2.VideoInfo{},
						Total: 0,
					},
				}, nil)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name: "Fail_list_error",
			mockSetup: func(m *MockVideoClient) {
				m.On("VideoList", mock.Anything, mock.MatchedBy(func(req *video2.VideoListReq) bool {
					return req.UserId != ""
				})).Return(&video2.VideoListResp{Code: consts.VideoDBSelectError}, nil)
			},
			wantCode: consts.VideoDBSelectError,
			wantErr:  true,
		},
		{
			name: "Fail_list_rpc_error",
			mockSetup: func(m *MockVideoClient) {
				m.On("VideoList", mock.Anything, mock.MatchedBy(func(req *video2.VideoListReq) bool {
					return req.UserId != ""
				})).Return(nil, assert.AnError)
			},
			wantCode: consts.VideoDBSelectError,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockVideoClient)
			Rpc.SetVideoClient(mockClient)
			tt.mockSetup(mockClient)

			c := ut.CreateUtRequestContext("GET", "/video/list?user_id=123", nil)
			c.Set("userid", "123")

			VideoList(context.Background(), c)

			assert.Equal(t, 200, c.Response.StatusCode())
			assertResponseCode(t, c, tt.wantCode)
			mockClient.AssertExpectations(t)
		})
	}
}

func TestVideoSearch(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(*MockVideoClient)
		wantCode  int32
		wantErr   bool
	}{
		{
			name: "Success_search",
			mockSetup: func(m *MockVideoClient) {
				m.On("VideoSearch", mock.Anything, mock.MatchedBy(func(req *video2.VideoSearchReq) bool {
					return req.KeyWord != ""
				})).Return(&video2.VideoSearchResp{
					Code: consts.Success,
					Data: &video2.VideoData{
						Items: []*video2.VideoInfo{},
						Total: 0,
					},
				}, nil)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name: "Fail_search_error",
			mockSetup: func(m *MockVideoClient) {
				m.On("VideoSearch", mock.Anything, mock.MatchedBy(func(req *video2.VideoSearchReq) bool {
					return req.KeyWord != ""
				})).Return(&video2.VideoSearchResp{Code: consts.VideoDBSelectError}, nil)
			},
			wantCode: consts.VideoDBSelectError,
			wantErr:  true,
		},
		{
			name: "Fail_search_rpc_error",
			mockSetup: func(m *MockVideoClient) {
				m.On("VideoSearch", mock.Anything, mock.MatchedBy(func(req *video2.VideoSearchReq) bool {
					return req.KeyWord != ""
				})).Return(nil, assert.AnError)
			},
			wantCode: consts.VideoDBSelectError,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockVideoClient)
			Rpc.SetVideoClient(mockClient)
			tt.mockSetup(mockClient)

			body, header := jsonBody(map[string]any{"keyword": "test"})
			c := ut.CreateUtRequestContext("POST", "/video/search", body, header...)

			VideoSearch(context.Background(), c)

			assert.Equal(t, 200, c.Response.StatusCode())
			assertResponseCode(t, c, tt.wantCode)
			mockClient.AssertExpectations(t)
		})
	}
}

func TestVideoPopular(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(*MockVideoClient)
		wantCode  int32
		wantErr   bool
	}{
		{
			name: "Success_popular",
			mockSetup: func(m *MockVideoClient) {
				m.On("VideoPopular", mock.Anything, mock.Anything).Return(&video2.VideoHotResp{
					Code: consts.Success,
					Data: &video2.VideoData{
						Items: []*video2.VideoInfo{},
						Total: 0,
					},
				}, nil)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name: "Fail_popular_error",
			mockSetup: func(m *MockVideoClient) {
				m.On("VideoPopular", mock.Anything, mock.Anything).Return(&video2.VideoHotResp{Code: consts.VideoDBSelectError}, nil)
			},
			wantCode: consts.VideoDBSelectError,
			wantErr:  true,
		},
		{
			name: "Fail_popular_rpc_error",
			mockSetup: func(m *MockVideoClient) {
				m.On("VideoPopular", mock.Anything, mock.Anything).Return(nil, assert.AnError)
			},
			wantCode: consts.VideoDBSelectError,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockVideoClient)
			Rpc.SetVideoClient(mockClient)
			tt.mockSetup(mockClient)

			c := ut.CreateUtRequestContext("GET", "/video/popular", nil)

			VideoPopular(context.Background(), c)

			assert.Equal(t, 200, c.Response.StatusCode())
			assertResponseCode(t, c, tt.wantCode)
			mockClient.AssertExpectations(t)
		})
	}
}

func TestVideoStream(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(*MockVideoClient)
		wantCode  int32
		wantErr   bool
	}{
		{
			name: "Success_stream",
			mockSetup: func(m *MockVideoClient) {
				m.On("VideoStream", mock.Anything, mock.Anything).Return(&video2.VideoStreamResp{
					Code: consts.Success,
					Data: &video2.VideoData{
						Items: []*video2.VideoInfo{},
						Total: 0,
					},
				}, nil)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name: "Fail_stream_error",
			mockSetup: func(m *MockVideoClient) {
				m.On("VideoStream", mock.Anything, mock.Anything).Return(&video2.VideoStreamResp{Code: consts.VideoDBSelectError}, nil)
			},
			wantCode: consts.VideoDBSelectError,
			wantErr:  true,
		},
		{
			name: "Fail_stream_rpc_error",
			mockSetup: func(m *MockVideoClient) {
				m.On("VideoStream", mock.Anything, mock.Anything).Return(nil, assert.AnError)
			},
			wantCode: consts.VideoDBSelectError,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockVideoClient)
			Rpc.SetVideoClient(mockClient)
			tt.mockSetup(mockClient)

			c := ut.CreateUtRequestContext("GET", "/video/feed", nil)

			VideoStream(context.Background(), c)

			assert.Equal(t, 200, c.Response.StatusCode())
			assertResponseCode(t, c, tt.wantCode)
			mockClient.AssertExpectations(t)
		})
	}
}

func TestVideoList_BindError(t *testing.T) {
	mockClient := new(MockVideoClient)
	Rpc.SetVideoClient(mockClient)
	mockClient.On("VideoList", mock.Anything, mock.Anything).Return(
		&video2.VideoListResp{
			Code: consts.Success,
			Data: &video2.VideoData{Items: []*video2.VideoInfo{}, Total: 0},
		}, nil,
	)

	c := ut.CreateUtRequestContext("GET", "/video/list", nil)
	c.Set("userid", "123")

	VideoList(context.Background(), c)

	assert.Equal(t, 200, c.Response.StatusCode())
	mockClient.AssertExpectations(t)
}

func TestVideoPopular_BindError(t *testing.T) {
	mockClient := new(MockVideoClient)
	Rpc.SetVideoClient(mockClient)
	mockClient.On("VideoPopular", mock.Anything, mock.Anything).Return(
		&video2.VideoHotResp{
			Code: consts.Success,
			Data: &video2.VideoData{Items: []*video2.VideoInfo{}, Total: 0},
		}, nil,
	)

	c := ut.CreateUtRequestContext("GET", "/video/popular", nil)

	VideoPopular(context.Background(), c)

	assert.Equal(t, 200, c.Response.StatusCode())
	mockClient.AssertExpectations(t)
}

func TestVideoStream_BindError(t *testing.T) {
	mockClient := new(MockVideoClient)
	Rpc.SetVideoClient(mockClient)
	mockClient.On("VideoStream", mock.Anything, mock.Anything).Return(
		&video2.VideoStreamResp{
			Code: consts.Success,
			Data: &video2.VideoData{Items: []*video2.VideoInfo{}, Total: 0},
		}, nil,
	)

	c := ut.CreateUtRequestContext("GET", "/video/feed", nil)

	VideoStream(context.Background(), c)

	assert.Equal(t, 200, c.Response.StatusCode())
	mockClient.AssertExpectations(t)
}
