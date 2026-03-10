package router

import (
	"Tiktok/biz/handler"
	"Tiktok/biz/middleware"

	"github.com/cloudwego/hertz/pkg/app/server"
)

func SetRouters(handler *handler.Handler, userHandler *handler.UserHandler) {
	h := server.Default(
		server.WithHostPorts(":8888"),
		server.WithMaxRequestBodySize(10*1024*1024),
	)
	defer h.Close()
	//注册、登录、用户信息、上传头像
	user := h.Group("/user")
	{
		user.GET("/info", middleware.AuthMiddleware, userHandler.UserInfo)
		user.POST("/login", userHandler.UserLogin)
		user.POST("/register", userHandler.UserRegister)
		user.PUT("/avatar/upload", middleware.AuthMiddleware, userHandler.UserAvatar)
	}
	authMfa := h.Group("/auth/mfa")
	{
		authMfa.GET("/qrcode", middleware.AuthMiddleware, handler.MfaQrcode)
		authMfa.POST("/bind", middleware.AuthMiddleware, handler.MfaBind)
	}
	//投稿、发布列表、搜索视频、热门排行榜
	video := h.Group("/video")
	video.Use(middleware.AuthMiddleware)
	{
		video.POST("/publish", handler.VideoPublish)
		video.GET("/list", handler.VideoList)
		video.POST("/search", handler.VideoSearch)
		video.GET("/popular", handler.VideoPopular)
	}
	//点赞操作、点赞列表、评论、评论列表、删除评论
	like := h.Group("/like")
	like.Use(middleware.AuthMiddleware)
	{
		like.POST("/action", handler.LikeAction)
		like.GET("/list", handler.LikeList)
	}
	comment := h.Group("/comment")
	comment.Use(middleware.AuthMiddleware)
	{
		comment.POST("/publish", handler.CommentPublish)
		comment.GET("/list", handler.CommentList)
		comment.DELETE("/delete", handler.CommentDelete)
	}
	//关注操作、关注列表、粉丝列表、好友列表
	h.POST("/relation/action", middleware.AuthMiddleware, handler.RelationAction)
	h.GET("/following/list", middleware.AuthMiddleware, handler.FollowingList)
	h.GET("/follower/list", middleware.AuthMiddleware, handler.FollowerList)
	h.GET("/friends/list", middleware.AuthMiddleware, handler.FriendList)
	h.Spin()
}
