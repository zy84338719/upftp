#!/bin/bash
#
# Build and publish packages to GitHub Releases
# Requires: goreleaser, gh CLI
#

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

info() { echo -e "${BLUE}[INFO]${NC} $1"; }
success() { echo -e "${GREEN}[✓]${NC} $1"; }
warn() { echo -e "${YELLOW}[!]${NC} $1"; }
error() { echo -e "${RED}[✗]${NC} $1"; exit 1; }

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
    
    success "All dependencies satisfied"
}

validate_environment() {
    info "Validating environment..."
    
    if [ -z "$GITHUB_TOKEN" ] && [ -z "$GH_TOKEN" ]; then
        warn "GITHUB_TOKEN not set. Release may fail."
    fi
    
    if [ -n "$(git status --porcelain)" ]; then
        error "Working directory is not clean. Commit or stash changes first."
    fi
    
    success "Environment validated"
}

run_tests() {
    info "Running tests..."
    go test ./... -v
    success "Tests passed"
}

build_snapshot() {
    info "Building snapshot..."
    goreleaser build --snapshot --clean
    success "Snapshot build complete"
}

release_local() {
    info "Running local release (no publish)..."
    goreleaser release --snapshot --clean
    success "Local release complete. Check dist/ directory."
}

release() {
    local skip_publish="$1"
    
    info "Running goreleaser..."
    
    if [ "$skip_publish" = true ]; then
        goreleaser release --skip=publish --clean
    else
        goreleaser release --clean
    fi
    
    success "Release complete!"
}

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
    echo ""
    echo "Options:"
    echo "  --skip-publish    Skip publishing to GitHub"
    echo ""
    echo "Examples:"
    echo "  $0 snapshot"
    echo "  $0 local"
    echo "  $0 release"
    echo "  $0 release --skip-publish"
    echo ""
}

main() {
    local command="${1:-help}"
    local skip_publish=false
    
    shift || true
    
    while [ $# -gt 0 ]; do
        case "$1" in
            --skip-publish)
                skip_publish=true
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
            goreleaser check
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

main "$@"
