package entity

import (
	"Tiktok/biz/model/user"
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

func (u *UserEntity) ToUserInfo() *user.UserInfo {
	return &user.UserInfo{
		ID:        u.Id,
		Username:  u.Username,
		AvatarURL: u.Avatar_url,
		CreatedAt: u.Created_at.String(),
		UpdatedAt: u.Updated_at.String(),
	}
}
