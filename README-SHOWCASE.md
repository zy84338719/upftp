# UPFTP Project Showcase

## 🚀 服务运行状态

```
╔════════════════════════════════════════════════════════════════════════════╗
║                                                                            ║
║   ██╗   ██╗██╗   ██╗████████╗███████╗██████╗                                ║
║   ██║   ██║██║   ██║╚══██╔══╝██╔════╝██╔══██╗                               ║
║   ██║   ██║██║   ██║   ██║   █████╗  ██████╔╝                                ║
║   ╚██╗ ██╔╝██║   ██║   ██║   ██╔══╝  ██╔══██╗                                ║
║    ╚████╔╝ ╚█████╔╝   ██║   ███████╗██║  ██║                                ║
║     ╚═══╝   ╚════╝    ╚═╝   ╚══════╝╚═╝  ╚═╝                                ║
║                                                                            ║
║   AI-First Lightweight File Sharing Server                                 ║
║   Version: v0.1.3                                                          ║
║                                                                            ║
╚════════════════════════════════════════════════════════════════════════════╝
```

## 📊 多协议支持

```
┌────────────────────────────────────────────────────────────────────────────┐
│ Server Status                                                               │
├────────────────────────────────────────────────────────────────────────────┤
│ ✅ HTTP Server      http://10.10.30.99:10002                                │
│ ✅ FTP Server       ftp://10.10.30.99:2123                                  │
│ ✅ WebDAV Server    http://10.10.30.99:8082                                 │
│ ✅ NFS Server       nfs://10.10.30.99:2051                                  │
│                                                                            │
│ 📁 Shared Directory: /tmp/upftp-share                                      │
│ 📤 File Upload: enabled (max size: 100 MB)                                 │
│ 👤 FTP Credentials: zhangyi / ****                                          │
└────────────────────────────────────────────────────────────────────────────┘
```

## 🖥️ Web 界面展示

```
┌────────────────────────────────────────────────────────────────────────────┐
│ UPFTP - AI-First File Server                               v0.1.3  [🌐 EN] │
├────────────────────────────────────────────────────────────────────────────┤
│ [📁 Files] [📤 Upload] [⚙️ Settings] [ℹ️ About]                            │
├────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  📁 documents/                                          📅 2026-03-31      │
│     📄 sample.txt                                       337 bytes         │
│     📄 data.json                                        215 bytes         │
│     📄 config.yaml                                      189 bytes         │
│                                                                             │
│  📁 images/                                             📅 2026-03-31      │
│     🖼️ logo.png                                        2.3 MB            │
│                                                                             │
│  📁 projects/                                           📅 2026-03-31      │
│     📦 source.tar.gz                                    1.5 MB            │
│                                                                             │
│  📁 videos/                                             📅 2026-03-31      │
│     🎬 demo.mp4                                         15.2 MB           │
│                                                                             │
│  📄 README.md                                           337 bytes         │
│                                                                             │
├────────────────────────────────────────────────────────────────────────────┤
│ 📊 Total: 7 items (19.4 MB)              [📤 Upload Files] [📁 New Folder] │
└────────────────────────────────────────────────────────────────────────────┘
```

## 💻 命令行界面

```bash
$ upftp -auto

Available network interfaces:
[0] 10.10.30.99
[1] 192.168.31.185

Starting UPFTP v0.1.3...
✓ HTTP server: http://10.10.30.99:10002
✓ FTP server: ftp://10.10.30.99:2123
✓ WebDAV server: http://10.10.30.99:8082
✓ NFS server: nfs://10.10.30.99:2051
✓ File upload: enabled (max size: 100 MB)

Server is running. Press Ctrl+C to stop.
```

## 📦 快速安装

### macOS (Homebrew)
```bash
brew tap zy84338719/tap
brew install upftp
```

### Windows (Scoop)
```powershell
scoop bucket add upftp https://github.com/zy84338719/scoop-bucket
scoop install upftp
```

### Linux
```bash
wget https://github.com/zy84338719/upftp/releases/download/v0.1.3/upftp_Linux_x86_64.tar.gz
tar -xzf upftp_Linux_x86_64.tar.gz
sudo mv upftp /usr/local/bin/
```

## 🎯 主要特性

- **🌐 现代化 Web 界面**: 响应式设计，支持移动设备
- **🚀 多协议支持**: HTTP, FTP, WebDAV, NFS, MCP
- **🎥 丰富的文件预览**: 图片、视频、音频、代码高亮
- **🌍 多语言支持**: 中英文自动切换
- **📤 文件上传**: 支持拖拽上传（最大 100MB）
- **🔒 安全认证**: 支持 HTTP 认证
- **🤖 MCP 集成**: 让 AI 助手直接访问文件

## 📸 创建真实截图

如果你想创建真实的 Web 界面截图：

```bash
# 1. 启动服务
upftp -auto -d /path/to/share

# 2. 在浏览器中打开
open http://localhost:10002

# 3. 使用以下工具截图：
# - macOS: Cmd+Shift+4 然后按空格键，点击浏览器窗口
# - 或使用 peekaboo: peekaboo image --mode window --app Safari
```

## 🔗 项目链接

- **GitHub**: https://github.com/zy84338719/upftp
- **Releases**: https://github.com/zy84338719/upftp/releases
- **Documentation**: https://github.com/zy84338719/upftp#readme
- **Issues**: https://github.com/zy84338719/upftp/issues

---

Made with ❤️ by [zy84338719](https://github.com/zy84338719)
