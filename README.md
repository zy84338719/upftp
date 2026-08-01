# upftp

> 一个快速、即开即用的 FTP 文件分享工具。单二进制,跨平台(Linux / macOS / Windows),
> 同时提供 **FTP 服务器**、**REST API** 和**极简 Web 界面**。

## 特性

- 🚀 **一行启动** — `upftp` 即可分享当前目录,自动打印访问地址
- 📡 **完整 FTP 协议** — 支持被动/主动模式(PASV/EPSV/PORT/EPRT)、上传下载、断点续传、MLSD
- 🌐 **REST API** — 程序可调用,支持动态创建/停止 FTP 会话、文件上传下载
- 🖥️ **极简 Web 界面** — 浏览器访问,文件列表、拖拽上传、一键下载
- 📱 **TUI + 二维码** — 终端选中文件即显示二维码,手机扫码下载
- 🔒 **沙箱隔离** — 所有路径强制限制在共享根目录内,防目录穿越
- 🔓 **默认匿名** — 即开即用;`-user/-pass` 启用密码认证
- 📦 **零运行时依赖** — 单二进制,核心代码 ~2500 行

## 安装

```bash
# 从源码
go install github.com/zy84338719/upftp@latest

# 或克隆后构建
git clone https://github.com/zy84338719/upftp.git
cd upftp && make build
```

## 快速开始

```bash
# 分享当前目录(匿名、默认端口 2121 / 8080)
upftp
#   📁 upftp 已启动
#   📡 FTP:   ftp://192.168.1.5:2121
#   🌐 Web:   http://192.168.1.5:8080
#   🔓 匿名访问

# 指定目录、端口、密码
upftp -d ~/share -p 2121 -user me -pass secret

# 只读模式(仅下载,禁止上传/删除)
upftp -d ~/share --read-only

# 不启动 TUI(纯服务模式)
upftp -d ~/share --tui=false
```

## 命令行参数

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `-d <dir>` | 共享目录 | 当前目录 |
| `-p <port>` | FTP 端口 | 2121 |
| `-http <port>` | HTTP/Web 端口 | 8080 |
| `-user <name>` | 用户名(留空=匿名) | 空 |
| `-pass <pass>` | 密码 | 空 |
| `-read-only` | 只读模式 | false |
| `-tui` | 启用终端界面 | true(TTY 时) |
| `-log <level>` | 日志级别 | info |
| `-version` | 打印版本 | |

配置文件查找顺序:`./upftp.yaml` → `~/.upftp/config.yaml` → `/etc/upftp/config.yaml`,命令行参数优先级最高。

## 三种使用形态

### 1. 人类使用 — FTP 客户端 / 浏览器 / TUI

```bash
# FTP 客户端(FileZilla、lftp、curl)
curl ftp://192.168.1.5:2121/hello.txt -o hello.txt

# 浏览器:访问 http://192.168.1.5:8080,拖拽上传/点击下载

# TUI:终端里选中文件显示二维码,手机扫码下载
upftp --tui
```

### 2. API — 直接文件操作

```bash
# 列出文件
curl http://localhost:8080/api/files?path=/

# 下载
curl -o x.zip "http://localhost:8080/api/download?path=/x.zip"

# 上传(raw body)
curl -X POST "http://localhost:8080/api/upload?path=/dest.txt" --data-binary @local.txt

# 上传(multipart,浏览器表单格式)
curl -X POST http://localhost:8080/api/upload -F file=@local.txt -F path=/dest.txt

# 新建文件夹 / 删除
curl -X POST "http://localhost:8080/api/mkdir?path=/newdir"
curl -X DELETE "http://localhost:8080/api/files?path=/old.txt"
```

### 3. API — 动态管理 FTP 会话

```bash
# 创建一个新的 FTP 会话(可指定不同目录/端口)
curl -X POST http://localhost:8080/api/sessions \
  -d '{"dir":"/tmp/another","port":2122,"anonymous":true}'
# → {"id":"s2","port":2122,"anonymous":true,...}

# 列出所有会话
curl http://localhost:8080/api/sessions

# 停止某个会话
curl -X DELETE http://localhost:8080/api/sessions/s2
```

## 认证

- **匿名模式**(默认):任何人可访问,适合局域网快速分享
- **密码模式**:`-user admin -pass secret`,FTP 需登录,HTTP 走 Basic Auth(同一份凭据)

## 架构

```
upftp 单二进制
├── FTP 服务器     internal/ftp/      手写协议,无外部依赖
├── HTTP API+Web   internal/http/     net/http + go:embed 单 HTML
├── TUI            internal/cli/      bubbletea
└── 共享核心       internal/core/     会话管理 / 文件操作 / 认证(FTP 与 HTTP 共用)
```

直接依赖仅 3 个:`yaml.v3`(配置)、`bubbletea`(TUI)、`go-qrcode`(二维码)。

## 开发

```bash
make test          # 运行测试
make vet           # 静态检查
make build-all     # 交叉编译全部平台
```

## 许可证

MIT
