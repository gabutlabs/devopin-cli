package command

import (
	"context"
	"fmt"
	"gabutlabs/devopin-cli/internal/config"
	"gabutlabs/devopin-cli/internal/worker_monitoring"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

var workerMonitoringCmd = &cobra.Command{
	Use:   "monitor-worker",
	Short: "Manage worker run on systemd",
	Long:  `Monitoring worker run on systemd by .service file`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			fmt.Println(err)
		}
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go worker_monitoring.RunWorkerMonitoring(ctx, cfg)
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		fmt.Println("Shutting down monitor worker command...")
	},
}
