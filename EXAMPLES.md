# UPFTP 使用示例

这里提供一些常见的使用场景和配置示例。

## 基础使用场景

### 1. 临时文件分享
在需要快速分享文件时：
```bash
# 在包含文件的目录中启动
cd /path/to/files
./upftp -auto

# 访问 http://your-ip:10000
```

### 2. 家庭媒体服务器
分享家庭照片和视频：
```bash
# 分享媒体文件夹，使用自定义端口
./upftp -p 8080 -d /home/user/Media -auto

# 支持在浏览器中预览图片和视频
```

### 3. 开发文件传输
在开发过程中传输文件：
```bash
# 启用FTP便于IDE或工具集成
./upftp -enable-ftp -p 8080 -ftp 2121 -d /project/dist
```

### 4. 局域网文件共享
在局域网中共享文件：
```bash
# 使用默认配置，手动选择网络接口
./upftp -d /shared/folder

# 选择正确的局域网IP地址
```

## 高级配置示例

### 1. 完整配置的媒体服务器
```bash
#!/bin/bash
# media-server.sh

./upftp \
  -p 8080 \
  -enable-ftp \
  -ftp 2121 \
  -user mediauser \
  -pass SecurePass123 \
  -d /home/user/MediaLibrary \
  -auto

echo "媒体服务器已启动"
echo "HTTP访问: http://$(hostname -I | awk '{print $1}'):8080"
echo "FTP访问: ftp://$(hostname -I | awk '{print $1}'):2121"
echo "FTP用户: mediauser / SecurePass123"
```

### 2. 安全的文件传输服务
```bash
#!/bin/bash
# secure-transfer.sh

# 生成随机密码
PASSWORD=$(openssl rand -base64 12)

./upftp \
  -p 9999 \
  -enable-ftp \
  -ftp 2122 \
  -user transfer \
  -pass "$PASSWORD" \
  -d /tmp/transfer \
  -auto

echo "安全传输服务已启动"
echo "临时密码: $PASSWORD"
echo "请将此密码安全地发送给接收方"
```

### 3. 多目录服务配置
```bash
#!/bin/bash
# multi-service.sh

# 启动文档服务器
./upftp -p 8081 -d /docs -auto &
DOC_PID=$!

# 启动媒体服务器
./upftp -p 8082 -d /media -auto &
MEDIA_PID=$!

# 启动下载服务器
./upftp -p 8083 -enable-ftp -d /downloads -auto &
DOWNLOAD_PID=$!

echo "多服务已启动:"
echo "文档服务: http://your-ip:8081"
echo "媒体服务: http://your-ip:8082" 
echo "下载服务: http://your-ip:8083 (含FTP)"

# 优雅退出
trap "kill $DOC_PID $MEDIA_PID $DOWNLOAD_PID" EXIT
wait
```

## 客户端访问示例

### Web浏览器访问
```
打开浏览器访问: http://server-ip:port
- 支持文件预览
- 支持搜索过滤
- 支持批量下载（ZIP）
- 移动端适配
```

### FTP客户端访问

#### 命令行FTP
```bash
# 连接到FTP服务器
ftp server-ip port

# 登录
Username: your-username
Password: your-password

# 常用命令
ls                  # 列出文件
get filename        # 下载文件
put filename        # 上传文件
mget *.txt         # 批量下载
binary             # 二进制模式
bye                # 退出
```

#### FileZilla配置
```
服务器: server-ip
端口: ftp-port
协议: FTP - 文件传输协议
加密: 只使用普通FTP（不安全）
登录类型: 正常
用户: your-username
密码: your-password
```

### 命令行下载

#### 使用curl
```bash
# 下载单个文件
curl -O http://server-ip:port/download/filename.txt

# 下载并重命名
curl -o newname.txt http://server-ip:port/download/filename.txt

# 显示进度
curl -# -O http://server-ip:port/download/largefile.zip

# 下载文件夹（自动打包为ZIP）
curl -O http://server-ip:port/download/foldername
```

#### 使用wget
```bash
# 下载单个文件
wget http://server-ip:port/download/filename.txt

# 下载并重命名
wget -O newname.txt http://server-ip:port/download/filename.txt

# 显示进度
wget --progress=bar http://server-ip:port/download/largefile.zip

# 递归下载（注意：这会下载整个网站结构）
wget -r -np http://server-ip:port/
```

## 常见问题和解决方案

### 1. 端口被占用
```bash
# 检查端口占用
netstat -tlnp | grep :10000

# 使用其他端口
./upftp -p 10001
```

### 2. 防火墙配置
```bash
# Ubuntu/Debian
sudo ufw allow 10000
sudo ufw allow 2121

# CentOS/RHEL
sudo firewall-cmd --permanent --add-port=10000/tcp
sudo firewall-cmd --permanent --add-port=2121/tcp
sudo firewall-cmd --reload
```

### 3. 权限问题
```bash
# 确保共享目录有正确权限
chmod 755 /path/to/share

# 确保执行文件有权限
chmod +x upftp
```

### 4. 网络接口选择
```bash
# 查看所有网络接口
ip addr show

# 手动选择正确的网络接口
./upftp  # 不使用 -auto 参数
```

## 性能优化

### 1. 大文件传输
```bash
# 对于大文件传输，建议使用FTP
./upftp -enable-ftp -ftp 2121

# FTP支持断点续传
```

### 2. 高并发访问
```bash
# 增加系统文件句柄限制
ulimit -n 65536

# 使用SSD存储提高IO性能
```

### 3. 网络优化
```bash
# 使用有线网络而非WiFi
# 确保网络带宽充足
# 避免网络拥塞时段
```

## 安全建议

### 1. 网络安全
```bash
# 使用强密码
./upftp -enable-ftp -user admin -pass $(openssl rand -base64 16)

# 限制访问IP（通过防火墙）
sudo ufw allow from 192.168.1.0/24 to any port 10000
```

### 2. 文件安全
```bash
# 使用只读目录
chmod 444 /path/to/readonly/files

# 创建专门的共享用户
sudo useradd -r -s /bin/false upftp-user
sudo chown -R upftp-user:upftp-user /path/to/share
```

### 3. 临时访问
```bash
# 创建临时共享脚本
#!/bin/bash
TEMP_DIR=$(mktemp -d)
cp /path/to/file "$TEMP_DIR"
./upftp -d "$TEMP_DIR" -auto &
SERVER_PID=$!

echo "临时服务器启动，60秒后自动关闭"
sleep 60
kill $SERVER_PID
rm -rf "$TEMP_DIR"
```

这些示例应该能帮助您在各种场景下有效地使用UPFTP。根据具体需求调整配置参数。
