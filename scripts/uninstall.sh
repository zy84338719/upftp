#!/bin/bash
#
# UPFTP Uninstaller Script
# https://github.com/zy84338719/upftp
#

set -e

REPO_OWNER="zy84338719"
REPO_NAME="upftp"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m'

info() { echo -e "${BLUE}[INFO]${NC} $1"; }
success() { echo -e "${GREEN}[✓]${NC} $1"; }
warn() { echo -e "${YELLOW}[!]${NC} $1"; }
error() { echo -e "${RED}[✗]${NC} $1"; }

check_sudo() {
    if [ "$EUID" -ne 0 ]; then
        if command -v sudo &> /dev/null; then
            SUDO="sudo"
        else
            error "This script requires sudo privileges"
            exit 1
        fi
    else
        SUDO=""
    fi
}

stop_service() {
    if command -v systemctl &> /dev/null; then
        if systemctl is-active --quiet upftp 2>/dev/null; then
            info "Stopping upftp service..."
            $SUDO systemctl stop upftp 2>/dev/null || true
        fi
        
        if systemctl is-enabled --quiet upftp 2>/dev/null; then
            info "Disabling upftp service..."
            $SUDO systemctl disable upftp 2>/dev/null || true
        fi
    fi
}

remove_homebrew() {
    if command -v brew &> /dev/null && brew list upftp &> /dev/null; then
        info "Removing via Homebrew..."
        brew uninstall upftp
        brew untap "${REPO_OWNER}/tap" 2>/dev/null || true
        success "Removed via Homebrew"
        return 0
    fi
    return 1
}

remove_apt() {
    if command -v apt-get &> /dev/null && dpkg -l upftp &> /dev/null; then
        info "Removing via APT..."
        $SUDO apt-get remove -y -qq upftp
        $SUDO rm -f /etc/apt/sources.list.d/upftp.list
        $SUDO apt-get update -qq
        success "Removed via APT"
        return 0
    fi
    return 1
}

remove_rpm() {
    if command -v rpm &> /dev/null && rpm -q upftp &> /dev/null; then
        info "Removing via RPM..."
        if command -v dnf &> /dev/null; then
            $SUDO dnf remove -y -q upftp
        else
            $SUDO yum remove -y -q upftp
        fi
        success "Removed via RPM"
        return 0
    fi
    return 1
}

remove_scoop() {
    if command -v scoop &> /dev/null && scoop list upftp &> /dev/null; then
        info "Removing via Scoop..."
        scoop uninstall upftp
        scoop bucket rm upftp 2>/dev/null || true
        success "Removed via Scoop"
        return 0
    fi
    return 1
}

remove_binary() {
    local locations=(
        "/usr/local/bin/upftp"
        "/usr/bin/upftp"
        "$HOME/bin/upftp"
        "$HOME/.local/bin/upftp"
    )
    
    for loc in "${locations[@]}"; do
        if [ -f "$loc" ]; then
            info "Removing binary from $loc..."
            rm -f "$loc"
            success "Removed binary from $loc"
            return 0
        fi
    done
    return 1
}

remove_service_files() {
    info "Removing service files..."
    
    if [ -f /etc/systemd/system/upftp.service ]; then
        $SUDO rm -f /etc/systemd/system/upftp.service
        $SUDO systemctl daemon-reload
        success "Removed systemd service file"
    fi
}

remove_config_files() {
    info "Removing configuration files..."
    
    local config_files=(
        "/etc/upftp"
        "/etc/default/upftp"
    )
    
    for file in "${config_files[@]}"; do
        if [ -e "$file" ]; then
            $SUDO rm -rf "$file"
        fi
    done
}

ask_remove_data() {
    local data_dir="/var/lib/upftp"
    
    if [ -d "$data_dir" ]; then
        echo ""
        warn "Data directory exists: $data_dir"
        read -p "Remove data directory? [y/N] " -n 1 -r
        echo ""
        
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            $SUDO rm -rf "$data_dir"
            success "Removed data directory"
        else
            info "Data directory preserved: $data_dir"
        fi
    fi
}

remove_user() {
    if getent passwd upftp &> /dev/null; then
        info "Removing upftp user..."
        $SUDO userdel upftp 2>/dev/null || true
    fi
    
    if getent group upftp &> /dev/null; then
        info "Removing upftp group..."
        $SUDO groupdel upftp 2>/dev/null || true
    fi
}

show_usage() {
    echo "UPFTP Uninstaller"
    echo ""
    echo "Usage: $0 [options]"
    echo ""
    echo "Options:"
    echo "  -h, --help      Show this help message"
    echo "  -y, --yes       Skip confirmation prompts"
    echo "  --purge         Remove all data including shared files"
    echo ""
}

main() {
    local auto_confirm=false
    local purge=false
    
    while [ $# -gt 0 ]; do
        case "$1" in
            -h|--help)
                show_usage
                exit 0
                ;;
            -y|--yes)
                auto_confirm=true
                ;;
            --purge)
                purge=true
                ;;
            *)
                warn "Unknown option: $1"
                ;;
        esac
        shift
    done
    
    echo -e "${CYAN}╔═══════════════════════════════════════╗${NC}"
    echo -e "${CYAN}║      UPFTP Uninstaller                ║${NC}"
    echo -e "${CYAN}╚═══════════════════════════════════════╝${NC}"
    echo ""
    
    check_sudo
    
    if [ "$auto_confirm" = false ]; then
        warn "This will remove UPFTP from your system."
        read -p "Continue? [y/N] " -n 1 -r
        echo ""
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            info "Aborted."
            exit 0
        fi
    fi
    
    echo ""
    stop_service
    
    remove_homebrew || \
    remove_apt || \
    remove_rpm || \
    remove_scoop || \
    remove_binary
    
    remove_service_files
    remove_config_files
    
    if [ "$purge" = true ]; then
        $SUDO rm -rf /var/lib/upftp 2>/dev/null || true
        remove_user
        success "Purged all data"
    else
        ask_remove_data
    fi
    
    echo ""
    success "UPFTP has been uninstalled successfully!"
    echo ""
}

main "$@"
