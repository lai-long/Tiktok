# Tiktok
## 飞书文档
https://my.feishu.cn/wiki/KjWKwoMYFiBMxxkBUIRcCJaonZg?from=from_copylink
## 技术栈

| 分类     | 技术选型          |
|----------|-------------------|
| web框架   | Hertz             |
| 数据库   | MySQL + sqlx      |
| 缓存     | Redis             |

## 实现接口
1、用户模块：注册、登陆、用户信息、上传头像、绑定mfa、获取mfa qrcode

2、视频模块：发布视频、视频列表、热门排行榜、搜索视频

3、互动模块：点赞操作、点赞列表、评论、评论列表、删除评论

4、社交模块：关注操作、关注列表、粉丝列表、好友列表

## 目录
    ├── biz
    │     ├── dao
    │     │     ├── db
    │     │     │     ├── dbInit.go
    │     │     │     ├── mfa.go
    │     │     │     ├── react.go
    │     │     │     ├── social.go
    │     │     │     ├── user.go
    │     │     │     └── video.go
    │     │     └── redis
    │     │         ├── redisInit.go
    │     │         └── video.go
    │     ├── handler
    │     │     ├── comment.go
    │     │     ├── like.go
    │     │     ├── mfa.go
    │     │     ├── social.go
    │     │     ├── user.go
    │     │     └── video.go
    │     ├── middleware
    │     │     └── auth.go
    │     ├── model
    │     │     ├── dto
    │     │     │     ├── comment.go
    │     │     │     ├── response.go
    │     │     │     ├── user.go
    │     │     │     └── video.go
    │     │     └── entity
    │     │         ├── comment.go
    │     │         ├── user_entity.go
    │     │         └── video.go
    │     ├── router
    │     │     └── routers.go
    │     └── service
    │         ├── comment.go
    │         ├── like.go
    │         ├── mfa.go
    │         ├── social.go
    │         ├── user.go
    │         └── video.go
    ├── Dockerfile
    ├── go.mod
    ├── go.sum
    ├── main.go
    ├── pkg
    │     ├── conf
    │     │     ├── config.go
    │     │     └── config.yaml
    │     ├── consts
    │     │     └── consts.go
    │     └── utils
    │         ├── checkimage.go
    │         ├── idgenerater.go
    │         ├── jwt.go
    │         └── password.go
    └── README.md

## 接口文档
k7wl3pn34m.apifox.cn

