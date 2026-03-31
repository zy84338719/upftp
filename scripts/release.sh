#!/bin/bash
#
# Build and publish packages to GitHub Releases
# Requires: goreleaser, gh CLI
#

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

# 函数定义
info() { echo -e "${BLUE}[INFO]${NC} $1"; }
success() { echo -e "${GREEN}[✓]${NC} $1"; }
warn() { echo -e "${YELLOW}[!]${NC} $1"; }
error() { echo -e "${RED}[✗]${NC} $1"; exit 1; }
important() { echo -e "${PURPLE}[IMPORTANT]${NC} $1"; }

# 检查依赖
check_dependencies() {
    info "Checking dependencies..."
    
    if ! command -v goreleaser &> /dev/null; then
        error "goreleaser is not installed. Install with: go install github.com/goreleaser/goreleaser/v2@latest"
    fi
    
    if ! command -v gh &> /dev/null; then
        warn "gh CLI not found. Some features may not work."
    fi
    
    if ! command -v git &> /dev/null; then
        error "git is not installed"
    fi
    
    if ! command -v go &> /dev/null; then
        error "Go is not installed"
    fi
    
    success "All dependencies satisfied"
}

# 验证环境
validate_environment() {
    info "Validating environment..."
    
    if [ -z "$GITHUB_TOKEN" ] && [ -z "$GH_TOKEN" ]; then
        warn "GITHUB_TOKEN not set. Release may fail."
    fi
    
    if [ -n "$(git status --porcelain)" ]; then
        error "Working directory is not clean. Commit or stash changes first."
    fi
    
    # 检查当前分支
    current_branch=$(git rev-parse --abbrev-ref HEAD)
    info "Current branch: $current_branch"
    
    success "Environment validated"
}

# 运行测试
run_tests() {
    info "Running tests..."
    
    # 运行单元测试
    go test ./... -v
    
    if [[ $? -ne 0 ]]; then
        error "Tests failed"
    fi
    
    success "Tests passed"
}

# 构建快照
build_snapshot() {
    info "Building snapshot..."
    
    # 清理旧的构建
    rm -rf dist/
    
    # 运行 goreleaser 构建快照
    goreleaser build --snapshot --clean
    
    if [[ $? -ne 0 ]]; then
        error "Snapshot build failed"
    fi
    
    success "Snapshot build complete"
    important "Check dist/ directory for built artifacts"
}

# 本地发布
release_local() {
    info "Running local release (no publish)..."
    
    # 清理旧的构建
    rm -rf dist/
    
    # 运行 goreleaser 本地发布
    goreleaser release --snapshot --clean
    
    if [[ $? -ne 0 ]]; then
        error "Local release failed"
    fi
    
    success "Local release complete"
    important "Check dist/ directory for release artifacts"
}

# 正式发布
release() {
    local skip_publish="$1"
    
    # 检查是否有版本标签
    if ! git describe --tags --exact-match &> /dev/null; then
        error "No exact version tag found. Please create a version tag first."
    fi
    
    current_version=$(git describe --tags --exact-match)
    info "Releasing version: $current_version"
    
    # 清理旧的构建
    rm -rf dist/
    
    info "Running goreleaser..."
    
    if [ "$skip_publish" = true ]; then
        goreleaser release --skip=publish --clean
    else
        goreleaser release --clean
    fi
    
    if [[ $? -ne 0 ]]; then
        error "Release failed"
    fi
    
    success "Release complete!"
    important "Check dist/ directory for release artifacts"
    important "Check GitHub Releases for published packages"
}

# 验证配置
validate_config() {
    info "Validating goreleaser configuration..."
    
    goreleaser check
    
    if [[ $? -ne 0 ]]; then
        error "Configuration validation failed"
    fi
    
    success "Configuration validated successfully"
}

# 显示版本信息
display_version() {
    info "UPFTP Release Script"
    info "Current version: $(git describe --tags --always 2>/dev/null || echo "dev")"
    info "Goreleaser version: $(goreleaser --version 2>&1 | head -n 1)"
    info "Git commit: $(git rev-parse --short HEAD 2>/dev/null || echo "unknown")"
}

# 显示帮助信息
show_usage() {
    echo "UPFTP Release Script"
    echo ""
    echo "Usage: $0 <command> [options]"
    echo ""
    echo "Commands:"
    echo "  test        Run tests only"
    echo "  snapshot    Build a snapshot (no version tag required)"
    echo "  local       Run local release (no publish)"
    echo "  release     Create and publish release (requires tag)"
    echo "  validate    Validate goreleaser config"
    echo "  version     Display version information"
    echo ""
    echo "Options:"
    echo "  --skip-publish    Skip publishing to GitHub"
    echo "  --help            Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 test"
    echo "  $0 snapshot"
    echo "  $0 local"
    echo "  $0 release"
    echo "  $0 release --skip-publish"
    echo "  $0 validate"
    echo "  $0 version"
    echo ""
}

# 主函数
main() {
    local command="${1:-help}"
    local skip_publish=false
    
    shift || true
    
    while [ $# -gt 0 ]; do
        case "$1" in
            --skip-publish)
                skip_publish=true
                ;;
            --help|-h)
                show_usage
                exit 0
                ;;
            *)
                warn "Unknown option: $1"
                ;;
        esac
        shift
    done
    
    case "$command" in
        test)
            check_dependencies
            run_tests
            ;;
        snapshot)
            check_dependencies
            build_snapshot
            ;;
        local)
            check_dependencies
            run_tests
            release_local
            ;;
        release)
            check_dependencies
            validate_environment
            run_tests
            release "$skip_publish"
            ;;
        validate)
            check_dependencies
            validate_config
            ;;
        version)
            display_version
            ;;
        help|--help|-h)
            show_usage
            ;;
        *)
            error "Unknown command: $command"
            show_usage
            exit 1
            ;;
    esac
}

# 运行主函数
main "$@"
