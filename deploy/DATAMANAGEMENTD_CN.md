# datamanagementd 部署说明（数据管理）

本文说明如何在宿主机部署 `datamanagementd`，并与主进程联动开启“数据管理”功能。

## 1. 关键约束

- 主进程固定探测路径：`/tmp/x-frame-datamanagement.sock`
- 仅当该 Unix Socket 可连通且 `Health` 成功时，后台“数据管理”才会启用
- `datamanagementd` 使用 SQLite 持久化元数据，不依赖主库

## 2. 宿主机构建与运行

```bash
cd /opt/x-frame-src/datamanagement
go build -o /opt/x-frame/datamanagementd ./cmd/datamanagementd

mkdir -p /var/lib/x-frame/datamanagement
chown -R x-frame:x-frame /var/lib/x-frame/datamanagement
```

手动启动示例：

```bash
/opt/x-frame/datamanagementd \
  -socket-path /tmp/x-frame-datamanagement.sock \
  -sqlite-path /var/lib/x-frame/datamanagement/datamanagementd.db \
  -version 1.0.0
```

## 3. systemd 托管（推荐）

仓库已提供示例服务文件：`deploy/x-frame-datamanagementd.service`

```bash
sudo cp deploy/x-frame-datamanagementd.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable --now x-frame-datamanagementd
sudo systemctl status x-frame-datamanagementd
```

查看日志：

```bash
sudo journalctl -u x-frame-datamanagementd -f
```

也可以使用一键安装脚本（自动安装二进制 + 注册 systemd）：

```bash
# 方式一：使用现成二进制
sudo ./deploy/install-datamanagementd.sh --binary /path/to/datamanagementd

# 方式二：从源码构建后安装
sudo ./deploy/install-datamanagementd.sh --source /path/to/x-frame
```

## 4. Docker 部署联动

若 `x-frame` 运行在 Docker 容器中，需要将宿主机 Socket 挂载到容器同路径：

```yaml
services:
  x-frame:
    volumes:
      - /tmp/x-frame-datamanagement.sock:/tmp/x-frame-datamanagement.sock
```

建议在 `docker-compose.override.yml` 中维护该挂载，避免覆盖主 compose 文件。

## 5. 依赖检查

`datamanagementd` 执行备份时依赖以下工具：

- `pg_dump`
- `redis-cli`
- `docker`（仅 `source_mode=docker_exec` 时）

缺失依赖会导致对应任务失败，并在任务详情中体现错误信息。
