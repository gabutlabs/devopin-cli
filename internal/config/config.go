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
	MonitorWorker MonitorWorkerConfig
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

type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	FromEmail    string
	ToEmails     []string
}

type NotifyConfig struct {
	Telegram TelegramConfig
	Email    EmailConfig
	Channels NotifyChannelsConfig
}

type NotifyChannelsConfig struct {
	Telegram bool
	Email    bool
}

type ServerConfig struct {
	Host string
}

type MonitorWorkerConfig struct {
	Interval       time.Duration
	ExcludeWorkers []string
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
	_ = viper.BindEnv("notify.email.smtp_host", "DEVOPIN_EMAIL_SMTP_HOST")
	_ = viper.BindEnv("notify.email.smtp_port", "DEVOPIN_EMAIL_SMTP_PORT")
	_ = viper.BindEnv("notify.email.smtp_user", "DEVOPIN_EMAIL_SMTP_USER")
	_ = viper.BindEnv("notify.email.smtp_password", "DEVOPIN_EMAIL_SMTP_PASSWORD")
	_ = viper.BindEnv("notify.email.from_email", "DEVOPIN_EMAIL_FROM_EMAIL")
	_ = viper.BindEnv("notify.email.to_emails", "DEVOPIN_EMAIL_TO_EMAILS")
	_ = viper.BindEnv("notify.channels.telegram", "DEVOPIN_NOTIFY_CHANNELS_TELEGRAM")
	_ = viper.BindEnv("notify.channels.email", "DEVOPIN_NOTIFY_CHANNELS_EMAIL")
	_ = viper.BindEnv("server.host", "DEVOPIN_SERVER_HOST")
	_ = viper.BindEnv("monitor_worker.interval", "DEVOPIN_MONITOR_WORKER_INTERVAL")
	_ = viper.BindEnv("monitor_worker.exclude_workers", "DEVOPIN_MONITOR_WORKER_EXCLUDE_WORKERS")

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
			Email: EmailConfig{
				SMTPHost:     viper.GetString("notify.email.smtp_host"),
				SMTPPort:     viper.GetInt("notify.email.smtp_port"),
				SMTPUser:     viper.GetString("notify.email.smtp_user"),
				SMTPPassword: viper.GetString("notify.email.smtp_password"),
				FromEmail:    viper.GetString("notify.email.from_email"),
				ToEmails:     viper.GetStringSlice("notify.email.to_emails"),
			},
			Channels: NotifyChannelsConfig{
				Telegram: viper.GetBool("notify.channels.telegram"),
				Email:    viper.GetBool("notify.channels.email"),
			},
		},
		Server: ServerConfig{
			Host: viper.GetString("server.host"),
		},
		MonitorWorker: MonitorWorkerConfig{
			Interval:       viper.GetDuration("monitor_worker.interval"),
			ExcludeWorkers: viper.GetStringSlice("monitor_worker.exclude_workers"),
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
	viper.SetDefault("resource_alert.interval", 1)
	viper.SetDefault("resource_alert.memory.max_percent", 90)
	viper.SetDefault("resource_alert.cpu.max_percent", 90)
	viper.SetDefault("resource_alert.disk.max_percent", 90)
	viper.SetDefault("notify.telegram.bot_token", "")
	viper.SetDefault("notify.telegram.chat_id", 0)
	viper.SetDefault("notify.telegram.chat_id", 0)
	viper.SetDefault("notify.email.smtp_host", "")
	viper.SetDefault("notify.email.smtp_port", 587)
	viper.SetDefault("notify.email.smtp_user", "")
	viper.SetDefault("notify.email.smtp_password", "")
	viper.SetDefault("notify.email.from_email", "")
	viper.SetDefault("notify.email.to_emails", []string{})
	viper.SetDefault("notify.channels.telegram", true)
	viper.SetDefault("notify.channels.email", false)
	viper.SetDefault("server.host", getHostName())
	viper.SetDefault("monitor_worker.interval", 1)
	viper.SetDefault("monitor_worker.exclude_workers", []string{"ua-auto-attach", "syslog"})

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

	// Validate notification channels
	if !cfg.Notify.Channels.Telegram && !cfg.Notify.Channels.Email {
		return errors.New("at least one notification channel (telegram or email) must be enabled")
	}

	if cfg.Notify.Channels.Telegram {
		if cfg.Notify.Telegram.BotToken == "" {
			return errors.New("notify.telegram.bot_token is required when telegram channel is enabled")
		}

		if cfg.Notify.Telegram.ChatID == 0 {
			return errors.New("notify.telegram.chat_id is required when telegram channel is enabled")
		}
	}

	if cfg.Notify.Channels.Email {
		if cfg.Notify.Email.SMTPHost == "" {
			return errors.New("notify.email.smtp_host is required when email channel is enabled")
		}

		if cfg.Notify.Email.SMTPPort <= 0 || cfg.Notify.Email.SMTPPort > 65535 {
			return errors.New("notify.email.smtp_port must be between 1-65535 when email channel is enabled")
		}

		if cfg.Notify.Email.SMTPUser == "" {
			return errors.New("notify.email.smtp_user is required when email channel is enabled")
		}

		if cfg.Notify.Email.SMTPPassword == "" {
			return errors.New("notify.email.smtp_password is required when email channel is enabled")
		}

		if cfg.Notify.Email.FromEmail == "" {
			return errors.New("notify.email.from_email is required when email channel is enabled")
		}

		if len(cfg.Notify.Email.ToEmails) == 0 {
			return errors.New("notify.email.to_emails is required when email channel is enabled")
		}
	}

	if cfg.MonitorWorker.Interval <= 0 {
		return errors.New("monitor_worker.interval must be > 0")
	}

	return nil
}
