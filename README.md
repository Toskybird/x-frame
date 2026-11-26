# X-Frame

可运行的 Go + Vue 基础脚手架：鉴权、用户、系统设置、安装向导与管理端布局壳，便于在此之上扩展业务。

边界与扩展方式见 [FRAMEWORK.md](FRAMEWORK.md)。本地开发见 [DEV_GUIDE.md](DEV_GUIDE.md)。

## 已包含

- HTTP 装配（Gin + Wire）、鉴权与用户管理
- 系统设置、安装向导、公告、仪表盘
- 管理端 / 用户端布局壳、i18n、通用组件
- PostgreSQL + Redis

## 技术栈

| 层 | 技术 |
|----|------|
| 后端 | Go、Gin、Ent、Wire、Viper、PostgreSQL、Redis |
| 前端 | Vue 3、Vite、TypeScript、Pinia、vue-i18n、Tailwind |
| 部署 | Docker Compose、systemd、Goreleaser |

Go module：`github.com/x-frame/x-frame`（请改成你自己的仓库路径）。

## 目录

```text
backend/     服务端（cmd / internal / ent / migrations）
frontend/    管理台与用户台
deploy/      Compose、安装脚本、示例配置
FRAMEWORK.md 框架边界与如何加业务
```

## 快速开始

```bash
# 基础设施（Postgres + Redis）
cd deploy && docker compose -f docker-compose.infra.yml up -d

# 后端
cd backend && go run ./cmd/server

# 前端
cd frontend && pnpm install && pnpm dev
```

首次访问会进入安装向导。默认库名等配置见 `deploy/.env.example` 与 `deploy/config.example.yaml`。

Windows 本机也可用绿色版 Postgres / Redis，见 `deploy/_internal/README.md` 与 [DEV_GUIDE.md](DEV_GUIDE.md)。

## 扩展业务

1. 后端：新增 `ent/schema` → `go generate ./ent` → repository / service / handler → 注册路由与 Wire。
2. 前端：加 `api/`、`views/`、路由与侧栏项。
3. 校验：`cd backend && go build ./...`；`cd frontend && pnpm exec vue-tsc --noEmit`。

更多约定见 [FRAMEWORK.md](FRAMEWORK.md)。
