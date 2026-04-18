package dao

func (m *MySQLdb) SaveMfaSecret(mfa string, userId string) error {
	sql := `UPDATE users SET mfa_secret = ? WHERE id = ?`
	_, err := m.db.Exec(sql, mfa, userId)
	return err
}
func (m *MySQLdb) GetMfaSecret(userId string) (string, error) {
	sql := `SELECT mfa_secret FROM users WHERE id = ?`
	var secret string
	err := m.db.Get(&secret, sql, userId)
	return secret, err
}
func (m *MySQLdb) MfaBindUpdate(userId string) error {
	sql := `UPDATE users SET mfa_enabled = true WHERE id = ?`
	_, err := m.db.Exec(sql, userId)
	return err
}
func (m *MySQLdb) CheckMfaBind(userId string) (int, error) {
	sql := `SELECT mfa_enabled FROM users WHERE id = ?`
	var ok int
	err := m.db.Get(&ok, sql, userId)
	return ok, err
}
