#!/bin/bash

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 函数定义
info() { echo -e "${BLUE}[INFO]${NC} $1"; }
success() { echo -e "${GREEN}[✓]${NC} $1"; }
warn() { echo -e "${YELLOW}[!]${NC} $1"; }
error() { echo -e "${RED}[✗]${NC} $1"; exit 1; }

# 变量默认值
RUN_NAME="upftp"
OUTPUT_DIR="output"
BINARY_DIR="${OUTPUT_DIR}/bin"
BUILD_MODE="release"
VERSION="$(git describe --tags --always 2>/dev/null || echo "dev")"
COMMIT="$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")"
BUILD_DATE="$(date -u +"%Y-%m-%dT%H:%M:%SZ")"
GO_VERSION="$(go version | awk '{print $3}')"

# 解析命令行参数
while [[ $# -gt 0 ]]; do
    case $1 in
        -o|--output)
            OUTPUT_DIR="$2"
            BINARY_DIR="${OUTPUT_DIR}/bin"
            shift 2
            ;;
        -n|--name)
            RUN_NAME="$2"
            shift 2
            ;;
        -m|--mode)
            BUILD_MODE="$2"
            shift 2
            ;;
        -h|--help)
            echo "UPFTP Build Script"
            echo ""
            echo "Usage: $0 [options]"
            echo ""
            echo "Options:"
            echo "  -o, --output <dir>    Output directory (default: output)"
            echo "  -n, --name <name>     Binary name (default: upftp)"
            echo "  -m, --mode <mode>     Build mode (debug|release) (default: release)"
            echo "  -h, --help            Show this help message"
            echo ""
            exit 0
            ;;
        *)
            warn "Unknown option: $1"
            shift
            ;;
    esac
done

# 检查 Go 环境
if ! command -v go &> /dev/null; then
    error "Go is not installed"
fi

# 清理并创建输出目录
info "Creating output directory structure..."
rm -rf "${OUTPUT_DIR}"
mkdir -p "${BINARY_DIR}"

# 设置构建标志
LDFLAGS="-X main.Version=${VERSION} -X main.LastCommit=${COMMIT} -X main.BuildDate=${BUILD_DATE}"

if [[ "${BUILD_MODE}" == "release" ]]; then
    info "Building in release mode..."
    LDFLAGS="${LDFLAGS} -s -w"
    CGO_ENABLED=0 go build -ldflags "${LDFLAGS}" -o "${BINARY_DIR}/${RUN_NAME}"
else
    info "Building in debug mode..."
    go build -ldflags "${LDFLAGS}" -o "${BINARY_DIR}/${RUN_NAME}"
fi

# 检查构建是否成功
if [[ $? -ne 0 ]]; then
    error "Build failed"
fi

# 复制必要文件
info "Copying necessary files..."
cp -r script/* "${OUTPUT_DIR}" 2>/dev/null
chmod +x "${OUTPUT_DIR}/bootstrap.sh" 2>/dev/null

# 显示构建信息
success "Build completed successfully!"
success "Binary: ${BINARY_DIR}/${RUN_NAME}"
success "Version: ${VERSION}"
success "Commit: ${COMMIT}"
success "Build Date: ${BUILD_DATE}"
success "Go Version: ${GO_VERSION}"

export PATH="${BINARY_DIR}:${PATH}"
info "You can now run the binary: ${RUN_NAME} --help"