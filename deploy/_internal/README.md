# `_internal` — 本机依赖离线包

本目录存放 **Windows 绿色版 PostgreSQL + Redis**（推荐），以及可选的 Docker/WSL 离线包。大文件默认不进 Git（见 `deploy/.gitignore`），用 U 盘/内网盘拷贝整个 `_internal`。

> 运行时硬依赖 PostgreSQL + Redis；应用仍在本机跑（`go run` / `pnpm dev`）。

## 推荐：绿色版（无需 Docker / WSL）

不解压安装器、不注册 Windows 服务、不改系统 PATH。进程模式启停，仅监听 `127.0.0.1`。

### 目录结构

```
_internal/
  portable/
    packages/                 # 下载的 zip
    pgsql/                    # 解压后的 PostgreSQL
    redis/                    # 解压后的 Redis
    data/postgres|redis/      # 数据目录
    logs/
    .env.portable.example
    .env.portable             # 本地凭据（gitignore）
  scripts/
    Download-PortableInfra.ps1
    Init-PortableInfra.ps1
    Start-PortableInfra.ps1
    Stop-PortableInfra.ps1
    PortableInfra.Common.ps1
  README.md
  manifest-portable.txt
```

### 有外网：下载并初始化

```powershell
cd deploy\_internal\scripts
.\Download-PortableInfra.ps1
.\Init-PortableInfra.ps1
.\Start-PortableInfra.ps1
```

内网机器：把整个 `_internal\portable`（含 `packages` 或已解压的 `pgsql`/`redis`）拷过去后：

```powershell
# 若只有 zip、尚未解压：
.\Download-PortableInfra.ps1 -SkipDownload
.\Init-PortableInfra.ps1
.\Start-PortableInfra.ps1
```

停止：

```powershell
.\Stop-PortableInfra.ps1
```

### 默认连接

| 服务 | 地址 | 凭据 |
|------|------|------|
| PostgreSQL | `127.0.0.1:5432` | user/password/db = `x-frame` / `x-frame` / `x-frame` |
| Redis | `127.0.0.1:6379` | 无密码 |

超级用户（仅维护用）：`postgres` / `postgres`（见 `.env.portable.example`）。

### 版本来源

| 组件 | 包 |
|------|-----|
| PostgreSQL 18 | EDB `postgresql-18.*-windows-x64-binaries.zip` |
| Redis 8 | [redis-windows](https://github.com/redis-windows/redis-windows) `Redis-*-Windows-x64-msys2.zip`（无 Service） |
| VC++ 运行库 | `vc_redist.x64.exe`（Download 脚本会一并拉取；PG 18 `initdb` 依赖较新的 MSVCP140） |

可用参数覆盖：`.\Download-PortableInfra.ps1 -PostgresUrl <url> -RedisVersion 8.10.3`。

若 `initdb` 在 post-bootstrap 阶段崩溃且事件查看器指向 `MSVCP140.dll`：安装 `portable\packages\vc_redist.x64.exe` 后必要时重启，再跑 `Init-PortableInfra.ps1`。

---

## 可选备选：Docker + WSL（本机引擎不可用时可忽略）

若日后要走容器，本目录仍可保存 Docker Desktop / WSL 安装包与镜像 tar。

```
_internal/
  windows/                 # Docker Desktop + WSL MSI
  docker-images/           # docker save 的镜像 tar
  scripts/
    Download-IntranetPackages.ps1
    Install-IntranetPrereqs.ps1
```

有外网准备 Windows 包：

```powershell
.\Download-IntranetPackages.ps1
```

引擎可用后导出/导入镜像并启动，见仓库 `deploy/docker-compose.infra.yml` 与 `deploy/scripts/`。

注意：

- `DockerDesktop-Installer-local-bootstrap.exe`（约几 MB）**不能**替代完整 `DockerDesktopInstaller.exe`。
- Docker Desktop 依赖可用的 WSL；WSL 损坏时请优先使用上方绿色版。
