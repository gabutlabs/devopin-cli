# Devopin CLI - Installation Guide

## Quick Install

### One-Liner Installation (Recommended)

Install the latest version directly from GitHub Releases:

```bash
# Using curl (Linux/macOS)
curl -fsSL https://github.com/gabutlabs/devopin-cli/releases/latest/download/install.sh | sudo bash

# Or using wget
wget -qO- https://github.com/gabutlabs/devopin-cli/releases/latest/download/install.sh | sudo bash
```

### Install Specific Version

```bash
curl -fsSL https://github.com/gabutlabs/devopin-cli/releases/download/v1.0.0/install.sh | sudo bash -s -- --version v1.0.0
```

## Installation Options

The installer supports several options:

```bash
# Show help
curl -fsSL .../install.sh | sudo bash -s -- --help

# Install specific version
curl -fsSL .../install.sh | sudo bash -s -- --version v1.0.0

# Skip systemd service installation
curl -fsSL .../install.sh | sudo bash -s -- --skip-service

# Skip configuration file creation
curl -fsSL .../install.sh | sudo bash -s -- --skip-config

# Force reinstallation
curl -fsSL .../install.sh | sudo bash -s -- --force
```

## Manual Installation

### 1. Download Binary

Visit the [Releases page](https://github.com/gabutlabs/devopin-cli/releases) and download the appropriate binary for your system:

- **Linux AMD64**: `devopin_vX.X.X_linux_amd64.tar.gz`
- **Linux ARM64**: `devopin_vX.X.X_linux_arm64.tar.gz`
- **Linux ARMv7**: `devopin_vX.X.X_linux_arm.tar.gz`
- **macOS AMD64**: `devopin_vX.X.X_darwin_amd64.tar.gz`
- **macOS ARM64**: `devopin_vX.X.X_darwin_arm64.tar.gz`

### 2. Extract and Install

```bash
# Extract the tarball
tar -xzf devopin_v1.0.0_linux_amd64.tar.gz

# Move binary to system path
sudo mv devopin /usr/local/bin/
sudo chmod +x /usr/local/bin/devopin
```

### 3. Verify Installation

```bash
devopin version
```

## Post-Installation Setup

### 1. Configure Devopin CLI

Create the configuration file:

```bash
sudo mkdir -p /etc/devopin
sudo cp /path/to/config.yaml.example /etc/devopin/config.yaml
sudo nano /etc/devopin/config.yaml
```

Edit the configuration and add your Telegram credentials:

```yaml
resource_alert:
  interval: 30s
  memory:
    max_percent: 90
  cpu:
    max_percent: 90
  disk:
    max_percent: 90

notify:
  telegram:
    bot_token: "YOUR_BOT_TOKEN"  # Get from @BotFather
    chat_id: 123456789            # Get from @userinfobot

server:
  host: ""  # Leave empty for auto-detect
```

### 2. Using Environment Variables (Alternative)

Instead of config file, you can use environment variables:

```bash
export DEVOPIN_TELEGRAM_BOT_TOKEN="YOUR_BOT_TOKEN"
export DEVOPIN_TELEGRAM_CHAT_ID="123456789"
export DEVOPIN_RESOURCE_ALERT_INTERVAL="30s"
export DEVOPIN_RESOURCE_ALERT_CPU_MAX_PERCENT="90"
export DEVOPIN_RESOURCE_ALERT_MEMORY_MAX_PERCENT="90"
export DEVOPIN_RESOURCE_ALERT_DISK_MAX_PERCENT="90"
```

### 3. Setup Systemd Service (Linux only)

The installer automatically sets up the systemd service. To manually configure:

```bash
# Copy service file
sudo cp scripts/devopin-resource-alert.service /etc/systemd/system/

# Reload systemd
sudo systemctl daemon-reload

# Enable service (start on boot)
sudo systemctl enable devopin-resource-alert

# Start service
sudo systemctl start devopin-resource-alert

# Check status
sudo systemctl status devopin-resource-alert
```

## Usage

### Running Manually

```bash
# Start resource monitoring
devopin resource-alert

# Check version
devopin version

# Uninstall
sudo devopin uninstall
```

### Using Systemd Service

```bash
# Start service
sudo systemctl start devopin-resource-alert

# Stop service
sudo systemctl stop devopin-resource-alert

# Restart service
sudo systemctl restart devopin-resource-alert

# Check status
sudo systemctl status devopin-resource-alert

# View logs
sudo journalctl -u devopin-resource-alert -f

# Enable auto-start on boot
sudo systemctl enable devopin-resource-alert

# Disable auto-start
sudo systemctl disable devopin-resource-alert
```

## Uninstallation

### Using the CLI

```bash
sudo devopin uninstall
```

### Manual Uninstallation

```bash
# Stop and disable service
sudo systemctl stop devopin-resource-alert
sudo systemctl disable devopin-resource-alert

# Remove service file
sudo rm /etc/systemd/system/devopin-resource-alert.service
sudo systemctl daemon-reload

# Remove binary
sudo rm /usr/local/bin/devopin

# Remove configuration (optional)
sudo rm -rf /etc/devopin
```

## Troubleshooting

### Service Won't Start

Check the logs:

```bash
sudo journalctl -u devopin-resource-alert -n 50 --no-pager
```

Common issues:
- Missing configuration file
- Invalid Telegram bot token
- Invalid Telegram chat ID

### Check Configuration

```bash
# Verify config file exists
sudo ls -la /etc/devopin/config.yaml

# Test run manually to see errors
sudo devopin resource-alert
```

### Permission Issues

Ensure the binary has execute permissions:

```bash
sudo chmod +x /usr/local/bin/devopin
```

## Building from Source

```bash
# Clone repository
git clone https://github.com/gabutlabs/devopin-cli.git
cd devopin-cli

# Build binary
go build -o devopin ./cmd/devopin

# Install to system path
sudo mv devopin /usr/local/bin/
```

## Requirements

- **Operating System**: Linux (systemd) or macOS
- **Architecture**: AMD64, ARM64, or ARMv7
- **Dependencies**: None (statically linked binary)
- **Network**: Required for Telegram notifications

## Getting Help

- [GitHub Issues](https://github.com/gabutlabs/devopin-cli/issues)
- [Documentation](https://github.com/gabutlabs/devopin-cli#readme)
