package mfa

import (
	"Tiktok/internal/mfa/service"
	mfa "Tiktok/kitex_gen/mfa"
	"context"
)

type MfaServiceImpl struct {
	mfaRepo *service.MfaRepo
}

func NewMfaService(mfaRepo *service.MfaRepo) *MfaServiceImpl {
	return &MfaServiceImpl{
		mfaRepo: mfaRepo,
	}
}

func (s *MfaServiceImpl) MfaQrcode(ctx context.Context, req *mfa.MfaQrcodeReq) (resp *mfa.MfaQrcodeResp, err error) {
	resp = &mfa.MfaQrcodeResp{}
	qrCode, secret, code, _ := s.mfaRepo.GenerateMfa(req.UserName, req.UserID)
	resp = &mfa.MfaQrcodeResp{
		Code: code,
		Data: &mfa.MfaData{
			Secret: secret,
			Qrcode: qrCode,
		},
	}
	return resp, nil
}

func (s *MfaServiceImpl) MfaBind(ctx context.Context, req *mfa.MfaBindReq) (resp *mfa.MfaBindResp, err error) {
	resp = &mfa.MfaBindResp{}
	var code int32
	switch req.Type {
	case "secret":
		code, _ = s.mfaRepo.MfaBindBySecret(req.Secret, req.UserID)
	case "qrcode":
		code, _ = s.mfaRepo.MfaBindByCode(req.MfaCode, req.UserID)
	default:
		code = 1
	}
	resp = &mfa.MfaBindResp{
		Code: code,
	}
	return resp, nil
}

func (s *MfaServiceImpl) MfaConfirm(ctx context.Context, req *mfa.MfaConfirmReq) (resp *mfa.MfaConfirmResp, err error) {
	resp = &mfa.MfaConfirmResp{}
	code, err := s.mfaRepo.MfaConfirm(req.QrCode, req.UserID)
	if err != nil {
		resp.Code = code
		return resp, err
	}
	resp.Code = code
	return resp, nil
}
