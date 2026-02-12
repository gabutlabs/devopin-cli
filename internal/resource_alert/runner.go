package resource_alert

import (
	"context"
	"fmt"
	"gabutlabs/devopin-cli/internal/config"
	"gabutlabs/devopin-cli/internal/notification"
	"time"
)

func checkResourceAlerts(cfg *config.Config, notif *notification.Notification) {
	collect, err := collector().GetSystemStats()
	if err != nil {
		fmt.Println(err)
	}
	// Implementasi runner untuk resource alert
	var message string
	if collect.CPU.UsagePercent > float64(cfg.ResourceAlert.CPU.MaxPercent) {
		fmt.Printf("CPU usage alert! Current usage: %.2f%%\n", collect.CPU.UsagePercent)
		message = notif.FormatResourceAlertMessage(cfg.Server.Host, "CPU", collect.CPU.UsagePercent, cfg.ResourceAlert.CPU.MaxPercent)
	}
	if collect.Memory.UsedPercent > float64(cfg.ResourceAlert.Memory.MaxPercent) {
		fmt.Printf("Memory usage alert! Current usage: %.2f%%\n", collect.Memory.UsedPercent)
		message = notif.FormatResourceAlertMessage(cfg.Server.Host, "Memory", collect.Memory.UsedPercent, cfg.ResourceAlert.Memory.MaxPercent)
	}
	if float64(collect.Disk.UsedPercent) > float64(cfg.ResourceAlert.Disk.MaxPercent) {
		fmt.Printf("Disk usage alert! Current usage: %.2f%%\n", float64((collect.Disk.Used/collect.Memory.Total)*100))
		message = notif.FormatResourceAlertMessage(cfg.Server.Host, "Disk", float64((collect.Disk.Used/collect.Memory.Total)*100), cfg.ResourceAlert.Disk.MaxPercent)
	}
	notif.SendTelegramAlert(message)
}

func RunResourceAlertRunner(ctx context.Context, cfg *config.Config) {
	ticker := time.NewTicker(cfg.ResourceAlert.Interval * time.Minute)
	defer ticker.Stop()
	notif := notification.NewNotification(ctx, cfg)
	checkResourceAlerts(cfg, notif)
	for {
		select {
		case <-ticker.C:
			checkResourceAlerts(cfg, notif)
		case <-ctx.Done():
			fmt.Println("Resource alert runner stopped.")
			return
		}
	}
}
