package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	ResourceAlert ResourceAlertConfig
	Notify        NotifyConfig
	Server        ServerConfig
}

type ResourceAlertConfig struct {
	Interval time.Duration
	Memory   ThresholdConfig
	CPU      ThresholdConfig
	Disk     ThresholdConfig
}

type ThresholdConfig struct {
	MaxPercent int
}

type TelegramConfig struct {
	BotToken string
	ChatID   int64
}

type NotifyConfig struct {
	Telegram TelegramConfig
}

type ServerConfig struct {
	Host string
}

func Load() (*Config, error) {
	setDefaults()

	// ENV support
	// Semua env var pakai prefix DEVOPIN_
	// Contoh: DEVOPIN_NOTIFY_TELEGRAM_CONFIG_BOT_TOKEN
	viper.SetEnvPrefix("DEVOPIN")
	viper.AutomaticEnv()

	// Eksplisit bind env ke key viper
	// supaya AutomaticEnv bisa override nested key dengan benar
	_ = viper.BindEnv("resource_alert.interval", "DEVOPIN_RESOURCE_ALERT_INTERVAL")
	_ = viper.BindEnv("resource_alert.memory.max_percent", "DEVOPIN_RESOURCE_ALERT_MEMORY_MAX_PERCENT")
	_ = viper.BindEnv("resource_alert.cpu.max_percent", "DEVOPIN_RESOURCE_ALERT_CPU_MAX_PERCENT")
	_ = viper.BindEnv("resource_alert.disk.max_percent", "DEVOPIN_RESOURCE_ALERT_DISK_MAX_PERCENT")
	_ = viper.BindEnv("notify.telegram.bot_token", "DEVOPIN_TELEGRAM_BOT_TOKEN")
	_ = viper.BindEnv("notify.telegram.chat_id", "DEVOPIN_TELEGRAM_CHAT_ID")
	_ = viper.BindEnv("server.host", "DEVOPIN_SERVER_HOST")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if isDev() {
		// Dev mode: cari config di root project (lokasi binary dijalankan)
		viper.AddConfigPath(".")
		viper.AddConfigPath("./config")
	} else {
		// Production: ambil dari /etc/devopin/
		viper.AddConfigPath("/etc/devopin/")
	}

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFound viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFound) {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		// Config file tidak ditemukan = ok, lanjut pakai env/defaults
	}

	cfg := &Config{
		ResourceAlert: ResourceAlertConfig{
			Interval: viper.GetDuration("resource_alert.interval"),
			Memory: ThresholdConfig{
				MaxPercent: viper.GetInt("resource_alert.memory.max_percent"),
			},
			CPU: ThresholdConfig{
				MaxPercent: viper.GetInt("resource_alert.cpu.max_percent"),
			},
			Disk: ThresholdConfig{
				MaxPercent: viper.GetInt("resource_alert.disk.max_percent"),
			},
		},
		Notify: NotifyConfig{
			Telegram: TelegramConfig{
				BotToken: viper.GetString("notify.telegram.bot_token"),
				ChatID:   viper.GetInt64("notify.telegram.chat_id"),
			},
		},
		Server: ServerConfig{
			Host: viper.GetString("server.host"),
		},
	}

	if err := validate(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// isDev cek apakah running di mode development
// Prioritas: APP_ENV env var → cek file .dev di project root
func isDev() bool {
	env := os.Getenv("APP_ENV")
	if env != "" {
		return env == "development" || env == "dev"
	}

	// Fallback: kalau ada file .dev atau devopin.yaml di direktori sekarang = dev mode
	_, errDev := os.Stat(".dev")
	_, errYaml := os.Stat("config.yaml")
	return !os.IsNotExist(errDev) || !os.IsNotExist(errYaml)
}

func setDefaults() {
	viper.SetDefault("resource_alert.interval", "30s")
	viper.SetDefault("resource_alert.memory.max_percent", 90)
	viper.SetDefault("resource_alert.cpu.max_percent", 90)
	viper.SetDefault("resource_alert.disk.max_percent", 90)
	viper.SetDefault("notify.telegram.bot_token", "")
	viper.SetDefault("notify.telegram.chat_id", 0)
	viper.SetDefault("notify.telegram.chat_id", 0)
	viper.SetDefault("server.host", getHostName())
}

func getHostName() string {
	path, err := os.Hostname()
	if err != nil {
		return "unknown-host"
	}
	return path
}
func validate(cfg *Config) error {
	if cfg.ResourceAlert.Interval <= 0 {
		return errors.New("resource_alert.interval must be > 0")
	}

	if cfg.ResourceAlert.Memory.MaxPercent <= 0 || cfg.ResourceAlert.Memory.MaxPercent > 100 {
		return errors.New("resource_alert.memory.max_percent must be between 1-100")
	}

	if cfg.ResourceAlert.CPU.MaxPercent <= 0 || cfg.ResourceAlert.CPU.MaxPercent > 100 {
		return errors.New("resource_alert.cpu.max_percent must be between 1-100")
	}

	if cfg.ResourceAlert.Disk.MaxPercent <= 0 || cfg.ResourceAlert.Disk.MaxPercent > 100 {
		return errors.New("resource_alert.disk.max_percent must be between 1-100")
	}

	if cfg.Notify.Telegram.BotToken == "" {
		return errors.New("notify.telegram.bot_token is required")
	}

	if cfg.Notify.Telegram.ChatID == 0 {
		return errors.New("notify.telegram.chat_id is required")
	}

	return nil
}
