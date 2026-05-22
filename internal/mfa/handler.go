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
	switch req.Type {
	case "secret":
		code, _ = s.mfaRepo.MfaBindBySecret(req.Secret, req.UserID)
	case "qrcode":
		code, _ = s.mfaRepo.MfaBindByCode(req.MfaCode, req.UserID)
	default:
		code = consts.MfaCodeFalse
	}
	resp = &mfa.MfaBindResp{
		Code: code,
	}
	return resp, nil
}

// MfaConfirm implements the MfaServiceImpl interface.
func (s *MfaServiceImpl) MfaConfirm(ctx context.Context, req *mfa.MfaConfirmReq) (resp *mfa.MfaConfirmResp, err error) {
	// TODO: Your code here...
	resp = &mfa.MfaConfirmResp{}
	code, err := s.mfaRepo.MfaConfirm(req.QrCode, req.UserID)
	if err != nil {
		resp.Code = code
		return resp, err
	}
	resp.Code = consts.Success
	return resp, nil
}
