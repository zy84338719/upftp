<h1 align="center">
  <br>
  <a href="https://github.com/zy84338719/upftp" alt="logo" ><img src="https://raw.githubusercontent.com/cloudreve/frontend/master/public/static/img/logo192.png" width="150"/></a>
  <br>
  upftp
  <br>
</h1>

<div align="center">
  <h4>一个现代化的跨平台文件共享服务器 | A modern cross-platform file sharing server</h4>
</div>

<p align="center">
  <a href="https://github.com/zy84338719/upftp/actions/workflows/build.yml">
    <img src="https://github.com/zy84338719/upftp/actions/workflows/build.yml/badge.svg?branch=main"
         alt="Build Status">
  </a>
  <a href="https://goreportcard.com/report/github.com/zy84338719/upftp">
    <img src="https://goreportcard.com/badge/github.com/zy84338719/upftp?style=flat-square">
  </a>
  <a href="https://github.com/zy84338719/upftp/releases">
    <img src="https://img.shields.io/github/v/release/zy84338719/upftp?include_prereleases&style=flat-square">
  </a>
  <a href="https://github.com/zy84338719/upftp/blob/main/LICENSE.txt">
    <img src="https://img.shields.io/github/license/zy84338719/upftp?style=flat-square">
  </a>
</p>

[English](#english) | [中文](#中文)

---

# 中文

## ✨ 主要特性

### 🌐 现代化Web界面
- 响应式设计，支持移动设备
- 实时文件搜索功能
- 直观的文件类型图标
- 优雅的预览模态框
- **🌍 多语言支持**: 自动检测浏览器语言，支持中英文切换

### 🎨 智能语言体验
- **自动语言检测**: 根据浏览器语言自动选择中文或英文界面
- **手动语言切换**: 一键切换中英文，实时生效无需刷新
- **语言偏好记忆**: 自动保存用户语言选择，下次访问时自动应用
- **完整界面翻译**: 所有文本元素均支持中英文显示

### 🎥 丰富的文件预览
- **图片**: JPG, PNG, GIF, SVG, WebP 等
- **视频**: MP4, AVI, MOV, WebM, MKV 等
- **音频**: MP3, WAV, FLAC, AAC, OGG 等  
- **文本/代码**: 支持语法高亮的代码预览
- **文档**: PDF, Office文档下载支持

### 🚀 双协议支持
- **HTTP服务器**: 现代Web界面，支持浏览器访问
- **FTP服务器**: 传统FTP协议，支持各种FTP客户端
- 独立端口配置，可单独启用

### 🔧 便捷的管理
- 交互式命令行界面
- 文件搜索和列表管理
- 下载链接生成
- 实时文件系统刷新

### 🌍 跨平台支持
- **Linux**: amd64, arm64
- **Windows**: amd64, 386  
- **macOS**: Intel, Apple Silicon

## 🚀 快速开始

### 一键安装脚本

最简单的安装方式，自动检测系统并选择最佳安装方法：

```bash
curl -fsSL https://install.upftp.dev | bash
```

### 包管理器安装

#### Ubuntu/Debian (APT)

```bash
# 添加UPFTP仓库
curl -fsSL https://apt.upftp.dev/key.gpg | sudo apt-key add -
echo "deb https://apt.upftp.dev stable main" | sudo tee /etc/apt/sources.list.d/upftp.list

# 安装
sudo apt update
sudo apt install upftp

# 启动服务
sudo systemctl start upftp
sudo systemctl enable upftp  # 开机自启
```

#### macOS (Homebrew)

```bash
# 添加tap并安装
brew tap zy84338719/tap
brew install upftp

# 启动服务
brew services start upftp
```

#### CentOS/RHEL/Fedora (RPM)

```bash
# 下载RPM包
wget https://github.com/zy84338719/upftp/releases/latest/download/upftp_linux_amd64.rpm

# 安装
sudo rpm -ivh upftp_linux_amd64.rpm
# 或者使用yum/dnf
sudo yum localinstall upftp_linux_amd64.rpm
```

### 手动下载和安装

从 [Releases页面](https://github.com/zy84338719/upftp/releases) 下载适合您系统的版本：

```bash
# Linux amd64
wget https://github.com/zy84338719/upftp/releases/latest/download/upftp_linux_amd64.tar.gz
tar -zxvf upftp_linux_amd64.tar.gz
chmod +x upftp_linux_amd64

# macOS (Intel)
wget https://github.com/zy84338719/upftp/releases/latest/download/upftp_darwin_amd64.tar.gz  
tar -zxvf upftp_darwin_amd64.tar.gz
chmod +x upftp_darwin_amd64

# macOS (Apple Silicon)
wget https://github.com/zy84338719/upftp/releases/latest/download/upftp_darwin_arm64.tar.gz
tar -zxvf upftp_darwin_arm64.tar.gz  
chmod +x upftp_darwin_arm64
```

### 基本使用

```bash
# 使用默认配置启动
./upftp

# 指定端口和目录
./upftp -p 8888 -d /path/to/share

# 启用FTP服务器
./upftp -enable-ftp -user admin -pass mypassword

# 自动选择网络接口（适合脚本使用）
./upftp -auto
```

### 完整参数

```bash
upftp [选项]

选项：
    -p <port>       HTTP服务器端口 (默认: 10000)
    -ftp <port>     FTP服务器端口 (默认: 2121)
    -d <dir>        共享目录 (默认: 当前目录)
    -auto           自动选择第一个可用网络接口
    -enable-ftp     启用FTP服务器
    -user <name>    FTP用户名 (默认: admin)
    -pass <pass>    FTP密码 (默认: admin)
    -h              显示帮助信息
```

### 访问方式

启动后可通过以下方式访问：

1. **Web浏览器**: `http://你的IP:端口`
2. **FTP客户端**: `ftp://你的IP:FTP端口`  
3. **命令行下载**: `curl -O http://你的IP:端口/download/文件名`

## 🛠️ 从源码构建

需要 Go 1.21 或更高版本：

```bash
# 克隆代码库
git clone https://github.com/zy84338719/upftp.git
cd upftp

# 构建当前平台
make build

# 构建所有平台
make build-all

# 创建发布包
make package

# 查看所有构建选项
make help
```

## 📖 详细文档

查看 [INSTALL.md](INSTALL.md) 获取完整的安装和使用说明。

---

# English

## ✨ Key Features

### 🌐 Modern Web Interface
- Responsive design with mobile support
- Real-time file search functionality  
- Intuitive file type icons
- Elegant preview modal dialogs
- **🌍 Multi-language Support**: Auto-detect browser language, supports Chinese/English switching

### 🎨 Smart Language Experience
- **Auto Language Detection**: Automatically selects Chinese or English based on browser language
- **Manual Language Switching**: One-click switch between Chinese/English, takes effect immediately without refresh
- **Language Preference Memory**: Automatically saves user language choice, applies on next visit
- **Complete Interface Translation**: All text elements support Chinese/English display

### 🎥 Rich File Preview
- **Images**: JPG, PNG, GIF, SVG, WebP, etc.
- **Videos**: MP4, AVI, MOV, WebM, MKV, etc.
- **Audio**: MP3, WAV, FLAC, AAC, OGG, etc.
- **Text/Code**: Syntax-highlighted code preview
- **Documents**: PDF, Office document download support

### 🚀 Dual Protocol Support  
- **HTTP Server**: Modern web interface for browser access
- **FTP Server**: Traditional FTP protocol for FTP clients
- Independent port configuration, can be enabled separately

### 🔧 Convenient Management
- Interactive command-line interface
- File search and listing management
- Download link generation
- Real-time file system refresh

### 🌍 Cross-Platform Support
- **Linux**: amd64, arm64
- **Windows**: amd64, 386
- **macOS**: Intel, Apple Silicon

## 🚀 Quick Start

### One-line Install Script

The easiest way to install, automatically detects your system and chooses the best installation method:

```bash
curl -fsSL https://install.upftp.dev | bash
```

### Package Manager Installation

#### Ubuntu/Debian (APT)

```bash
# Add UPFTP repository
curl -fsSL https://apt.upftp.dev/key.gpg | sudo apt-key add -
echo "deb https://apt.upftp.dev stable main" | sudo tee /etc/apt/sources.list.d/upftp.list

# Install
sudo apt update
sudo apt install upftp

# Start service
sudo systemctl start upftp
sudo systemctl enable upftp  # Enable on boot
```

#### macOS (Homebrew)

```bash
# Add tap and install
brew tap zy84338719/tap
brew install upftp

# Start service
brew services start upftp
```

#### CentOS/RHEL/Fedora (RPM)

```bash
# Download RPM package
wget https://github.com/zy84338719/upftp/releases/latest/download/upftp_linux_amd64.rpm

# Install
sudo rpm -ivh upftp_linux_amd64.rpm
# Or use yum/dnf
sudo yum localinstall upftp_linux_amd64.rpm
```

### Manual Download and Install

Download the appropriate version for your system from the [Releases page](https://github.com/zy84338719/upftp/releases):

```bash
# Linux amd64
wget https://github.com/zy84338719/upftp/releases/latest/download/upftp_linux_amd64.tar.gz
tar -zxvf upftp_linux_amd64.tar.gz
chmod +x upftp_linux_amd64

# macOS (Intel)  
wget https://github.com/zy84338719/upftp/releases/latest/download/upftp_darwin_amd64.tar.gz
tar -zxvf upftp_darwin_amd64.tar.gz
chmod +x upftp_darwin_amd64

# macOS (Apple Silicon)
wget https://github.com/zy84338719/upftp/releases/latest/download/upftp_darwin_arm64.tar.gz  
tar -zxvf upftp_darwin_arm64.tar.gz
chmod +x upftp_darwin_arm64
```

### Basic Usage

```bash
# Start with default configuration
./upftp

# Specify port and directory  
./upftp -p 8888 -d /path/to/share

# Enable FTP server
./upftp -enable-ftp -user admin -pass mypassword

# Auto-select network interface (suitable for scripts)
./upftp -auto
```

### Full Options

```bash
upftp [options]

Options:
    -p <port>       HTTP server port (default: 10000)
    -ftp <port>     FTP server port (default: 2121)  
    -d <dir>        Share directory (default: current directory)
    -auto           Automatically select first available network interface
    -enable-ftp     Enable FTP server
    -user <name>    FTP username (default: admin)
    -pass <pass>    FTP password (default: admin)
    -h              Show help message
```

### Access Methods

After startup, you can access via:

1. **Web Browser**: `http://your-ip:port`
2. **FTP Client**: `ftp://your-ip:ftp-port`
3. **Command Line**: `curl -O http://your-ip:port/download/filename`

## 🛠️ Build from Source

Requires Go 1.21 or higher:

```bash
# Clone repository
git clone https://github.com/zy84338719/upftp.git
cd upftp

# Build for current platform
make build

# Build for all platforms  
make build-all

# Create release packages
make package

# View all build options
make help
```

## 📖 Documentation

See [INSTALL.md](INSTALL.md) for complete installation and usage instructions.

## 📝 License

[MIT](LICENSE.txt)

---

> GitHub [@zy84338719](https://github.com/zy84338719) &nbsp;&middot;&nbsp;
> Twitter [@murphyyi](https://twitter.com/murphyyi) &nbsp;&middot;&nbsp;
> Website [murphyyi.com](https://murphyyi.com)
