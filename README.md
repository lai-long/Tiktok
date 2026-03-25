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
    用户模块：
    [x] 注册。
    [x] 登陆。(通过json返回access_token和refresh_token)
    [x] 用户信息。
    [x] 上传头像。
    [x] 绑定mfa。
    [x] 获取mfa qrcode。
    [x] 刷新token。
    视频模块：
    [] 视频流
    [x] 发布视频。
    [x] 视频列表。
    [x] 热门排行榜。
    [x] 搜索视频。
    互动模块：
    [x] 点赞操作。
    [x] 点赞列表。
    [x] 评论。
    [x] 评论列表。
    [x] 删除评论。
    社交模块：
    [] websocket聊天(仅实现一对一聊天)
    [x] 关注操作。
    [x] 关注列表。
    [x] 粉丝列表。
    [x] 好友列表。
    [x] 添加好友
    [x] 删除好友

## 目录
    ├── biz
    │   ├── dao
    │   │   ├── db
    │   │   │   ├── dbInit.go
    │   │   │   ├── mfa.go
    │   │   │   ├── react.go
    │   │   │   ├── social.go
    │   │   │   ├── user.go
    │   │   │   ├── video.go
    │   │   │   └── websocket.go
    │   │   └── re
    │   │       ├── redisInit.go
    │   │       ├── user.go
    │   │       ├── video.go
    │   │       └── websocket.go
    │   ├── handler
    │   │   ├── comment.go
    │   │   ├── like.go
    │   │   ├── mfa.go
    │   │   ├── social.go
    │   │   ├── user.go
    │   │   ├── video.go
    │   │   └── websocket.go
    │   ├── middleware
    │   │   └── auth.go
    │   ├── model
    │   │   ├── dto
    │   │   │   ├── comment.go
    │   │   │   ├── response.go
    │   │   │   ├── user.go
    │   │   │   ├── video.go
    │   │   │   └── websocket.go
    │   │   └── entity
    │   │       ├── comment.go
    │   │       ├── user_entity.go
    │   │       └── video.go
    │   ├── router
    │   │   └── routers.go
    │   └── service
    │       ├── comment.go
    │       ├── like.go
    │       ├── mfa.go
    │       ├── social.go
    │       ├── user.go
    │       ├── video.go
    │       └── websocket.go
    ├── Dockerfile
    ├── go.mod
    ├── go.sum
    ├── main.go
    ├── pkg
    │   ├── conf
    │   │   ├── config.go
    │   │   └── config.yaml
    │   ├── consts
    │   │   └── consts.go
    │   └── utils
    │       └── utils.go
    └── README.md

## 接口文档
k7wl3pn34m.apifox.cn



