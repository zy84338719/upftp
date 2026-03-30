# UPFTP 项目架构改进与重构计划 v2

> 基于 2026-03-30 全面代码审查（第二轮）| **当前项目编译失败，必须先修复**
> 
> 前一版改进计划执行了一半就中断了：新目录 `dal/biz/model/middleware/protocol/` 已创建，但关键的 handler 层（`route.go`、`api.go`、`page.go`）被 AI 写成了语法错误+逻辑混乱的废代码，导致整个项目无法编译。本计划基于**真实现状**重新制定。

导航：[1. 现状分析](#1-项目现状分析) | [2. 问题清单](#2-核心问题清单) | [3. 目标架构](#3-目标架构设计) | [4. 执行阶段](#4-执行阶段) | [5. 安全修复](#5-安全修复清单) | [6. 测试计划](#6-测试计划)

---

## 1. 项目现状分析

### 1.1 当前目录结构（实际状态）

```
upftp/
├── main.go                                    # 入口（110 行）✅ 正常
├── internal/
│   ├── auth/session.go                        # Session 管理（121 行）✅ 已修复协程泄漏
│   ├── biz/
│   │   ├── file_service.go                    # 文件业务逻辑（177 行）✅ 正常
│   │   └── dir_service.go                     # 目录业务逻辑（117 行）✅ 正常
│   ├── cli/                                   # TUI 层 ✅ 正常
│   ├── conf/conf.go                           # 配置管理（345 行）✅ 正常（但 Port 仍是 string）
│   ├── dal/
│   │   ├── dal.go                             # 接口定义（28 行）✅ 正常
│   │   ├── file.go                            # 文件操作实现（90 行）✅ 正常
│   │   ├── dir.go                             # 空文件 ✅ 正常
│   │   └── path.go                            # 统一路径安全（57 行）✅ 正常
│   ├── filehandler/types.go                   # 🗑️ 旧代码（278 行）—— 与 model/ 完全重复
│   ├── handler/
│   │   ├── auth.go                            # ✅ 正常（151 行）
│   │   ├── settings.go                        # ✅ 正常（111 行），密码已脱敏
│   │   ├── download.go                        # ⚠️ 仍引用 filehandler 包，未用 biz
│   │   ├── fileops.go                         # ⚠️ 仍引用 filehandler 包，未用 biz
│   │   ├── route.go                           # ❌ 完全损坏（248 行垃圾代码）
│   │   ├── api.go                             # ❌ 语法错误（import 双引号）
│   │   └── page.go                            # ❌ 完全损坏（68 行垃圾代码）
│   ├── logger/logger.go                       # ✅ 正常
│   ├── middleware/
│   │   ├── auth.go                            # ✅ 正常（47 行）
│   │   ├── cors.go                            # ✅ 正常（23 行），但配置过于宽松
│   │   └── recovery.go                        # ✅ 正常（22 行）
│   ├── model/
│   │   ├── file_type.go                       # ✅ 正常（251 行）
│   │   ├── file_info.go                       # ✅ 正常（47 行），含 LegacyFileInfo 适配
│   │   └── tree_node.go                       # ✅ 正常（9 行）
│   ├── network/network.go                     # ✅ 正常
│   └── protocol/
│       ├── ftp/
│       │   ├── server.go                      # ✅ 正常（434 行），已用 dal.PathResolver
│       │   ├── dataconn.go                    # ✅ 正常
│       │   ├── transfer.go                    # ⚠️ 仍直接调 os 操作，未完全委托 biz
│       │   └── fileops.go                     # ⚠️ 仍直接调 os 操作，未完全委托 biz
│       ├── http/server.go                     # ✅ 正常（70 行）
│       └── mcp/
│           ├── server.go                      # ⚠️ 有独立 isPathSafe()，未用 dal
│           ├── filesystem.go                  # ⚠️ 仍直接调 os 操作，未用 biz
│           ├── fileops.go                     # ⚠️ 仍直接调 os 操作，未用 biz
│           └── serverctl.go                   # ✅ 基本正常
```

### 1.2 编译状态

```
❌ go build ./...  → FAILED
   internal/handler/api.go:10:5: expected ';', found github
   （连带导致所有依赖 handler 包的模块编译失败）
```

### 1.3 已完成 vs 未完成

| 组件 | 旧计划目标 | 当前状态 | 需要动作 |
|------|-----------|---------|---------|
| `dal/` 接口 + 实现 | ✅ | ✅ 已完成 | 无需改动 |
| `biz/` 业务逻辑层 | ✅ | ✅ 已完成 | 无需改动 |
| `model/` 数据模型 | ✅ | ✅ 已完成 | 需删除重复的 `filehandler/` |
| `middleware/` 中间件 | ✅ | ✅ 已完成 | CORS 需收紧配置 |
| `conf/` 配置管理 | ✅ | ✅ 已完成 | Port 类型仍为 string（小问题） |
| `auth/` Session 修复 | ✅ | ✅ 已完成 | 无需改动 |
| `handler/route.go` | ✅ | ❌ **垃圾代码** | **必须完全重写** |
| `handler/api.go` | ✅ | ❌ **语法错误** | **必须完全重写** |
| `handler/page.go` | ✅ | ❌ **垃圾代码** | **必须完全重写** |
| HTTP handler 用 biz | ✅ | ⚠️ 未完成 | fileops/download 仍用旧方式 |
| FTP 用 biz | ✅ | ⚠️ 未完成 | transfer/fileops 仍直接 os |
| MCP 用 biz | ✅ | ⚠️ 未完成 | filesystem/fileops 仍直接 os |
| 删除 `filehandler/` | ✅ | ⚠️ 未完成 | 仍存在，仍被引用 |
| MCP 统一路径安全 | ✅ | ⚠️ 未完成 | 仍有独立 `isPathSafe()` |

---

## 2. 核心问题清单

### 2.1 🔴 Critical — 阻塞性问题

| # | 问题 | 位置 | 影响 |
|---|------|------|------|
| C1 | **项目无法编译** | `handler/api.go:10` 语法错误 | 整个项目不可用 |
| C2 | **`handler/route.go` 是 AI 幻觉垃圾代码** | 248 行无意义代码片段拼接 | 路由注册完全缺失 |
| C3 | **`handler/page.go` 是 AI 幻觉垃圾代码** | 68 行无意义代码片段拼接 | 首页渲染完全缺失 |
| C4 | **`handler/api.go` 几乎为空** | 仅 16 行，API 功能缺失 | 文件列表/目录树/QR Code API 缺失 |
| C5 | **文件操作仍重复 3 处** | FTP/HTTP/MCP 各自直接调 os | 改一处忘另两处，维护成本极高 |
| C6 | **路径安全仍重复 3 处** | `filehandler.IsPathSafe` / `mcp.isPathSafe` / `dal.PathResolver` | 校验标准不一致 |

### 2.2 🟠 High — 架构问题

| # | 问题 | 位置 | 影响 |
|---|------|------|------|
| H1 | **HTTP handler 未用 biz 层** | `handler/fileops.go`, `handler/download.go` | 绕过了大小校验和统一错误处理 |
| H2 | **FTP 协议层未用 biz 层** | `protocol/ftp/transfer.go`, `protocol/ftp/fileops.go` | 同上 |
| H3 | **MCP 协议层未用 biz 层** | `protocol/mcp/filesystem.go`, `protocol/mcp/fileops.go` | 同上 |
| H4 | **`filehandler/` 包未删除** | `handler/fileops.go:14`, `handler/download.go:15` | 与 `model/` 包完全重复 |
| H5 | **CORS 配置过于宽松** | `middleware/cors.go:11` `Allow-Origin: *` | 允许任意域跨域请求 |
| H6 | **配置 Port 仍为 string 类型** | `conf/conf.go:14` `Port string yaml:"port"` | 需要运行时解析，类型不安全 |
| H7 | **函数内定义匿名 struct** | `handler/fileops.go`, `handler/settings.go` 等 | struct 定义分散，不可复用 |

### 2.3 🟡 Medium — 次要问题

| # | 问题 | 位置 |
|---|------|------|
| M1 | MCP `isPathSafe()` 独立实现 | `protocol/mcp/server.go:86-94` |
| M2 | FTP `handleLIST`/`handleMLSD` 仍用 `os.ReadDir` | `protocol/ftp/transfer.go:37,83` |
| M3 | FTP `handleSTOR` 不校验上传大小 | `protocol/ftp/transfer.go:155` |
| M4 | HTTP handler 没有使用 biz 层的 `FileService` 注入 | `handler/fileops.go` 直接 os 操作 |
| M5 | Zip 下载中 `defer file.Close()` 在 Walk 闭包内 | `handler/download.go:83`（已修复为直接 Close） |
| M6 | `dal/dir.go` 空文件 | 可删除或保留占位 |

---

## 3. 目标架构设计

### 3.1 目标：协议层全部委托 biz 层

```
                    ┌──────────────────────────────────┐
                    │          main.go                 │
                    │  初始化 conf/dal/biz → 注入各协议  │
                    └──────────┬───────────────────────┘
                               │
              ┌────────────────┼────────────────┐
              ▼                ▼                 ▼
    ┌─────────────────┐ ┌──────────────┐ ┌──────────────┐
    │ protocol/http   │ │ protocol/ftp │ │ protocol/mcp │
    │ (Hertz handler) │ │ (FTP 命令)    │ │ (MCP tools)  │
    │  仅做HTTP编解码  │ │  仅做FTP编解码│ │  仅做MCP编解码│
    └────────┬────────┘ └──────┬───────┘ └──────┬───────┘
             │                 │                 │
             └────────┬────────┘─────────────────┘
                      ▼
            ┌──────────────────┐
            │   biz/           │  ← 核心服务层（唯一业务逻辑入口）
            │   FileService    │     Upload/Download/Delete/Rename/...
            └────────┬─────────┘
                     │
            ┌────────┴─────────┐
            │   dal/           │  ← 数据访问层（唯一文件系统入口）
            │   FileStore      │     os.ReadFile / os.WriteFile / ...
            │   PathResolver   │     统一路径安全校验
            └──────────────────┘
```

### 3.2 核心原则

| 原则 | 规则 |
|------|------|
| **单一入口** | 所有文件操作必须通过 `biz.FileService`，协议层禁止直接调 `os.*` |
| **统一路径安全** | 所有路径校验必须通过 `dal.PathResolver`，禁止各自实现 `isPathSafe` |
| **依赖注入** | `biz.FileService` 在 `main.go` 创建，通过 `SetFileService()` 注入到各协议层 |
| **无全局状态** | handler 不依赖全局变量，通过闭包或 setter 获取服务实例 |

---

## 4. 执行阶段

### Phase 0: 紧急修复 — 恢复编译（预计 0.5 天）

**目标：** 让项目重新编译通过

**步骤：**

#### 0.1 重写 `handler/route.go`

当前的 `route.go` 是 248 行 AI 幻觉垃圾代码。需要完全重写为：

```go
package handler

import (
    "embed"
    
    "github.com/cloudwego/hertz/pkg/route"
    "github.com/zy84338719/upftp/internal/middleware"
)

var (
    templates   embed.FS
    serverInfo  *ServerInfo
    fileSvc     *biz.FileService
)

type ServerInfo struct { ... }

func SetServerInfo(ip string, httpPort, ftpPort int, root string) { ... }
func SetFileService(svc *biz.FileService) { fileSvc = svc }

func RegisterRoutes(r *route.RouterGroup) {
    // 公开路由
    r.GET("/login", HandleLoginPage)
    r.POST("/api/login", HandleLogin)
    r.POST("/api/logout", HandleLogout)
    r.GET("/api/settings", HandleGetSettings)
    r.POST("/api/settings/language", HandleSetLanguage)
    r.POST("/api/settings/http-auth", HandleSetHTTPAuth)
    r.POST("/api/settings/ftp", HandleSetFTP)

    // 需要认证的路由
    auth := r.Group("/", middleware.AuthMiddleware(conf.AppConfig, GetSessionManager()))
    auth.GET("/", HandleIndexPage)
    auth.GET("/api/info", HandleServerInfo)
    auth.GET("/api/files", HandleFileListAPI)
    auth.GET("/api/tree", HandleDirectoryTree)
    auth.GET("/api/qrcode", HandleQRCode)
    auth.POST("/api/upload", handleUpload)
    auth.POST("/api/create-folder", handleCreateFolder)
    auth.POST("/api/delete", handleDelete)
    auth.POST("/api/rename", handleRename)
    auth.GET("/download/:path", handleDownload)
    auth.GET("/preview/:path", handlePreview)
    auth.GET("/files/:path", handleFiles)

    // 静态资源
    r.Static("/static", "./static")
}
```

#### 0.2 重写 `handler/api.go`

需要包含完整功能：

- `HandleServerInfo` — 服务器信息 API
- `HandleFileListAPI` — 文件列表 API（使用 `biz.FileService.ListFiles`）
- `HandleDirectoryTree` — 目录树 API（使用 `biz.FileService.BuildTree`）
- `HandleQRCode` — 二维码 API

#### 0.3 重写 `handler/page.go`

首页渲染，使用 `biz.FileService.ListFiles` 列出根目录文件。

#### 0.4 删除 `filehandler/` 包

所有引用改为 `model` 包 + `dal.PathResolver`：

```bash
# 需要修改的文件：
handler/fileops.go    → 用 fileSvc 替换 filehandler.IsPathSafe + 直接 os 操作
handler/download.go   → 用 fileSvc 替换 filehandler.IsPathSafe + GetFileType + GetMimeType
# 删除：
rm -rf internal/filehandler/
```

**验证：** `go build ./...` 通过

### Phase 1: HTTP handler 层委托 biz（预计 1 天）

**目标：** 所有 HTTP handler 通过 `biz.FileService` 操作文件

**前置条件：** Phase 0 完成

#### 1.1 重构 `handler/fileops.go`

```go
// 改前：直接调 os
func handleUpload(ctx context.Context, c *app.RequestContext) {
    if !filehandler.IsPathSafe(targetPath) { ... }
    targetDir := path.Join(conf.AppConfig.Root, targetPath)
    os.MkdirAll(targetDir, 0755)
    dst, _ := os.Create(dstPath)
    io.Copy(dst, file)
}

// 改后：委托 biz
func handleUpload(ctx context.Context, c *app.RequestContext) {
    if !conf.AppConfig.Upload.Enabled { ... }
    for _, fileHeader := range files {
        file, _ := fileHeader.Open()
        if err := fileSvc.Upload(
            path.Join(targetPath, fileHeader.Filename),
            file,
            fileHeader.Size,
        ); err != nil {
            logger.Error("Upload error: %v", err)
            continue
        }
    }
}
```

同理重构 `handleCreateFolder`、`handleDelete`、`handleRename`。

#### 1.2 重构 `handler/download.go`

```go
// 改前：
func handleDownload(ctx context.Context, c *app.RequestContext) {
    if !filehandler.IsPathSafe(filename) { ... }
    filePath := path.Join(conf.AppConfig.Root, filename)
    // 直接 os.Stat + c.File
}

// 改后：
func handleDownload(ctx context.Context, c *app.RequestContext) {
    fullPath, err := fileSvc.Download(filename)
    if err != nil { ... } // 包含路径安全校验
    // c.File(fullPath)
}
```

#### 1.3 验证

- `go build ./...` 通过
- 手动测试上传/下载/删除/重命名

### Phase 2: FTP 协议层委托 biz（预计 1 天）

**目标：** FTP 文件操作全部通过 `biz.FileService`

**前置条件：** Phase 1 完成

#### 2.1 重构 `protocol/ftp/transfer.go`

```go
// 改前：
func (s *FTPServer) handleLIST(...) {
    dirPath, _ := s.pathResolver.ResolvePath(client.cwd)
    files, _ := os.ReadDir(dirPath)
    // 手工格式化
}

// 改后：
func (s *FTPServer) handleLIST(...) {
    files, err := s.fileSvc.ListFiles(client.cwd)
    if err != nil { ... }
    // 格式化输出 files
}
```

`handleRETR`、`handleSTOR`、`handleAPPE` 同理委托。

#### 2.2 重构 `protocol/ftp/fileops.go`

`handleMKD` → `fileSvc.CreateFolder()`
`handleRMD` → `fileSvc.Delete()`（目录）
`handleDELE` → `fileSvc.Delete()`（文件）
`handleRNTO` → `fileSvc.Rename()`
`handleSIZE` → `fileSvc.Stat()`
`handleMDTM` → `fileSvc.Stat()`

#### 2.3 注入 `biz.FileService` 到 FTP

```go
// server.go
func (s *FTPServer) SetFileService(svc *biz.FileService) {
    s.fileSvc = svc
}
```

#### 2.4 验证

- `go build ./...` 通过
- 手动测试 FTP 连接/上传/下载

### Phase 3: MCP 协议层委托 biz（预计 1 天）

**目标：** MCP 工具全部通过 `biz.FileService`，删除独立 `isPathSafe()`

**前置条件：** Phase 2 完成

#### 3.1 重构 `protocol/mcp/server.go`

- 删除 `isPathSafe()` 方法（86-94 行）
- 注入 `*biz.FileService`
- 所有文件操作改为调用 `fileSvc`

#### 3.2 重构 `protocol/mcp/filesystem.go`

| MCP 方法 | 当前 | 目标 |
|---------|------|------|
| `handleListFiles` | `os.ReadDir` | `fileSvc.ListFiles()` |
| `handleGetFileInfo` | `os.Stat` | `fileSvc.Stat()` |
| `handleReadFile` | `os.ReadFile` + 独立安全校验 | `fileSvc.ReadFileContent()` |
| `handleWriteFile` | `os.WriteFile` + 独立安全校验 | `fileSvc.WriteFileContent()` |
| `handleDownloadFile` | `os.ReadFile` | `fileSvc.ReadFileContent()` |
| `handleSearchFiles` | `filepath.Walk` | `fileSvc.SearchFiles()` |
| `handleGetDirectoryTree` | `buildTree` + `os.ReadDir` | `fileSvc.BuildTree()` |

#### 3.3 重构 `protocol/mcp/fileops.go`

| MCP 方法 | 当前 | 目标 |
|---------|------|------|
| `handleUploadFile` | base64 解码 + `os.WriteFile` | `fileSvc.WriteFileContent()` |
| `handleDeleteFile` | `os.RemoveAll` | `fileSvc.Delete()` |
| `handleRenameFile` | `os.Rename` | `fileSvc.Rename()` |
| `handleMoveFile` | `os.Rename` | `fileSvc.Move()` |
| `handleCopyFile` | `io.Copy` | `fileSvc.Copy()` |
| `handleCreateDirectory` | `os.MkdirAll` | `fileSvc.CreateFolder()` |

#### 3.4 验证

- `go build ./...` 通过
- MCP 工具调用测试

### Phase 4: CORS 收紧 + 清理（预计 0.5 天）

**目标：** 修复安全问题、清理冗余代码

**前置条件：** Phase 3 完成

#### 4.1 CORS 配置收紧

```go
// 改前：允许所有来源
c.Header("Access-Control-Allow-Origin", "*")

// 改后：可配置化，默认仅允许局域网
func CORSMiddleware(allowedOrigins []string) app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        origin := string(c.GetHeader("Origin"))
        if isAllowedOrigin(origin, allowedOrigins) {
            c.Header("Access-Control-Allow-Origin", origin)
        }
        // ...
    }
}
```

#### 4.2 清理工作

- 删除 `var _ = os.DevNull`（`protocol/mcp/server.go:115`）
- 删除 `ftpJoin()`、`statPath()` 辅助函数（`protocol/ftp/transfer.go:13-19`）如果已不再使用
- 确认 `dal/dir.go` 是否需要保留
- 全局搜索 `os.ReadDir` / `os.ReadFile` / `os.Create` / `os.Open` / `os.Remove` / `os.Rename` / `os.MkdirAll` 确认仅在 `dal/` 中出现

#### 4.3 验证

- `go build ./...` 通过
- `go vet ./...` 无警告
- `grep -rn "os\.ReadDir\|os\.ReadFile\|os\.Create\|os\.Remove\|os\.Rename\|os\.MkdirAll" internal/ --include="*.go" | grep -v "internal/dal/"` 应无结果

### Phase 5: 补全测试（预计 1 天）

**前置条件：** Phase 4 完成

#### 5.1 已有测试

- `internal/dal/dal_test.go` — 需确认覆盖度
- `internal/biz/file_service_test.go` — 需确认覆盖度

#### 5.2 需要新增的测试

```
internal/dal/
  └── path_test.go          # 路径穿越攻击向量 + 正常路径

internal/biz/
  └── file_service_test.go  # Upload 大小校验 / Delete / Rename / 路径穿越拒绝
```

#### 5.3 安全测试向量

```go
traversalVectors := []string{
    "..",
    "../",
    "../../etc/passwd",
    "sub/../../../etc",
    "/../../../etc",
    "..\\..\\windows\\system32",  // Windows 风格
}
```

**验证：** `go test ./...` 全部通过

---

## 5. 安全修复清单

| # | 问题 | 严重度 | 修复位置 | 阶段 |
|---|------|--------|----------|------|
| 1 | **项目无法编译** | 🔴 Critical | `handler/route.go`, `api.go`, `page.go` | Phase 0 |
| 2 | **AI 生成垃圾代码** | 🔴 Critical | `handler/route.go`, `page.go` | Phase 0 |
| 3 | **filehandler 重复包** | 🟠 High | 删除 `internal/filehandler/` | Phase 0 |
| 4 | **HTTP handler 绕过 biz** | 🟠 High | `handler/fileops.go`, `download.go` | Phase 1 |
| 5 | **FTP 绕过 biz + 无大小校验** | 🟠 High | `protocol/ftp/transfer.go` | Phase 2 |
| 6 | **MCP 独立 isPathSafe()** | 🟠 High | `protocol/mcp/server.go:86-94` | Phase 3 |
| 7 | **MCP 绕过 biz** | 🟠 High | `protocol/mcp/filesystem.go`, `fileops.go` | Phase 3 |
| 8 | **CORS Allow-Origin: \*** | 🟡 Medium | `middleware/cors.go` | Phase 4 |

---

## 6. 测试计划

每个 Phase 完成后的验证标准：

| Phase | 验证命令 | 预期 |
|-------|---------|------|
| 0 | `go build ./...` | ✅ 编译通过 |
| 1 | `go build ./...` + 手动测试 HTTP 上传/下载/删除 | ✅ |
| 2 | `go build ./...` + 手动测试 FTP 连接 | ✅ |
| 3 | `go build ./...` + MCP 工具调用 | ✅ |
| 4 | `go vet ./...` + `grep` 检查无 os 直接调用 | ✅ |
| 5 | `go test ./...` | ✅ 全部通过 |

---

## 7. 时间估算

| 阶段 | 描述 | 预估时间 | 状态 |
|------|------|----------|------|
| Phase 0 | **紧急修复：恢复编译** | 0.5 天 | ⏳ 待开始 |
| Phase 1 | HTTP handler 委托 biz | 1 天 | ⏳ 待开始 |
| Phase 2 | FTP 协议层委托 biz | 1 天 | ⏳ 待开始 |
| Phase 3 | MCP 协议层委托 biz | 1 天 | ⏳ 待开始 |
| Phase 4 | CORS 收紧 + 清理 | 0.5 天 | ⏳ 待开始 |
| Phase 5 | 补全测试 | 1 天 | ⏳ 待开始 |
| **总计** | | **5 天** | |

---

## 8. 关键决策记录

| 决策 | 选择 | 理由 |
|------|------|------|
| route.go/page.go/api.go 重写方式 | 完全重写 | 当前是 AI 幻觉垃圾代码，无法修补 |
| filehandler/ 包处理 | 直接删除，全局替换为 model/ + dal | 已有完整替代品 |
| FTP/HTTP/MCP 委托顺序 | HTTP → FTP → MCP | HTTP 最常用且最容易测试 |
| CORS 策略 | 可配置白名单 | 默认仅允许局域网，保持便捷性 |
| biz.FileService 注入方式 | setter 方法 | 最小改动，避免改动 main.go 启动流程 |

---

*本计划基于 2026-03-30 代码真实状态制定，前一版计划因执行中断留下了大量损坏文件。*
