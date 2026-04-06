package service

import (
	"Tiktok/pkg/consts"

	"github.com/pkg/errors"
	"github.com/pquerna/otp/totp"
)

type MfaDatabase interface {
	SaveMfaSecret(mfa string, userId string) error
	GetMfaSecret(userId string) (string, error)
	MfaBindUpdate(userId string) error
	CheckMfaBind(userId string) (error, int)
}

func (s *UserService) GenerateMfa(username string, userId string) (string, string, int32, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Tk",
		AccountName: username,
	})
	if err != nil {
		return "", "", consts.MfaGenerateError, errors.Wrap(err, "->generate mfa totp.GenerateMfa error")
	}
	secret := key.Secret()
	err = s.mfaDb.SaveMfaSecret(secret, userId)
	if err != nil {
		return "", "", consts.UserDBUpdateError, errors.Wrap(err, "->generate mfa save MFA error")
	}
	return key.URL(), secret, consts.Success, nil
}
func (s *UserService) MfaBindByCode(code string, userId string) (int32, error) {
	secret, err := s.mfaDb.GetMfaSecret(userId)
	if err != nil {
		return consts.UserDBSelectError, errors.Wrap(err, "->mfa bind by code get mfa secret error")
	}
	valid := totp.Validate(code, secret)
	if !valid {
		return consts.MfaCodeFalse, nil
	}
	err = s.mfaDb.MfaBindUpdate(userId)
	if err != nil {
		return consts.UserDBUpdateError, errors.Wrap(err, "->mfa bind by code update MFA error")
	}
	return consts.Success, nil
}
func (s *UserService) MfaBindBySecret(secret string, userId string) (int32, error) {
	dbSecret, err := s.mfaDb.GetMfaSecret(userId)
	if err != nil {
		return consts.UserDBSelectError, errors.Wrap(err, "->mfa bind by secret get mfa secret error")
	}
	if dbSecret != secret {
		return consts.MfaCodeFalse, nil
	}
	err = s.mfaDb.MfaBindUpdate(userId)
	if err != nil {
		return consts.UserDBUpdateError, errors.Wrap(err, "->mfa bind by secret update MFA error")
	}
	return consts.Success, nil
}
