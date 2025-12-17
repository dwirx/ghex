#!/bin/bash
# GHEX Uninstaller Script
# Usage: curl -sSL https://raw.githubusercontent.com/dwirx/ghex/main/scripts/uninstall.sh | bash

set -e

BINARY_NAME="ghex"
INSTALL_DIR="/usr/local/bin"
CONFIG_DIR_PRIMARY="$HOME/.config/ghe"
CONFIG_DIR_LEGACY="$HOME/.config/github-switch"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_banner() {
    echo -e "${RED}"
    echo "  ██████╗ ██╗  ██╗███████╗██╗  ██╗"
    echo " ██╔════╝ ██║  ██║██╔════╝╚██╗██╔╝"
    echo " ██║  ███╗███████║█████╗   ╚███╔╝ "
    echo " ██║   ██║██╔══██║██╔══╝   ██╔██╗ "
    echo " ╚██████╔╝██║  ██║███████╗██╔╝ ██╗"
    echo "  ╚═════╝ ╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝"
    echo -e "${NC}"
    echo "GHEX Uninstaller"
    echo ""
}

info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

confirm() {
    local prompt="$1"
    local default="${2:-n}"
    
    if [ "$default" = "y" ]; then
        prompt="$prompt [Y/n]: "
    else
        prompt="$prompt [y/N]: "
    fi
    
    read -p "$prompt" response
    response=${response:-$default}
    
    case "$response" in
        [yY][eE][sS]|[yY]) return 0 ;;
        *) return 1 ;;
    esac
}

remove_binary() {
    local binary_path="${INSTALL_DIR}/${BINARY_NAME}"
    
    if [ -f "$binary_path" ]; then
        info "Removing binary: $binary_path"
        
        if [ -w "$binary_path" ] || [ -w "$INSTALL_DIR" ]; then
            rm -f "$binary_path"
            success "Binary removed"
        else
            info "Elevated permissions required..."
            if sudo rm -f "$binary_path"; then
                success "Binary removed"
            else
                error "Failed to remove binary. Try manually: sudo rm $binary_path"
                return 1
            fi
        fi
    else
        warn "Binary not found at $binary_path"
    fi
}

remove_config() {
    local removed=false
    
    if [ -d "$CONFIG_DIR_PRIMARY" ]; then
        info "Removing config directory: $CONFIG_DIR_PRIMARY"
        rm -rf "$CONFIG_DIR_PRIMARY"
        success "Config directory removed: $CONFIG_DIR_PRIMARY"
        removed=true
    fi
    
    if [ -d "$CONFIG_DIR_LEGACY" ]; then
        info "Removing legacy config directory: $CONFIG_DIR_LEGACY"
        rm -rf "$CONFIG_DIR_LEGACY"
        success "Legacy config directory removed: $CONFIG_DIR_LEGACY"
        removed=true
    fi
    
    if [ "$removed" = false ]; then
        warn "No config directories found"
    fi
}

show_preview() {
    echo ""
    info "The following will be removed:"
    echo ""
    
    local binary_path="${INSTALL_DIR}/${BINARY_NAME}"
    if [ -f "$binary_path" ]; then
        echo "  Binary: $binary_path"
    else
        echo "  Binary: (not found)"
    fi
    
    if [ -d "$CONFIG_DIR_PRIMARY" ]; then
        echo "  Config: $CONFIG_DIR_PRIMARY"
    fi
    
    if [ -d "$CONFIG_DIR_LEGACY" ]; then
        echo "  Legacy Config: $CONFIG_DIR_LEGACY"
    fi
    
    echo ""
}

main() {
    print_banner
    
    # Parse arguments
    local purge=false
    local force=false
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            --purge|-p)
                purge=true
                shift
                ;;
            --force|-f)
                force=true
                shift
                ;;
            --help|-h)
                echo "Usage: uninstall.sh [OPTIONS]"
                echo ""
                echo "Options:"
                echo "  --purge, -p    Remove config files as well"
                echo "  --force, -f    Skip confirmation prompts"
                echo "  --help, -h     Show this help message"
                exit 0
                ;;
            *)
                warn "Unknown option: $1"
                shift
                ;;
        esac
    done
    
    show_preview
    
    # Confirm uninstallation
    if [ "$force" = false ]; then
        if ! confirm "Do you want to uninstall GHEX?"; then
            info "Uninstallation cancelled"
            exit 0
        fi
    fi
    
    # Remove binary
    remove_binary
    
    # Handle config removal
    if [ "$purge" = true ]; then
        remove_config
    elif [ "$force" = false ]; then
        echo ""
        if confirm "Do you want to remove configuration files as well?"; then
            remove_config
        else
            info "Configuration files preserved"
        fi
    fi
    
    echo ""
    success "GHEX has been uninstalled!"
    echo ""
    echo "Thank you for using GHEX! 👋"
}

main "$@"
