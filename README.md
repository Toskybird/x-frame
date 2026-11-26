# X-Frame

个人基础框架（Go + Vue），从 Sub2API 改造而来。

- **框架核心**：鉴权、用户、系统设置、安装向导、管理端布局壳
- **示例业务**：原 AI 网关相关页面（账号、渠道、计费、用量等）完整保留，侧栏标注「示例」，供后续页面设计对照

边界与裁剪方式见 [FRAMEWORK.md](FRAMEWORK.md)。本地开发见 [DEV_GUIDE.md](DEV_GUIDE.md)。

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
FRAMEWORK.md 核心 vs 示例边界
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

## 设计新页面时

1. 对照侧栏「示例业务」里相近的列表/表单页（如 `AccountsView`、`UsersView`）。
2. 在框架核心路由下新增自己的 `views` + `api` + 后端 handler/service。
3. 稳定后按 [FRAMEWORK.md](FRAMEWORK.md) 裁剪不再需要的示例模块。
