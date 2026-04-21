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
# Run server (with embedded frontend) — 本地开发必须用此方式
# 注意顺序：必须先编译前端，再启动后端（后端编译时嵌入前端产物）
cd frontend && pnpm build && cd ../backend && go run -tags embed ./cmd/server/

# Run server (backend only, no frontend)
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
cd frontend && pnpm install && pnpm build   # build for embedding

# Build backend binary
cd backend && CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o bin/server ./cmd/server

# 一键编译前端 + 后端（嵌入前端）并部署到服务器（macOS/Linux 直接运行；Windows 请使用 Git Bash）
bash deploy/build-and-deploy.sh
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

## Release 发布流程

```bash
# 1. 提交代码并推送到 GitHub
git push origin main

# 2. 检查 CI 是否通过（仓库: gghyoo/sub2api-plus）
gh run list --limit 3 --repo gghyoo/sub2api-plus
# 如有失败，查看日志修复后重新推送
gh run view <run-id> --repo gghyoo/sub2api-plus --log-failed

# 3. CI 通过后，更新版本号
# 编辑 backend/cmd/server/VERSION，如 0.2.2 → 0.2.3
# 注意：仅在发布 release 时才变更此文件

# 4. 提交版本号并推送
git add backend/cmd/server/VERSION
git commit -m "chore: sync VERSION to x.x.x [skip ci]"
git push origin main

# 5. 创建带注释的 tag 并推送（触发 release workflow 构建 Docker 镜像）
#    重要：必须使用 annotated tag（-a），release.yml 才能正确读取 tag message 作为 Release Notes
git tag -a vx.x.x -m "变更说明..."
git push origin vx.x.x

# 6. 等待 release workflow 完成后，GitHub Release 会自动创建
#    检查 release workflow 状态：
gh run list --limit 3 --repo gghyoo/sub2api-plus --workflow=release.yml
```

### 注意事项

- `[skip ci]` 避免版本号提交触发额外的 CI 运行
- **tag 必须推送到 GitHub**：只有推送到远程的 `v*` tag 才会触发 `release.yml` workflow
- 使用 annotated tag（`git tag -a`），tag message 会自动成为 GitHub Release 的发布说明
- 如果 CI 失败，先修复问题再重新推送，不要跳过 CI 直接发布
- tag 推送后 release.yml 会自动：构建前端 → 编译多架构二进制 → 构建 Docker 镜像 → 推送到 GHCR + DockerHub → 创建 GitHub Release → 发送 Telegram 通知

## 测试部署（二进制热更新）

快速将本地修改部署到测试服务器（`ggh@10.10.200.21`），用于功能验证：

```bash
# 一键完成：编译前端 → 编译后端（嵌入前端）→ 上传二进制 → 替换容器内文件 → 重启 → 健康检查
bash deploy/build-and-deploy.sh
```

脚本默认连接测试服务器 `ggh@10.10.200.21`。可通过环境变量覆盖：

```bash
SERVER_HOST=192.168.160.21 SERVER_USER=ggh bash deploy/build-and-deploy.sh
```

**平台说明**：
- macOS / Linux：直接运行
- Windows：请通过 **Git Bash**（随 Git for Windows 安装）或 **WSL** 执行

脚本会交叉编译 `GOOS=linux GOARCH=amd64` 的二进制并嵌入前端产物，通过 `scp` 上传后执行 `docker cp` + `docker restart` 完成热更新。

**注意**：热更新仅在当前容器生命周期内有效。`docker compose down && up` 会恢复为镜像内版本。

## 生产服务器更新

生产服务器（`192.168.160.145`）只通过 `docker pull` 拉取最新镜像更新：

```bash
ssh ggh@192.168.160.145
cd /home/ggh/dockers/sub2api-plus
docker compose pull && docker compose up -d
```

## Local Dev Environment

| Service | Connection |
|---------|-----------|
| PostgreSQL 16 | `127.0.0.1:5432`, user=`sub2api`, password=`sub2api`, dbname=`sub2api` |
| Redis | `127.0.0.1:6379`, no password |
