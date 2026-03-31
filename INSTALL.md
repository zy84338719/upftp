# UPFTP - 跨平台文件共享服务器

一个轻量级的文件共享服务器，支持HTTP浏览和FTP访问，提供现代化的Web界面和丰富的文件预览功能。

## 🚀 特性

### 核心功能
- **跨平台支持**: Linux、macOS、Windows 三平台原生支持
- **多协议服务**: HTTP Web界面 + FTP服务器 + WebDAV服务器 + NFS服务器
- **文件预览**: 支持图片、视频、音频、文本和代码文件预览
- **文件夹下载**: 自动打包为ZIP文件下载
- **搜索功能**: 实时文件名和类型搜索
- **命令行界面**: 交互式命令行管理界面，带有 ASCII 艺术标志横幅
- **AI集成**: 支持 Model Context Protocol (MCP)，可与 AI 助手集成

### 支持的文件类型

#### 🖼️ 图片预览
- JPG, JPEG, PNG, GIF, BMP, WebP, SVG, ICO

#### 🎥 视频预览
- MP4, AVI, MOV, WMV, FLV, WebM, MKV, M4V

#### 🎵 音频预览
- MP3, WAV, FLAC, AAC, OGG, WMA, M4A

#### 📝 文本/代码预览
- TXT, MD, JSON, XML, YAML, CSV, LOG
- Go, JavaScript, TypeScript, HTML, CSS, Python
- Java, C++, C, PHP, Ruby, Rust, Shell, SQL

#### 📄 文档支持
- PDF, DOC, DOCX, XLS, XLSX, PPT, PPTX
- *注意: Office文档提供下载，不提供在线预览*

## 📦 安装方法

### 方法1: 一键安装（推荐）

#### macOS / Linux
```bash
curl -fsSL https://zy84338719.github.io/upftp/install.sh | bash
```

#### Windows (PowerShell)
```powershell
# 使用 Scoop
scoop bucket add zy84338719 https://github.com/zy84338719/scoop-bucket
scoop install upftp
```

### 方法2: 包管理器

#### macOS (Homebrew)
```bash
brew tap zy84338719/tap
brew install upftp
```

#### Ubuntu / Debian (APT)
```bash
# 添加 APT 仓库
echo "deb [trusted=yes] https://zy84338719.github.io/upftp/apt stable main" | sudo tee /etc/apt/sources.list.d/upftp.list

# 安装
sudo apt update
sudo apt install upftp

# 启动服务
sudo systemctl start upftp
sudo systemctl enable upftp
```

#### CentOS / RHEL / Fedora (RPM)
```bash
# 下载 RPM 包
sudo yum install -y https://github.com/zy84338719/upftp/releases/latest/download/upftp-1.0.0-1.x86_64.rpm

# 或使用 dnf
sudo dnf install -y https://github.com/zy84338719/upftp/releases/latest/download/upftp-1.0.0-1.x86_64.rpm
```

### 方法3: 下载预编译版本

从 [Releases](https://github.com/zy84338719/upftp/releases) 页面下载适合您系统的版本：

#### Linux (amd64)
```bash
wget https://github.com/zy84338719/upftp/releases/latest/download/upftp_linux_x86_64.tar.gz
tar -xzf upftp_linux_x86_64.tar.gz
chmod +x upftp
sudo mv upftp /usr/local/bin/
```

#### Linux (arm64)
```bash
wget https://github.com/zy84338719/upftp/releases/latest/download/upftp_linux_arm64.tar.gz
tar -xzf upftp_linux_arm64.tar.gz
chmod +x upftp
sudo mv upftp /usr/local/bin/
```

#### macOS (Intel / Apple Silicon)
```bash
# 使用 Homebrew (推荐)
brew tap zy84338719/tap
brew install upftp

# 或手动下载
# Intel
wget https://github.com/zy84338719/upftp/releases/latest/download/upftp_darwin_x86_64.tar.gz
# Apple Silicon
wget https://github.com/zy84338719/upftp/releases/latest/download/upftp_darwin_arm64.tar.gz

tar -xzf upftp_darwin_*.tar.gz
chmod +x upftp
sudo mv upftp /usr/local/bin/
```

#### Windows
下载 `upftp_windows_x86_64.zip`，解压后将 `upftp.exe` 放入 PATH 目录。

### 方法4: 从源码编译

需要 Go 1.21 或更高版本：

```bash
git clone https://github.com/zy84338719/upftp.git
cd upftp
make build
sudo make install
```

## 🔧 使用方法

### 基本使用

```bash
# 使用默认配置启动 (端口 10000，当前目录)
upftp

# 指定端口和目录
upftp -p 8888 -d /path/to/share

# 自动选择网络接口
upftp -auto

# 启用FTP服务器
upftp -enable-ftp

# 启用WebDAV服务器
upftp -enable-webdav

# 启用NFS服务器
upftp -enable-nfs

# 启用MCP服务器 (AI集成)
upftp -enable-mcp

# 完整配置
upftp -p 8080 -enable-ftp -ftp 2121 -enable-webdav -webdav 8081 -enable-nfs -nfs 2049 -user admin -pass secret -d /home/shared
```

### 完整参数列表

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `-p <port>` | HTTP服务器端口 | 10000 |
| `-ftp <port>` | FTP服务器端口 | 2121 |
| `-webdav <port>` | WebDAV服务器端口 | 8080 |
| `-nfs <port>` | NFS服务器端口 | 2049 |
| `-d <dir>` | 共享目录 | 当前目录 |
| `-auto` | 自动选择网络接口 | false |
| `-enable-ftp` | 启用FTP服务器 | false |
| `-enable-webdav` | 启用WebDAV服务器 | false |
| `-enable-nfs` | 启用NFS服务器 | false |
| `-enable-mcp` | 启用MCP服务器 (AI集成) | false |
| `-user <name>` | FTP用户名 | admin |
| `-pass <pass>` | FTP密码 | admin |
| `-h` | 显示帮助信息 | - |

## 🌐 访问方式

### Web界面
启动后访问显示的 URL（如 `http://192.168.1.100:10000`）

### FTP客户端
```bash
ftp 192.168.1.100 2121
# 用户名: admin
# 密码: admin
```

### WebDAV客户端
```bash
# 使用 curl 测试
curl http://192.168.1.100:8080

# 在文件管理器中访问
# Windows: 映射网络驱动器 -> http://192.168.1.100:8080
# macOS: 前往 -> 连接服务器 -> http://192.168.1.100:8080
# Linux: 挂载 WebDAV
mount -t davfs http://192.168.1.100:8080 /mnt/webdav
```

### NFS客户端
```bash
# Linux/macOS
mount -t nfs 192.168.1.100:/share /mnt/nfs

# Windows
# 映射网络驱动器 -> \\192.168.1.100\share
```

### 命令行下载
```bash
# 下载文件
curl -O http://192.168.1.100:10000/download/filename
wget http://192.168.1.100:10000/download/filename

# 下载并重命名
curl -o newname.txt http://192.168.1.100:10000/download/filename.txt
wget -O newname.txt http://192.168.1.100:10000/download/filename.txt

# 显示进度
curl -# -O http://192.168.1.100:10000/download/largefile.zip
wget --progress=bar http://192.168.1.100:10000/download/largefile.zip

# 下载文件夹（自动打包为ZIP）
curl -O http://192.168.1.100:10000/download/foldername
wget http://192.168.1.100:10000/download/foldername
```

## 🎛️ 系统服务 (Linux)

使用包管理器安装后，可以配置为系统服务：

```bash
# 启动服务
sudo systemctl start upftp

# 开机自启
sudo systemctl enable upftp

# 查看状态
sudo systemctl status upftp

# 查看日志
sudo journalctl -u upftp -f
```

服务配置文件位于 `/etc/systemd/system/upftp.service`

## 🗑️ 卸载

### macOS (Homebrew)
```bash
brew uninstall upftp
brew untap zy84338719/tap
```

### Ubuntu / Debian
```bash
sudo apt remove upftp
sudo rm /etc/apt/sources.list.d/upftp.list
```

### 通用卸载脚本
```bash
curl -fsSL https://zy84338719.github.io/upftp/uninstall.sh | bash
```

## 🔧 开发构建

```bash
# 构建当前平台
make build

# 构建所有平台
make build-all

# 运行测试
make test

# 创建本地发布（测试）
make release-snapshot

# 验证配置
make validate
```

## 🔒 安全说明

1. 建议更改默认FTP密码
2. HTTP和FTP都是明文协议，建议在可信网络使用
3. 注意共享目录的权限设置

## 📄 许可证

[MIT License](LICENSE.txt)

---

> GitHub [@zy84338719](https://github.com/zy84338719)
