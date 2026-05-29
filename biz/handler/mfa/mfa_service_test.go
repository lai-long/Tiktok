package mfa

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"Tiktok/pkg/consts"

	Rpc "Tiktok/biz/rpc"
	mfa2 "Tiktok/kitex_gen/mfa"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/cloudwego/kitex/client/callopt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func jsonBody(v any) (*ut.Body, []ut.Header) {
	b, _ := json.Marshal(v)
	return &ut.Body{Body: bytes.NewReader(b), Len: len(b)},
		[]ut.Header{{Key: "Content-Type", Value: "application/json"}}
}

func assertResponseCode(t *testing.T, c *app.RequestContext, wantCode int32) {
	var result map[string]interface{}
	assert.NoError(t, json.Unmarshal(c.Response.Body(), &result))
	base := result["base"].(map[string]interface{})
	assert.Equal(t, wantCode, int32(base["code"].(float64)))
}

type MockMfaClient struct {
	mock.Mock
}

func (m *MockMfaClient) MfaQrcode(ctx context.Context, req *mfa2.MfaQrcodeReq, callOptions ...callopt.Option) (*mfa2.MfaQrcodeResp, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*mfa2.MfaQrcodeResp), args.Error(1)
}

func (m *MockMfaClient) MfaBind(ctx context.Context, req *mfa2.MfaBindReq, callOptions ...callopt.Option) (*mfa2.MfaBindResp, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*mfa2.MfaBindResp), args.Error(1)
}

func (m *MockMfaClient) MfaConfirm(ctx context.Context, req *mfa2.MfaConfirmReq, callOptions ...callopt.Option) (*mfa2.MfaConfirmResp, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*mfa2.MfaConfirmResp), args.Error(1)
}

func TestMfaQrcode(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(*MockMfaClient)
		wantCode  int32
		wantErr   bool
	}{
		{
			name: "Success_qrcode",
			mockSetup: func(m *MockMfaClient) {
				m.On("MfaQrcode", mock.Anything, mock.MatchedBy(func(req *mfa2.MfaQrcodeReq) bool {
					return req.UserID != "" && req.UserName != ""
				})).Return(&mfa2.MfaQrcodeResp{
					Code: consts.Success,
					Data: &mfa2.MfaData{
						Secret: "test_secret",
						Qrcode: "test_qrcode",
					},
				}, nil)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name: "Fail_qrcode_error",
			mockSetup: func(m *MockMfaClient) {
				m.On("MfaQrcode", mock.Anything, mock.MatchedBy(func(req *mfa2.MfaQrcodeReq) bool {
					return req.UserID != "" && req.UserName != ""
				})).Return(&mfa2.MfaQrcodeResp{Code: consts.MfaCodeFalse}, nil)
			},
			wantCode: consts.MfaCodeFalse,
			wantErr:  true,
		},
		{
			name: "Fail_qrcode_rpc_error",
			mockSetup: func(m *MockMfaClient) {
				m.On("MfaQrcode", mock.Anything, mock.MatchedBy(func(req *mfa2.MfaQrcodeReq) bool {
					return req.UserID != "" && req.UserName != ""
				})).Return(nil, assert.AnError)
			},
			wantCode: consts.MfaCodeFalse,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockMfaClient)
			Rpc.SetMfaClient(mockClient)
			tt.mockSetup(mockClient)

			c := ut.CreateUtRequestContext("GET", "/auth/mfa/qrcode", nil)
			c.Set("userid", "123")
			c.Set("username", "testuser")

			MfaQrcode(context.Background(), c)

			assert.Equal(t, 200, c.Response.StatusCode())
			assertResponseCode(t, c, tt.wantCode)
			mockClient.AssertExpectations(t)
		})
	}
}

func TestMfaBind(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(*MockMfaClient)
		wantCode  int32
		wantErr   bool
	}{
		{
			name: "Fail_bind_invalid_request",
			mockSetup: func(m *MockMfaClient) {
				// No RPC call when both Secret and Code are empty
			},
			wantCode: consts.UserReqValidError,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockMfaClient)
			Rpc.SetMfaClient(mockClient)
			tt.mockSetup(mockClient)

			body, header := jsonBody(map[string]interface{}{})
			c := ut.CreateUtRequestContext("POST", "/auth/mfa/bind", body, header...)
			c.Set("userid", "123")

			MfaBind(context.Background(), c)

			assert.Equal(t, 200, c.Response.StatusCode())
			assertResponseCode(t, c, tt.wantCode)
			mockClient.AssertExpectations(t)
		})
	}
}

// TestMfaQrcode_BindError - GET 请求，无 required 字段，RPC 仍会被调用
func TestMfaQrcode_BindError(t *testing.T) {
	mockClient := new(MockMfaClient)
	Rpc.SetMfaClient(mockClient)
	mockClient.On("MfaQrcode", mock.Anything, mock.Anything).Return(&mfa2.MfaQrcodeResp{Code: consts.Success, Data: &mfa2.MfaData{Secret: "test", Qrcode: "test"}}, nil)

	c := ut.CreateUtRequestContext("GET", "/auth/mfa/qrcode", nil)
	c.Set("userid", "123")
	c.Set("username", "testuser")

	MfaQrcode(context.Background(), c)

	assert.Equal(t, 200, c.Response.StatusCode())
	mockClient.AssertExpectations(t)
}
