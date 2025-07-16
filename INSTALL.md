# UPFTP - 跨平台文件共享服务器

一个轻量级的文件共享服务器，支持HTTP浏览和FTP访问，提供现代化的Web界面和丰富的文件预览功能。

## 🚀 特性

### 核心功能
- **跨平台支持**: Linux、macOS、Windows 三平台原生支持
- **双协议服务**: HTTP Web界面 + FTP服务器
- **文件预览**: 支持图片、视频、音频、文本和代码文件预览
- **文件夹下载**: 自动打包为ZIP文件下载
- **搜索功能**: 实时文件名和类型搜索
- **命令行界面**: 交互式命令行管理界面

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

### 方法1: 下载预编译版本

从 [Releases](https://github.com/zy84338719/upftp/releases) 页面下载适合您系统的版本：

#### Linux (amd64)
```bash
wget https://github.com/zy84338719/upftp/releases/latest/download/upftp_VERSION_linux_amd64.tar.gz
tar -zxvf upftp_VERSION_linux_amd64.tar.gz
chmod +x upftp_linux_amd64
./upftp_linux_amd64
```

#### Linux (arm64)
```bash
wget https://github.com/zy84338719/upftp/releases/latest/download/upftp_VERSION_linux_arm64.tar.gz
tar -zxvf upftp_VERSION_linux_arm64.tar.gz
chmod +x upftp_linux_arm64
./upftp_linux_arm64
```

#### macOS (Intel)
```bash
wget https://github.com/zy84338719/upftp/releases/latest/download/upftp_VERSION_darwin_amd64.tar.gz
tar -zxvf upftp_VERSION_darwin_amd64.tar.gz
chmod +x upftp_darwin_amd64
./upftp_darwin_amd64
```

#### macOS (Apple Silicon)
```bash
wget https://github.com/zy84338719/upftp/releases/latest/download/upftp_VERSION_darwin_arm64.tar.gz
tar -zxvf upftp_VERSION_darwin_arm64.tar.gz
chmod +x upftp_darwin_arm64
./upftp_darwin_arm64
```

#### Windows
下载 `upftp_VERSION_windows_amd64.zip` 或 `upftp_VERSION_windows_386.zip`，解压后运行 `upftp_windows_amd64.exe`

### 方法2: 从源码编译

需要 Go 1.21 或更高版本：

```bash
# 克隆仓库
git clone https://github.com/zy84338719/upftp.git
cd upftp

# 编译当前平台版本
make build

# 或编译所有平台版本
make build-all

# 创建发布包
make package
```

## 🔧 使用方法

### 基本使用

```bash
# 使用默认配置启动 (端口 10000，当前目录)
./upftp

# 指定端口和目录
./upftp -p 8888 -d /path/to/share

# 自动选择网络接口
./upftp -auto

# 启用FTP服务器
./upftp -enable-ftp

# 指定FTP端口和认证信息
./upftp -enable-ftp -ftp 2121 -user admin -pass secretpass
```

### 完整参数列表

```bash
选项：
  -p <port>       HTTP服务器端口 (默认: 10000)
  -ftp <port>     FTP服务器端口 (默认: 2121)  
  -d <dir>        共享目录 (默认: 当前目录)
  -auto           自动选择第一个可用网络接口
  -enable-ftp     启用FTP服务器 (默认: 关闭)
  -user <name>    FTP用户名 (默认: admin)
  -pass <pass>    FTP密码 (默认: admin)
  -h              显示帮助信息
```

### 使用示例

#### 1. 基本文件共享
```bash
# 共享当前目录，使用默认端口
./upftp

# 访问: http://你的IP:10000
```

#### 2. 高级配置
```bash
# 启用FTP，自定义端口和认证
./upftp -p 8080 -enable-ftp -ftp 2121 -user myuser -pass mypass -d /home/user/shared

# HTTP访问: http://你的IP:8080
# FTP访问: ftp://你的IP:2121 (用户名: myuser, 密码: mypass)
```

#### 3. 自动模式（适合脚本使用）
```bash
# 自动选择网络接口，无需手动选择
./upftp -auto -enable-ftp
```

## 🌐 访问方式

### Web界面访问
1. 启动服务器后，打开浏览器访问显示的URL
2. 支持文件浏览、预览、下载
3. 可搜索文件和文件夹
4. 支持键盘快捷键 (Ctrl+F 聚焦搜索，ESC 关闭预览)

### FTP客户端访问
```bash
# 命令行FTP客户端
ftp 你的IP
# 输入用户名和密码

# FileZilla等图形FTP客户端
服务器: 你的IP
端口: 2121 (或自定义端口)
用户名: admin (或自定义)
密码: admin (或自定义)
```

### 命令行下载
```bash
# 使用curl下载
curl -O http://你的IP:10000/download/文件名

# 使用wget下载
wget http://你的IP:10000/download/文件名

# 下载文件夹(ZIP格式)
curl -O http://你的IP:10000/download/文件夹名
```

## 🎛️ 命令行界面

服务器启动后提供交互式命令行界面：

```
Commands:
  [1] Search files        - 搜索文件
  [2] List all files      - 列出所有文件  
  [3] Show download examples - 显示下载示例
  [4] Refresh file list   - 刷新文件列表
  [5] FTP connection info - FTP连接信息 (如果启用)
  [q] Quit server        - 退出服务器
```

## 🔧 开发构建

### 构建命令

```bash
# 显示构建信息
make debugInfo

# 下载依赖
make deps

# 构建当前平台
make build

# 构建所有平台
make build-all

# 构建特定平台
make build-linux     # Linux (amd64 + arm64)
make build-windows   # Windows (amd64 + 386)
make build-darwin    # macOS (amd64 + arm64)

# 创建发布包
make package

# 运行测试
make test

# 清理构建文件
make clean

# 开发模式运行
make dev

# 运行并启用FTP
make run-ftp

# 代码格式化
make fmt

# 代码检查 (需要golangci-lint)
make lint
```

## 🔒 安全说明

1. **生产环境使用**:
   - 建议更改默认FTP用户名和密码
   - 考虑使用防火墙限制访问
   - 注意共享目录的权限设置

2. **网络安全**:
   - HTTP和FTP都是明文协议
   - 在不可信网络中使用时请谨慎
   - 建议在局域网环境使用

## 📝 更新日志

### v2.0.0 (最新)
- ✨ 全新现代化Web界面
- 🎥 支持视频和音频预览
- 🌐 添加FTP服务器支持
- 🔍 实时搜索功能
- 📱 移动端适配
- 🎨 更好的文件类型识别和图标
- ⚡ 提升性能和稳定性
- 🔧 增强的命令行界面
- 📦 跨平台预编译版本

### v1.x
- 基础HTTP文件服务
- 简单的文件浏览和下载
- 基础文件预览

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

[MIT License](LICENSE.txt)

---

> GitHub [@zy84338719](https://github.com/zy84338719) &nbsp;&middot;&nbsp;
> Twitter [@murphyyi](https://twitter.com/murphyyi)
> Index: [murphyyi](https://murphyyi.com)
