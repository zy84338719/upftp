<h1 align="center">
  <br>
  <a href="https://github.com/zy84338719/upftp" alt="logo" ><img src="https://raw.githubusercontent.com/cloudreve/frontend/master/public/static/img/logo192.png" width="150"/></a>
  <br>
  upftp
  <br>
</h1>

<div align="center">
  <h4>一个轻量级的文件共享服务器 | A lightweight file sharing server</h4>
</div>

<p align="center">
  <a href="https://github.com/zy84338719/upftp">
    <img src="https://github.com/zy84338719/upftp/actions/workflows/build.yml/badge.svg?branch=main"
         alt="Build Status">
  </a>
  <a href="https://codecov.io/gh/zy84338719/upftp">
    <img src="https://img.shields.io/codecov/c/github/zy84338719/upftp?style=flat-square">
  </a>
  <a href="https://goreportcard.com/report/github.com/zy84338719/upftp">
    <img src="https://goreportcard.com/badge/github.com/zy84338719/upftp?style=flat-square">
  </a>
  <a href="https://github.com/zy84338719/upftp/releases">
    <img src="https://img.shields.io/github/v/release/zy84338719/upftp?include_prereleases&style=flat-square">
  </a>
</p>

[English](#english) | [中文](#中文)

# 中文

## :sparkles: 特性

* 快速启动：一键启动文件共享服务
* 简单易用：支持网页浏览和命令行操作
* 文件预览：支持图片和文本文件在线预览
* 文件夹支持：可以浏览和下载整个文件夹
* 多种下载方式：支持浏览器、curl、wget等下载方式
* 无需配置：自动检测网络接口，快速部署

## :hammer_and_wrench: 使用方法

```bash
upftp [options]

选项：
    -p <port>    指定端口号（默认：10000）
    -d <dir>     指定共享目录（默认：当前目录）
    -auto        自动选择第一个可用网络接口
```

### 示例：
```bash
# 使用默认配置启动
upftp

# 指定端口和目录
upftp -p 8888 -d /path/to/share

# 自动选择网络接口
upftp -auto
```

## :rocket: 快速安装

### 使用预编译版本

从 [Releases](https://github.com/zy84338719/upftp/releases) 页面下载适合您系统的版本。

```bash
# 解压
tar -zxvf upftp_VERSION_OS_ARCH.tar.gz

# 添加执行权限
chmod +x ./upftp

# 运行
./upftp
```

### 从源码安装

需要 Go 1.16 或更高版本。

```bash
go install github.com/zy84338719/upftp@latest
```

## :gear: 从源码构建

```bash
# 克隆仓库
git clone https://github.com/zy84338719/upftp.git

# 获取版本信息
export COMMIT_SHA=$(git rev-parse --short HEAD)
export VERSION=$(git describe --tags)

# 编译
go build -a -ldflags " -X \"main.Version=$(VERSION)\" -X \"main.LastCommit=$(COMMIT_SHA)\" " -o upftp
```

---

# English

## :sparkles: Features

* Quick Start: Launch file sharing service with one command
* Easy to Use: Support both web interface and command line operations
* File Preview: Preview images and text files online
* Directory Support: Browse and download entire directories
* Multiple Download Methods: Support browser, curl, wget downloads
* Zero Configuration: Auto-detect network interfaces for quick deployment

## :hammer_and_wrench: Usage

```bash
upftp [options]

Options:
    -p <port>    Specify port number (default: 10000)
    -d <dir>     Specify share directory (default: current directory)
    -auto        Automatically select first available network interface
```

### Examples:
```bash
# Start with default configuration
upftp

# Specify port and directory
upftp -p 8888 -d /path/to/share

# Auto-select network interface
upftp -auto
```

## :rocket: Quick Installation

### Using Pre-built Binaries

Download the appropriate version for your system from the [Releases](https://github.com/zy84338719/upftp/releases) page.

```bash
# Extract
tar -zxvf upftp_VERSION_OS_ARCH.tar.gz

# Add execute permission
chmod +x ./upftp

# Run
./upftp
```

### Install from Source

Requires Go 1.16 or higher.

```bash
go install github.com/zy84338719/upftp@latest
```

## :gear: Build from Source

```bash
# Clone repository
git clone https://github.com/zy84338719/upftp.git

# Get version info
export COMMIT_SHA=$(git rev-parse --short HEAD)
export VERSION=$(git describe --tags)

# Build
go build -a -ldflags " -X \"main.Version=$(VERSION)\" -X \"main.LastCommit=$(COMMIT_SHA)\" " -o upftp
```

## :scroll: License

[MIT](https://github.com/zy84338719/upftp/blob/main/LICENSE.txt)

---
> GitHub [@zy84338719](https://github.com/zy84338719) &nbsp;&middot;&nbsp;
> Twitter [@murphyyi](https://twitter.com/murphyyi)
> index: [murphyyi](https://murphyyi.com)
