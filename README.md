<h1 align="center">
  <br>
  <a href="https://cloudreve.org/" alt="logo" ><img src="https://raw.githubusercontent.com/cloudreve/frontend/master/public/static/img/logo192.png" width="150"/></a>
  <br>
  upftp
  <br>
</h1>

<h4 align="center">支持快速建立简易ftp下载站点.</h4>

<p align="center">
  <a href="https://github.com/zy84338719/upftp">
    <img src="https://github.com/zy84338719/upftp/actions/workflows/build.yml/badge.svg?branch=main"
         alt="travis">
  </a>
  <a href="https://codecov.io/gh/zy84338719/upftp"><img src="https://img.shields.io/codecov/c/github/zy84338719/upftp?style=flat-square"></a>
  <a href="https://goreportcard.com/report/github.com/zy84338719/upftp">
      <img src="https://goreportcard.com/badge/github.com/zy84338719/upftp?style=flat-square">
  </a>
  <a href="https://github.com/zy84338719/upftp/releases">
    <img src="https://img.shields.io/github/v/release/zy84338719/upftp?include_prereleases&style=flat-square">
  </a>
</p>

<p align="center">
  <a href="#scroll-许可证">许可证</a>
</p>

![Screenshot](https://raw.githubusercontent.com/zy84338719/upftp/master/img.png)

## :sparkles: 特性

* :cloud: 支持本机任意目录进行ftp服务启动 用于构建简易ftp站点

## :hammer_and_wrench: 使用方法
```bash
    upftp 
        -p 8888 端口 默认 10000 
        -d /opt 目录 默认 ./ 当前目录
```

## :hammer_and_wrench: 部署

下载适用于您目标机器操作系统、CPU架构的主程序，直接运行即可。

```shell
# 解压程序包
tar -zxvf upftp_VERSION_OS_ARCH.tar.gz

# 赋予执行权限
chmod +x ./upftp

# 启动 upftp
./upftp
```

以上为最简单的部署示例，您可以参考 [文档 - 起步](https://docs.cloudreve.org/) 进行更为完善的部署。

## :gear: 构建

自行构建前需要拥有 `Go >= 1.11`等必要依赖。

#### 克隆代码

```shell
git clone https://github.com/zy84338719/upftp.git
```

go 1.16
```bash
    go install github.com/zy84338719/upftp
```

#### 编译项目

```shell
# 获得当前版本号、Commit
export COMMIT_SHA=$(git rev-parse --short HEAD)
export VERSION=$(git describe --tags)

# 开始编译
go build -a -ldflags " -X \"main.Version=$(VERSION)\" -X \"main.LastCommit=$(COMMIT_SHA)\" " -o upftp
```

你也可以使用项目根目录下的`build.sh`快速开始构建：

```shell
make build
```

## :alembic: 技术栈

* [Go ](https://golang.org/) + [Gin](https://github.com/gin-gonic/gin)

## :scroll: 许可证

GPL V3

---
> GitHub [@zy84338719](https://github.com/zy84338719) &nbsp;&middot;&nbsp;
> Twitter [@murphyyi](https://twitter.com/murphyyi)