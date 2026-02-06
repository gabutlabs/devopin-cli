package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	ResourceAlert ResourceAlertConfig
}

type ResourceAlertConfig struct {
	Interval time.Duration
	Memory   MemoryConfig
	CPU      MemoryConfig
	Disk     MemoryConfig
	Notify   NotifyConfig
}

type MemoryConfig struct {
	MaxPercent int
}

type NotifyConfig struct {
	TelegramWebhook string
}

func Load() (*Config, error) {
	setDefaults()

	// ENV support
	viper.SetEnvPrefix("DEVOPIN")
	viper.AutomaticEnv()

	// YAML (optional)
	viper.SetConfigName("devopin")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/devopin")

	if err := viper.ReadInConfig(); err != nil {
		// yaml not found = dev mode (OK)
		var configFileNotFound viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFound) {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	cfg := &Config{
		ResourceAlert: ResourceAlertConfig{
			Interval: viper.GetDuration("resource_alert.interval"),
			Memory: MemoryConfig{
				MaxPercent: viper.GetInt("resource_alert.memory.max_percent"),
			},
			CPU: MemoryConfig{
				MaxPercent: viper.GetInt("resource_alert.cpu.max_percent"),
			},
			Disk: MemoryConfig{
				MaxPercent: viper.GetInt("resource_alert.disk.max_percent"),
			},
			Notify: NotifyConfig{
				TelegramWebhook: viper.GetString("resource_alert.notify.telegram_webhook"),
			},
		},
	}

	if err := validate(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func setDefaults() {
	viper.SetDefault("resource_alert.interval", "30s")
	viper.SetDefault("resource_alert.memory.max_percent", 90)
	viper.SetDefault("resource_alert.cpu.max_percent", 90)
	viper.SetDefault("resource_alert.disk.max_percent", 90)
	viper.SetDefault("resource_alert.notify.telegram_webhook", "http://telegram.webhook.url")
}

func validate(cfg *Config) error {
	if cfg.ResourceAlert.Interval <= 0 {
		return errors.New("resource_alert.interval must be > 0")
	}

	if cfg.ResourceAlert.Memory.MaxPercent <= 0 || cfg.ResourceAlert.Memory.MaxPercent > 100 {
		return errors.New("resource_alert.memory.max_percent must be between 1-100")
	}

	if cfg.ResourceAlert.Notify.TelegramWebhook == "" {
		return errors.New("resource_alert.notify.telegram_webhook is required")
	}

	return nil
}
