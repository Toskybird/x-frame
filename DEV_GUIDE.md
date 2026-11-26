# X-Frame 开发指南

本仓库是个人基础框架 **X-Frame**（由 Sub2API 改造）。架构分层与裁剪说明见 [FRAMEWORK.md](FRAMEWORK.md)。

| 项 | 值 |
|----|-----|
| Go module | `github.com/x-frame/x-frame` |
| 后端入口 | `backend/cmd/server` |
| 前端 | `frontend/`（pnpm） |
| 配置示例 | `deploy/config.example.yaml` |

## 本地开发

1. 启动 Postgres / Redis（**推荐绿色版**，无需 Docker）：
   ```powershell
   cd deploy\_internal\scripts
   .\Download-PortableInfra.ps1   # 首次 / 有网刷新
   .\Init-PortableInfra.ps1       # 首次初始化数据目录
   .\Start-PortableInfra.ps1
   ```
   说明见 [deploy/_internal/README.md](deploy/_internal/README.md)。若本机 Docker 可用，也可用 `deploy/docker-compose.infra.yml`。
2. 后端：`cd backend && go run ./cmd/server`
3. 前端：`cd frontend && pnpm install && pnpm dev`
4. 改 Ent schema 后：`cd backend && go generate ./ent`
5. 改 Wire 依赖后：`cd backend && go generate ./cmd/server`

## 本地凭据（可按环境修改）

| 项 | 常见本地值 |
|----|------------|
| Postgres | `127.0.0.1:5432`，user/password/dbname = `x-frame` |
| Redis | `127.0.0.1:6379`，无密码 |
| 前端包管理 | **pnpm**（不要用 npm） |

## 约定

- 新业务优先挂在「框架核心」导航下；对照「示例业务」页面做 UI。
- 不要手改 `backend/ent/` 生成代码（除 `schema/`）。
- 产品名、二进制、配置路径统一为 `x-frame` / `X-Frame`。
