package router

import (
	"Tiktok/biz/handler"
	"Tiktok/biz/middleware"

	"github.com/cloudwego/hertz/pkg/app/server"
)

func SetRouters(commentHandler *handler.CommentHandler, userHandler *handler.UserHandler, videoHandler *handler.VideoHandler, socialHandler *handler.SocialHandler, likesHandler *handler.LikesHandler, websocketHandler *handler.WebsocketSever) {
	h := server.Default(
		server.WithHostPorts(":8888"),
		server.WithMaxRequestBodySize(10*1024*1024),
	)
	defer h.Close()
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
	video := h.Group("/video")
	video.Use(middleware.AuthMiddleware)
	{
		video.POST("/publish", videoHandler.VideoPublish)
		video.GET("/list", videoHandler.VideoList)
		video.POST("/search", videoHandler.VideoSearch)
		video.GET("/popular", videoHandler.VideoPopular)
	}
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
	h.POST("/relation/action", middleware.AuthMiddleware, socialHandler.RelationAction)
	h.GET("/following/list", middleware.AuthMiddleware, socialHandler.FollowingList)
	h.GET("/follower/list", middleware.AuthMiddleware, socialHandler.FollowerList)
	h.GET("/friends/list", middleware.AuthMiddleware, socialHandler.FriendList)
	h.POST("/friend/add", middleware.AuthMiddleware, socialHandler.AddFriend)
	h.DELETE("/friend/delete", middleware.AuthMiddleware, socialHandler.DeleteFriend)

	h.GET("/ws", middleware.AuthMiddleware, websocketHandler.WebSocketHandler)
	h.Spin()
}
