# UPFTP 架构优化 - The Implementation Plan (Decomposed and Prioritized Task List)

## [ ] Task 1: 修复后端编译错误
- **Priority**: P0
- **Depends On**: None
- **Description**: 
  - 修复 `biz/service/mcp/serverctl.go` 中引用不存在的 `biz/service/http` 包的问题
  - 检查并移除对该包的依赖或创建必要的实现
- **Acceptance Criteria Addressed**: [AC-1]
- **Test Requirements**:
  - `programmatic` TR-1.1: `go build ./...` 编译通过，无错误
- **Notes**: 这是最高优先级任务，必须先完成才能进行其他工作

## [ ] Task 2: 创建类型定义和常量模块
- **Priority**: P0
- **Depends On**: None
- **Description**: 
  - 提取 `app.ts` 中的接口定义到单独的类型文件
  - 提取魔法数字和字符串为常量
  - 创建 `web/src/types/` 目录存放类型定义
  - 创建 `web/src/constants/` 目录存放常量
- **Acceptance Criteria Addressed**: [AC-4]
- **Test Requirements**:
  - `human-judgement` TR-2.1: 所有类型定义集中在 types 目录
  - `human-judgement` TR-2.2: 所有常量集中在 constants 目录，无魔法值
- **Notes**: 为后续重构奠定基础

## [ ] Task 3: 创建统一的 API 层
- **Priority**: P0
- **Depends On**: Task 2
- **Description**: 
  - 创建 `web/src/api/` 目录
  - 封装所有 API 调用，包括文件、设置、服务等
  - 统一错误处理和响应格式
  - 提供类型安全的 API 函数
- **Acceptance Criteria Addressed**: [AC-3, AC-5]
- **Test Requirements**:
  - `human-judgement` TR-3.1: 所有 API 调用通过 api 层
  - `human-judgement` TR-3.2: API 层有完整的类型定义
  - `human-judgement` TR-3.3: 有统一的错误处理机制
- **Notes**: 这是前端重构的核心

## [ ] Task 4: 模块化拆分 Store
- **Priority**: P1
- **Depends On**: Task 2, Task 3
- **Description**: 
  - 将 `app.ts` 拆分为多个 store：
    - `file.store.ts`: 文件列表、目录树、上传下载等
    - `settings.store.ts`: 语言、HTTP 认证等设置
    - `services.store.ts`: FTP、WebDAV、NFS、MCP 服务配置
    - `ui.store.ts`: 加载状态、搜索查询等 UI 状态
  - 每个 store 使用 API 层进行数据交互
  - 保持功能完全兼容
- **Acceptance Criteria Addressed**: [AC-2, AC-3, AC-4, AC-5]
- **Test Requirements**:
  - `human-judgement` TR-4.1: Store 按功能拆分，每个职责单一
  - `human-judgement` TR-4.2: 每个 store 文件 &lt; 200 行
  - `human-judgement` TR-4.3: 所有功能正常工作
- **Notes**: 逐步迁移，保持功能可用

## [ ] Task 5: 提取工具函数
- **Priority**: P1
- **Depends On**: Task 2
- **Description**: 
  - 提取 `getFileIcon`、`getFileType`、`getDownloadUrl`、`copyLink` 等工具函数
  - 创建 `web/src/utils/` 目录
  - 按功能分类存放工具函数
- **Acceptance Criteria Addressed**: [AC-2]
- **Test Requirements**:
  - `human-judgement` TR-5.1: 工具函数集中在 utils 目录
  - `human-judgement` TR-5.2: 工具函数有类型定义和文档
- **Notes**: 提高代码复用性和可测试性

## [ ] Task 6: 清理和优化
- **Priority**: P2
- **Depends On**: Task 1-5
- **Description**: 
  - 删除不再使用的代码和文件
  - 运行 lint 和 format 确保代码质量
  - 更新类型定义文档
- **Acceptance Criteria Addressed**: [AC-2, AC-4]
- **Test Requirements**:
  - `programmatic` TR-6.1: `npm run lint` 通过
  - `programmatic` TR-6.2: `npm run format` 通过
  - `programmatic` TR-6.3: `npm run type-check` 通过
- **Notes**: 确保代码质量

## [ ] Task 7: 完整功能测试
- **Priority**: P1
- **Depends On**: Task 1-6
- **Description**: 
  - 手动测试所有现有功能
  - 确保前端功能完全兼容
  - 确保后端功能正常
- **Acceptance Criteria Addressed**: [AC-1, AC-2, AC-3, AC-4, AC-5]
- **Test Requirements**:
  - `human-judgement` TR-7.1: 所有现有功能正常工作
  - `human-judgement` TR-7.2: UI 行为无变化
- **Notes**: 确保重构不破坏现有功能
