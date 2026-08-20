#!/bin/sh
# upftp 一行安装脚本 —— linux / macOS
#
# 用法:
#   curl -fsSL https://github.com/zy84338719/upftp/raw/main/scripts/install.sh | bash
#
# 行为:自动检测 OS/架构,从 GitHub Releases 下载对应二进制,
#       校验 SHA256,安装到 ~/.local/bin(或 /usr/local/bin)。
#       无需 Go 环境,无需 sudo(默认装到用户目录)。

set -e

# --- 配置 ---
OWNER="zy84338719"
REPO="upftp"
INSTALL_DIR="${UPFTP_INSTALL_DIR:-$HOME/.local/bin}"

# --- 终端着色 ---
if [ -t 1 ]; then
  BLUE='\033[36m'; GREEN='\033[32m'; YELLOW='\033[33m'; RED='\033[31m'; BOLD='\033[1m'; RESET='\033[0m'
else
  BLUE=''; GREEN=''; YELLOW=''; RED=''; BOLD=''; RESET=''
fi
info()  { printf "${BLUE}▸${RESET} %s\n" "$1"; }
ok()    { printf "${GREEN}✓${RESET} %s\n" "$1"; }
warn()  { printf "${YELLOW}⚠${RESET} %s\n" "$1"; }
error() { printf "${RED}✗${RESET} %s\n" "$1" >&2; }

# --- 检测操作系统 ---
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$OS" in
  linux)  OS=linux ;;
  darwin) OS=darwin ;;
  *) error "不支持的操作系统: $(uname -s) (仅支持 linux / macOS)"; exit 1 ;;
esac

case "$ARCH" in
  x86_64|amd64)  ARCH=amd64 ;;
  aarch64|arm64) ARCH=arm64 ;;
  i386|i686)     ARCH=386 ;;   # 兜底(若未来构建含 386)
  *) error "不支持的架构: $ARCH"; exit 1 ;;
esac

# --- 确定最新版本 ---
info "查询最新版本..."
# 优先读取命令行参数指定的版本:$1
TAG="${1:-}"
if [ -z "$TAG" ]; then
  TAG=$(curl -fsSL "https://api.github.com/repos/${OWNER}/${REPO}/releases/latest" \
        | grep '"tag_name"' | head -1 | sed -E 's/.*"([^"]+)".*/\1/')
fi
if [ -z "$TAG" ]; then
  error "无法确定最新版本。请检查网络,或手动指定版本: install.sh v1.0.0"
  exit 1
fi
VERSION="${TAG#v}"  # 去掉前导 v
ok "目标版本: ${BOLD}${TAG}${RESET} (${OS}/${ARCH})"

# --- 拼接下载地址 ---
# goreleaser 产物命名:upftp_VERSION_OS_ARCH.tar.gz
ASSET="upftp_${VERSION}_${OS}_${ARCH}.tar.gz"
CHECKSUMS="checksums.txt"
DOWNLOAD_URL="https://github.com/${OWNER}/${REPO}/releases/download/${TAG}/${ASSET}"
CHECKSUM_URL="https://github.com/${OWNER}/${REPO}/releases/download/${TAG}/${CHECKSUMS}"

# --- 临时目录 ---
TMPDIR=$(mktemp -d)
trap 'rm -rf "$TMPDIR"' EXIT

# --- 下载 ---
info "下载 ${ASSET}..."
if ! curl -fsSL "$DOWNLOAD_URL" -o "${TMPDIR}/${ASSET}"; then
  error "下载失败: $DOWNLOAD_URL"
  error "可能该平台未发布二进制。检查: https://github.com/${OWNER}/${REPO}/releases/tag/${TAG}"
  exit 1
fi

# --- 校验 SHA256(可选:校验失败仅警告,不阻断,避免 macOS 无 sha256sum 的硬失败)---
if command -v sha256sum >/dev/null 2>&1 || command -v shasum >/dev/null 2>&1; then
  info "校验完整性..."
  if curl -fsSL "$CHECKSUM_URL" -o "${TMPDIR}/${CHECKSUMS}"; then
    EXPECTED=$(grep "${ASSET}" "${TMPDIR}/${CHECKSUMS}" | awk '{print $1}')
    if [ -n "$EXPECTED" ]; then
      if command -v sha256sum >/dev/null 2>&1; then
        ACTUAL=$(sha256sum "${TMPDIR}/${ASSET}" | awk '{print $1}')
      else
        ACTUAL=$(shasum -a 256 "${TMPDIR}/${ASSET}" | awk '{print $1}')
      fi
      if [ "$EXPECTED" = "$ACTUAL" ]; then
        ok "SHA256 校验通过"
      else
        error "SHA256 校验失败! 期望 ${EXPECTED}, 实际 ${ACTUAL}"
        exit 1
      fi
    else
      warn "checksums.txt 中未找到 ${ASSET},跳过校验"
    fi
  else
    warn "无法下载 checksums.txt,跳过校验"
  fi
else
  warn "系统无 sha256sum/shasum,跳过校验"
fi

# --- 解压 ---
info "解压..."
tar -xzf "${TMPDIR}/${ASSET}" -C "${TMPDIR}"
# 二进制在归档根目录
BINARY="${TMPDIR}/upftp"
if [ ! -f "$BINARY" ]; then
  # 兜底:某些归档可能把二进制放在子目录
  BINARY=$(find "$TMPDIR" -name upftp -type f | head -1)
fi
if [ -z "$BINARY" ] || [ ! -f "$BINARY" ]; then
  error "解压后未找到 upftp 二进制"
  exit 1
fi
chmod +x "$BINARY"

# --- 安装 ---
# 若 INSTALL_DIR 不可写(如 /usr/local/bin),尝试 sudo
if [ ! -d "$INSTALL_DIR" ]; then
  mkdir -p "$INSTALL_DIR" 2>/dev/null || {
    warn "无法创建 $INSTALL_DIR,尝试使用 sudo"
    sudo mkdir -p "$INSTALL_DIR"
  }
fi

if [ -w "$INSTALL_DIR" ]; then
  mv "$BINARY" "${INSTALL_DIR}/upftp"
else
  warn "$INSTALL_DIR 不可写,使用 sudo 安装"
  sudo mv "$BINARY" "${INSTALL_DIR}/upftp"
fi
chmod +x "${INSTALL_DIR}/upftp" 2>/dev/null || sudo chmod +x "${INSTALL_DIR}/upftp"
ok "已安装到 ${BOLD}${INSTALL_DIR}/upftp${RESET}"

# --- PATH 检查 ---
case ":$PATH:" in
  *":${INSTALL_DIR}:"*) ;;
  *)
    echo
    warn "${INSTALL_DIR} 不在 PATH 中。请将以下内容加入 shell 配置(~/.bashrc 或 ~/.zshrc):"
    printf '    export PATH="%s:$PATH"\n' "$INSTALL_DIR"
    echo
    ;;
esac

# --- 验证 ---
VERSION_INSTALLED=$("${INSTALL_DIR}/upftp" -version 2>/dev/null || echo "unknown")
echo
printf "${GREEN}${BOLD}  ✓ upftp 安装成功!${RESET}  ${VERSION_INSTALLED}\n"
echo
printf "  快速开始:\n"
printf "    upftp                  # 分享当前目录\n"
printf "    upftp -d ~/share       # 分享指定目录\n"
echo
