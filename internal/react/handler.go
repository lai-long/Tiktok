package react

import (
	react "Tiktok/kitex_gen/react"
	"context"
)

// LikeServiceImpl implements the last service interface defined in the IDL.
type LikeServiceImpl struct{}

// LikeAction implements the LikeServiceImpl interface.
func (s *LikeServiceImpl) LikeAction(ctx context.Context, req *react.LikeActionReq) (resp *react.LikeActionResp, err error) {
	// TODO: Your code here...
	return
}

// LikeList implements the LikeServiceImpl interface.
func (s *LikeServiceImpl) LikeList(ctx context.Context, req *react.LikeListReq) (resp *react.LikeListResp, err error) {
	// TODO: Your code here...
	return
}
