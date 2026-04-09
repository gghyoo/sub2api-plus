# Sub2API Plus

<div align="center">

[![Go](https://img.shields.io/badge/Go-1.26.2-00ADD8.svg)](https://golang.org/)
[![Vue](https://img.shields.io/badge/Vue-3.4+-4FC08D.svg)](https://vuejs.org/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED.svg)](https://www.docker.com/)

**AI API Gateway Platform — Enhanced Fork of [Sub2API](https://github.com/Wei-Shaw/sub2api)**

</div>

---

## About

Sub2API Plus is an enhanced fork of [Wei-Shaw/sub2api](https://github.com/Wei-Shaw/sub2api) (upstream v0.1.110). For full documentation, deployment guides, and upstream features, please refer to the [upstream README](https://github.com/Wei-Shaw/sub2api).

This fork focuses on homepage enhancements and deployment simplification.

## Plus Features

### Homepage Model Pricing Table

Displays model pricing info (input/output/cache prices) with platform and endpoint color-coded badges. Prices are cached for performance.

### Model Guide Tip Card

A customizable green info card below the pricing table for model selection guidance. Supports Markdown and HTML, configurable from **Admin Panel → System Settings → Site Settings → Model Guide Tip**.

### Tool Usage Guide

Quick start guides for Claude Code, Codex CLI, Cursor, etc. with one-click config copy and dynamic base URL.

### Endpoint Color Coding

Visual differentiation of endpoints in the model table:

| Endpoint | Color |
|----------|-------|
| `/v1/messages` (Anthropic) | Orange |
| `/v1/responses` (OpenAI) | Indigo |
| `/v1/chat/completions` (OpenAI) | Green |
| `/v1beta/*` (Gemini) | Blue |
| `/antigravity/*` | Purple |

### Pricing Data Source

Switched from LiteLLM to OpenRouter for model pricing data, with DB → OpenRouter → omit fallback chain.

### OpenAI Endpoint Optimization

`/v1/chat/completions` removed from default OpenAI routing; only `/v1/responses` retained for Codex CLI compatibility.

### DockerHub Only

Docker images published to DockerHub only (`gghyoo/sub2api-plus`). GHCR removed.

## Deployment

Refer to the [upstream deployment guide](https://github.com/Wei-Shaw/sub2api#deployment), replacing image references with `gghyoo/sub2api-plus`.

### Docker Quick Start

```bash
mkdir -p sub2api-deploy && cd sub2api-deploy
curl -sSL https://raw.githubusercontent.com/gghyoo/sub2api-plus/main/deploy/docker-deploy.sh | bash
docker compose up -d
```

### Build from Source

```bash
git clone https://github.com/gghyoo/sub2api-plus.git
cd sub2api-plus/frontend && pnpm install && pnpm build
cd ../backend && go build -tags embed -o sub2api ./cmd/server
```

## Changelog

### v0.2.4
- New: Model Guide tip card below pricing table (Markdown/HTML, configurable from admin)
- New: @tailwindcss/typography for better Markdown rendering

### v0.2.3
- New: Endpoint color coding (Anthropic/OpenAI/Gemini/Antigravity) in model table

### v0.2.0 - v0.2.2
- New: Homepage model pricing table with platform & endpoint columns, model ID copy
- New: Tool usage guide (Claude Code, Codex CLI, Cursor, etc.)
- New: `/v1/models` public pricing with DB → OpenRouter → omit fallback
- Changed: OpenAI endpoint only retains `/v1/responses`
- Changed: Pricing data source from LiteLLM to OpenRouter
- Changed: Docker images published to DockerHub only
- Fixed: Update check and release links pointing to gghyoo/sub2api-plus
- Fixed: Stream timeout test timing race assertion

## License

MIT License
