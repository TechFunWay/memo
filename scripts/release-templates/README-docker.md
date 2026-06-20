# Docker 部署方式

适用于 Docker 或 Docker Compose 环境。

## 镜像信息

- 镜像名：`techfunways/memo`
- 支持架构：`linux/amd64`、`linux/arm64`

## 快速启动

```bash
docker run -d \
  --name memo \
  -p 8905:8905 \
  -v ./data:/app/data \
  techfunways/memo:latest
```

启动后浏览器打开 http://localhost:8905，第一个注册的用户自动成为超级管理员。

## Docker 启动参数详解

```bash
docker run -d \
  --name memo \
  -p 8905:8905 \
  -v ./data:/app/data \
  -e TZ=Asia/Shanghai \
  --restart unless-stopped \
  techfunways/memo:latest
```

### 端口映射 `-p`

```bash
# 将宿主机的 8080 端口映射到容器的 8905 端口
-p 8080:8905
```

### 数据持久化 `-v`

```bash
# 将宿主机目录挂载到容器内的 /app/data
-v /path/on/host/data:/app/data
```

容器内的数据目录结构：
```
/app/data/
├── db/
│   └── database.db   # SQLite 数据库
├── upload/           # 上传文件
└── logs/             # 运行日志
```

### 环境变量 `-e`

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `PORT` | 8905 | 服务监听端口（容器内部，一般无需修改） |
| `DATA_DIR` | /app/data | 数据存储目录 |
| `STATIC_DIR` | /app/www | 前端静态文件目录 |
| `TZ` | UTC | 时区，建议设为 `Asia/Shanghai` |

### 重启策略 `--restart`

```bash
--restart always           # 总是重启
--restart unless-stopped   # 除非手动停止，否则自动重启（推荐）
--restart on-failure       # 仅在退出码非零时重启
```

### 资源限制（可选）

```bash
--memory=256m          # 内存上限
--cpus=1               # CPU 核心数限制
```

## Docker Compose 部署

创建 `docker-compose.yaml`：

```yaml
version: '3.8'
services:
  memo:
    image: techfunways/memo:latest
    container_name: memo
    ports:
      - "8905:8905"
    volumes:
      - ./data:/app/data
    environment:
      - PORT=8905
      - TZ=Asia/Shanghai
    restart: unless-stopped
```

启动：

```bash
docker compose up -d
```

查看日志：

```bash
docker compose logs -f
```

停止：

```bash
docker compose down
```

## 升级方法

### Docker run 方式

```bash
# 拉取最新镜像
docker pull techfunways/memo:latest

# 停止并删除旧容器
docker stop memo
docker rm memo

# 用同样的参数启动新容器（数据卷不变，数据不会丢失）
docker run -d \
  --name memo \
  -p 8905:8905 \
  -v ./data:/app/data \
  techfunways/memo:latest
```

### Docker Compose 方式

```bash
docker compose pull
docker compose up -d
```

> 数据库会自动迁移到新版本，无需手动操作。升级前建议备份数据目录。
