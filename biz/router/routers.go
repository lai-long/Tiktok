package router

import (
	"Tiktok/biz/handler"
	"Tiktok/biz/middleware"

	"github.com/cloudwego/hertz/pkg/app/server"
)

func SetRouters() {
	h := server.Default(
		server.WithHostPorts("localhost:8888"),
	)
	defer h.Close()
	//注册、登录、用户信息、上传头像
	user := h.Group("/user")
	{
		user.GET("/info", middleware.AuthMiddleware, handler.UserInfo)
		user.POST("/login", handler.UserLogin)
		user.POST("/register", handler.UserRegister)
		user.POST("/avatar/upload", handler.UserAvatar)
	}
	//投稿、发布列表、搜索视频、热门排行榜
	video := h.Group("/video")
	video.Use(middleware.AuthMiddleware)
	{
		video.POST("/publish", handler.VideoPublish)
		video.GET("/list", handler.VideoList)
		video.POST("/search", handler.VideoSearch)
		//video.GET("/popular", handler.VideoPopular)
	}
	//点赞操作、点赞列表、评论、评论列表、删除评论
	like := h.Group("/like")

	{
		like.POST("/action", handler.LikeAction)
		like.GET("/list", handler.LikeList)
	}
	comment := h.Group("/comment")
	{
		comment.POST("/publish", handler.CommentPublish)
		comment.GET("/list", handler.CommentList)
		comment.DELETE("/delete", handler.CommentDelete)
	}
	//关注操作、关注列表、粉丝列表、好友列表
	h.POST("/relation/action", handler.RelationAction)
	h.GET("/following/list", handler.FollowingList)
	h.GET("/follower/list", handler.FollowerList)
	h.GET("/friend/list", handler.FriendList)
	h.Spin()
}
