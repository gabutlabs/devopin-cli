package command

import (
	"fmt"
	"gabutlabs/devopin-cli/internal/config"
	"gabutlabs/devopin-cli/internal/resource_alert"

	"github.com/spf13/cobra"
)

var resourceAlertCmd = &cobra.Command{
	Use:   "resource-alert",
	Short: "Manage resource alerts",
	Long:  `Monitoring resource like a cpu usage, disk usage and memory usage`,
	Run: func(cmd *cobra.Command, args []string) {
		// Implementation for managing resource alerts goes here
		cfg, err := config.Load()
		if err != nil {
			fmt.Println(err)
		}
		resource_alert.RunResourceAlertRunner(cfg)
	},
}
