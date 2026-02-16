# Devopin CLI

[![Release](https://img.shields.io/github/v/release/gabutlabs/devopin-cli)](https://github.com/gabutlabs/devopin-cli/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/gabutlabs/devopin-cli)](https://goreportcard.com/report/github.com/gabutlabs/devopin-cli)
[![License](https://img.shields.io/github/license/gabutlabs/devopin-cli)](LICENSE)

A command-line application for monitoring system resources (CPU, Memory, Disk) and sending real-time alerts via Telegram.

## Features

- 🖥️ **Real-time Monitoring**: Monitor CPU, Memory, and Disk usage
- 📱 **Telegram Alerts**: Get instant notifications when resource usage exceeds thresholds
- ⚙️ **Configurable**: Customizable thresholds and check intervals
- 🔧 **Systemd Integration**: Run as a background service on Linux
- 🌐 **Cross-Platform**: Supports Linux (AMD64, ARM64, ARMv7) and macOS

## Quick Install

### One-Liner Installation

```bash
# Install latest version (Linux/macOS)
curl -fsSL https://github.com/gabutlabs/devopin-cli/releases/latest/download/install.sh | sudo bash
```

### Install Specific Version

```bash
curl -fsSL https://github.com/gabutlabs/devopin-cli/releases/download/v1.0.0/install.sh | sudo bash -s -- --version v1.0.0
```

### Using Homebrew (macOS)

```bash
# Coming soon
```

For detailed installation instructions, see [INSTALL.md](INSTALL.md).

## Usage

### Start Resource Monitoring

```bash
# Run manually
devopin resource-alert

# Run as systemd service (Linux)
sudo systemctl start devopin-resource-alert
sudo systemctl enable devopin-resource-alert  # Auto-start on boot
```

### Check Service Status

```bash
sudo systemctl status devopin-resource-alert
```

### View Logs

```bash
# Follow logs in real-time
sudo journalctl -u devopin-resource-alert -f

# View last 50 lines
sudo journalctl -u devopin-resource-alert -n 50
```

## Configuration

### Configuration File

Create `/etc/devopin/config.yaml` (production) or `config.yaml` (development):

```yaml
resource_alert:
  interval: 30s  # Check every 30 seconds
  memory:
    max_percent: 90  # Alert when memory > 90%
  cpu:
    max_percent: 90  # Alert when CPU > 90%
  disk:
    max_percent: 90  # Alert when disk > 90%

notify:
  telegram:
    bot_token: "YOUR_BOT_TOKEN"  # Get from @BotFather
    chat_id: 123456789            # Get from @userinfobot

server:
  host: ""  # Leave empty for auto-detect
```

### Environment Variables

Alternatively, use environment variables (takes precedence over config file):

```bash
export DEVOPIN_TELEGRAM_BOT_TOKEN="your_bot_token"
export DEVOPIN_TELEGRAM_CHAT_ID="your_chat_id"
export DEVOPIN_RESOURCE_ALERT_INTERVAL="30s"
export DEVOPIN_RESOURCE_ALERT_CPU_MAX_PERCENT="90"
export DEVOPIN_RESOURCE_ALERT_MEMORY_MAX_PERCENT="90"
export DEVOPIN_RESOURCE_ALERT_DISK_MAX_PERCENT="90"
```

## Getting Telegram Credentials

1. **Get Bot Token**:
   - Open Telegram and search for `@BotFather`
   - Send `/newbot` and follow the instructions
   - Copy the bot token

2. **Get Chat ID**:
   - Open Telegram and search for `@userinfobot`
   - Start a chat and send any message
   - Copy your chat ID

## Commands

```bash
# Show version
devopin version

# Start resource monitoring
devopin resource-alert

# Uninstall
sudo devopin uninstall

# Show help
devopin --help
```

## Systemd Service Commands

```bash
# Start service
sudo systemctl start devopin-resource-alert

# Stop service
sudo systemctl stop devopin-resource-alert

# Restart service
sudo systemctl restart devopin-resource-alert

# Enable auto-start
sudo systemctl enable devopin-resource-alert

# Disable auto-start
sudo systemctl disable devopin-resource-alert

# Check status
sudo systemctl status devopin-resource-alert

# View logs
sudo journalctl -u devopin-resource-alert -f
```

## Uninstall

```bash
# Using the CLI
sudo devopin uninstall

# Or manually
sudo systemctl stop devopin-resource-alert
sudo systemctl disable devopin-resource-alert
sudo rm /etc/systemd/system/devopin-resource-alert.service
sudo rm /usr/local/bin/devopin
sudo rm -rf /etc/devopin
```

## Building from Source

```bash
# Clone repository
git clone https://github.com/gabutlabs/devopin-cli.git
cd devopin-cli

# Build binary
go build -o devopin ./cmd/devopin

# Run
./devopin resource-alert
```

## Project Structure

```
devopin-cli/
├── cmd/
│   └── devopin/
│       ├── main.go
│       └── command/
│           ├── root.go
│           ├── resource-alert.go
│           ├── version.go
│           └── uninstall.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── notification/
│   │   └── notif.go
│   └── resource_alert/
│       ├── resource-monitoring.go
│       └── runner.go
├── scripts/
│   ├── install.sh
│   └── devopin-resource-alert.service
└── config/
    └── config.yaml.example
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Support

- 📖 [Installation Guide](INSTALL.md)
- 🐛 [Issue Tracker](https://github.com/gabutlabs/devopin-cli/issues)
- 💬 [Discussions](https://github.com/gabutlabs/devopin-cli/discussions)
