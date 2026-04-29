# Tiktok

## 技术栈

| 分类     | 技术选型          |
|----------|-------------------|
| web框架   | Hertz             |
| 数据库   | MySQL + sqlx      |
| 缓存     | Redis             |

## 实现接口
### 用户模块：

    [x] 注册。
    [x] 登陆。(通过json返回access_token和refresh_token)
    [x] 用户信息。
    [x] 上传头像。
    [x] 绑定mfa。
    [x] 获取mfa qrcode。
    [x] 刷新token。
### 视频模块：

    [] 视频流
    [x] 发布视频。
    [x] 视频列表。
    [x] 热门排行榜。
    [x] 搜索视频。
### 互动模块：

    [x] 点赞操作。
    [x] 点赞列表。
    [x] 评论。
    [x] 评论列表。
    [x] 删除评论。
### 社交模块：

    [ ] websocket聊天(仅实现一对一聊天,在线群聊)
    [x] 关注操作。
    [x] 关注列表。
    [x] 粉丝列表。
    [x] 好友列表。
    [x] 添加好友
    [x] 删除好友

## 目录
    ├── biz
    │   ├── dal
    │   │   ├── cache
    │   │   └── dao
    │   ├── entity
    │   ├── handler
    │   │   ├── chat
    │   │   ├── mfa
    │   │   ├── react
    │   │   ├── social
    │   │   ├── user
    │   │   └── video
    │   ├── middleware
    │   ├── model
    │   │   ├── ai
    │   │   ├── api
    │   │   ├── chat
    │   │   ├── common
    │   │   ├── mfa
    │   │   ├── react
    │   │   ├── social
    │   │   ├── user
    │   │   └── video
    │   ├── router
    │   │   ├── chat
    │   │   ├── mfa
    │   │   ├── react
    │   │   ├── social
    │   │   ├── user
    │   │   └── video
    │   └── service
    │       ├── ai
    │       ├── comment
    │       ├── like
    │       ├── mfa
    │       ├── social
    │       ├── user
    │       ├── video
    │       └── websocket
    ├── docs
    ├── idl
    ├── mcp_service
    │   └── tools
    ├── pkg
    │   ├── config
    │   ├── consts
    │   └── utils
    └── script
## 接口文档
k7wl3pn34m.apifox.cn



