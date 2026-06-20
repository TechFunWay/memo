# 二进制直接启动方式

适用于 Linux、macOS、Windows 系统，无需安装 Docker 或任何运行时环境。

## 文件结构

```
memo-{{VERSION}}-<平台>/
├── memo             # 主程序（Windows 下为 memo.exe）
├── start.sh         # 启动脚本（Windows 下为 start.bat）
└── www/             # 前端静态资源（由程序自动提供）
```

## 快速启动

### Linux / macOS

```bash
# 解压
tar -xzf techfunway-memo-{{VERSION}}-linux-amd64.tar.gz
cd techfunway-memo-{{VERSION}}-linux-amd64

# 方式一：使用启动脚本
./start.sh

# 方式二：直接运行
./memo
```

### Windows

解压后双击 `start.bat`，或在命令行中运行：

```cmd
memo.exe
```

### 访问

启动后浏览器打开：http://localhost:8905

第一个注册的用户自动成为超级管理员。

## 启动参数详解

```bash
# 指定端口（默认 8905）
./memo -port 8080

# 指定数据目录（默认 ./data）
./memo -data-dir /path/to/data

# 组合使用
./memo -port 8080 -data-dir /mnt/disk1/memo-data
```

## 环境变量配置

所有参数均支持通过环境变量设置，优先级高于命令行参数：

| 环境变量 | 默认值 | 说明 |
|----------|--------|------|
| `PORT` | 8905 | 服务监听端口 |
| `DATA_DIR` | ./data | 数据存储目录（数据库、日志、上传文件均在此目录下） |
| `STATIC_DIR` | ./www | 前端静态文件目录（通常无需设置） |

### 环境变量使用示例

```bash
# Linux / macOS
export PORT=8080
export DATA_DIR=/mnt/disk1/memo-data
./memo

# Windows (CMD)
set PORT=8080
set DATA_DIR=D:\memo-data
memo.exe
```

## 数据目录结构

```
data/
├── db/
│   └── database.db   # SQLite 数据库文件
├── upload/           # 上传的图片文件
└── logs/             # 运行日志
```

> `database.db` 是核心数据文件，请定期备份。升级时只需替换 `memo` 二进制文件，数据库会自动迁移。

## 开机自启（Linux systemd）

创建 `/etc/systemd/system/memo.service`：

```ini
[Unit]
Description=Memo Service
After=network.target

[Service]
Type=simple
User=your-user
WorkingDirectory=/opt/memo
ExecStart=/opt/memo/memo -port 8905 -data-dir /opt/memo/data
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl daemon-reload
sudo systemctl enable memo
sudo systemctl start memo
```

## 升级方法

1. 下载新版本的压缩包
2. 停止当前运行的 memo 进程
3. 备份数据目录（可选，建议操作）
4. 替换 `memo` 二进制文件
5. 重新启动 memo

数据库会自动迁移到新版本，无需手动操作。
