package service

import (
	"Tiktok/biz/dao/db"
	"Tiktok/pkg/consts"
	"log"

	"github.com/pquerna/otp/totp"
)

func GenerateMfa(username string, userId string) (bool, string, string, int, string) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Tk",
		AccountName: username,
	})
	if err != nil {
		log.Println("Generate MFA err:", err)
		return false, "", "", consts.CodeMfaError, "Mfa Generate err"
	}
	secret := key.Secret()
	err = db.SaveMfaSecret(secret, userId)
	if err != nil {
		log.Println("Generate MFA err:", err)
		return false, "", "", consts.CodeDBUpdateError, "Mfa Generate err"
	}
	return true, key.URL(), secret, consts.CodeSuccess, "Mfa Generate success"
}
func MfaBindByCode(code string, userId string) (int, string) {
	secret, err := db.GetMfaSecret(userId)
	if err != nil {
		log.Println("Get MFA secret err:", err)
		return consts.CodeDBSelectError, "GetMfaSecret err:"
	}
	valid := totp.Validate(code, secret)
	if !valid {
		return consts.CodeMfaError, "GetMfaSecret err:"
	}
	err = db.MfaBindUpdate(userId)
	if err != nil {
		log.Println("MfaBindUpdate err:", err)
		return consts.CodeDBUpdateError, "MfaBindUpdate err:"
	}
	return consts.CodeSuccess, "MfaBindUpdate success"
}
func MfaBindBySecret(secret string, userId string) (int, string) {
	dbSecret, err := db.GetMfaSecret(userId)
	if err != nil {
		log.Println("Get MFA secret err:", err)
		return consts.CodeDBSelectError, "GetMfaSecret err:"
	}
	if dbSecret != secret {
		return consts.CodeMfaError, "MfaSecret false err:"
	}
	err = db.MfaBindUpdate(userId)
	if err != nil {
		log.Println("MfaBindUpdate err:", err)
		return consts.CodeDBUpdateError, "MfaBindUpdate err:"
	}
	return consts.CodeSuccess, "MfaBindUpdate success"
}
