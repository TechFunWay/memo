# 飞牛 NAS 安装与升级

适用于飞牛 NAS（fnOS）系统，通过 FPK 应用包安装。

## 安装

### 方式一：应用中心手动安装

1. 下载对应架构的 FPK 包：
   - `techfunway-memo-v{{VERSION}}-fnos-amd64.fpk` — x86_64 架构
   - `techfunway-memo-v{{VERSION}}-fnos-arm64.fpk` — ARM64 架构

2. 打开飞牛 NAS 桌面 → **应用中心**

3. 点击右上角 **手动安装** 按钮

4. 在弹出的文件选择窗口中选择下载的 `.fpk` 文件

5. 阅读应用信息，点击 **安装**

6. 在配置向导中设置：
   - **应用端口**：默认 `8905`，如需修改请填入其他端口号

7. 点击 **完成**，等待安装完成

### 方式二：命令行安装

```bash
# 上传 FPK 包到 NAS，然后执行
fnapp install techfunway-memo-v{{VERSION}}-fnos-amd64.fpk
```

## 首次使用

1. 安装完成后，在应用中心找到 **Memo 备忘录**，点击 **打开**

2. 或在浏览器访问：`http://<你的NAS IP>:8905`

3. 注册第一个账号，该账号自动成为超级管理员

## 配置说明

### 应用端口

安装时可通过向导配置应用端口，默认 `8905`。如需更改，在安装向导的"应用端口"输入框中填写新端口号。

## 数据存储

飞牛 NAS 上的数据默认存储在应用数据目录：

```
/app/data/db/database.db     # SQLite 数据库
/app/data/upload/            # 上传的图片文件
/app/data/logs/              # 运行日志
```

> 飞牛 NAS 的应用数据目录通常位于 `/vol1/@appdata/techfunway-memo/` 或通过 `fnapp` 管理的对应路径。

## 升级

### 方式一：应用中心升级

当有新版本发布时：

1. 下载新版本的 `.fpk` 包
2. 打开 **应用中心**
3. 点击 **手动安装**
4. 选择新的 `.fpk` 文件
5. 系统会自动识别为升级，保留原有数据

### 方式二：命令行升级

```bash
fnapp upgrade techfunway-memo-v{{VERSION}}-fnos-amd64.fpk
```

> 升级后数据库会自动迁移到新版本，数据不会丢失。建议在升级前手动备份数据目录。

## 卸载

### 应用中心卸载

1. 打开 **应用中心**
2. 找到 **Memo 备忘录**
3. 点击 **卸载**
4. 确认卸载

### 命令行卸载

```bash
fnapp uninstall techfunway-memo
```

> 卸载会删除应用数据，请提前备份重要笔记。

## 常见问题

### 如何重置管理员密码？

通过 SSH 登录 NAS，进入应用目录，执行：

```bash
# 暂不支持 CLI 管理工具
# 可通过数据库直接操作
```

### 如何备份数据？

备份数据目录（含 `database.db`、`upload/`、`logs/`）即可：

```bash
cp -r /vol1/@appdata/techfunway-memo /backup/memo-backup-$(date +%Y%m%d)
```

### 如何更换端口？

如果安装时未指定正确端口，需要卸载后重新安装并在配置向导中设置新端口。
