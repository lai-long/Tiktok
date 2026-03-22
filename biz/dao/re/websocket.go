package re

import (
	"context"
	"log"
	"time"
)

func GetMsgCountsBYClientID(ctx context.Context, id string) string {
	s, err := rdb.Get(ctx, id).Result()
	if err != nil {
		log.Println(err)
		return ""
	}
	return s
}
func MSgCountIncr(ctx context.Context, id string) {
	rdb.Incr(ctx, id)
	rdb.Expire(ctx, id, time.Hour*600)
}
