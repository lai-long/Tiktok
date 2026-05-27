# Tiktok

基于 Go + CloudWeGo 生态的短视频平台，采用微服务架构。

## 技术栈

| 分类 | 技术选型 |
|------|---------|
| Web 框架 | Hertz |
| RPC 框架 | Kitex + etcd |
| 数据库 | MySQL + sqlx |
| 缓存 | Redis (go-redis) |
| 消息协议 | Protobuf |
| 认证 | JWT (Access + Refresh Token) |
| 双因素认证 | TOTP |
| 文件存储 | 七牛云 |
| AI 对话 | Bifrost (MCP) + MiniMax |
| 限流熔断 | Sentinel |
| 日志 | Zap + Lumberjack |
| WebSocket | Gorilla WebSocket |
| 配置管理 | Viper (热更新) |
| 容器化 | Docker + docker-compose |


## 服务列表

| 服务 | 端口 | 协议 | 发现 |
|------|------|------|------|
| API Gateway | 8888 | HTTP | — |
| WebSocket | 8881 | HTTP+WS | — |
| userService | — | Kitex | etcd |
| videoService | — | Kitex | etcd |
| reactService | — | Kitex | etcd |
| socialService | — | Kitex | etcd |
| mfaService | — | Kitex | etcd |

## 功能模块

### 用户
- 注册 / 登录（JWT 双 Token）
- 用户信息 / 上传头像（七牛云）
- Token 刷新
- MFA 绑定（TOTP 二维码）

### 视频
- 视频流 / 发布 / 列表
- 热门排行榜
- 搜索视频

### 互动
- 点赞 / 取消点赞 / 点赞列表
- 评论 / 删除评论 / 评论列表

### 社交
- 关注/取消关注、关注/粉丝/好友列表
- 添加/删除好友

### WebSocket 聊天
- 私信、群聊
- 15s 心跳 / 30s 超时
- 离线消息、历史记录分页
- AI 对话（`@AI` 触发，MCP + MiniMax）

## 快速开始

### 环境要求
Go 1.26+, MySQL, Redis, etcd

```bash
# 数据库初始化
mysql -u root -p < config/init.sql

# 配置（支持环境变量覆盖）
cp config/config.example.yaml config/config.yaml

# 构建所有服务
./build.sh

# 启动 RPC 服务（各开终端）
go run ./cmd/user/
go run ./cmd/video/
go run ./cmd/react/
go run ./cmd/social/
go run ./cmd/mfa/

# 启动网关和 WS
go run ./cmd/api/    # :8888
go run ./cmd/ws/     # :8881
```

### 环境变量

| 变量 | 说明 |
|------|------|
| `MYSQL_PASSWORD` | MySQL 密码 |
| `REDIS_PASSWORD` | Redis 密码 |
| `JWT_ACCESS_SECRET` | JWT 密钥 |
| `JWT_REFRESH_SECRET` | Refresh 密钥 |
| `OPENAI_API_KEY` | MiniMax API Key |
| `QINIU_ACCESS_KEY` | 七牛云 Access |
| `QINIU_SECRET_KEY` | 七牛云 Secret |
| `ETCD_ADDR` | etcd 地址 (默认 127.0.0.1:2379) |
| `CONFIG_PATH` | 配置文件目录 |
| `SENTINEL_PATH` | Sentinel 规则文件目录 |
| `ENV_PATH` | .env 文件路径 (默认 .env) |

### Docker 部署

```bash
# 构建并启动所有服务
cd docker && docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f api

# 停止
docker-compose down
```

### Makefile

```bash
make build       # 编译所有服务
make docker-up   # docker-compose 启动
make docker-down # 停止容器
make test        # 运行测试
make lint        # 代码检查
```

### 测试

```bash
go test -race -count=1 -coverprofile=coverage.out ./...
```

## 项目结构

```
├── biz/                  # HTTP API 层 (Hertz)
│   ├── handler/          # 请求处理器
│   ├── middleware/       # JWT 认证、日志
│   ├── model/            # 请求/响应模型 (pb 生成)
│   ├── router/           # 路由注册
│   └── rpc/              # RPC 客户端
├── cmd/                  # 服务入口
│   ├── api/ ws/ user/ video/ react/ social/ mfa/
├── internal/             # RPC 服务实现
│   ├── user/video/react/social/mfa/ws/
│   │   └── handler.go + service/   # 各模块处理器+业务逻辑
│   └── middleware/       # Sentinel 限流
├── idl/                  # Protobuf IDL
├── kitex_gen/            # Kitex 生成代码
├── pkg/
│   ├── config/           # Viper 配置 (含 sentinel.yaml)
│   ├── consts/           # 错误码 (格式: XYYZZZ)
│   ├── dal/dao/          # MySQL (sqlx)
│   ├── dal/cache/        # Redis 缓存
│   ├── entity/           # 数据库实体
│   ├── logger/           # Zap 日志
│   └── utils/            # JWT、bcrypt、七牛云上传
├── docker/               # Docker 部署
│   ├── Dockerfile        # 多阶段构建 (SERVICE arg)
│   ├── docker-compose.yml # 全服务编排
│   ├── config/            # Docker 环境配置文件
│   ├── env/               # 基础设施环境变量
│   └── script/
├── mcp_service/          # MCP 工具服务 (天气查询等)
├── Makefile              # 便捷命令
└── build.sh
```

## 错误码规范

格式 `XYYZZZ`：X=模块(1用户/2视频/3互动/4社交/5WS/6MFA)，YY=层级(00通用/01请求/02DB)，ZZZ=具体错误。使用 `GetErrorCodeMsg()` 获取错误信息。

## Sentinel 限流熔断

规则配置在 `config/sentinel.yaml`，支持流控（快速失败/排队等待）和熔断（慢调用/错误比例/错误计数）。

## 配置热更新

基于 Viper + fsnotify，配置文件变更自动热加载。
