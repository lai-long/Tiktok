package social

import (
	social "Tiktok/kitex_gen/social"
	"Tiktok/pkg/consts"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockSocialService struct {
	mock.Mock
}

func (m *mockSocialService) RelationAction(ctx context.Context, toUserId string, actionType string, userId string) (int32, error) {
	args := m.Called(ctx, toUserId, actionType, userId)
	return args.Get(0).(int32), args.Error(1)
}

func (m *mockSocialService) FollowingList(ctx context.Context, userId string, pageNum int64, pageSize int64) (int32, []*social.UserInfo, error) {
	args := m.Called(ctx, userId, pageNum, pageSize)
	return args.Get(0).(int32), args.Get(1).([]*social.UserInfo), args.Error(2)
}

func (m *mockSocialService) FollowerList(ctx context.Context, userId string, pageNum int64, pageSize int64) (int32, []*social.UserInfo, error) {
	args := m.Called(ctx, userId, pageNum, pageSize)
	return args.Get(0).(int32), args.Get(1).([]*social.UserInfo), args.Error(2)
}

func (m *mockSocialService) FriendList(ctx context.Context, userId string, pageNum int64, pageSize int64) (int32, []*social.UserInfo, error) {
	args := m.Called(ctx, userId, pageNum, pageSize)
	return args.Get(0).(int32), args.Get(1).([]*social.UserInfo), args.Error(2)
}

func TestSocialServiceImpl_RelationAction(t *testing.T) {
	tests := []struct {
		name     string
		mock     func(*mockSocialService)
		wantCode int32
		wantErr  bool
	}{
		{
			name: "success",
			mock: func(m *mockSocialService) {
				m.On("RelationAction", mock.Anything, "to", "1", "from").Return(consts.Success, nil)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name: "db_error",
			mock: func(m *mockSocialService) {
				m.On("RelationAction", mock.Anything, "to", "1", "from").Return(consts.SocialDBInsertError, errors.New("db down"))
			},
			wantCode: consts.SocialDBInsertError,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := new(mockSocialService)
			tt.mock(m)
			h := NewSocialServiceImpl(m)
			resp, err := h.RelationAction(context.Background(), &social.RelationActionReq{ToUserId: "to", ActionType: "1", UserId: "from"})
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

func TestSocialServiceImpl_FollowingList(t *testing.T) {
	tests := []struct {
		name     string
		mock     func(*mockSocialService)
		wantCode int32
		wantErr  bool
	}{
		{
			name: "success",
			mock: func(m *mockSocialService) {
				m.On("FollowingList", mock.Anything, "uid", int64(1), int64(10)).Return(consts.Success, []*social.UserInfo{{ID: "fid"}}, nil)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name: "db_error",
			mock: func(m *mockSocialService) {
				m.On("FollowingList", mock.Anything, "uid", int64(1), int64(10)).Return(consts.SocialDBSelectError, []*social.UserInfo{}, errors.New("db down"))
			},
			wantCode: consts.SocialDBSelectError,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := new(mockSocialService)
			tt.mock(m)
			h := NewSocialServiceImpl(m)
			resp, err := h.FollowingList(context.Background(), &social.FollowingListReq{UserId: "uid", PageNum: 1, PageSize: 10})
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

func TestSocialServiceImpl_FollowerList(t *testing.T) {
	tests := []struct {
		name     string
		mock     func(*mockSocialService)
		wantCode int32
		wantErr  bool
	}{
		{
			name: "success",
			mock: func(m *mockSocialService) {
				m.On("FollowerList", mock.Anything, "uid", int64(1), int64(10)).Return(consts.Success, []*social.UserInfo{{ID: "fid"}}, nil)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name: "db_error",
			mock: func(m *mockSocialService) {
				m.On("FollowerList", mock.Anything, "uid", int64(1), int64(10)).Return(consts.SocialDBSelectError, []*social.UserInfo{}, errors.New("db down"))
			},
			wantCode: consts.SocialDBSelectError,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := new(mockSocialService)
			tt.mock(m)
			h := NewSocialServiceImpl(m)
			resp, err := h.FollowerList(context.Background(), &social.FollowerListReq{UserId: "uid", PageNum: 1, PageSize: 10})
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

func TestSocialServiceImpl_FriendList(t *testing.T) {
	tests := []struct {
		name     string
		mock     func(*mockSocialService)
		wantCode int32
		wantErr  bool
	}{
		{
			name: "success",
			mock: func(m *mockSocialService) {
				m.On("FriendList", mock.Anything, "uid", int64(1), int64(10)).Return(consts.Success, []*social.UserInfo{{ID: "fid"}}, nil)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name: "db_error",
			mock: func(m *mockSocialService) {
				m.On("FriendList", mock.Anything, "uid", int64(1), int64(10)).Return(consts.SocialDBSelectError, []*social.UserInfo{}, errors.New("db down"))
			},
			wantCode: consts.SocialDBSelectError,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := new(mockSocialService)
			tt.mock(m)
			h := NewSocialServiceImpl(m)
			resp, err := h.FriendList(context.Background(), &social.FriendListReq{UserId: "uid", PageNum: 1, PageSize: 10})
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
