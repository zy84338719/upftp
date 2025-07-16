# UPFTP 包管理器安装指南

本文档详细说明如何设置和使用UPFTP的各种包管理器安装方案。

## 📦 支持的包管理器

- **APT** (Debian/Ubuntu)
- **RPM** (CentOS/RHEL/Fedora)  
- **Homebrew** (macOS)
- **一键安装脚本** (跨平台)

## 🔧 APT 仓库设置

### 服务器端设置

1. **准备Web服务器**
```bash
# 安装Nginx或Apache
sudo apt install nginx

# 创建仓库目录
sudo mkdir -p /var/www/apt/{pool/main,dists/stable/main/binary-{amd64,arm64}}
sudo chown -R www-data:www-data /var/www/apt
```

2. **配置Nginx**
```nginx
server {
    listen 80;
    server_name apt.upftp.dev;  # 替换为你的域名
    root /var/www/apt;
    
    location / {
        autoindex on;
        autoindex_exact_size off;
        autoindex_localtime on;
    }
    
    location ~ \.deb$ {
        add_header Content-Type application/vnd.debian.binary-package;
    }
}
```

3. **发布包**
```bash
# 使用GoReleaser自动发布
git tag v1.2.0
git push origin v1.2.0

# 或手动发布
make build-packages
make publish-apt
```

### 客户端安装

```bash
# 添加仓库
curl -fsSL https://apt.upftp.dev/key.gpg | sudo apt-key add -
echo "deb https://apt.upftp.dev stable main" | sudo tee /etc/apt/sources.list.d/upftp.list

# 安装
sudo apt update
sudo apt install upftp

# 管理服务
sudo systemctl start upftp
sudo systemctl enable upftp
sudo systemctl status upftp
```

## 🍺 Homebrew Tap 设置

### 创建Homebrew Tap

1. **创建tap仓库**
```bash
# 在GitHub上创建仓库: homebrew-tap
git clone https://github.com/zy84338719/homebrew-tap.git
cd homebrew-tap
mkdir Formula
```

2. **生成formula**
```bash
# 使用脚本生成
make generate-brew-formula

# 复制到tap仓库
cp upftp.rb ../homebrew-tap/Formula/
```

3. **发布到tap**
```bash
cd homebrew-tap
git add Formula/upftp.rb
git commit -m "Add upftp formula"
git push origin main
```

### 客户端安装

```bash
# 添加tap
brew tap zy84338719/tap

# 安装
brew install upftp

# 管理服务
brew services start upftp
brew services stop upftp
brew services restart upftp
```

## 🚀 一键安装脚本

### 脚本功能

- 自动检测操作系统和架构
- 选择最佳安装方法
- 支持多种安装源
- 错误处理和回滚

### 使用方法

```bash
# 自动安装
curl -fsSL https://install.upftp.dev | bash

# 指定安装方法
curl -fsSL https://install.upftp.dev | bash -s apt
curl -fsSL https://install.upftp.dev | bash -s brew
curl -fsSL https://install.upftp.dev | bash -s download
```

### 自定义安装脚本

可以将安装脚本托管在自己的服务器上：

```bash
# 下载脚本
wget https://raw.githubusercontent.com/zy84338719/upftp/main/scripts/install.sh

# 自定义配置
vim install.sh  # 修改APT_REPO_URL等变量

# 托管在你的服务器
cp install.sh /var/www/html/install.sh
```

## 🔒 GPG 签名设置

### 生成GPG密钥

```bash
# 生成密钥
gpg --full-generate-key

# 导出公钥
gpg --armor --export your-email@example.com > key.gpg

# 上传到服务器
cp key.gpg /var/www/apt/
```

### 签名发布

```bash
# 在发布脚本中设置
export GPG_KEY_ID="your-key-id"
make publish-apt
```

## 📊 仓库监控

### 访问统计

```bash
# Nginx日志分析
sudo tail -f /var/log/nginx/access.log | grep "\.deb"

# 包下载统计
sudo grep "\.deb" /var/log/nginx/access.log | wc -l
```

### 自动化监控

```bash
#!/bin/bash
# 监控脚本

LOG_FILE="/var/log/nginx/access.log"
WEBHOOK_URL="https://hooks.slack.com/your-webhook"

# 检查最近1小时的下载量
DOWNLOADS=$(grep "$(date -d '1 hour ago' '+%d/%b/%Y:%H')" $LOG_FILE | grep "\.deb" | wc -l)

if [ $DOWNLOADS -gt 100 ]; then
    curl -X POST -H 'Content-type: application/json' \
        --data "{\"text\":\"High download activity: $DOWNLOADS downloads in the last hour\"}" \
        $WEBHOOK_URL
fi
```

## 🛠️ 故障排除

### 常见问题

1. **APT更新失败**
```bash
# 检查网络连接
ping apt.upftp.dev

# 检查仓库配置
cat /etc/apt/sources.list.d/upftp.list

# 手动更新
sudo apt update --allow-insecure-repositories
```

2. **Homebrew安装失败**
```bash
# 清理缓存
brew cleanup

# 重新添加tap
brew untap zy84338719/tap
brew tap zy84338719/tap

# 强制重新安装
brew reinstall upftp
```

3. **权限问题**
```bash
# 检查文件权限
ls -la /usr/local/bin/upftp

# 修复权限
sudo chmod +x /usr/local/bin/upftp
```

### 日志调试

```bash
# 检查系统日志
sudo journalctl -u upftp

# 检查服务状态
sudo systemctl status upftp

# 查看安装日志
cat /var/log/dpkg.log | grep upftp
```

## 📈 性能优化

### CDN加速

```bash
# 使用CloudFlare或其他CDN加速下载
# 配置CDN指向你的APT仓库

# 或使用GitHub Releases作为CDN
DOWNLOAD_URL="https://github.com/zy84338719/upftp/releases/latest/download/"
```

### 缓存策略

```nginx
# Nginx缓存配置
location ~ \.(deb|rpm)$ {
    expires 7d;
    add_header Cache-Control "public, immutable";
}
```

## 📋 维护清单

### 定期任务

- [ ] 检查仓库磁盘空间
- [ ] 清理旧版本包
- [ ] 更新GPG密钥
- [ ] 监控下载统计
- [ ] 备份仓库数据

### 发布检查

- [ ] 测试安装脚本
- [ ] 验证包完整性
- [ ] 检查依赖关系
- [ ] 测试服务启动
- [ ] 更新文档

---

更多详细信息请访问 [GitHub仓库](https://github.com/zy84338719/upftp)。
