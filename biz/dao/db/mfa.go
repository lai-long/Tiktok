package db

func SaveMfaSecret(mfa string, userId string) error {
	sql := `UPDATE users SET mfa_secret = ? WHERE id = ?`
	_, err := db.Exec(sql, mfa, userId)
	return err
}
func GetMfaSecret(userId string) (string, error) {
	sql := `SELECT mfa_secret FROM users WHERE id = ?`
	var secret string
	err := db.Get(&secret, sql, userId)
	return secret, err
}
func MfaBindUpdate(userId string) error {
	sql := `UPDATE users SET mfa_enabled = true WHERE id = ?`
	_, err := db.Exec(sql, userId)
	return err
}
func CheckMfaBind(userId string) (error, int) {
	sql := `SELECT mfa_enabled FROM users WHERE id = ?`
	var ok int
	err := db.Get(&ok, sql, userId)
	return err, ok
}
