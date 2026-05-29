package service

import (
	"Tiktok/pkg/consts"
	"Tiktok/pkg/entity"
	"context"
	"errors"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockVideoRedis struct {
	mock.Mock
}

func (m *MockVideoRedis) VideoHotSet(ctx context.Context, key string, member interface{}, score float64) error {
	args := m.Called(ctx, key, member, score)
	return args.Error(0)
}

func (m *MockVideoRedis) VideoHotGet(ctx context.Context, key string, pageNum int64, pageSize int64) ([]redis.Z, error) {
	args := m.Called(ctx, key, pageNum, pageSize)
	return args.Get(0).([]redis.Z), args.Error(1)
}

func (m *MockVideoRedis) VideoInfoSet(ctx context.Context, videoID string, video *entity.VideoEntity) error {
	args := m.Called(ctx, videoID, video)
	return args.Error(0)
}

func (m *MockVideoRedis) VideoInfoGet(ctx context.Context, videoID string) (*entity.VideoEntity, error) {
	args := m.Called(ctx, videoID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.VideoEntity), args.Error(1)
}

type MockVideoDb struct {
	mock.Mock
}

func (m *MockVideoDb) CreateVideo(ctx context.Context, entity entity.VideoEntity) error {
	args := m.Called(ctx, entity)
	return args.Error(0)
}

func (m *MockVideoDb) GetVideoByUserID(ctx context.Context, userId string, pageSize int64, pageNum int64) ([]entity.VideoEntity, error) {
	args := m.Called(ctx, userId, pageSize, pageNum)
	return args.Get(0).([]entity.VideoEntity), args.Error(1)
}

func (m *MockVideoDb) GetVideoByKeyWord(ctx context.Context, keyword string, pageNum int64, pageSize int64) ([]entity.VideoEntity, error) {
	args := m.Called(ctx, keyword, pageNum, pageSize)
	return args.Get(0).([]entity.VideoEntity), args.Error(1)
}

func (m *MockVideoDb) GetVideoByVideoId(ctx context.Context, videoId string) (entity.VideoEntity, error) {
	args := m.Called(ctx, videoId)
	return args.Get(0).(entity.VideoEntity), args.Error(1)
}

func (m *MockVideoDb) GetVideoStream(ctx context.Context) ([]entity.VideoEntity, error) {
	args := m.Called(ctx)
	return args.Get(0).([]entity.VideoEntity), args.Error(1)
}

func (m *MockVideoDb) GetVideoByIds(ctx context.Context, ids []string) ([]entity.VideoEntity, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]entity.VideoEntity), args.Error(1)
}

func TestVideoPublish(t *testing.T) {
	tests := []struct {
		name        string
		title       string
		description string
		url         string
		coverURL    string
		userID      string
		mockSetup   func(*MockVideoRedis, *MockVideoDb)
		wantCode    int32
		wantErr     bool
	}{
		{
			name:        "Success_publish",
			title:       "title",
			description: "description",
			url:         "http://example.com/video.mp4",
			coverURL:    "qiniu://cover/xxx.jpg",
			userID:      "userID",
			mockSetup: func(mr *MockVideoRedis, md *MockVideoDb) {
				mr.On("VideoHotSet", mock.Anything, "videoHot", mock.Anything, mock.Anything).Return(nil)
				md.On("CreateVideo", mock.Anything, mock.MatchedBy(func(e entity.VideoEntity) bool {
					return e.UserID != ""
				})).Return(nil)
				mr.On("VideoInfoSet", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name:        "Fail_redis_set_error",
			title:       "title",
			description: "description",
			url:         "http://example.com/video.mp4",
			coverURL:    "",
			userID:      "userID",
			mockSetup: func(mr *MockVideoRedis, md *MockVideoDb) {
				mr.On("VideoHotSet", mock.Anything, "videoHot", mock.Anything, mock.Anything).Return(errors.New("redis error"))
			},
			wantCode: consts.VideoRedisSetError,
			wantErr:  true,
		},
		{
			name:        "Fail_db_insert_error",
			title:       "title",
			description: "description",
			url:         "http://example.com/video.mp4",
			coverURL:    "",
			userID:      "userID",
			mockSetup: func(mr *MockVideoRedis, md *MockVideoDb) {
				mr.On("VideoHotSet", mock.Anything, "videoHot", mock.Anything, mock.Anything).Return(nil)
				md.On("CreateVideo", mock.Anything, mock.MatchedBy(func(e entity.VideoEntity) bool {
					return e.UserID != ""
				})).Return(errors.New("db error"))
			},
			wantCode: consts.VideoDBInsertError,
			wantErr:  true,
		},
		{
			name:        "Success_empty_url",
			title:       "title",
			description: "description",
			url:         "",
			coverURL:    "",
			userID:      "userID",
			mockSetup: func(mr *MockVideoRedis, md *MockVideoDb) {
				mr.On("VideoHotSet", mock.Anything, "videoHot", mock.Anything, mock.Anything).Return(nil)
				md.On("CreateVideo", mock.Anything, mock.MatchedBy(func(e entity.VideoEntity) bool {
					return e.UserID != ""
				})).Return(nil)
				mr.On("VideoInfoSet", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name:        "Success_empty_title",
			title:       "",
			description: "description",
			url:         "http://example.com/video.mp4",
			coverURL:    "",
			userID:      "userID",
			mockSetup: func(mr *MockVideoRedis, md *MockVideoDb) {
				mr.On("VideoHotSet", mock.Anything, "videoHot", mock.Anything, mock.Anything).Return(nil)
				md.On("CreateVideo", mock.Anything, mock.MatchedBy(func(e entity.VideoEntity) bool {
					return e.UserID != ""
				})).Return(nil)
				mr.On("VideoInfoSet", mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRedis := new(MockVideoRedis)
			mockDb := new(MockVideoDb)
			tt.mockSetup(mockRedis, mockDb)
			videoRepo := NewVideoRepo(mockDb, mockRedis)
			code, err := videoRepo.VideoPublish(context.Background(), tt.title, tt.description, tt.url, tt.coverURL, tt.userID)
			assert.Equal(t, tt.wantCode, code)
			assert.Equal(t, tt.wantErr, err != nil)
			mockRedis.AssertExpectations(t)
			mockDb.AssertExpectations(t)
		})
	}
}

func TestVideoList(t *testing.T) {
	tests := []struct {
		name      string
		userId    string
		pageSize  int64
		pageNum   int64
		mockSetup func(*MockVideoDb)
		wantCode  int32
		wantErr   bool
	}{
		{
			name:     "Success_list",
			userId:   "userID",
			pageSize: 10,
			pageNum:  1,
			mockSetup: func(m *MockVideoDb) {
				m.On("GetVideoByUserID", mock.Anything, "userID", int64(10), int64(1)).Return([]entity.VideoEntity{
					{ID: "1", Title: "title1"},
					{ID: "2", Title: "title2"},
				}, nil)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name:     "Fail_db_error",
			userId:   "userID",
			pageSize: 10,
			pageNum:  1,
			mockSetup: func(m *MockVideoDb) {
				m.On("GetVideoByUserID", mock.Anything, "userID", int64(10), int64(1)).Return([]entity.VideoEntity{}, errors.New("db error"))
			},
			wantCode: consts.VideoDBSelectError,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDb := new(MockVideoDb)
			tt.mockSetup(mockDb)
			videoRepo := NewVideoRepo(mockDb, nil)
			code, _, err := videoRepo.VideoList(context.Background(), tt.userId, tt.pageSize, tt.pageNum)
			assert.Equal(t, tt.wantCode, code)
			assert.Equal(t, tt.wantErr, err != nil)
			mockDb.AssertExpectations(t)
		})
	}
}

func TestVideoSearch(t *testing.T) {
	tests := []struct {
		name      string
		keyword   string
		pageNum   int64
		pageSize  int64
		mockSetup func(*MockVideoDb)
		wantCode  int32
		wantErr   bool
	}{
		{
			name:     "Success_search",
			keyword:  "test",
			pageNum:  1,
			pageSize: 10,
			mockSetup: func(m *MockVideoDb) {
				m.On("GetVideoByKeyWord", mock.Anything, "test", int64(1), int64(10)).Return([]entity.VideoEntity{
					{ID: "1", Title: "test video"},
				}, nil)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name:     "Fail_db_error",
			keyword:  "test",
			pageNum:  1,
			pageSize: 10,
			mockSetup: func(m *MockVideoDb) {
				m.On("GetVideoByKeyWord", mock.Anything, "test", int64(1), int64(10)).Return([]entity.VideoEntity{}, errors.New("db error"))
			},
			wantCode: consts.VideoDBSelectError,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDb := new(MockVideoDb)
			tt.mockSetup(mockDb)
			videoRepo := NewVideoRepo(mockDb, nil)
			code, _, err := videoRepo.VideoSearch(context.Background(), tt.keyword, tt.pageNum, tt.pageSize)
			assert.Equal(t, tt.wantCode, code)
			assert.Equal(t, tt.wantErr, err != nil)
			mockDb.AssertExpectations(t)
		})
	}
}

func TestVideoPopular(t *testing.T) {
	tests := []struct {
		name      string
		pageNum   int64
		pageSize  int64
		mockSetup func(*MockVideoRedis, *MockVideoDb)
		wantCode  int32
		wantErr   bool
	}{
		{
			name:     "Success_popular",
			pageNum:  1,
			pageSize: 10,
			mockSetup: func(mr *MockVideoRedis, md *MockVideoDb) {
				mr.On("VideoHotGet", mock.Anything, "videoHot", int64(1), int64(10)).Return([]redis.Z{
					{Score: 100, Member: "video1"},
					{Score: 90, Member: "video2"},
				}, nil)
				mr.On("VideoInfoGet", mock.Anything, "video1").Return(&entity.VideoEntity{ID: "video1", Title: "title1"}, nil)
				mr.On("VideoInfoGet", mock.Anything, "video2").Return(&entity.VideoEntity{ID: "video2", Title: "title2"}, nil)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name:     "Fail_redis_get_error",
			pageNum:  1,
			pageSize: 10,
			mockSetup: func(mr *MockVideoRedis, md *MockVideoDb) {
				mr.On("VideoHotGet", mock.Anything, "videoHot", int64(1), int64(10)).Return([]redis.Z{}, errors.New("redis error"))
			},
			wantCode: consts.VideoRedisGetError,
			wantErr:  true,
		},
		{
			name:     "Success_cache_miss_db_fallback",
			pageNum:  1,
			pageSize: 10,
			mockSetup: func(mr *MockVideoRedis, md *MockVideoDb) {
				mr.On("VideoHotGet", mock.Anything, "videoHot", int64(1), int64(10)).Return([]redis.Z{
					{Score: 100, Member: "video1"},
				}, nil)
				mr.On("VideoInfoGet", mock.Anything, "video1").Return(nil, errors.New("cache miss"))
				md.On("GetVideoByVideoId", mock.Anything, "video1").Return(entity.VideoEntity{ID: "video1", Title: "title1"}, nil)
				mr.On("VideoInfoSet", mock.Anything, "video1", mock.Anything).Return(nil)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name:     "Fail_db_fallback_error",
			pageNum:  1,
			pageSize: 10,
			mockSetup: func(mr *MockVideoRedis, md *MockVideoDb) {
				mr.On("VideoHotGet", mock.Anything, "videoHot", int64(1), int64(10)).Return([]redis.Z{
					{Score: 100, Member: "video1"},
				}, nil)
				mr.On("VideoInfoGet", mock.Anything, "video1").Return(nil, errors.New("cache miss"))
				md.On("GetVideoByVideoId", mock.Anything, "video1").Return(entity.VideoEntity{}, errors.New("db error"))
			},
			wantCode: consts.VideoDBSelectError,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRedis := new(MockVideoRedis)
			mockDb := new(MockVideoDb)
			tt.mockSetup(mockRedis, mockDb)
			videoRepo := NewVideoRepo(mockDb, mockRedis)
			code, _, err := videoRepo.VideoPopular(context.Background(), tt.pageNum, tt.pageSize)
			assert.Equal(t, tt.wantCode, code)
			assert.Equal(t, tt.wantErr, err != nil)
			mockRedis.AssertExpectations(t)
			mockDb.AssertExpectations(t)
		})
	}
}

func TestVideoStream(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(*MockVideoDb)
		wantCode  int32
		wantErr   bool
	}{
		{
			name: "Success_stream",
			mockSetup: func(m *MockVideoDb) {
				m.On("GetVideoStream", mock.Anything).Return([]entity.VideoEntity{
					{ID: "1", Title: "stream video 1"},
					{ID: "2", Title: "stream video 2"},
				}, nil)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name: "Fail_db_error",
			mockSetup: func(m *MockVideoDb) {
				m.On("GetVideoStream", mock.Anything).Return([]entity.VideoEntity{}, errors.New("db error"))
			},
			wantCode: consts.VideoDBSelectError,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDb := new(MockVideoDb)
			tt.mockSetup(mockDb)
			videoRepo := NewVideoRepo(mockDb, nil)
			code, _, err := videoRepo.VideoStream(context.Background())
			assert.Equal(t, tt.wantCode, code)
			assert.Equal(t, tt.wantErr, err != nil)
			mockDb.AssertExpectations(t)
		})
	}
}
