package mfa

import (
	"Tiktok/internal/mfa/service"
	mfa "Tiktok/kitex_gen/mfa"
	"Tiktok/pkg/consts"
	"context"

	"github.com/alibaba/sentinel-golang/api"
)

// MfaServiceImpl implements the last service interface defined in the IDL.
type MfaServiceImpl struct {
	mfaRepo *service.MfaRepo
}

func NewMfaService(mfaRepo *service.MfaRepo) *MfaServiceImpl {
	return &MfaServiceImpl{
		mfaRepo: mfaRepo,
	}
}

// MfaQrcode implements the MfaServiceImpl interface.
func (s *MfaServiceImpl) MfaQrcode(ctx context.Context, req *mfa.MfaQrcodeReq) (resp *mfa.MfaQrcodeResp, err error) {
	// TODO: Your code here...
	resp = &mfa.MfaQrcodeResp{}
	entry, blockErr := api.Entry("/auth/mfa/qrcode")
	if blockErr != nil {
		resp.Code = consts.SentinelBlock
		return resp, nil
	}
	defer entry.Exit()
	qrCode, secret, code, err := s.mfaRepo.GenerateMfa(req.UserName, req.UserID)
	resp = &mfa.MfaQrcodeResp{
		Code: code,
		Data: &mfa.MfaData{
			Secret: secret,
			Qrcode: qrCode,
		},
	}
	return resp, err
}

// MfaBind implements the MfaServiceImpl interface.
func (s *MfaServiceImpl) MfaBind(ctx context.Context, req *mfa.MfaBindReq) (resp *mfa.MfaBindResp, err error) {
	// TODO: Your code here...
	resp = &mfa.MfaBindResp{}
	entry, blockErr := api.Entry("/auth/mfa/bind")
	if blockErr != nil {
		resp.Code = consts.SentinelBlock
		return resp, nil
	}
	defer entry.Exit()
	var code int32
	if req.Type == "secret" {
		code, err = s.mfaRepo.MfaBindBySecret(req.Secret, req.UserID)
	} else if req.Type == "qrcode" {
		code, err = s.mfaRepo.MfaBindByCode(req.MfaCode, req.UserID)
	} else {
		code = consts.MfaCodeFalse
		err = nil
	}
	resp = &mfa.MfaBindResp{
		Code: code,
	}
	return resp, err
}
