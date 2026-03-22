package router

import (
	"Tiktok/biz/handler"
	"Tiktok/biz/middleware"

	"github.com/cloudwego/hertz/pkg/app/server"
)

func SetRouters(commentHandler *handler.CommentHandler, userHandler *handler.UserHandler, videoHandler *handler.VideoHandler, socialHandler *handler.SocialHandler, likesHandler *handler.LikesHandler) {
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
		user.POST("/refresh", userHandler.RefreshToken)
	}
	authMfa := h.Group("/auth/mfa")
	{
		authMfa.GET("/qrcode", middleware.AuthMiddleware, userHandler.MfaQrcode)
		authMfa.POST("/bind", middleware.AuthMiddleware, userHandler.MfaBind)
	}
	//投稿、发布列表、搜索视频、热门排行榜
	video := h.Group("/video")
	video.Use(middleware.AuthMiddleware)
	{
		video.POST("/publish", videoHandler.VideoPublish)
		video.GET("/list", videoHandler.VideoList)
		video.POST("/search", videoHandler.VideoSearch)
		video.GET("/popular", videoHandler.VideoPopular)
	}
	//点赞操作、点赞列表、评论、评论列表、删除评论
	like := h.Group("/like")
	like.Use(middleware.AuthMiddleware)
	{
		like.POST("/action", likesHandler.LikeAction)
		like.GET("/list", likesHandler.LikeList)
	}
	comment := h.Group("/comment")
	comment.Use(middleware.AuthMiddleware)
	{
		comment.POST("/publish", commentHandler.CommentPublish)
		comment.GET("/list", commentHandler.CommentList)
		comment.DELETE("/delete", commentHandler.CommentDelete)
	}
	//关注操作、关注列表、粉丝列表、好友列表
	h.POST("/relation/action", middleware.AuthMiddleware, socialHandler.RelationAction)
	h.GET("/following/list", middleware.AuthMiddleware, socialHandler.FollowingList)
	h.GET("/follower/list", middleware.AuthMiddleware, socialHandler.FollowerList)
	h.GET("/friends/list", middleware.AuthMiddleware, socialHandler.FriendList)

	h.GET("/ws", middleware.AuthMiddleware, handler.WebSocketHandler)
	h.Spin()
}
