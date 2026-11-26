# X-Frame 框架边界

X-Frame 已裁成**可运行的基础脚手架**：保留鉴权、用户、设置、安装向导与核心管理页；示例业务（账号池 / 网关 / 计费等）已从代码中移除。

Go module：`github.com/x-frame/x-frame`（可按自己的仓库改）。

## 分层

```text
┌─────────────────────────────────────────────┐
│  你的业务（在此扩展）                          │
│  新路由 · 新页面 · 新 Ent schema · 新服务     │
├─────────────────────────────────────────────┤
│  框架核心（已保留）                            │
│  HTTP 装配 · 鉴权 · 用户 · 设置 · 安装向导    │
│  布局壳 · i18n · 通用组件 · Postgres/Redis    │
│  仪表盘 · 公告 · 系统设置 / 备份              │
└─────────────────────────────────────────────┘
```

## 框架核心

| 区域 | 路径 / 说明 |
|------|-------------|
| 服务入口 | `backend/cmd/server`（Gin + Wire） |
| 配置 | `backend/internal/config` |
| HTTP / 中间件 | `backend/internal/server` |
| 鉴权与用户 | `routes/auth.go`、`handler/auth*`、`service/auth*`、`service/user*` |
| 设置 / 安装向导 | `service/setting*`、`internal/setup`、`frontend/src/views/setup` |
| 通用包 | `backend/internal/pkg`（logger、response、pagination 等） |
| 前端壳 | `layout/`、`router/`、`stores/`、`i18n/`、`views/auth` |
| 核心管理页 | 仪表盘、用户、公告、系统设置（含备份） |

侧栏仅保留**框架核心**分组。

## 已裁掉的示例业务

以下模块已从前后端删除（历史 SQL 迁移保留，表可空置）：

- 账号 / 渠道 / 分组 / 代理
- API Key / 用量统计
- 订阅 / 兑换码 / 优惠码
- 运维监控看板
- 网关转发（Claude / OpenAI / Gemini 等客户端与路由）

## 在脚手架上加业务

1. 后端：新增 `ent/schema` → `go generate ./ent` → 写 repository / service / handler → 注册 routes → 更新 Wire（`go generate ./cmd/server`）。
2. 前端：加 `api/` 客户端、`views/` 页面、路由与侧栏项。
3. 跑通：`cd backend && go build ./...`；`cd frontend && pnpm exec vue-tsc --noEmit`。

## 命名约定

| 项 | 值 |
|----|-----|
| 产品名 | X-Frame |
| 二进制 / 服务 | `x-frame` |
| 配置目录 | `/etc/x-frame` |
| 安装目录 | `/opt/x-frame` |
| 前端包名 | `x-frame-frontend` |
| 调试环境变量 | `XFRAME_DEBUG_*` |
| localStorage | `xframe_locale` |
| WebSocket 子协议 | `x-frame-admin` |

历史 SQL 迁移里若仍出现 `sub2api` 或已删表示例字符串，属于旧数据兼容，不要为改名去改已执行过的迁移语义。
