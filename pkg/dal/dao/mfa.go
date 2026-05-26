package dao

import "context"

func (m *MySQLdb) SaveMfaSecret(ctx context.Context, mfa string, userId string) error {
	sql := `UPDATE users SET mfa_secret = ? WHERE id = ?`
	_, err := m.db.ExecContext(ctx, sql, mfa, userId)
	return err
}
func (m *MySQLdb) GetMfaSecret(ctx context.Context, userId string) (string, error) {
	sql := `SELECT mfa_secret FROM users WHERE id = ?`
	var secret string
	err := m.db.GetContext(ctx, &secret, sql, userId)
	return secret, err
}
func (m *MySQLdb) MfaBindUpdate(ctx context.Context, userId string) error {
	sql := `UPDATE users SET mfa_enabled = true WHERE id = ?`
	_, err := m.db.ExecContext(ctx, sql, userId)
	return err
}
func (m *MySQLdb) CheckMfaBind(ctx context.Context, userId string) (int, error) {
	sql := `SELECT mfa_enabled FROM users WHERE id = ?`
	var ok int
	err := m.db.GetContext(ctx, &ok, sql, userId)
	return ok, err
}
