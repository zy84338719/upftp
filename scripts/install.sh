#!/bin/bash
#
# UPFTP Universal Installer Script
# https://github.com/zy84338719/upftp
#
# Usage:
#   curl -fsSL https://zy84338719.github.io/upftp/install.sh | bash
#   curl -fsSL https://zy84338719.github.io/upftp/install.sh | bash -s -- --help
#

set -e

REPO_OWNER="zy84338719"
REPO_NAME="upftp"
REPO_URL="https://github.com/${REPO_OWNER}/${REPO_NAME}"
GITHUB_PAGES_URL="https://${REPO_OWNER}.github.io/${REPO_NAME}"
APT_REPO_URL="${GITHUB_PAGES_URL}/apt"
BREW_TAP="${REPO_OWNER}/tap"
SCOOP_BUCKET="${REPO_OWNER}/scoop-bucket"

VERSION=""
OS=""
ARCH=""
INSTALL_DIR="/usr/local/bin"
TMP_DIR=""

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m'

print_banner() {
    echo -e "${CYAN}"
    echo "╔═══════════════════════════════════════════════════╗"
    echo "║                                                   ║"
    echo "║     ██╗   ██╗████████╗██████╗ ███████╗██████╗    ║"
    echo "║     ██║   ██║╚══██╔══╝██╔══██╗██╔════╝██╔══██╗   ║"
    echo "║     ██║   ██║   ██║   ██████╔╝█████╗  ██████╔╝   ║"
    echo "║     ╚██╗ ██╔╝   ██║   ██╔══██╗██╔══╝  ██╔══██╗   ║"
    echo "║      ╚████╔╝    ██║   ██║  ██║███████╗██║  ██║   ║"
    echo "║       ╚═══╝     ╚═╝   ╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝   ║"
    echo "║                                                   ║"
    echo "║     Modern Cross-Platform File Sharing Server    ║"
    echo "║                                                   ║"
    echo "╚═══════════════════════════════════════════════════╝"
    echo -e "${NC}"
}

info() { echo -e "${BLUE}[INFO]${NC} $1"; }
success() { echo -e "${GREEN}[✓]${NC} $1"; }
warn() { echo -e "${YELLOW}[!]${NC} $1"; }
error() { echo -e "${RED}[✗]${NC} $1"; }

cleanup() {
    if [ -n "$TMP_DIR" ] && [ -d "$TMP_DIR" ]; then
        rm -rf "$TMP_DIR"
    fi
}

trap cleanup EXIT

detect_os() {
    case "$(uname -s)" in
        Darwin*)  OS="darwin" ;;
        Linux*)   OS="linux" ;;
        MINGW*|MSYS*|CYGWIN*|Windows_NT*) OS="windows" ;;
        FreeBSD*) OS="freebsd" ;;
        *) error "Unsupported OS: $(uname -s)"; exit 1 ;;
    esac
}

detect_arch() {
    case "$(uname -m)" in
        x86_64|amd64) ARCH="amd64" ;;
        aarch64|arm64) ARCH="arm64" ;;
        armv7l|armv7) ARCH="arm" ;;
        i386|i686) ARCH="386" ;;
        *) error "Unsupported architecture: $(uname -m)"; exit 1 ;;
    esac
}

get_latest_version() {
    local version
    
    if command -v curl &> /dev/null; then
        version=$(curl -fsSL "https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/releases/latest" 2>/dev/null | grep '"tag_name"' | sed -E 's/.*"v([^"]+)".*/\1/')
    elif command -v wget &> /dev/null; then
        version=$(wget -qO- "https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/releases/latest" 2>/dev/null | grep '"tag_name"' | sed -E 's/.*"v([^"]+)".*/\1/')
    fi
    
    if [ -z "$version" ]; then
        error "Failed to get latest version from GitHub"
        exit 1
    fi
    
    VERSION="$version"
}

check_sudo() {
    if [ "$EUID" -ne 0 ] && [ "$OS" != "windows" ]; then
        if ! command -v sudo &> /dev/null; then
            error "This script requires sudo privileges but sudo is not available"
            exit 1
        fi
        SUDO="sudo"
    else
        SUDO=""
    fi
}

check_dependencies() {
    local deps="curl"
    
    if [ "$OS" = "linux" ]; then
        if command -v apt-get &> /dev/null; then
            $SUDO apt-get update -qq
            $SUDO apt-get install -y -qq curl ca-certificates
        elif command -v yum &> /dev/null; then
            $SUDO yum install -y -q curl ca-certificates
        elif command -v dnf &> /dev/null; then
            $SUDO dnf install -y -q curl ca-certificates
        elif command -v pacman &> /dev/null; then
            $SUDO pacman -Sy --noconfirm curl ca-certificates
        fi
    fi
}

install_via_homebrew() {
    if [ "$OS" != "darwin" ]; then
        return 1
    fi
    
    if ! command -v brew &> /dev/null; then
        return 1
    fi
    
    info "Installing via Homebrew..."
    
    brew tap "$BREW_TAP" 2>/dev/null || true
    brew install upftp
    
    success "Installed via Homebrew"
    return 0
}

install_via_apt() {
    if [ "$OS" != "linux" ]; then
        return 1
    fi
    
    if ! command -v apt-get &> /dev/null; then
        return 1
    fi
    
    info "Installing via APT..."
    
    $SUDO apt-get update -qq
    $SUDO apt-get install -y -qq curl gnupg ca-certificates
    
    local list_file="/etc/apt/sources.list.d/upftp.list"
    
    if [ -f "$list_file" ]; then
        $SUDO rm -f "$list_file"
    fi
    
    echo "deb [trusted=yes] ${APT_REPO_URL} stable main" | $SUDO tee "$list_file" > /dev/null
    
    $SUDO apt-get update -qq
    $SUDO apt-get install -y -qq upftp
    
    success "Installed via APT"
    return 0
}

install_via_rpm() {
    if [ "$OS" != "linux" ]; then
        return 1
    fi
    
    local pkg_manager=""
    if command -v dnf &> /dev/null; then
        pkg_manager="dnf"
    elif command -v yum &> /dev/null; then
        pkg_manager="yum"
    else
        return 1
    fi
    
    info "Installing via $pkg_manager (RPM)..."
    
    local rpm_url="${REPO_URL}/releases/download/v${VERSION}/upftp_${VERSION}_linux_${ARCH}.rpm"
    local rpm_file="/tmp/upftp-${VERSION}.rpm"
    
    curl -fsSL -o "$rpm_file" "$rpm_url"
    $SUDO $pkg_manager install -y -q "$rpm_file"
    rm -f "$rpm_file"
    
    success "Installed via $pkg_manager"
    return 0
}

install_via_scoop() {
    if [ "$OS" != "windows" ]; then
        return 1
    fi
    
    if ! command -v scoop &> /dev/null; then
        return 1
    fi
    
    info "Installing via Scoop..."
    
    scoop bucket add upftp "https://github.com/${SCOOP_BUCKET}"
    scoop install upftp
    
    success "Installed via Scoop"
    return 0
}

install_via_binary() {
    info "Installing binary directly..."
    
    local ext="tar.gz"
    local archive_name="upftp_${OS}_${ARCH}.${ext}"
    
    if [ "$OS" = "windows" ]; then
        ext="zip"
        archive_name="upftp_windows_${ARCH}.zip"
    fi
    
    local download_url="${REPO_URL}/releases/download/v${VERSION}/${archive_name}"
    
    TMP_DIR=$(mktemp -d)
    cd "$TMP_DIR"
    
    info "Downloading v${VERSION} for ${OS}/${ARCH}..."
    
    if ! curl -fsSL --progress-bar -o "upftp.${ext}" "$download_url"; then
        error "Failed to download from $download_url"
        exit 1
    fi
    
    info "Extracting..."
    
    if [ "$ext" = "zip" ]; then
        if ! command -v unzip &> /dev/null; then
            error "unzip is required but not installed"
            exit 1
        fi
        unzip -q "upftp.${ext}"
    else
        tar -xzf "upftp.${ext}"
    fi
    
    local binary="upftp"
    [ "$OS" = "windows" ] && binary="upftp.exe"
    
    if [ ! -f "$binary" ]; then
        error "Binary not found in archive"
        exit 1
    fi
    
    if [ "$OS" = "windows" ]; then
        INSTALL_DIR="$HOME/bin"
        mkdir -p "$INSTALL_DIR"
        cp "$binary" "$INSTALL_DIR/"
        chmod +x "$INSTALL_DIR/$binary"
    else
        $SUDO mkdir -p "$INSTALL_DIR"
        $SUDO cp "$binary" "$INSTALL_DIR/upftp"
        $SUDO chmod +x "$INSTALL_DIR/upftp"
    fi
    
    success "Installed to ${INSTALL_DIR}/${binary}"
    return 0
}

setup_systemd_service() {
    if [ "$OS" != "linux" ]; then
        return
    fi
    
    if ! command -v systemctl &> /dev/null; then
        return
    fi
    
    info "Setting up systemd service..."
    
    $SUDO useradd -r -s /bin/false upftp 2>/dev/null || true
    $SUDO mkdir -p /var/lib/upftp
    $Sudo chown upftp:upftp /var/lib/upftp 2>/dev/null || true
    
    success "Systemd service available. Run: sudo systemctl enable --now upftp"
}

verify_installation() {
    local binary="upftp"
    [ "$OS" = "windows" ] && binary="upftp.exe"
    
    if command -v "$binary" &> /dev/null; then
        local installed_version
        installed_version=$("$binary" -h 2>&1 | head -1 || echo "installed")
        success "Verification successful: $binary is available"
        return 0
    fi
    
    if [ -f "${INSTALL_DIR}/${binary}" ]; then
        if [[ ":$PATH:" != *":${INSTALL_DIR}:"* ]]; then
            warn "${INSTALL_DIR} is not in your PATH"
            echo ""
            echo "Add it to your PATH by running:"
            if [ -n "$ZSH_VERSION" ]; then
                echo "  echo 'export PATH=\"\$PATH:${INSTALL_DIR}\"' >> ~/.zshrc && source ~/.zshrc"
            elif [ -n "$BASH_VERSION" ]; then
                echo "  echo 'export PATH=\"\$PATH:${INSTALL_DIR}\"' >> ~/.bashrc && source ~/.bashrc"
            else
                echo "  export PATH=\"\$PATH:${INSTALL_DIR}\""
            fi
        fi
        return 0
    fi
    
    return 1
}

print_post_install() {
    echo ""
    echo -e "${BOLD}Installation completed!${NC}"
    echo ""
    echo "Quick Start:"
    echo "  upftp -auto -d /path/to/share"
    echo ""
    echo "Options:"
    echo "  -p <port>       HTTP port (default: 10000)"
    echo "  -ftp <port>     FTP port (default: 2121)"
    echo "  -d <dir>        Share directory"
    echo "  -auto           Auto-select network interface"
    echo "  -enable-ftp     Enable FTP server"
    echo ""
    if [ "$OS" = "linux" ] && command -v systemctl &> /dev/null; then
        echo "System Service:"
        echo "  sudo systemctl start upftp"
        echo "  sudo systemctl enable upftp"
        echo ""
    fi
    echo "Documentation: ${REPO_URL}"
    echo ""
}

show_usage() {
    echo "UPFTP Installer v1.0.0"
    echo ""
    echo "Usage: curl -fsSL <url> | bash [options]"
    echo ""
    echo "Options:"
    echo "  -h, --help          Show this help message"
    echo "  -v, --version VER   Install specific version"
    echo "  -m, --method M      Force installation method (brew, apt, rpm, scoop, binary)"
    echo "  -d, --dir DIR       Installation directory (default: /usr/local/bin)"
    echo "  --no-service        Skip systemd service setup"
    echo ""
    echo "Examples:"
    echo "  curl -fsSL https://zy84338719.github.io/upftp/install.sh | bash"
    echo "  curl -fsSL https://zy84338719.github.io/upftp/install.sh | bash -s -- -v 1.0.0"
    echo "  curl -fsSL https://zy84338719.github.io/upftp/install.sh | bash -s -- -m binary"
    echo ""
}

main() {
    local method="auto"
    local specific_version=""
    local setup_service=true
    
    while [ $# -gt 0 ]; do
        case "$1" in
            -h|--help)
                show_usage
                exit 0
                ;;
            -v|--version)
                shift
                specific_version="$1"
                ;;
            -m|--method)
                shift
                method="$1"
                ;;
            -d|--dir)
                shift
                INSTALL_DIR="$1"
                ;;
            --no-service)
                setup_service=false
                ;;
            *)
                warn "Unknown option: $1"
                ;;
        esac
        shift
    done
    
    print_banner
    
    detect_os
    detect_arch
    
    if [ -n "$specific_version" ]; then
        VERSION="$specific_version"
        info "Installing specified version: v${VERSION}"
    else
        get_latest_version
        info "Latest version: v${VERSION}"
    fi
    
    info "Platform: ${OS}/${ARCH}"
    info "Install directory: ${INSTALL_DIR}"
    echo ""
    
    check_sudo
    check_dependencies
    
    local installed=false
    
    if [ "$method" = "auto" ]; then
        case "$OS" in
            darwin)
                install_via_homebrew && installed=true
                ;;
            linux)
                install_via_apt && installed=true || \
                install_via_rpm && installed=true
                ;;
            windows)
                install_via_scoop && installed=true
                ;;
        esac
        
        if [ "$installed" = false ]; then
            install_via_binary && installed=true
        fi
    else
        case "$method" in
            brew|homebrew) install_via_homebrew && installed=true ;;
            apt|deb|debian) install_via_apt && installed=true ;;
            rpm|yum|dnf) install_via_rpm && installed=true ;;
            scoop) install_via_scoop && installed=true ;;
            binary|download) install_via_binary && installed=true ;;
            *)
                error "Unknown installation method: $method"
                exit 1
                ;;
        esac
    fi
    
    if [ "$installed" = false ]; then
        error "Installation failed"
        exit 1
    fi
    
    if [ "$setup_service" = true ] && [ "$method" != "brew" ] && [ "$method" != "scoop" ]; then
        setup_systemd_service
    fi
    
    verify_installation
    print_post_install
}

main "$@"
