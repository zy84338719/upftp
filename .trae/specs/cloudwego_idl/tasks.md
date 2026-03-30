# upftp - CloudWeGo IDL 重构项目 - 实现计划

## [ ] 任务 1: 环境搭建和依赖安装
- **Priority**: P0
- **Depends On**: None
- **Description**:
  - 安装 Go 1.24+
  - 安装 CloudWeGo hz 工具链
  - 安装必要的依赖库（zap 日志库等）
- **Acceptance Criteria Addressed**: AC-1, AC-2
- **Test Requirements**:
  - `programmatic` TR-1.1: 验证 Go 版本 >= 1.24
  - `programmatic` TR-1.2: 验证 hz 工具安装成功
  - `programmatic` TR-1.3: 验证依赖库安装成功
- **Notes**: 确保环境变量配置正确

## [ ] 任务 2: 编写 CloudWeGo IDL 接口定义
- **Priority**: P0
- **Depends On**: 任务 1
- **Description**:
  - 设计文件服务接口
  - 设计认证服务接口
  - 设计系统管理接口
  - 编写完整的 IDL 文件
- **Acceptance Criteria Addressed**: AC-1
- **Test Requirements**:
  - `human-judgment` TR-2.1: IDL 文件结构清晰，符合 CloudWeGo 规范
  - `human-judgment` TR-2.2: 接口定义完整，覆盖所有必要功能
- **Notes**: 参考原有项目的功能需求，确保接口设计合理

## [ ] 任务 3: 使用 hz 工具生成代码
- **Priority**: P0
- **Depends On**: 任务 2
- **Description**:
  - 使用 hz 工具根据 IDL 文件生成代码
  - 检查生成的代码结构和内容
  - 调整生成的代码以适应项目需求
- **Acceptance Criteria Addressed**: AC-2
- **Test Requirements**:
  - `programmatic` TR-3.1: 生成的代码无语法错误
  - `programmatic` TR-3.2: 生成的代码符合 CloudWeGo 规范
- **Notes**: 注意生成代码的目录结构

## [x] 任务 4: 重构项目目录结构
- **Priority**: P0
- **Depends On**: 任务 3
- **Description**:
  - 设计新的项目目录结构
  - 组织生成的代码和自定义代码
  - 确保模块化设计
- **Acceptance Criteria Addressed**: AC-3
- **Test Requirements**:
  - `human-judgment` TR-4.1: 目录结构清晰合理
  - `human-judgment` TR-4.2: 模块划分明确
- **Notes**: 参考 CloudWeGo 推荐的项目结构

## [x] 任务 5: 集成 zap 日志库
- **Priority**: P1
- **Depends On**: 任务 4
- **Description**:
  - 安装 zap 日志库
  - 配置日志系统
  - 在各个模块中使用统一的日志接口
- **Acceptance Criteria Addressed**: AC-4
- **Test Requirements**:
  - `programmatic` TR-5.1: 日志系统初始化成功
  - `programmatic` TR-5.2: 日志输出格式正确
- **Notes**: 确保日志级别可配置

## [x] 任务 6: 实现 HTTP 服务
- **Priority**: P0
- **Depends On**: 任务 5
- **Description**:
  - 基于生成的代码实现 HTTP 服务
  - 实现文件浏览和下载功能
  - 实现 Web 界面
- **Acceptance Criteria Addressed**: AC-5, AC-10
- **Test Requirements**:
  - `programmatic` TR-6.1: HTTP 服务启动成功
  - `programmatic` TR-6.2: 文件浏览功能正常
  - `programmatic` TR-6.3: 文件下载功能正常
  - `human-judgment` TR-6.4: Web 界面显示正常
- **Notes**: 使用 CloudWeGo Hertz 框架

## [/] 任务 7: 实现 FTP 服务
- **Priority**: P0
- **Depends On**: 任务 5
- **Description**:
  - 实现 FTP 协议服务器
  - 支持 FTP 客户端连接
  - 实现文件上传和下载功能
- **Acceptance Criteria Addressed**: AC-6
- **Test Requirements**:
  - `programmatic` TR-7.1: FTP 服务启动成功
  - `programmatic` TR-7.2: FTP 客户端可连接
  - `programmatic` TR-7.3: 文件传输功能正常
- **Notes**: 可考虑使用成熟的 FTP 库

## [ ] 任务 8: 实现 WebDAV 服务
- **Priority**: P1
- **Depends On**: 任务 5
- **Description**:
  - 实现 WebDAV 协议服务器
  - 支持 WebDAV 客户端连接
  - 实现基本的文件操作功能
- **Acceptance Criteria Addressed**: AC-7
- **Test Requirements**:
  - `programmatic` TR-8.1: WebDAV 服务启动成功
  - `programmatic` TR-8.2: WebDAV 客户端可连接
  - `programmatic` TR-8.3: 基本文件操作功能正常
- **Notes**: 可考虑使用成熟的 WebDAV 库

## [ ] 任务 9: 实现 NFS 服务
- **Priority**: P1
- **Depends On**: 任务 5
- **Description**:
  - 实现 NFS 协议服务器
  - 支持 NFS 客户端挂载
  - 实现文件共享功能
- **Acceptance Criteria Addressed**: AC-8
- **Test Requirements**:
  - `programmatic` TR-9.1: NFS 服务启动成功
  - `programmatic` TR-9.2: NFS 客户端可挂载
  - `programmatic` TR-9.3: 文件共享功能正常
- **Notes**: 注意 NFS 协议的实现复杂度

## [ ] 任务 10: 实现 TUI 界面
- **Priority**: P1
- **Depends On**: 任务 5
- **Description**:
  - 实现命令行 TUI 界面
  - 支持文件管理操作
  - 支持服务状态管理
  - 保持原有交互模式
- **Acceptance Criteria Addressed**: AC-9
- **Test Requirements**:
  - `programmatic` TR-10.1: TUI 界面启动成功
  - `human-judgment` TR-10.2: TUI 界面操作流畅，保持原有交互模式
  - `programmatic` TR-10.3: 文件管理功能正常
- **Notes**: 可使用成熟的 TUI 库，参考原有实现保持交互模式

## [ ] 任务 11: 实现命令行参数和配置文件处理
- **Priority**: P0
- **Depends On**: 任务 4
- **Description**:
  - 实现命令行参数解析
  - 实现配置文件加载和解析
  - 保持现有的命令行参数和配置文件格式
- **Acceptance Criteria Addressed**: AC-11, AC-12
- **Test Requirements**:
  - `programmatic` TR-11.1: 命令行参数解析正确
  - `programmatic` TR-11.2: 配置文件加载正确
  - `programmatic` TR-11.3: 保持与原有格式兼容
- **Notes**: 参考原有实现，确保参数和配置文件格式保持不变