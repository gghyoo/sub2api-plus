# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What is this project?

Sub2API is an **AI API Gateway Platform** for subscription quota distribution. It reverse-proxies upstream AI services (Anthropic Claude, OpenAI, Google Gemini, Antigravity) through platform-generated API Keys, handling authentication, token-level billing, intelligent account scheduling/load-balancing, concurrency control, and rate limiting.

- Go module: `github.com/Wei-Shaw/sub2api`, Go 1.26.1
- Stack: Go backend (Ent ORM + Gin) + Vue 3 frontend (pnpm)
- Database: PostgreSQL 16 + Redis
- License: MIT

## Build & Development Commands

All commands run from `backend/` unless noted.

```bash
# Run server
cd backend && go run ./cmd/server/

# Unit tests (build tag: unit)
go test -tags=unit ./...

# Integration tests (build tag: integration, uses testcontainers)
go test -tags=integration ./...

# Lint (golangci-lint v2)
golangci-lint run ./...

# Regenerate Ent ORM code after schema changes
go generate ./ent

# Regenerate Wire DI code
go generate ./cmd/server

# Full generate (Ent + Wire)
go generate ./ent && go generate ./cmd/server

# Frontend (must use pnpm, NOT npm)
cd frontend && pnpm install && pnpm dev   # dev server
cd frontend && pnpm build                  # production build

# Build backend binary
cd backend && CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o bin/server ./cmd/server
```

## Architecture

### Layered Architecture (enforced by depguard)

```
handler → service → repository
  ↓         ↓          ↓
 NO repo  NO repo    DB/Redis
```

- **handler** (`internal/handler/`) — HTTP handlers, must only import service
- **service** (`internal/service/`) — business logic, must only import repository interfaces
- **repository** (`internal/repository/`) — data access, implements service-layer interfaces
- Depguard linter will fail if handler imports repository directly

### Dependency Injection: Google Wire

- ProviderSets in each package's `wire.go`: config, repository, service, handler, server/middleware
- Generated output: `cmd/server/wire_gen.go`
- After adding/removing providers: `go generate ./cmd/server`

### Key Subsystems

- **Gateway routing** (`internal/server/routes/gateway.go`): `/v1/messages` (Claude), `/v1/chat/completions` (OpenAI), `/v1/responses` (OpenAI Responses), `/v1beta/models/*` (Gemini), `/antigravity/v1/*` — auto-routed by API Key's group platform
- **Account scheduling** (`internal/service/`): priority-based with sticky sessions, top-K load balancing, concurrency slots, outbox-based snapshot rebuilding
- **Auth middleware** (`internal/server/middleware/`): JWT (admin/user sessions) + API Key auth with L1 (in-memory) + L2 (Redis) caching, singleflight
- **Background services**: All have `Start()`/`Stop()` lifecycle, started via Wire — token refresh, account/subscription expiry, ops metrics/aggregation/alerts, billing cache, email queue, scheduled tests, backups

### Database: Ent ORM

- Schema definitions (source of truth): `backend/ent/schema/` — 25 entities
- Key entities: Account (upstream credentials), User, APIKey, Group (platform + billing + scheduling), UsageLog, UserSubscription, Proxy
- After modifying any schema: `go generate ./ent` then commit generated `ent/` files
- Ent client + raw `*sql.DB` both available

### Configuration

- YAML config file (`config.yaml`), loaded via Viper with env var overrides (`_` separator, e.g. `DATABASE_HOST` → `database.host`)
- Two run modes: `standard` (full SaaS with billing) and `simple` (no billing)
- First-run setup wizard if no config found
- Full example: `deploy/config.example.yaml`

## Critical Rules

### Ent Schema Changes
1. Edit `backend/ent/schema/*.go`
2. Run `go generate ./ent`
3. Commit the generated `ent/` directory — generated code must be in git

### Interface Changes
When adding methods to a service interface, ALL test stubs in `internal/testutil/stubs.go` and any mocks must be updated — Go will fail to compile otherwise.

### Frontend Dependencies
Always use **pnpm**. After changing `package.json`, run `pnpm install` and commit `pnpm-lock.yaml`. Never use npm — leftover `node_modules` from npm will conflict.

### Build Tags for Embedding
- Production Docker build uses `-tags embed` to embed frontend into Go binary
- Local dev runs without the tag (frontend served separately by Vite dev server)

## Test Patterns

- **Unit tests**: `//go:build unit` tag, run with `go test -tags=unit ./...`
- **Integration tests**: `//go:build integration` tag, use testcontainers (PostgreSQL + Redis)
- **E2E tests**: `//go:build e2e` tag, in `internal/integration/`
- Test stubs: `internal/testutil/stubs.go` provides zero-value interface implementations
- Test fixtures and HTTP helpers: `internal/testutil/fixtures.go`, `internal/testutil/httptest.go`

## CI Workflows

- **backend-ci.yml**: unit tests → integration tests → golangci-lint v2.7
- **security-scan.yml**: govulncheck + gosec + pnpm audit (runs on push + weekly)
- **release.yml**: triggered on `v*` tags, multi-stage Docker build via GoReleaser, publishes to GHCR + DockerHub

## Deploy to Server (Binary Replacement)

部署服务器: `ggh@192.168.160.145`，容器名: `sub2api-plus`，docker目录: `/home/ggh/dockers/sub2api-plus`

容器内二进制路径为 `/app/sub2api-plus`（由 `docker-entrypoint.sh` 启动），**不是** `/app/server`。

### 步骤

```bash
# 1. 构建前端
cd frontend && pnpm build

# 2. 交叉编译 Linux 二进制（嵌入前端）
cd backend && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags embed -ldflags="-s -w" -trimpath -o bin/server ./cmd/server

# 3. 上传二进制到服务器
scp backend/bin/server ggh@192.168.160.145:/tmp/sub2api-server

# 4. 确保容器处于运行状态（如果已运行则跳过）
ssh ggh@192.168.160.145 "cd /home/ggh/dockers/sub2api-plus && docker compose up -d --no-deps sub2api"

# 5. 趁容器还在运行，复制二进制并设置可执行权限
ssh ggh@192.168.160.145 "docker cp /tmp/sub2api-server sub2api-plus:/app/sub2api-plus && docker exec sub2api-plus chmod +x /app/sub2api-plus"

# 6. 重启容器
ssh ggh@192.168.160.145 "cd /home/ggh/dockers/sub2api-plus && docker compose restart sub2api"

# 7. 查看日志确认启动成功
ssh ggh@192.168.160.145 "docker logs sub2api-plus --tail 10"
```

### 注意事项

- `docker cp` 复制的文件**不会继承**原有文件的执行权限，必须 `docker exec chmod +x`
- 容器停止后无法 `docker exec`，所以必须先启动容器（即使 crash 循环中也行），再 cp + chmod
- 如果容器已完全停止无法 exec，可用 `docker compose run --rm --entrypoint sh sub2api -c 'cp /app/data/sub2api-server /app/sub2api-plus && chmod +x /app/sub2api-plus'`（需先将二进制放到 `./data/` 挂载目录）
- 访问地址: `http://192.168.160.145:8081`（映射到容器内 8080 端口）

## Local Dev Environment

| Service | Connection |
|---------|-----------|
| PostgreSQL 16 | `127.0.0.1:5432`, user=`sub2api`, password=`sub2api`, dbname=`sub2api` |
| Redis | `127.0.0.1:6379`, no password |
