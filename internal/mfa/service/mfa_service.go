package service

import (
	"Tiktok/pkg/consts"
	"Tiktok/pkg/utils"
	"context"

	"github.com/pkg/errors"
	"github.com/pquerna/otp/totp"
)

type MfaDatabase interface {
	SaveMfaSecret(ctx context.Context, mfa string, userId string) error
	GetMfaSecret(ctx context.Context, userId string) (string, error)
	MfaBindUpdate(ctx context.Context, userId string) error
	CheckMfaBind(ctx context.Context, userId string) (int, error)
}

type MfaRepo struct {
	mfaDb MfaDatabase
}

func NewMfaRepo(mfaDb MfaDatabase) *MfaRepo {
	return &MfaRepo{mfaDb: mfaDb}
}

func (s *MfaRepo) GenerateMfa(ctx context.Context, username string, userId string) (string, string, int32, error) {
	defer utils.TrackTime(ctx, "GenerateMfa")()
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Tk",
		AccountName: username,
	})
	if err != nil {
		return "", "", consts.MfaGenerateError, errors.Wrap(err, "->generate mfa totp.GenerateMfa error")
	}
	secret := key.Secret()
	err = s.mfaDb.SaveMfaSecret(ctx, secret, userId)
	if err != nil {
		return "", "", consts.UserDBUpdateError, errors.Wrap(err, "->generate mfa save MFA error")
	}
	return key.URL(), secret, consts.Success, nil
}

func (s *MfaRepo) MfaBindByCode(ctx context.Context, code string, userId string) (int32, error) {
	defer utils.TrackTime(ctx, "MfaBindByCode")()
	secret, err := s.mfaDb.GetMfaSecret(ctx, userId)
	if err != nil {
		return consts.MfaDBSelectError, errors.Wrap(err, "->mfa bind by code get mfa secret error")
	}
	valid := totp.Validate(code, secret)
	if !valid {
		return consts.MfaCodeFalse, nil
	}
	err = s.mfaDb.MfaBindUpdate(ctx, userId)
	if err != nil {
		return consts.UserDBUpdateError, errors.Wrap(err, "->mfa bind by code update MFA error")
	}
	return consts.Success, nil
}

func (s *MfaRepo) MfaBindBySecret(ctx context.Context, secret string, userId string) (int32, error) {
	defer utils.TrackTime(ctx, "MfaBindBySecret")()
	dbSecret, err := s.mfaDb.GetMfaSecret(ctx, userId)
	if err != nil {
		return consts.MfaDBSelectError, errors.Wrap(err, "->mfa bind by secret get mfa secret error")
	}
	if dbSecret != secret {
		return consts.MfaCodeFalse, nil
	}
	err = s.mfaDb.MfaBindUpdate(ctx, userId)
	if err != nil {
		return consts.UserDBUpdateError, errors.Wrap(err, "->mfa bind by secret update MFA error")
	}
	return consts.Success, nil
}

func (s *MfaRepo) MfaConfirm(ctx context.Context, mfaCode string, userID string) (int32, error) {
	defer utils.TrackTime(ctx, "MfaConfirm")()
	isBind, err := s.mfaDb.CheckMfaBind(ctx, userID)
	if err != nil {
		return consts.MfaDBSelectError, errors.Wrap(err, "->check mfa bind error")
	}
	if isBind != 0 {
		if mfaCode == "" {
			return consts.MfaReqValidError, nil
		}
		mfaSecret, err := s.mfaDb.GetMfaSecret(ctx, userID)
		if err != nil {
			return consts.MfaDBSelectError, errors.Wrap(err, "->mfa confirm mfa secret error")
		}
		if !totp.Validate(mfaCode, mfaSecret) {
			return consts.MfaCodeFalse, nil
		}
	}
	return consts.Success, nil
}
