# upftp - CloudWeGo IDL 重构项目 - 产品需求文档

## 概述
- **Summary**: 使用 CloudWeGo IDL 协议约束接口，通过 hz 工具链生成接口代码，全面重写 upftp 项目，优化项目结构和代码质量，同时保持原有的交互模式。
- **Purpose**: 提升项目的可维护性、可扩展性和性能，采用现代化的接口定义和代码生成方式，同时保持用户熟悉的交互体验。
- **Target Users**: 开发者和系统管理员，需要快速部署和使用文件共享服务的用户。

## 目标
- 使用 CloudWeGo IDL 定义所有接口
- 通过 hz 工具链生成接口代码
- 重构项目目录结构，采用模块化设计
- 集成成熟的日志库（如 zap）
- 支持多种协议：HTTP、FTP、WebDAV、NFS
- 提供 TUI 和 Web 界面
- 保持原有的核心功能和交互模式，同时提升性能和可靠性

## 非目标（范围外）
- 保持与原有代码的兼容性
- 维持旧的项目结构和代码风格
- 保留已有的自定义实现的协议服务器
- 改变现有的命令行参数和配置文件格式

## 背景与上下文
- 原项目是一个跨平台的文件共享服务器，支持 HTTP、FTP 和 MCP 协议
- 项目当前使用标准库和自定义实现，缺乏统一的接口定义
- 需要采用现代化的接口定义和代码生成方式，提升项目质量

## 功能需求
- **FR-1**: 使用 CloudWeGo IDL 定义所有服务接口
- **FR-2**: 通过 hz 工具链生成接口代码
- **FR-3**: 重构项目目录结构，采用模块化设计
- **FR-4**: 集成成熟的日志库（如 zap）
- **FR-5**: 支持 HTTP 协议的文件共享和管理
- **FR-6**: 支持 FTP 协议的文件传输
- **FR-7**: 支持 WebDAV 协议
- **FR-8**: 支持 NFS 协议
- **FR-9**: 提供 TUI 界面进行管理，保持原有交互模式
- **FR-10**: 提供 Web 界面进行文件管理
- **FR-11**: 保持现有的命令行参数和配置文件格式
- **FR-12**: 保持现有的服务启动和管理方式

## 非功能需求
- **NFR-1**: 性能优化，支持并发访问
- **NFR-2**: 高可靠性，具备错误处理和恢复机制
- **NFR-3**: 可扩展性，支持插件式架构
- **NFR-4**: 安全性，提供基本的认证和授权机制
- **NFR-5**: 可配置性，支持通过配置文件调整系统行为

## 约束
- **技术**: Go 1.24+，CloudWeGo Kitex，CloudWeGo Hertz
- **依赖**: CloudWeGo IDL，hz 工具链，zap 日志库
- **平台**: 跨平台支持（Linux、macOS、Windows）

## 假设
- 项目将完全重写，不考虑与原有代码的兼容性
- 使用 CloudWeGo 生态系统的工具和库
- 保持原有的核心功能和用户体验

## 验收标准

### AC-1: CloudWeGo IDL 接口定义
- **Given**: 项目环境已搭建
- **When**: 编写 IDL 文件定义所有服务接口
- **Then**: IDL 文件应包含所有必要的服务和方法定义
- **Verification**: `human-judgment`

### AC-2: 代码生成
- **Given**: IDL 文件已编写完成
- **When**: 使用 hz 工具链生成接口代码
- **Then**: 生成的代码应符合 CloudWeGo 规范，无语法错误
- **Verification**: `programmatic`

### AC-3: 项目目录结构重构
- **Given**: 代码生成完成
- **When**: 调整项目目录结构
- **Then**: 目录结构应清晰合理，符合模块化设计原则
- **Verification**: `human-judgment`

### AC-4: 日志库集成
- **Given**: 项目结构已调整
- **When**: 集成 zap 日志库
- **Then**: 系统应使用 zap 进行日志记录，日志格式统一
- **Verification**: `programmatic`

### AC-5: HTTP 服务
- **Given**: 核心功能已实现
- **When**: 启动 HTTP 服务
- **Then**: 服务应正常运行，支持文件浏览和下载
- **Verification**: `programmatic`

### AC-6: FTP 服务
- **Given**: 核心功能已实现
- **When**: 启动 FTP 服务
- **Then**: 服务应正常运行，支持 FTP 客户端连接和文件传输
- **Verification**: `programmatic`

### AC-7: WebDAV 服务
- **Given**: 核心功能已实现
- **When**: 启动 WebDAV 服务
- **Then**: 服务应正常运行，支持 WebDAV 客户端连接
- **Verification**: `programmatic`

### AC-8: NFS 服务
- **Given**: 核心功能已实现
- **When**: 启动 NFS 服务
- **Then**: 服务应正常运行，支持 NFS 客户端挂载
- **Verification**: `programmatic`

### AC-9: TUI 界面
- **Given**: 核心功能已实现
- **When**: 启动 TUI 界面
- **Then**: 界面应正常显示，支持文件管理操作，保持原有交互模式
- **Verification**: `human-judgment`

### AC-10: Web 界面
- **Given**: 核心功能已实现
- **When**: 访问 Web 界面
- **Then**: 界面应正常显示，支持文件管理操作
- **Verification**: `human-judgment`

### AC-11: 命令行参数
- **Given**: 项目已构建
- **When**: 使用现有命令行参数启动服务
- **Then**: 服务应正常启动，参数应被正确解析
- **Verification**: `programmatic`

### AC-12: 配置文件格式
- **Given**: 项目已构建
- **When**: 使用现有格式的配置文件启动服务
- **Then**: 服务应正常启动，配置应被正确加载
- **Verification**: `programmatic`

## 未解决问题
- [ ] 具体的 IDL 接口设计细节
- [ ] WebDAV 和 NFS 协议的具体实现方式
- [ ] 各协议服务的配置参数设计
- [ ] 性能测试和优化策略