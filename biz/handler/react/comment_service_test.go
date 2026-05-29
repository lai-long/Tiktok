package react

import (
	"context"
	"testing"

	"Tiktok/pkg/consts"

	Rpc "Tiktok/biz/rpc"
	react2 "Tiktok/kitex_gen/react"

	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/cloudwego/kitex/client/callopt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// jsonBody and assertResponseCode are in like_service_test.go
type MockCommentClient struct {
	mock.Mock
}

func (m *MockCommentClient) CommentPublish(
	ctx context.Context, req *react2.CommentPublishReq, callOptions ...callopt.Option,
) (*react2.CommentPublishResp, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*react2.CommentPublishResp), args.Error(1)
}

func (m *MockCommentClient) CommentList(
	ctx context.Context, req *react2.CommentListReq, callOptions ...callopt.Option,
) (*react2.CommentListResp, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*react2.CommentListResp), args.Error(1)
}

func (m *MockCommentClient) CommentDelete(
	ctx context.Context, req *react2.CommentDeleteReq, callOptions ...callopt.Option,
) (*react2.CommentDeleteResp, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*react2.CommentDeleteResp), args.Error(1)
}

func TestCommentPublish(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(*MockCommentClient)
		wantCode  int32
		wantErr   bool
	}{
		{
			name: "Success_publish",
			mockSetup: func(m *MockCommentClient) {
				m.On("CommentPublish", mock.Anything, mock.MatchedBy(func(req *react2.CommentPublishReq) bool {
					return req.UserID != ""
				})).Return(&react2.CommentPublishResp{Code: consts.Success}, nil)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name: "Fail_publish_error",
			mockSetup: func(m *MockCommentClient) {
				m.On("CommentPublish", mock.Anything, mock.MatchedBy(func(req *react2.CommentPublishReq) bool {
					return req.UserID != ""
				})).Return(&react2.CommentPublishResp{Code: consts.ReactDBSelectError}, nil)
			},
			wantCode: consts.ReactDBSelectError,
			wantErr:  true,
		},
		{
			name: "Fail_publish_rpc_error",
			mockSetup: func(m *MockCommentClient) {
				m.On("CommentPublish", mock.Anything, mock.MatchedBy(func(req *react2.CommentPublishReq) bool {
					return req.UserID != ""
				})).Return(nil, assert.AnError)
			},
			wantCode: consts.ReactError,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockCommentClient)
			Rpc.SetCommentClient(mockClient)
			tt.mockSetup(mockClient)

			body, header := jsonBody(map[string]interface{}{})
			c := ut.CreateUtRequestContext("POST", "/comment/publish", body, header...)
			c.Set("userid", "123")

			CommentPublish(context.Background(), c)

			assert.Equal(t, 200, c.Response.StatusCode())
			assertResponseCode(t, c, tt.wantCode)
			mockClient.AssertExpectations(t)
		})
	}
}

func TestCommentList(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(*MockCommentClient)
		wantCode  int32
		wantErr   bool
	}{
		{
			name: "Success_list",
			mockSetup: func(m *MockCommentClient) {
				m.On("CommentList", mock.Anything, mock.Anything).Return(&react2.CommentListResp{
					Code: consts.Success,
					Data: &react2.CommentData{
						Items: []*react2.CommentInfo{},
						Total: 0,
					},
				}, nil)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name: "Fail_list_error",
			mockSetup: func(m *MockCommentClient) {
				m.On("CommentList", mock.Anything, mock.Anything).Return(&react2.CommentListResp{Code: consts.ReactDBSelectError}, nil)
			},
			wantCode: consts.ReactDBSelectError,
			wantErr:  true,
		},
		{
			name: "Fail_list_rpc_error",
			mockSetup: func(m *MockCommentClient) {
				m.On("CommentList", mock.Anything, mock.Anything).Return(nil, assert.AnError)
			},
			wantCode: consts.ReactDBSelectError,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockCommentClient)
			Rpc.SetCommentClient(mockClient)
			tt.mockSetup(mockClient)

			c := ut.CreateUtRequestContext("GET", "/comment/list", nil)

			CommentList(context.Background(), c)

			assert.Equal(t, 200, c.Response.StatusCode())
			assertResponseCode(t, c, tt.wantCode)
			mockClient.AssertExpectations(t)
		})
	}
}

func TestCommentDelete(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(*MockCommentClient)
		wantCode  int32
		wantErr   bool
	}{
		{
			name: "Success_delete",
			mockSetup: func(m *MockCommentClient) {
				m.On("CommentDelete", mock.Anything, mock.MatchedBy(func(req *react2.CommentDeleteReq) bool {
					return req.UserID != ""
				})).Return(&react2.CommentDeleteResp{Code: consts.Success}, nil)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name: "Fail_delete_error",
			mockSetup: func(m *MockCommentClient) {
				m.On("CommentDelete", mock.Anything, mock.MatchedBy(func(req *react2.CommentDeleteReq) bool {
					return req.UserID != ""
				})).Return(&react2.CommentDeleteResp{Code: consts.ReactDBSelectError}, nil)
			},
			wantCode: consts.ReactDBSelectError,
			wantErr:  true,
		},
		{
			name: "Fail_delete_rpc_error",
			mockSetup: func(m *MockCommentClient) {
				m.On("CommentDelete", mock.Anything, mock.MatchedBy(func(req *react2.CommentDeleteReq) bool {
					return req.UserID != ""
				})).Return(nil, assert.AnError)
			},
			wantCode: consts.ReactError,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockCommentClient)
			Rpc.SetCommentClient(mockClient)
			tt.mockSetup(mockClient)

			body, header := jsonBody(map[string]interface{}{})
			c := ut.CreateUtRequestContext("DELETE", "/comment/delete", body, header...)
			c.Set("userid", "123")

			CommentDelete(context.Background(), c)

			assert.Equal(t, 200, c.Response.StatusCode())
			assertResponseCode(t, c, tt.wantCode)
			mockClient.AssertExpectations(t)
		})
	}
}

// TestCommentList_BindError - GET 请求，无 required 字段，RPC 仍会被调用
func TestCommentList_BindError(t *testing.T) {
	mockClient := new(MockCommentClient)
	Rpc.SetCommentClient(mockClient)
	mockClient.On("CommentList", mock.Anything, mock.Anything).Return(
		&react2.CommentListResp{
			Code: consts.Success,
			Data: &react2.CommentData{Items: []*react2.CommentInfo{}, Total: 0},
		}, nil,
	)

	c := ut.CreateUtRequestContext("GET", "/comment/list", nil)

	CommentList(context.Background(), c)

	assert.Equal(t, 200, c.Response.StatusCode())
	mockClient.AssertExpectations(t)
}
