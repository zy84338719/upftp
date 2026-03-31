# 📸 UPFTP 项目运行截图指南

## 🚀 当前运行状态

**服务已成功启动！**

```bash
✅ HTTP Server:   http://localhost:10002
✅ FTP Server:    ftp://localhost:2123
✅ WebDAV Server: http://localhost:8082
✅ NFS Server:    nfs://localhost:2051

📁 共享目录: /tmp/upftp/demo-files
👤 FTP用户名: zhangyi
```

**浏览器已自动打开**: http://localhost:10002

---

## 📸 如何创建项目运行截图

### 方法一：使用系统截图工具

#### 1. Web 界面截图

**全屏截图**:
- **macOS**: `Cmd + Shift + 3`
- **Windows**: `Win + PrintScreen`
- **Linux**: `PrintScreen`

**窗口截图** (推荐):
- **macOS**: `Cmd + Shift + 4` → 按空格键 → 点击浏览器窗口
- **Windows**: `Win + Shift + S` → 选择窗口
- **Linux**: 使用 `gnome-screenshot -w`

**区域截图**:
- **macOS**: `Cmd + Shift + 4` → 拖动选择区域
- **Windows**: `Win + Shift + S` → 选择区域
- **Linux**: 使用 `flameshot gui`

#### 2. 命令行界面截图

打开终端，执行以下命令，然后截图：

```bash
cd /tmp/upftp
./upftp -h
```

显示帮助信息后截图。

### 方法二：使用命令行工具

#### macOS (使用 peekaboo)

```bash
# 安装 peekaboo (如果未安装)
brew install steipete/tap/peekaboo

# 授权屏幕录制
# 打开系统偏好设置 → 安全性与隐私 → 屏幕录制 → 添加终端

# 截取窗口
peekaboo image --mode window --app Safari --path /tmp/upftp-web-ui.png

# 截取全屏
peekaboo image --mode screen --screen-index 0 --path /tmp/upftp-fullscreen.png
```

#### Linux (使用 gnome-screenshot)

```bash
# 安装
sudo apt-get install gnome-screenshot

# 截取窗口
gnome-screenshot -w -f /tmp/upftp-web-ui.png

# 截取全屏
gnome-screenshot -f /tmp/upftp-fullscreen.png
```

### 方法三：使用浏览器开发者工具

1. 在浏览器中打开 `http://localhost:10002`
2. 按 `F12` 或 `Cmd+Option+I` 打开开发者工具
3. 按 `Cmd+Shift+P` (Mac) 或 `Ctrl+Shift+P` (Windows/Linux)
4. 输入 `screenshot` 并选择：
   - `Capture full size screenshot` - 完整页面截图
   - `Capture screenshot` - 可见区域截图
   - `Capture node screenshot` - 特定元素截图

---

## 🎯 推荐的截图场景

### 1. 主界面文件浏览器

**显示内容**:
- 文件列表（文件夹、文件）
- 文件大小、日期信息
- 上传按钮、新建文件夹按钮

**操作步骤**:
1. 确保浏览器已打开 http://localhost:10002
2. 按 `Cmd+Shift+4` (Mac) 或 `Win+Shift+S` (Windows)
3. 选择浏览器窗口区域
4. 保存为 `upftp-web-ui.png`

### 2. 命令行启动界面

**显示内容**:
- ASCII Logo
- 服务器启动信息
- 多协议监听状态

**操作步骤**:
1. 打开终端
2. 执行: `cd /tmp/upftp && ./upftp -h`
3. 截图终端窗口
4. 保存为 `upftp-cli.png`

### 3. 服务器状态面板

**显示内容**:
- 所有服务的运行状态
- 监听地址和端口
- 共享目录信息

**操作步骤**:
1. 在浏览器中访问 http://localhost:10002
2. 点击 "Settings" 或 "About" 按钮
3. 截图状态面板
4. 保存为 `upftp-status.png`

### 4. 文件预览功能

**显示内容**:
- 图片预览
- 代码高亮预览
- 视频播放

**操作步骤**:
1. 在 Web 界面中点击一个文件（如 README.md）
2. 等待预览窗口弹出
3. 截图预览窗口
4. 保存为 `upftp-preview.png`

---

## 📁 截图文件命名建议

```
screenshots/
├── upftp-web-ui-main.png      # 主界面文件浏览器
├── upftp-cli-start.png        # 命令行启动界面
├── upftp-file-preview.png     # 文件预览功能
├── upftp-settings.png         # 设置页面
├── upftp-mobile.png           # 移动端响应式界面
└── upftp-upload.png           # 文件上传界面
```

---

## 🎨 截图美化建议

1. **使用深色主题**: 如果 Web 界面支持，切换到深色主题更美观
2. **添加示例文件**: 在共享目录中添加一些真实的示例文件
3. **使用高分辨率**: 确保截图清晰度足够（推荐 2x 或 Retina）
4. **裁剪多余区域**: 去除浏览器地址栏、书签栏等无关区域
5. **添加标注**: 使用工具添加箭头、文字说明等

---

## 📤 将截图添加到 README

创建完成后，将截图添加到 GitHub README:

```markdown
# UPFTP

## Screenshots

![Web Interface](screenshots/upftp-web-ui-main.png)
*Modern web interface for file browsing*

![CLI](screenshots/upftp-cli-start.png)
*Easy-to-use command line interface*

![File Preview](screenshots/upftp-file-preview.png)
*Rich file preview with syntax highlighting*
```

---

## 🔄 重新启动服务

如果服务停止了，可以使用以下命令重新启动：

```bash
cd /tmp/upftp

# 停止旧进程
pkill -f "upftp.*10002"

# 启动新进程
./upftp -auto -d /tmp/upftp/demo-files -p 10002 -ftp 2123 -webdav 8082 -nfs 2051

# 在浏览器中打开
open http://localhost:10002
```

---

## 📝 当前示例文件

我已经为你创建了以下示例文件：

```
/tmp/upftp/demo-files/
├── README.md              # 项目说明文档
├── documents/
│   ├── sample.txt         # 示例文本文件
│   ├── data.json          # JSON 数据文件
│   └── config.yaml        # 配置文件示例
├── images/
│   └── logo.png           # Logo 图片占位符
└── projects/
    ├── app.zip            # 应用打包文件占位符
    └── source.tar.gz      # 源码压缩包占位符
```

这些文件会在 Web 界面中显示，让截图更加真实和完整。

---

## ⚠️ 注意事项

1. **屏幕录制权限**: 如果使用 peekaboo 或其他命令行截图工具，需要授予终端屏幕录制权限
2. **端口占用**: 确保 10002、2123、8082、2051 端口没有被其他程序占用
3. **文件权限**: 确保共享目录中的文件有正确的读取权限
4. **浏览器缓存**: 如果页面显示不正常，尝试硬刷新（Cmd+Shift+R）

---

## 🎉 完成！

现在你可以在浏览器中查看 Web 界面，并使用系统截图工具或上述任何方法创建项目运行截图。

截图完成后，记得将它们添加到 GitHub 仓库的 `screenshots/` 目录，并更新 README.md 文件。

祝截图顺利！📸✨
