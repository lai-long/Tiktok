package service

import (
	"Tiktok/pkg/consts"
	"log"

	"github.com/pquerna/otp/totp"
)

func (s *UserService) GenerateMfa(username string, userId string) (bool, string, string, int, string) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Tk",
		AccountName: username,
	})
	if err != nil {
		log.Println("Generate MFA err:", err)
		return false, "", "", consts.CodeMfaError, "Mfa Generate err"
	}
	secret := key.Secret()
	err = s.mfaDb.SaveMfaSecret(secret, userId)
	if err != nil {
		log.Println("Generate MFA err:", err)
		return false, "", "", consts.CodeDBUpdateError, "Mfa Generate err"
	}
	return true, key.URL(), secret, consts.CodeSuccess, "Mfa Generate success"
}
func (s *UserService) MfaBindByCode(code string, userId string) (int, string) {
	secret, err := s.mfaDb.GetMfaSecret(userId)
	if err != nil {
		log.Println("Get MFA secret err:", err)
		return consts.CodeDBSelectError, "GetMfaSecret err:"
	}
	valid := totp.Validate(code, secret)
	if !valid {
		return consts.CodeMfaError, "GetMfaSecret err:"
	}
	err = s.mfaDb.MfaBindUpdate(userId)
	if err != nil {
		log.Println("MfaBindUpdate err:", err)
		return consts.CodeDBUpdateError, "MfaBindUpdate err:"
	}
	return consts.CodeSuccess, "MfaBindUpdate success"
}
func (s *UserService) MfaBindBySecret(secret string, userId string) (int, string) {
	dbSecret, err := s.mfaDb.GetMfaSecret(userId)
	if err != nil {
		log.Println("Get MFA secret err:", err)
		return consts.CodeDBSelectError, "GetMfaSecret err:"
	}
	if dbSecret != secret {
		return consts.CodeMfaError, "MfaSecret false err:"
	}
	err = s.mfaDb.MfaBindUpdate(userId)
	if err != nil {
		log.Println("MfaBindUpdate err:", err)
		return consts.CodeDBUpdateError, "MfaBindUpdate err:"
	}
	return consts.CodeSuccess, "MfaBindUpdate success"
}
