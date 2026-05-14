# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a TikTok-like application using a dual-framework architecture:
- **Hertz** (HTTP web framework) for REST API handlers
- **Kitex** (RPC framework) + etcd for service registration and discovery

## Commands

```bash
# Build all services
./build.sh

# Run linter
golangci-lint run

# Run a specific test
go test ./pkg/utils/... -v

# Generate code from IDL (requires kitex/thriftgo)
cd idl && thriftgo -g go:... *.proto

# Generate Hertz handler code
hz update -mod Tiktok -idl idl/api.thrift
```

## Architecture

### Dual Framework Structure

Each module (user, video, mfa) exists in two layers:

| Layer | Framework | Purpose | Entry Point |
|-------|-----------|---------|-------------|
| `cmd/*/main.go` | Hertz | HTTP API server | Port 8x90 |
| `internal/*/handler.go` | Kitex | RPC service | Port 889x |

The `biz/` directory contains Hertz handlers (business logic for HTTP), while `internal/` contains Kitex service implementations (RPC business logic).

### Module Organization

```
cmd/
  ├── mfa/main.go       # MFA RPC service (Kitex, port 8890)
  ├── user/main.go      # User RPC service (Kitex)
  └── video/main.go     # Video RPC service (Kitex)

internal/
  ├── mfa/              # MFA Kitex handler + service
  ├── user/             # User Kitex handler + service
  └── video/            # Video Kitex handler + service

biz/
  ├── handler/          # HTTP handlers (receives calls from Hertz)
  ├── service/          # Business logic for HTTP layer
  ├── model/            # Request/Response models
  ├── router/           # Hertz route definitions
  ├── middleware/       # Hertz middleware
  ├── rpc/              # Kitex RPC clients (calls internal/ services)
  └── dal/              # Data access layer (dao + cache)

idl/                   # Protobuf IDL definitions
kitex_gen/            # Generated Kitex code from IDL
pkg/
  ├── config/          # Configuration (config.yaml, sentinel.yaml)
  ├── consts/          # Constants
  └── utils/           # Utility functions
```

### Service Communication

HTTP requests → Hertz handlers (`biz/handler/`) → Kitex RPC clients (`biz/rpc/`) → Kitex services (`internal/*/handler.go`)

### Key Dependencies

- **Database**: MySQL with `sqlx`
- **Cache**: Redis
- **Rate Limiting**: Sentinel
- **Service Discovery**: etcd (for Kitex services)
- **Authentication**: JWT (access + refresh tokens)

### Configuration

- Environment variables in `.env` (see `.env.example`)
- Application config: `pkg/config/config.yaml`
- Rate limit rules: `pkg/config/sentinel.yaml`
- Database schema: `pkg/config/init.sql`