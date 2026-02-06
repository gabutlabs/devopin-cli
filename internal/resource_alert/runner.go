package resource_alert

import (
	"fmt"
	"gabutlabs/devopin-cli/internal/config"
)

func RunResourceAlertRunner(cfg *config.Config) {
	collect, err := collector().GetSystemStats()
	if err != nil {
		fmt.Println(err)
	}
	// Implementasi runner untuk resource alert
	if collect.CPU.UsagePercent > float64(cfg.ResourceAlert.CPU.MaxPercent) {
		fmt.Printf("CPU usage alert! Current usage: %.2f%%\n", collect.CPU.UsagePercent)
	}
	if collect.Memory.UsedPercent > float64(cfg.ResourceAlert.Memory.MaxPercent) {
		fmt.Printf("Memory usage alert! Current usage: %.2f%%\n", collect.Memory.UsedPercent)
	}
	if float64(collect.Disk.UsedPercent) > float64(cfg.ResourceAlert.Disk.MaxPercent) {
		fmt.Printf("Disk usage alert! Current usage: %.2f%%\n", float64((collect.Disk.Used/collect.Memory.Total)*100))
	}
}
