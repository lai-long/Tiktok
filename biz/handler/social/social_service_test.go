package social

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"Tiktok/pkg/consts"

	Rpc "Tiktok/biz/rpc"
	social2 "Tiktok/kitex_gen/social"

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

type MockSocialClient struct {
	mock.Mock
}

func (m *MockSocialClient) RelationAction(
	ctx context.Context, req *social2.RelationActionReq, callOptions ...callopt.Option,
) (*social2.RelationActionResp, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*social2.RelationActionResp), args.Error(1)
}

func (m *MockSocialClient) FollowingList(
	ctx context.Context, req *social2.FollowingListReq, callOptions ...callopt.Option,
) (*social2.FollowingListResp, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*social2.FollowingListResp), args.Error(1)
}

func (m *MockSocialClient) FollowerList(
	ctx context.Context, req *social2.FollowerListReq, callOptions ...callopt.Option,
) (*social2.FollowerListResp, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*social2.FollowerListResp), args.Error(1)
}

func (m *MockSocialClient) FriendList(
	ctx context.Context, req *social2.FriendListReq, callOptions ...callopt.Option,
) (*social2.FriendListResp, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*social2.FriendListResp), args.Error(1)
}

func TestRelationAction(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(*MockSocialClient)
		wantCode  int32
		wantErr   bool
	}{
		{
			name: "Success_follow",
			mockSetup: func(m *MockSocialClient) {
				m.On("RelationAction", mock.Anything, mock.MatchedBy(func(req *social2.RelationActionReq) bool {
					return req.UserId != ""
				})).Return(&social2.RelationActionResp{Code: consts.Success}, nil)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name: "Success_unfollow",
			mockSetup: func(m *MockSocialClient) {
				m.On("RelationAction", mock.Anything, mock.MatchedBy(func(req *social2.RelationActionReq) bool {
					return req.UserId != ""
				})).Return(&social2.RelationActionResp{Code: consts.Success}, nil)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name: "Fail_action_error",
			mockSetup: func(m *MockSocialClient) {
				m.On("RelationAction", mock.Anything, mock.MatchedBy(func(req *social2.RelationActionReq) bool {
					return req.UserId != ""
				})).Return(&social2.RelationActionResp{Code: consts.SocialDBSelectError}, nil)
			},
			wantCode: consts.SocialDBSelectError,
			wantErr:  true,
		},
		{
			name: "Fail_action_rpc_error",
			mockSetup: func(m *MockSocialClient) {
				m.On("RelationAction", mock.Anything, mock.MatchedBy(func(req *social2.RelationActionReq) bool {
					return req.UserId != ""
				})).Return(nil, assert.AnError)
			},
			wantCode: consts.SocialReqValueError,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockSocialClient)
			Rpc.SetSocialClient(mockClient)
			tt.mockSetup(mockClient)

			body, header := jsonBody(map[string]interface{}{})
			c := ut.CreateUtRequestContext("POST", "/relation/action", body, header...)
			c.Set("userid", "123")

			RelationAction(context.Background(), c)

			assert.Equal(t, 200, c.Response.StatusCode())
			assertResponseCode(t, c, tt.wantCode)
			mockClient.AssertExpectations(t)
		})
	}
}

func TestFollowingList(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(*MockSocialClient)
		wantCode  int32
		wantErr   bool
	}{
		{
			name: "Success_list",
			mockSetup: func(m *MockSocialClient) {
				m.On("FollowingList", mock.Anything, mock.MatchedBy(func(req *social2.FollowingListReq) bool {
					return req.UserId != ""
				})).Return(&social2.FollowingListResp{
					Code: consts.Success,
					Data: &social2.SocialData{
						Items: []*social2.UserInfo{},
						Total: 0,
					},
				}, nil)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name: "Fail_list_error",
			mockSetup: func(m *MockSocialClient) {
				m.On("FollowingList", mock.Anything, mock.MatchedBy(func(req *social2.FollowingListReq) bool {
					return req.UserId != ""
				})).Return(&social2.FollowingListResp{Code: consts.SocialDBSelectError}, nil)
			},
			wantCode: consts.SocialDBSelectError,
			wantErr:  true,
		},
		{
			name: "Fail_list_rpc_error",
			mockSetup: func(m *MockSocialClient) {
				m.On("FollowingList", mock.Anything, mock.MatchedBy(func(req *social2.FollowingListReq) bool {
					return req.UserId != ""
				})).Return(nil, assert.AnError)
			},
			wantCode: consts.SocialDBSelectError,
			wantErr:  true,
		},
		{
			name: "Fail_list_nil_data",
			mockSetup: func(m *MockSocialClient) {
				m.On("FollowingList", mock.Anything, mock.MatchedBy(func(req *social2.FollowingListReq) bool {
					return req.UserId != ""
				})).Return(&social2.FollowingListResp{Code: consts.Success, Data: nil}, nil)
			},
			wantCode: consts.SocialDBSelectError,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockSocialClient)
			Rpc.SetSocialClient(mockClient)
			tt.mockSetup(mockClient)

			c := ut.CreateUtRequestContext("GET", "/following/list?user_id=123", nil)
			c.Set("userid", "123")

			FollowingList(context.Background(), c)

			assert.Equal(t, 200, c.Response.StatusCode())
			assertResponseCode(t, c, tt.wantCode)
			mockClient.AssertExpectations(t)
		})
	}
}

func TestFollowerList(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(*MockSocialClient)
		wantCode  int32
		wantErr   bool
	}{
		{
			name: "Success_list",
			mockSetup: func(m *MockSocialClient) {
				m.On("FollowerList", mock.Anything, mock.MatchedBy(func(req *social2.FollowerListReq) bool {
					return req.UserId != ""
				})).Return(&social2.FollowerListResp{
					Code: consts.Success,
					Data: &social2.SocialData{
						Items: []*social2.UserInfo{},
						Total: 0,
					},
				}, nil)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name: "Fail_list_error",
			mockSetup: func(m *MockSocialClient) {
				m.On("FollowerList", mock.Anything, mock.MatchedBy(func(req *social2.FollowerListReq) bool {
					return req.UserId != ""
				})).Return(&social2.FollowerListResp{Code: consts.SocialDBSelectError}, nil)
			},
			wantCode: consts.SocialDBSelectError,
			wantErr:  true,
		},
		{
			name: "Fail_list_rpc_error",
			mockSetup: func(m *MockSocialClient) {
				m.On("FollowerList", mock.Anything, mock.MatchedBy(func(req *social2.FollowerListReq) bool {
					return req.UserId != ""
				})).Return(nil, assert.AnError)
			},
			wantCode: consts.SocialDBSelectError,
			wantErr:  true,
		},
		{
			name: "Fail_list_nil_data",
			mockSetup: func(m *MockSocialClient) {
				m.On("FollowerList", mock.Anything, mock.MatchedBy(func(req *social2.FollowerListReq) bool {
					return req.UserId != ""
				})).Return(&social2.FollowerListResp{Code: consts.Success, Data: nil}, nil)
			},
			wantCode: consts.SocialDBSelectError,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockSocialClient)
			Rpc.SetSocialClient(mockClient)
			tt.mockSetup(mockClient)

			c := ut.CreateUtRequestContext("GET", "/follower/list?user_id=123", nil)
			c.Set("userid", "123")

			FollowerList(context.Background(), c)

			assert.Equal(t, 200, c.Response.StatusCode())
			assertResponseCode(t, c, tt.wantCode)
			mockClient.AssertExpectations(t)
		})
	}
}

func TestFriendList(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(*MockSocialClient)
		wantCode  int32
		wantErr   bool
	}{
		{
			name: "Success_list",
			mockSetup: func(m *MockSocialClient) {
				m.On("FriendList", mock.Anything, mock.MatchedBy(func(req *social2.FriendListReq) bool {
					return req.UserId != ""
				})).Return(&social2.FriendListResp{
					Code: consts.Success,
					Data: &social2.SocialData{
						Items: []*social2.UserInfo{},
						Total: 0,
					},
				}, nil)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name: "Fail_list_error",
			mockSetup: func(m *MockSocialClient) {
				m.On("FriendList", mock.Anything, mock.MatchedBy(func(req *social2.FriendListReq) bool {
					return req.UserId != ""
				})).Return(&social2.FriendListResp{Code: consts.SocialDBSelectError}, nil)
			},
			wantCode: consts.SocialDBSelectError,
			wantErr:  true,
		},
		{
			name: "Fail_list_rpc_error",
			mockSetup: func(m *MockSocialClient) {
				m.On("FriendList", mock.Anything, mock.MatchedBy(func(req *social2.FriendListReq) bool {
					return req.UserId != ""
				})).Return(nil, assert.AnError)
			},
			wantCode: consts.SocialDBSelectError,
			wantErr:  true,
		},
		{
			name: "Fail_list_nil_data",
			mockSetup: func(m *MockSocialClient) {
				m.On("FriendList", mock.Anything, mock.MatchedBy(func(req *social2.FriendListReq) bool {
					return req.UserId != ""
				})).Return(&social2.FriendListResp{Code: consts.Success, Data: nil}, nil)
			},
			wantCode: consts.SocialDBSelectError,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockSocialClient)
			Rpc.SetSocialClient(mockClient)
			tt.mockSetup(mockClient)

			c := ut.CreateUtRequestContext("GET", "/friend/list?user_id=123", nil)
			c.Set("userid", "123")

			FriendList(context.Background(), c)

			assert.Equal(t, 200, c.Response.StatusCode())
			assertResponseCode(t, c, tt.wantCode)
			mockClient.AssertExpectations(t)
		})
	}
}
