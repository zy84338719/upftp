# UPFTP 架构优化 - Product Requirement Document

## Overview
- **Summary**: 对 UPFTP 项目进行全面架构重构和优化，修复编译错误，提升代码质量、可维护性和可扩展性。包括后端 Go 代码和前端 Vue 代码的优化。
- **Purpose**: 解决当前项目存在的编译错误、架构混乱、职责不清、代码重复等问题，建立清晰的分层架构，提升项目整体质量。
- **Target Users**: 项目开发者和维护者

## Goals
- 修复后端编译错误，确保项目可以正常构建和运行
- 建立清晰的后端分层架构（Handler → Service → DAL）
- 优化前端代码结构，提升可维护性和开发体验
- 减少代码重复，统一业务逻辑入口
- 提升代码质量，遵循最佳实践

## Non-Goals (Out of Scope)
- 不添加新的业务功能
- 不重构 CLI 界面代码
- 不进行大规模 UI 重构
- 不改变现有的 API 接口契约

## Background & Context
当前项目存在以下主要问题：

1. **后端编译错误**：`biz/service/mcp/serverctl.go` 引用了不存在的包
2. **前端代码结构混乱**：`app.ts` 一个文件包含 300+ 行代码，职责不清
3. **缺少 API 层封装**：前端直接在 store 中使用 fetch
4. **缺少类型安全**：前端有很多魔法数字和字符串
5. **错误处理不完善**：缺少统一的错误处理机制

## Functional Requirements
- **FR-1**: 修复后端编译错误，项目可以正常 `go build ./...`
- **FR-2**: 前端 store 模块化拆分，按功能领域分离
- **FR-3**: 建立统一的前端 API 层封装
- **FR-4**: 建立统一的类型定义和常量
- **FR-5**: 建立统一的错误处理机制

## Non-Functional Requirements
- **NFR-1**: 前端代码结构清晰，单一职责
- **NFR-2**: 代码遵循 TypeScript 最佳实践
- **NFR-3**: 提升开发体验，便于后续维护和扩展
- **NFR-4**: 保持现有功能完全兼容

## Constraints
- **Technical**: 
  - 后端使用 Go 1.24+
  - 前端使用 Vue 3 + TypeScript + Pinia
  - 保持 API 接口兼容性
- **Business**: 不改变现有用户功能体验
- **Dependencies**: 不引入新的第三方依赖

## Assumptions
- 现有 API 接口契约保持不变
- 用户功能需求保持不变
- 技术栈保持不变

## Acceptance Criteria

### AC-1: 后端编译修复
- **Given**: 当前项目存在编译错误
- **When**: 执行 `go build ./...`
- **Then**: 编译成功通过，无错误
- **Verification**: `programmatic`

### AC-2: 前端 Store 模块化
- **Given**: 现有 app.ts 包含 300+ 行代码
- **When**: 完成重构
- **Then**: Store 按功能拆分为多个模块（文件管理、设置、服务等），每个模块职责单一
- **Verification**: `human-judgment`

### AC-3: 统一 API 层
- **Given**: 前端直接在 store 中使用 fetch
- **When**: 完成重构
- **Then**: 所有 API 调用通过统一的 API 层封装，有类型定义和错误处理
- **Verification**: `human-judgment`

### AC-4: 类型安全
- **Given**: 代码中有魔法数字和字符串
- **When**: 完成重构
- **Then**: 所有常量和类型都有统一定义，避免魔法值
- **Verification**: `human-judgment`

### AC-5: 错误处理机制
- **Given**: 缺少统一的错误处理
- **When**: 完成重构
- **Then**: 建立统一的错误处理机制，包括用户提示和日志
- **Verification**: `human-judgment`

## Open Questions
- 无明确问题需要用户确认
