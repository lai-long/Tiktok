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
| AI 对话 | Bifrost (MCP) + OpenAI 兼容 API（MiniMax / DeepSeek / OpenAI 等） |
| 限流熔断 | Sentinel |
| 日志 | Zap + Lumberjack → Filebeat → Elasticsearch → Kibana |
| WebSocket | Gorilla WebSocket |
| 配置管理 | Viper (热更新) |
| 容器化 | Docker + docker-compose |

## 服务列表

| 服务 | 端口 | 协议 | 服务发现 |
|------|------|------|---------|
| API Gateway | 8888 | HTTP | — |
| WebSocket | 8881 | HTTP+WS | — |
| userService | — | Kitex RPC | etcd |
| videoService | — | Kitex RPC | etcd |
| reactService | — | Kitex RPC | etcd |
| socialService | — | Kitex RPC | etcd |
| mfaService | — | Kitex RPC | etcd |

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
- AI 对话（`@AI` 触发，MCP + Bifrost，支持所有 OpenAI 兼容 API）


## 架构概览

```
                    ┌──────────────┐
                    │   Client     │
                    └──────┬───────┘
                           │ HTTP
                    ┌──────┴───────┐
                    │ API Gateway  │ :8888
                    │   (Hertz)    │
                    └──────┬───────┘
                           │ Kitex RPC (etcd 服务发现)
          ┌────────┬───────┼───────┬────────┐
          │        │       │       │        │
       ┌──┴──┐ ┌──┴──┐ ┌──┴──┐ ┌──┴──┐ ┌──┴──┐
       │user │ │video│ │react│ │social│ │ mfa │
       └──┬──┘ └──┬──┘ └──┬──┘ └──┬──┘ └──┬──┘
          │       │       │       │       │
          └───────┴───────┼───────┴───────┘
                          │
               ┌──────────┼──────────┐
               │          │          │
          ┌────┴───┐ ┌────┴───┐ ┌────┴───┐
          │ MySQL  │ │ Redis  │ │  etcd  │
          └────────┘ └────────┘ └────────┘

     ┌──────────────┐
     │ WebSocket    │ :8881
     │  (Gorilla)   │
     └──────────────┘
```


### 环境要求

- Go 1.26+
- Docker
- MySQL, Redis, etcd（本地开发时可用 Docker 启动）

### 1. 配置文件

```bash
# 复制配置文件
cp config/config.example.yaml config/config.yaml
cp .env.example .env

# 编辑 .env，填入你的密钥和数据库密码
vim .env
```

### 2. 本地开发

```bash
# 启动基础设施（MySQL, Redis, etcd）
make run-infra

# 数据库初始化
mysql -u root -p < config/init.sql

# 启动所有服务
make run
```

或者手动逐个启动：

```bash
go run ./cmd/user/
go run ./cmd/video/
go run ./cmd/react/
go run ./cmd/social/
go run ./cmd/mfa/
go run ./cmd/api/    # :8888
go run ./cmd/ws/     # :8881
```

### 3. Docker 部署

```bash
# 构建并启动所有服务
make docker-up
# 或
cd docker && docker compose up -d

# 查看服务状态
docker compose ps

# 查看日志
docker logs -f tiktok-api
docker logs -f tiktok-user

# 重新构建并启动
make docker-rebuild
# 或
cd docker && docker compose up -d --build

# 停止所有服务
make docker-down
# 或
cd docker && docker compose down
```

## Makefile 命令

| 命令 | 说明 |
|------|------|
| `make help` | 查看所有可用命令 |
| `make run-infra` | 启动基础设施（MySQL, Redis, etcd） |
| `make run` | 本地启动所有服务（需先启动基础设施） |
| `make docker-up` | Docker 启动所有服务 |
| `make docker-rebuild` | 重新构建镜像并启动 |
| `make docker-down` | 停止所有 Docker 容器 |
| `make test` | 运行测试 |

## 环境变量

| 变量 | 说明             | 默认值 |
|------|----------------|--------|
| `MYSQL_PASSWORD` | MySQL 密码       | — |
| `REDIS_PASSWORD` | Redis 密码       | — |
| `JWT_ACCESS_SECRET` | JWT 密钥         | — |
| `JWT_REFRESH_SECRET` | Refresh 密钥     | — |
| `OPENAI_API_KEY` | OpenAI 兼容 API Key（MiniMax / DeepSeek / OpenAI 等） | — |
| `QINIU_ACCESS_KEY` | 七牛云 Access Key | — |
| `QINIU_SECRET_KEY` | 七牛云 Secret Key | — |
| `ETCD_ADDR` | etcd 地址        | `127.0.0.1:2379` |
| `CONFIG_PATH` | 配置文件目录         | `./config` |
| `SENTINEL_PATH` | Sentinel 规则目录  | `./config` |
| `ENV_PATH` | .env 文件路径      | `./.env` |


## 测试

```bash
make test
go test -race -count=1 -coverprofile=coverage.out ./...
```

## 项目结构

```
├── biz/                  # HTTP API 层 (Hertz)
│   ├── handler/          # 请求处理器
│   ├── middleware/       # JWT 认证、请求日志
│   ├── model/            # 请求/响应模型 (pb 生成)
│   ├── router/           # 路由注册
│   └── rpc/              # RPC 客户端 (etcd 服务发现)
├── cmd/                  # 服务入口
│   ├── api/              # API 网关
│   ├── ws/               # WebSocket 服务
│   └── user/ video/ react/ social/ mfa/   # Kitex RPC 服务
├── config/               # 配置文件
│   ├── config.yaml       # 主配置
│   ├── config.example.yaml
│   ├── init.sql          # 数据库初始化
│   ├── sentinel.yaml     # Sentinel 限流规则
│   └── filebeat.yml      # Filebeat 日志采集配置
├── docker/               # Docker 部署
│   ├── Dockerfile        # 多阶段构建 (SERVICE arg)
│   └── docker-compose.yml
├── idl/                  # Protobuf IDL 定义
├── internal/             # RPC 服务实现
│   ├── config/           # Viper 配置加载 (含热更新)
│   ├── middleware/       # Sentinel 限流
│   └── user/ video/ react/ social/ mfa/ ws/
│       └── handler.go + service/   # 各模块处理器+业务逻辑
├── kitex_gen/            # Kitex 生成代码
├── logs/                 # 本地日志输出目录
├── mcp_service/          # MCP 工具服务 (天气查询等)
├── pkg/                  # 公共库
│   ├── consts/           # 错误码 (格式: XYYZZZ)
│   ├── dal/dao/          # MySQL (sqlx)
│   ├── dal/cache/        # Redis 缓存
│   ├── entity/           # 数据库实体
│   ├── logger/           # Zap 日志 (双路输出)
│   └── utils/            # JWT、bcrypt、七牛云上传
├── third_party/          # 第三方依赖
├── Makefile              # 便捷命令
└── .env.example          # 环境变量模板
```

## 错误码规范

格式 `XYYZZZ`（6 位数字）：

| 位 | 含义 | 值 |
|----|------|-----|
| X | 类型 | 1=客户端错误, 2=服务端错误 |
| YY | 模块 | 01=用户, 02=视频, 03=互动, 04=社交, 05=WS, 06=MFA |
| ZZZ | 具体错误 | 自增编号 |

示例：`101001` = 请求层(1) + 用户模块(01) + 用户名已存在(001) |

使用 `GetErrorCodeMsg()` 获取错误信息。

## Sentinel 限流熔断

规则配置在 `config/sentinel.yaml`，支持：
- **流控**：快速失败 / 排队等待
- **熔断**：慢调用比例 / 错误比例 / 错误计数

## AI 对话

基于 Bifrost 框架，支持所有 OpenAI 兼容 API。通过 `config.yaml` 的 `api.base_url` 切换模型提供商：

```yaml
api:
  base_url: "https://api.minimaxi.com"   # MiniMax
  # base_url: "https://api.deepseek.com"  # DeepSeek
  # base_url: "https://api.openai.com"    # OpenAI
```

环境变量 `OPENAI_API_KEY` 对应填入对应平台的 API Key。用户在聊天中发送 `@AI` 触发 AI 对话，支持 MCP 工具调用（如天气查询）。

## 配置热更新

基于 Viper + fsnotify，`config/config.yaml` 变更后自动热加载，无需重启服务。
