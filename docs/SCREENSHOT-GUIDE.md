# 📸 UPFTP 项目运行截图指南

## ✅ 服务已成功启动！

```bash
✅ HTTP Server:   http://localhost:10003
✅ FTP Server:    ftp://localhost:2122
✅ WebDAV Server: http://localhost:8083
✅ NFS Server:    nfs://localhost:2050

📁 共享目录: /tmp/upftp/demo-files
👤 用户名: zhangyi
```

**浏览器已自动打开**: http://localhost:10003

---

## 🎯 推荐的截图场景

### 1. 主界面 - 文件浏览器
**展示内容**:
- 现代化的 Web 界面
- 文件和文件夹列表
- 文件大小、修改时间信息
- 上传、新建文件夹按钮

**操作步骤**:
1. 确保浏览器已打开 http://localhost:10003
2. 等待页面完全加载（看到文件列表）
3. 按 `Cmd+Shift+4` (Mac) 或 `Win+Shift+S` (Windows)
4. 按空格键，点击浏览器窗口
5. 保存为 `screenshots/upftp-web-ui-main.png`

### 2. 文件预览功能
**展示内容**:
- 点击文件查看预览效果
- 代码高亮、图片预览等

**操作步骤**:
1. 在 Web 界面点击 `README.md` 文件
2. 等待预览窗口弹出
3. 截图预览窗口
4. 保存为 `screenshots/upftp-file-preview.png`

### 3. 命令行界面
**展示内容**:
- 帮助信息
- ASCII Logo
- 服务器启动日志

**操作步骤**:
1. 打开终端
2. 执行: `cd /tmp/upftp && ./upftp -h`
3. 截图终端窗口
4. 保存为 `screenshots/upftp-cli-help.png`

或者:
1. 在终端执行: `./upftp -auto`
2. 等待服务器启动
3. 截图启动日志
4. 保存为 `screenshots/upftp-cli-start.png`

### 4. 文件夹浏览
**展示内容**:
- 点击 `documents` 文件夹
- 显示文件夹内容

**操作步骤**:
1. 在 Web 界面点击 `documents` 文件夹
2. 等待文件夹内容加载
3. 截图浏览器窗口
4. 保存为 `screenshots/upftp-folder-browse.png`

### 5. 响应式设计（移动端）
**展示内容**:
- 移动端界面
- 响应式布局

**操作步骤**:
1. 打开浏览器开发者工具 (F12)
2. 切换到移动设备模拟模式
3. 选择 iPhone 或 iPad
4. 截图模拟的移动界面
5. 保存为 `screenshots/upftp-mobile.png`

---

## 🛠️ 截图方法

### macOS
```bash
# 窗口截图（推荐）
Cmd + Shift + 4 → 按空格键 → 点击浏览器窗口

# 区域截图
Cmd + Shift + 4 → 拖动选择区域

# 全屏截图
Cmd + Shift + 3
```

### Windows
```bash
# 窗口截图
Win + Shift + S → 选择窗口

# 区域截图
Win + Shift + S → 选择区域

# 全屏截图
Win + PrintScreen
```

### Linux
```bash
# 使用 gnome-screenshot
gnome-screenshot -w  # 窗口
gnome-screenshot -a  # 全屏

# 使用 flameshot
flameshot gui
```

### 浏览器开发者工具
```bash
# 完整页面截图（推荐）
1. 按 F12 打开开发者工具
2. 按 Cmd+Shift+P (Mac) 或 Ctrl+Shift+P
3. 输入 "screenshot"
4. 选择 "Capture full size screenshot"
5. 保存图片
```

---

## 📁 项目结构建议

创建 `screenshots/` 目录存放截图：

```
upftp/
├── screenshots/
│   ├── upftp-web-ui-main.png      # 主界面
│   ├── upftp-file-preview.png     # 文件预览
│   ├── upftp-cli-start.png        # 命令行启动
│   ├── upftp-folder-browse.png    # 文件夹浏览
│   └── upftp-mobile.png           # 移动端界面
└── README.md
```

---

## 📝 添加到 README

截图完成后，更新 README.md：

```markdown
## 📸 Screenshots

### Web Interface

![Web Interface](screenshots/upftp-web-ui-main.png)
*Modern web interface with file browser and preview capabilities*

### File Preview

![File Preview](screenshots/upftp-file-preview.png)
*Rich file preview with syntax highlighting and image viewing*

### Command Line

![CLI](screenshots/upftp-cli-start.png)
*Easy-to-use command line interface with multi-protocol support*
```

---

## 🎨 緻加示例文件（可选）

为了让截图更丰富，可以添加更多示例文件：

```bash
# 创建示例图片文件
cd /tmp/upftp/demo-files/images
curl -o logo.png "https://via.placeholder.com/150"
curl -o banner.jpg "https://via.placeholder.com/800x200"

# 创建示例文档
cd ../documents
cat > tutorial.md << 'EOF'
# UPFTP Tutorial

This is a comprehensive tutorial for UPFTP...

## Installation
...
EOF

# 创建示例代码文件
cat > config.yaml << 'EOF'
server:
  http_port: 10003
  ftp_port: 2122
  webdav_port: 8083
  nfs_port: 2050

upload:
  max_size: 100MB
  enabled: true
EOF

# 刷新浏览器查看新文件
```

---

## 🔄 重启服务

如果需要重启服务：

```bash
# 停止服务
pkill -f "upftp.*10003"

# 重新启动
cd /tmp/upftp
./upftp -auto -d /tmp/upftp/demo-files -p 10003

# 在浏览器中打开
open http://localhost:10003
```

---

## ✨ 灯箱提示

1. **等待加载**: 确保页面完全加载后再截图
2. **高分辨率**: 使用 Retina 或 2x 分辨率
3. **避免敏感信息**: 确保不包含密码等敏感信息
4. **裁剪美化**: 去除浏览器地址栏、书签栏等
5. **使用 PNG 格式**: 保持高质量
6. **文件命名**: 使用有意义的文件名

---

## 🎉 完成！

现在你已经准备好创建专业的项目展示截图了！

1. ✅ **服务正在运行**: http://localhost:10003
2. ✅ **浏览器已打开**: 可以直接看到 Web 界面
3. ✅ **示例文件已创建**: README.md, documents/, images/, projects/
4. ✅ **截图指南已提供**: 按照上述步骤创建截图

创建截图后，记得：
- 将它们添加到 `screenshots/` 目录
- 更新 README.md 添加截图展示
- 提交到 GitHub 仓库

祝截图顺利！📸✨

---

## 📊 当前服务信息

```json
{
  "project": "UPFTP",
  "version": "0.1.3",
  "http_port": 10003,
  "ftp_port": 2122,
  "webdav_port": 8083,
  "nfs_port": 2050,
  "shared_directory": "/tmp/upftp/demo-files",
  "username": "zhangyi",
  "status": "running"
}
```

**当前运行的文件列表**:
```
/tmp/upftp/demo-files/
├── README.md (176 B)
├── documents/ (96 B)
│   └── sample.txt
├── images/ (96 B)
│   └── logo.png
└── projects/ (96 B)
    ├── app.zip
    └── source.tar.gz
```

---

## ⚠️ 注意事项

如果遇到问题：
1. **404 错误**: 确保 `templates/` 目录存在并包含前端文件
2. **端口占用**: 检查 10003 端口是否被占用
3. **权限问题**: 确保有读取共享目录的权限
4. **截图工具**: macOS 需要授权屏幕录制

需要帮助？ 请提供具体的错误信息！🔧
