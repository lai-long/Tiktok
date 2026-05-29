package react

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"Tiktok/pkg/consts"

	Rpc "Tiktok/biz/rpc"
	react2 "Tiktok/kitex_gen/react"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/cloudwego/kitex/client/callopt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// jsonBody 辅助函数
func jsonBody(v any) (*ut.Body, []ut.Header) {
	b, _ := json.Marshal(v)
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

type MockLikeClient struct {
	mock.Mock
}

func (m *MockLikeClient) LikeAction(ctx context.Context, req *react2.LikeActionReq, callOptions ...callopt.Option) (*react2.LikeActionResp, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*react2.LikeActionResp), args.Error(1)
}

func (m *MockLikeClient) LikeList(ctx context.Context, req *react2.LikeListReq, callOptions ...callopt.Option) (*react2.LikeListResp, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*react2.LikeListResp), args.Error(1)
}

func TestLikeAction(t *testing.T) {
	tests := []struct {
		name      string
		setBody   bool
		mockSetup func(*MockLikeClient)
		wantCode  int32
		wantErr   bool
	}{
		{
			name:    "Fail_action_target_type_empty",
			setBody: false, // TargetType will be empty, handler returns early
			mockSetup: func(m *MockLikeClient) {
				// No RPC call expected when TargetType is empty
			},
			wantCode: consts.ReactReqValueError,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockLikeClient)
			Rpc.SetLikeClient(mockClient)
			tt.mockSetup(mockClient)

			body, header := jsonBody(map[string]interface{}{})
			c := ut.CreateUtRequestContext("POST", "/like/action", body, header...)
			c.Set("userid", "123")

			LikeAction(context.Background(), c)

			assert.Equal(t, 200, c.Response.StatusCode())
			assertResponseCode(t, c, tt.wantCode)
			mockClient.AssertExpectations(t)
		})
	}
}

func TestLikeList(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(*MockLikeClient)
		wantCode  int32
		wantErr   bool
	}{
		{
			name: "Success_list_empty",
			mockSetup: func(m *MockLikeClient) {
				m.On("LikeList", mock.Anything, mock.MatchedBy(func(req *react2.LikeListReq) bool {
					return req.UserId != ""
				})).Return(&react2.LikeListResp{
					Code: consts.Success,
					Data: &react2.LikeVideoData{
						VideoIds: []string{},
						Total:    0,
					},
				}, nil)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name: "Fail_list_error",
			mockSetup: func(m *MockLikeClient) {
				m.On("LikeList", mock.Anything, mock.MatchedBy(func(req *react2.LikeListReq) bool {
					return req.UserId != ""
				})).Return(&react2.LikeListResp{Code: consts.ReactDBSelectError}, nil)
			},
			wantCode: consts.ReactDBSelectError,
			wantErr:  true,
		},
		{
			name: "Fail_list_rpc_error",
			mockSetup: func(m *MockLikeClient) {
				m.On("LikeList", mock.Anything, mock.MatchedBy(func(req *react2.LikeListReq) bool {
					return req.UserId != ""
				})).Return(nil, assert.AnError)
			},
			wantCode: consts.ReactDBSelectError,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockLikeClient)
			Rpc.SetLikeClient(mockClient)
			tt.mockSetup(mockClient)

			c := ut.CreateUtRequestContext("GET", "/like/list?user_id=123", nil)

			LikeList(context.Background(), c)

			assert.Equal(t, 200, c.Response.StatusCode())
			assertResponseCode(t, c, tt.wantCode)
			mockClient.AssertExpectations(t)
		})
	}
}

// TestLikeList_BindError - GET 请求，无 required 字段，RPC 仍会被调用
func TestLikeList_BindError(t *testing.T) {
	mockClient := new(MockLikeClient)
	Rpc.SetLikeClient(mockClient)
	mockClient.On("LikeList", mock.Anything, mock.Anything).Return(
		&react2.LikeListResp{
			Code: consts.Success,
			Data: &react2.LikeVideoData{VideoIds: []string{}, Total: 0},
		}, nil,
	)

	c := ut.CreateUtRequestContext("GET", "/like/list", nil)

	LikeList(context.Background(), c)

	assert.Equal(t, 200, c.Response.StatusCode())
	mockClient.AssertExpectations(t)
}
