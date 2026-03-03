package entity

import (
	"database/sql"
	"time"
)

type UserEntity struct {
	Id         string       `db:"id"`
	Username   string       `db:"username"`
	Password   string       `db:"password"`
	Avatar_url string       `db:"avatar_url"`
	Created_at time.Time    `db:"created_at"`
	Updated_at time.Time    `db:"updated_at"`
	Deleted_at sql.NullTime `db:"deleted_at"`
	MfaSecret  string       `db:"mfa_secret"`
	MfaEnabled bool         `db:"mfa_enabled"`
}
