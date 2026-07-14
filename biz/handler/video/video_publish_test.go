package video

import (
	"bytes"
	"context"
	"errors"
	"io"
	"mime/multipart"
	"testing"

	"Tiktok/pkg/consts"
	"Tiktok/pkg/utils"

	Rpc "Tiktok/biz/rpc"
	video2 "Tiktok/kitex_gen/video"

	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type stubUploader struct {
	imageKey string
	imageErr error
	fileKey  string
	fileErr  error
}

func (s *stubUploader) UploadImage(context.Context, *multipart.FileHeader, string) (string, error) {
	return s.imageKey, s.imageErr
}

func (s *stubUploader) UploadFile(context.Context, io.Reader, string, string) (string, error) {
	return s.fileKey, s.fileErr
}

func setUploadProvider(t *testing.T, p utils.UploadProvider) {
	old := uploadProvider
	uploadProvider = p
	t.Cleanup(func() { uploadProvider = old })
}

func buildMultipartBody(fields map[string]string, files map[string][]byte) (*ut.Body, []ut.Header) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for k, v := range fields {
		if err := w.WriteField(k, v); err != nil {
			panic(err)
		}
	}
	for name, content := range files {
		ext := ".bin"
		if name == "cover" {
			ext = ".png"
		} else if name == "data" {
			ext = ".mp4"
		}
		part, err := w.CreateFormFile(name, name+ext)
		if err != nil {
			panic(err)
		}
		if _, err := part.Write(content); err != nil {
			panic(err)
		}
	}
	if err := w.Close(); err != nil {
		panic(err)
	}
	return &ut.Body{Body: bytes.NewReader(b.Bytes()), Len: b.Len()},
		[]ut.Header{{Key: "Content-Type", Value: w.FormDataContentType()}}
}

func TestVideoPublish(t *testing.T) {
	tests := []struct {
		name      string
		fields    map[string]string
		files     map[string][]byte
		uploader  *stubUploader
		mockSetup func(*MockVideoClient)
		userID    string
		wantCode  int32
		wantErr   bool
	}{
		{
			name:     "Success_video_only",
			fields:   map[string]string{"title": "test title"},
			files:    map[string][]byte{"data": []byte("fake video")},
			uploader: &stubUploader{fileKey: "qiniu://video/key"},
			mockSetup: func(m *MockVideoClient) {
				m.On("VideoPublish", mock.Anything, mock.MatchedBy(func(req *video2.VideoPublishReq) bool {
					return req.Title == "test title" && req.VideoURL == "qiniu://video/key" && req.UserID == "userID"
				})).Return(&video2.VideoPublishResp{Code: consts.Success}, nil)
			},
			userID:   "userID",
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name:   "Success_video_with_cover",
			fields: map[string]string{"title": "test title", "description": "desc"},
			files:  map[string][]byte{"data": []byte("fake video"), "cover": []byte("fake cover")},
			uploader: &stubUploader{
				fileKey:  "qiniu://video/key",
				imageKey: "qiniu://cover/key",
			},
			mockSetup: func(m *MockVideoClient) {
				m.On("VideoPublish", mock.Anything, mock.MatchedBy(func(req *video2.VideoPublishReq) bool {
					return req.Title == "test title" && req.CoverURL == "qiniu://cover/key"
				})).Return(&video2.VideoPublishResp{Code: consts.Success}, nil)
			},
			userID:   "userID",
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name:      "Fail_bind_error",
			fields:    map[string]string{},
			files:     map[string][]byte{},
			uploader:  &stubUploader{},
			mockSetup: func(m *MockVideoClient) {},
			userID:    "userID",
			wantCode:  consts.VideoReqValidError,
			wantErr:   false,
		},
		{
			name:      "Fail_missing_data",
			fields:    map[string]string{"title": "test title"},
			files:     map[string][]byte{},
			uploader:  &stubUploader{},
			mockSetup: func(m *MockVideoClient) {},
			userID:    "userID",
			wantCode:  consts.VideoReqValidError,
			wantErr:   false,
		},
		{
			name:      "Fail_user_id_empty",
			fields:    map[string]string{"title": "test title"},
			files:     map[string][]byte{"data": []byte("fake video")},
			uploader:  &stubUploader{},
			mockSetup: func(m *MockVideoClient) {},
			userID:    "",
			wantCode:  consts.VideoReqValidError,
			wantErr:   false,
		},
		{
			name:     "Fail_video_upload_error",
			fields:   map[string]string{"title": "test title"},
			files:    map[string][]byte{"data": []byte("fake video")},
			uploader: &stubUploader{fileErr: errors.New("upload failed")},
			mockSetup: func(m *MockVideoClient) {
			},
			userID:   "userID",
			wantCode: consts.FileError,
			wantErr:  false,
		},
		{
			name:   "Fail_cover_upload_error",
			fields: map[string]string{"title": "test title"},
			files:  map[string][]byte{"data": []byte("fake video"), "cover": []byte("fake cover")},
			uploader: &stubUploader{
				fileKey:  "qiniu://video/key",
				imageErr: errors.New("cover upload failed"),
			},
			mockSetup: func(m *MockVideoClient) {
			},
			userID:   "userID",
			wantCode: consts.FileError,
			wantErr:  false,
		},
		{
			name:     "Fail_rpc_error",
			fields:   map[string]string{"title": "test title"},
			files:    map[string][]byte{"data": []byte("fake video")},
			uploader: &stubUploader{fileKey: "qiniu://video/key"},
			mockSetup: func(m *MockVideoClient) {
				m.On("VideoPublish", mock.Anything, mock.Anything).Return(nil, assert.AnError)
			},
			userID:   "userID",
			wantCode: consts.VideoReqValidError,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setUploadProvider(t, tt.uploader)

			mockClient := new(MockVideoClient)
			Rpc.SetVideoClient(mockClient)
			tt.mockSetup(mockClient)

			body, header := buildMultipartBody(tt.fields, tt.files)
			c := ut.CreateUtRequestContext("POST", "/video/publish", body, header...)
			c.Set("userid", tt.userID)

			VideoPublish(context.Background(), c)

			assert.Equal(t, 200, c.Response.StatusCode())
			assertResponseCode(t, c, tt.wantCode)
			mockClient.AssertExpectations(t)
		})
	}
}
