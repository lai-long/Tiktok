package utils

import "github.com/rs/xid"

func IdGenerate() string {
	return xid.New().String()
}
