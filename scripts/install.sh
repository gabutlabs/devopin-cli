#!/bin/bash

set -e

# =============================================================================
# Devopin CLI Installer
# =============================================================================
# Installation script for Devopin CLI with systemd service support.
# Downloads the latest release from GitHub and sets up the resource-alert service.
#
# Usage:
#   curl -fsSL https://github.com/gabutlabs/devopin-cli/releases/latest/download/install.sh | sudo bash
#   OR
#   curl -fsSL https://raw.githubusercontent.com/gabutlabs/devopin-cli/main/scripts/install.sh | bash -s -- --version v1.0.0
# =============================================================================

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
GITHUB_REPO="gabutlabs/devopin-cli"
BINARY_NAME="devopin"
SERVICE_NAME="devopin-resource-alert"
INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="/etc/devopin"
SYSTEMD_DIR="/etc/systemd/system"

# Variables
VERSION=""
SKIP_SERVICE=false
SKIP_CONFIG=false
FORCE_REINSTALL=false

# =============================================================================
# Helper Functions
# =============================================================================

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

die() {
    log_error "$1"
    exit 1
}

# =============================================================================
# Argument Parsing
# =============================================================================

parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -v|--version)
                VERSION="$2"
                shift 2
                ;;
            --skip-service)
                SKIP_SERVICE=true
                shift
                ;;
            --skip-config)
                SKIP_CONFIG=true
                shift
                ;;
            -f|--force)
                FORCE_REINSTALL=true
                shift
                ;;
            -h|--help)
                show_help
                exit 0
                ;;
            *)
                die "Unknown option: $1. Use --help for usage information."
                ;;
        esac
    done
}

show_help() {
    cat << EOF
Devopin CLI Installer

Usage:
  install.sh [OPTIONS]

Options:
  -v, --version VERSION    Install specific version (e.g., v1.0.0)
                           Default: latest release
  --skip-service           Skip systemd service installation
  --skip-config            Skip configuration file creation
  -f, --force              Force reinstallation even if already installed
  -h, --help               Show this help message

Examples:
  # Install latest version
  curl -fsSL https://github.com/gabutlabs/devopin-cli/releases/latest/download/install.sh | sudo bash

  # Install specific version
  curl -fsSL https://github.com/gabutlabs/devopin-cli/releases/download/v1.0.0/install.sh | sudo bash -s -- --version v1.0.0

  # Install without systemd service
  curl -fsSL .../install.sh | sudo bash -s -- --skip-service

  # Install to custom location (requires manual PATH update)
  curl -fsSL .../install.sh | bash -s -- --skip-service
EOF
}

# =============================================================================
# System Checks
# =============================================================================

check_system() {
    log_info "Checking system compatibility..."

    # Check if running as root (required for system-wide installation)
    if [[ $EUID -ne 0 ]] && [[ "$INSTALL_DIR" == "/usr/local/bin" ]]; then
        die "This script must be run as root (sudo) for system-wide installation."
    fi

    # Detect OS and architecture
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)

    case $ARCH in
        x86_64)
            ARCH="amd64"
            ;;
        aarch64|arm64)
            ARCH="arm64"
            ;;
        armv7l)
            ARCH="arm"
            ;;
        *)
            die "Unsupported architecture: $ARCH"
            ;;
    esac

    case $OS in
        linux)
            OS="linux"
            ;;
        darwin)
            OS="darwin"
            log_warning "macOS detected. systemd service will not be installed."
            SKIP_SERVICE=true
            ;;
        *)
            die "Unsupported operating system: $OS"
            ;;
    esac

    log_info "Detected: ${OS}/${ARCH}"
}

check_dependencies() {
    log_info "Checking dependencies..."

    if ! command -v curl &> /dev/null; then
        die "curl is required but not installed."
    fi

    if ! command -v jq &> /dev/null; then
        log_warning "jq not found. Installing..."
        if command -v apt-get &> /dev/null; then
            apt-get update -qq && apt-get install -y -qq jq
        elif command -v yum &> /dev/null; then
            yum install -y -q jq
        elif command -v brew &> /dev/null; then
            brew install jq
        else
            die "Cannot install jq automatically. Please install jq manually."
        fi
    fi
}

# =============================================================================
# Version Resolution
# =============================================================================

resolve_version() {
    if [[ -n "$VERSION" ]]; then
        log_info "Using specified version: $VERSION"
        return
    fi

    log_info "Fetching latest version from GitHub..."
    
    # Try to get latest release from GitHub API
    LATEST_RELEASE=$(curl -fsSL "https://api.github.com/repos/${GITHUB_REPO}/releases/latest" 2>/dev/null || echo "")
    
    if [[ -n "$LATEST_RELEASE" ]]; then
        VERSION=$(echo "$LATEST_RELEASE" | jq -r '.tag_name')
        log_info "Latest version: $VERSION"
    else
        die "Failed to fetch latest version. Please specify a version with --version"
    fi
}

# =============================================================================
# Download and Install Binary
# =============================================================================

download_binary() {
    local asset_name="${BINARY_NAME}_${VERSION}_${OS}_${ARCH}.tar.gz"
    local download_url="https://github.com/${GITHUB_REPO}/releases/download/${VERSION}/${asset_name}"
    local temp_dir=$(mktemp -d)
    local tarball="${temp_dir}/${asset_name}"

    log_info "Downloading ${asset_name}..."

    # Check if URL exists
    if ! curl -fsSL -o /dev/null "$download_url" 2>/dev/null; then
        # Try alternative naming (without OS)
        asset_name="${BINARY_NAME}_${VERSION}_${ARCH}.tar.gz"
        download_url="https://github.com/${GITHUB_REPO}/releases/download/${VERSION}/${asset_name}"
        
        if ! curl -fsSL -o /dev/null "$download_url" 2>/dev/null; then
            die "Binary not found for ${OS}/${ARCH}. Available assets may not include your platform."
        fi
    fi

    # Download the tarball
    if ! curl -fsSL "$download_url" -o "$tarball"; then
        die "Failed to download binary from $download_url"
    fi

    # Extract binary
    log_info "Extracting binary..."
    tar -xzf "$tarball" -C "$temp_dir"

    # Find the binary (might be in a subdirectory)
    local binary_path=$(find "$temp_dir" -type f -name "${BINARY_NAME}" -executable | head -n 1)
    
    if [[ -z "$binary_path" ]]; then
        die "Binary '${BINARY_NAME}' not found in archive"
    fi

    # Install binary
    log_info "Installing binary to ${INSTALL_DIR}..."
    mkdir -p "$INSTALL_DIR"
    cp "$binary_path" "$INSTALL_DIR/${BINARY_NAME}"
    chmod +x "$INSTALL_DIR/${BINARY_NAME}"

    # Cleanup
    rm -rf "$temp_dir"

    log_success "Binary installed successfully!"
}

# =============================================================================
# Configuration Setup
# =============================================================================

setup_config() {
    if [[ "$SKIP_CONFIG" == true ]]; then
        log_info "Skipping configuration setup..."
        return
    fi

    log_info "Setting up configuration..."

    mkdir -p "$CONFIG_DIR"

    local config_file="${CONFIG_DIR}/config.yaml"

    if [[ -f "$config_file" ]] && [[ "$FORCE_REINSTALL" != true ]]; then
        log_warning "Configuration file already exists. Skipping..."
        return
    fi

    # Create example config
    cat > "$config_file" << 'EOF'
# Devopin CLI Configuration
# =========================
# Copy this file and edit the values as needed.
# Environment variables can also be used (prefix with DEVOPIN_)

# Resource Alert Settings
resource_alert:
  interval: 30s  # Check interval
  memory:
    max_percent: 90  # Alert when memory usage exceeds this percentage
  cpu:
    max_percent: 90  # Alert when CPU usage exceeds this percentage
  disk:
    max_percent: 90  # Alert when disk usage exceeds this percentage

# Notification Settings
notify:
  telegram:
    bot_token: ""  # Your Telegram Bot Token (get from @BotFather)
    chat_id: 0     # Your Telegram Chat ID (get from @userinfobot)

# Server Settings
server:
  host: ""  # Leave empty to auto-detect hostname
EOF

    chmod 644 "$config_file"
    log_success "Configuration file created at ${config_file}"
    log_warning "Please edit ${config_file} and add your Telegram credentials!"
}

# =============================================================================
# Systemd Service Setup
# =============================================================================

setup_systemd_service() {
    if [[ "$SKIP_SERVICE" == true ]]; then
        log_info "Skipping systemd service setup..."
        return
    fi

    if [[ "$OS" != "linux" ]]; then
        log_info "Skipping systemd service (not on Linux)..."
        return
    fi

    log_info "Setting up systemd service..."

    # Create service file
    cat > "${SYSTEMD_DIR}/${SERVICE_NAME}.service" << EOF
[Unit]
Description=Devopin CLI Resource Alert Monitor
Documentation=https://github.com/${GITHUB_REPO}
After=network.target
Wants=network.target

[Service]
Type=simple
ExecStart=${INSTALL_DIR}/${BINARY_NAME} resource-alert
Restart=on-failure
RestartSec=10
StandardOutput=journal
StandardError=journal
SyslogIdentifier=${SERVICE_NAME}

# Security hardening
NoNewPrivileges=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/etc/devopin
PrivateTmp=true

# Environment
Environment=PATH=/usr/local/bin:/usr/bin:/bin

[Install]
WantedBy=multi-user.target
EOF

    chmod 644 "${SYSTEMD_DIR}/${SERVICE_NAME}.service"

    # Reload systemd daemon
    log_info "Reloading systemd daemon..."
    systemctl daemon-reload

    # Enable service (but don't start yet)
    log_info "Enabling ${SERVICE_NAME} service..."
    systemctl enable "${SERVICE_NAME}"

    log_success "Systemd service installed and enabled!"
    log_info "Start the service with: sudo systemctl start ${SERVICE_NAME}"
    log_info "Check status with: sudo systemctl status ${SERVICE_NAME}"
}

# =============================================================================
# Post-Installation
# =============================================================================

print_post_install_message() {
    echo ""
    echo "=============================================="
    echo "  Devopin CLI Installation Complete!"
    echo "=============================================="
    echo ""
    echo "Binary location: ${INSTALL_DIR}/${BINARY_NAME}"
    echo "Config location: ${CONFIG_DIR}/config.yaml"
    echo ""
    echo "Next steps:"
    echo "  1. Edit config file:"
    echo "     sudo nano ${CONFIG_DIR}/config.yaml"
    echo ""
    if [[ "$SKIP_SERVICE" != true ]] && [[ "$OS" == "linux" ]]; then
        echo "  2. Start the resource-alert service:"
        echo "     sudo systemctl start ${SERVICE_NAME}"
        echo ""
        echo "  3. Check service status:"
        echo "     sudo systemctl status ${SERVICE_NAME}"
        echo ""
        echo "  4. Enable auto-start on boot (if not already):"
        echo "     sudo systemctl enable ${SERVICE_NAME}"
        echo ""
        echo "  5. View logs:"
        echo "     sudo journalctl -u ${SERVICE_NAME} -f"
    else
        echo "  2. Run manually:"
        echo "     ${BINARY_NAME} resource-alert"
    fi
    echo ""
    echo "Uninstall:"
    echo "  sudo ${INSTALL_DIR}/${BINARY_NAME} uninstall"
    echo "  OR manually remove:"
    echo "    sudo rm ${INSTALL_DIR}/${BINARY_NAME}"
    echo "    sudo rm ${SYSTEMD_DIR}/${SERVICE_NAME}.service"
    echo "    sudo rm -rf ${CONFIG_DIR}"
    echo ""
    echo "=============================================="
}

# =============================================================================
# Main Installation Flow
# =============================================================================

main() {
    echo ""
    echo "=============================================="
    echo "  Devopin CLI Installer"
    echo "=============================================="
    echo ""

    parse_args "$@"
    check_system
    check_dependencies
    resolve_version
    
    echo ""
    log_info "Starting installation..."
    echo ""

    download_binary
    setup_config
    setup_systemd_service
    
    print_post_install_message
}

# Run main function with all arguments
main "$@"
