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

    [x] 视频流
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

    [x] websocket聊天(仅实现一对一聊天,在线群聊)
    [x] 关注操作。
    [x] 关注列表。
    [x] 粉丝列表。
    [x] 好友列表。
    [x] 添加好友
    [x] 删除好友

## 目录
    ├── biz                     # HTTP层 (Hertz handlers)
    │   ├── handler             # HTTP请求处理
    │   ├── middleware          # 中间件 (auth认证)
    │   ├── model               # 请求/响应模型
    │   ├── router              # 路由注册
    │   └── rpc                 # RPC客户端
    ├── cmd                     # 服务入口
    │   ├── api                 # HTTP API服务 (:8888)
    │   ├── mfa/react/social/user/video  # RPC服务 (Kitex)
    │   └── ws                  # WebSocket服务 (:8881)
    ├── idl                     # protobuf IDL定义
    │   └── api
    ├── internal                # RPC层 (Kitex handlers)
    │   ├── mfa/react/social/user/video/ws
    │   │   ├── handler.go      # Kitex服务实现
    │   │   └── service/        # 业务逻辑
    │   └── ws
    ├── kitex_gen               # 生成的Kitex代码
    │   ├── mfa/react/social/user/video
    ├── pkg                     # 公共包
    │   ├── config              # 配置
    │   ├── consts             # 常量
    │   ├── dal                # 数据访问层 (dao/cache)
    │   ├── entity             # 实体
    │   └── utils              # 工具函数
    ├── Dockerfile, go.mod, go.sum, README.md
## 接口文档
k7wl3pn34m.apifox.cn