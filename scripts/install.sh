#!/bin/bash

# UPFTP Universal Installer Script
# Supports multiple installation methods

set -e

REPO_URL="https://github.com/zy84338719/upftp"
APT_REPO_URL="https://apt.upftp.dev"  # Replace with your actual APT repository
BREW_TAP="zy84338719/tap"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Detect OS and architecture
detect_os() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)
    
    case $ARCH in
        x86_64) ARCH="amd64" ;;
        aarch64|arm64) ARCH="arm64" ;;
        i386|i686) ARCH="386" ;;
        *) 
            print_error "Unsupported architecture: $ARCH"
            exit 1
            ;;
    esac
    
    print_status "Detected OS: $OS, Architecture: $ARCH"
}

# Install via APT (Debian/Ubuntu)
install_via_apt() {
    print_status "Installing UPFTP via APT..."
    
    # Check if running as root or with sudo
    if [[ $EUID -ne 0 ]]; then
        SUDO="sudo"
    else
        SUDO=""
    fi
    
    # Add repository
    print_status "Adding UPFTP repository..."
    $SUDO apt-get update
    $SUDO apt-get install -y curl gnupg
    
    # Add GPG key (if you have one)
    # curl -fsSL $APT_REPO_URL/key.gpg | $SUDO apt-key add -
    
    # Add repository
    echo "deb $APT_REPO_URL stable main" | $SUDO tee /etc/apt/sources.list.d/upftp.list
    
    # Update and install
    $SUDO apt-get update
    $SUDO apt-get install -y upftp
    
    print_success "UPFTP installed successfully via APT!"
    print_status "Start service: sudo systemctl start upftp"
    print_status "Enable on boot: sudo systemctl enable upftp"
}

# Install via Homebrew (macOS)
install_via_brew() {
    print_status "Installing UPFTP via Homebrew..."
    
    # Check if Homebrew is installed
    if ! command -v brew &> /dev/null; then
        print_error "Homebrew is not installed. Please install it first:"
        print_status "Visit https://brew.sh for installation instructions"
        exit 1
    fi
    
    # Add tap and install
    brew tap $BREW_TAP
    brew install upftp
    
    print_success "UPFTP installed successfully via Homebrew!"
    print_status "Start service: brew services start upftp"
}

# Install via direct download
install_via_download() {
    print_status "Installing UPFTP via direct download..."
    
    # Get latest release version
    VERSION=$(curl -s https://api.github.com/repos/zy84338719/upftp/releases/latest | grep '"tag_name"' | cut -d'"' -f4)
    
    if [ -z "$VERSION" ]; then
        print_error "Failed to get latest version"
        exit 1
    fi
    
    print_status "Latest version: $VERSION"
    
    # Determine file format
    if [ "$OS" = "windows" ]; then
        EXT="zip"
    else
        EXT="tar.gz"
    fi
    
    # Download URL
    FILENAME="upftp_${OS}_${ARCH}.${EXT}"
    DOWNLOAD_URL="$REPO_URL/releases/download/$VERSION/$FILENAME"
    
    print_status "Downloading: $DOWNLOAD_URL"
    
    # Create temporary directory
    TMP_DIR=$(mktemp -d)
    cd "$TMP_DIR"
    
    # Download and extract
    curl -fsSL "$DOWNLOAD_URL" -o "$FILENAME"
    
    if [ "$EXT" = "zip" ]; then
        unzip -q "$FILENAME"
    else
        tar -xzf "$FILENAME"
    fi
    
    # Install binary
    if [ "$OS" = "windows" ]; then
        BINARY="upftp.exe"
        INSTALL_DIR="$HOME/bin"
    else
        BINARY="upftp"
        INSTALL_DIR="/usr/local/bin"
    fi
    
    # Check permissions and install
    if [ -w "$INSTALL_DIR" ] || [ "$EUID" -eq 0 ]; then
        cp "$BINARY" "$INSTALL_DIR/"
        chmod +x "$INSTALL_DIR/$BINARY"
    else
        sudo cp "$BINARY" "$INSTALL_DIR/"
        sudo chmod +x "$INSTALL_DIR/$BINARY"
    fi
    
    # Cleanup
    cd - > /dev/null
    rm -rf "$TMP_DIR"
    
    print_success "UPFTP installed successfully to $INSTALL_DIR!"
}

# Show usage information
show_usage() {
    echo "UPFTP Installation Script"
    echo ""
    echo "Usage: $0 [method]"
    echo ""
    echo "Installation methods:"
    echo "  apt      - Install via APT repository (Debian/Ubuntu)"
    echo "  brew     - Install via Homebrew (macOS)"
    echo "  download - Install via direct download (all platforms)"
    echo "  auto     - Auto-detect best method (default)"
    echo ""
    echo "Examples:"
    echo "  $0 apt"
    echo "  $0 brew"
    echo "  $0 download"
    echo "  curl -fsSL https://install.upftp.dev | bash"
}

# Auto-detect installation method
auto_install() {
    detect_os
    
    case $OS in
        linux)
            if command -v apt-get &> /dev/null; then
                print_status "Detected Debian/Ubuntu system, using APT"
                install_via_apt
            elif command -v yum &> /dev/null || command -v dnf &> /dev/null; then
                print_warning "RPM-based system detected, falling back to direct download"
                install_via_download
            else
                print_warning "Package manager not detected, using direct download"
                install_via_download
            fi
            ;;
        darwin)
            if command -v brew &> /dev/null; then
                print_status "Detected macOS with Homebrew, using Homebrew"
                install_via_brew
            else
                print_warning "Homebrew not found, using direct download"
                install_via_download
            fi
            ;;
        *)
            print_status "Using direct download for $OS"
            install_via_download
            ;;
    esac
}

# Main execution
main() {
    METHOD=${1:-auto}
    
    case $METHOD in
        apt)
            detect_os
            install_via_apt
            ;;
        brew)
            detect_os
            install_via_brew
            ;;
        download)
            detect_os
            install_via_download
            ;;
        auto)
            auto_install
            ;;
        help|--help|-h)
            show_usage
            exit 0
            ;;
        *)
            print_error "Unknown installation method: $METHOD"
            show_usage
            exit 1
            ;;
    esac
    
    echo ""
    print_success "Installation completed!"
    print_status "Run 'upftp -h' to see usage options"
    print_status "Visit $REPO_URL for documentation"
}

# Run main function
main "$@"
