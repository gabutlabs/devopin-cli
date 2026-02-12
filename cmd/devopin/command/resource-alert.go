package command

import (
	"context"
	"fmt"
	"gabutlabs/devopin-cli/internal/config"
	"gabutlabs/devopin-cli/internal/resource_alert"
	"os"
	"os/signal"
	"syscall"

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
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go resource_alert.RunResourceAlertRunner(ctx, cfg)
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		fmt.Println("Shutting down resource alert command...")
	},
}
