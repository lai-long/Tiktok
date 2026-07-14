package user

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
	user2 "Tiktok/kitex_gen/user"

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

func buildMultipartBody(files map[string][]byte) (*ut.Body, []ut.Header) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for name, content := range files {
		part, err := w.CreateFormFile(name, name+".png")
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

func TestMfaQrcode(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(*MockUserClient)
		wantCode  int32
		wantErr   bool
	}{
		{
			name: "Success_qrcode",
			mockSetup: func(m *MockUserClient) {
				m.On("MfaQrcode", mock.Anything, mock.Anything).Return(&user2.MfaQrcodeResp{
					Code: consts.Success,
					Data: &user2.MfaData{Secret: "secret", Qrcode: "qrcode"},
				}, nil)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name: "Fail_rpc_error",
			mockSetup: func(m *MockUserClient) {
				m.On("MfaQrcode", mock.Anything, mock.Anything).Return(nil, assert.AnError)
			},
			wantCode: consts.MfaCodeFalse,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockUserClient)
			Rpc.SetUserClient(mockClient)
			tt.mockSetup(mockClient)

			c := ut.CreateUtRequestContext("GET", "/auth/mfa/qrcode", nil)
			c.Set("userid", "userID")
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
		body      map[string]any
		mockSetup func(*MockUserClient)
		wantCode  int32
		wantErr   bool
	}{
		{
			name: "Success_bind_by_secret",
			body: map[string]any{"secret": "testsecret"},
			mockSetup: func(m *MockUserClient) {
				m.On("MfaBind", mock.Anything, mock.MatchedBy(func(req *user2.MfaBindReq) bool {
					return req.Secret == "testsecret" && req.Type == "secret"
				})).Return(&user2.MfaBindResp{Code: consts.Success}, nil)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name: "Success_bind_by_code",
			body: map[string]any{"code": "123456"},
			mockSetup: func(m *MockUserClient) {
				m.On("MfaBind", mock.Anything, mock.MatchedBy(func(req *user2.MfaBindReq) bool {
					return req.MfaCode == "123456" && req.Type == "qrcode"
				})).Return(&user2.MfaBindResp{Code: consts.Success}, nil)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name:      "Fail_missing_secret_and_code",
			body:      map[string]any{},
			mockSetup: func(m *MockUserClient) {},
			wantCode:  consts.UserReqValidError,
			wantErr:   false,
		},
		{
			name: "Fail_rpc_error",
			body: map[string]any{"secret": "testsecret"},
			mockSetup: func(m *MockUserClient) {
				m.On("MfaBind", mock.Anything, mock.Anything).Return(nil, assert.AnError)
			},
			wantCode: consts.MfaCodeFalse,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockUserClient)
			Rpc.SetUserClient(mockClient)
			tt.mockSetup(mockClient)

			body, header := jsonBody(tt.body)
			c := ut.CreateUtRequestContext("POST", "/auth/mfa/bind", body, header...)
			c.Set("userid", "userID")
			c.Set("username", "testuser")

			MfaBind(context.Background(), c)

			assert.Equal(t, 200, c.Response.StatusCode())
			assertResponseCode(t, c, tt.wantCode)
			mockClient.AssertExpectations(t)
		})
	}
}

func TestUserAvatar(t *testing.T) {
	tests := []struct {
		name        string
		files       map[string][]byte
		uploader    *stubUploader
		mockSetup   func(*MockUserClient)
		wantCode    int32
		wantErr     bool
		wantDataNil bool
	}{
		{
			name:     "Success_avatar_upload",
			files:    map[string][]byte{"data": []byte("fake image content")},
			uploader: &stubUploader{imageKey: "qiniu://avatar/key"},
			mockSetup: func(m *MockUserClient) {
				m.On("UserAvatar", mock.Anything, mock.MatchedBy(func(req *user2.UserAvatarReq) bool {
					return req.AvatarURL == "qiniu://avatar/key" && req.UserID == "userID"
				})).Return(&user2.UserAvatarResp{
					Code: consts.Success,
					Data: &user2.UserInfo{ID: "userID", Username: "testuser"},
				}, nil)
			},
			wantCode: consts.Success,
			wantErr:  false,
		},
		{
			name:        "Fail_missing_file",
			files:       map[string][]byte{},
			uploader:    &stubUploader{},
			mockSetup:   func(m *MockUserClient) {},
			wantCode:    consts.UserReqValidError,
			wantErr:     false,
			wantDataNil: true,
		},
		{
			name:     "Fail_upload_error",
			files:    map[string][]byte{"data": []byte("fake image content")},
			uploader: &stubUploader{imageErr: errors.New("upload failed")},
			mockSetup: func(m *MockUserClient) {
			},
			wantCode:    consts.FileError,
			wantErr:     false,
			wantDataNil: true,
		},
		{
			name:     "Fail_invalid_image",
			files:    map[string][]byte{"data": []byte("fake image content")},
			uploader: &stubUploader{imageKey: ""},
			mockSetup: func(m *MockUserClient) {
			},
			wantCode:    consts.ImageFalse,
			wantErr:     false,
			wantDataNil: true,
		},
		{
			name:     "Fail_rpc_error",
			files:    map[string][]byte{"data": []byte("fake image content")},
			uploader: &stubUploader{imageKey: "qiniu://avatar/key"},
			mockSetup: func(m *MockUserClient) {
				m.On("UserAvatar", mock.Anything, mock.Anything).Return(nil, assert.AnError)
			},
			wantCode:    consts.UserReqValidError,
			wantErr:     true,
			wantDataNil: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setUploadProvider(t, tt.uploader)

			mockClient := new(MockUserClient)
			Rpc.SetUserClient(mockClient)
			tt.mockSetup(mockClient)

			body, header := buildMultipartBody(tt.files)
			c := ut.CreateUtRequestContext("PUT", "/user/avatar/upload", body, header...)
			c.Set("userid", "userID")

			UserAvatar(context.Background(), c)

			assert.Equal(t, 200, c.Response.StatusCode())
			assertResponseCode(t, c, tt.wantCode)
			mockClient.AssertExpectations(t)
		})
	}
}
