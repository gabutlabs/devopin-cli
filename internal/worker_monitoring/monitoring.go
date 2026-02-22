package worker_monitoring

import (
	"context"
	"gabutlabs/devopin-cli/internal/config"
	"gabutlabs/devopin-cli/internal/utils"
	"time"
)

func RunWorkerMonitoring(ctx context.Context, cfg *config.Config) {

	interval := cfg.MonitorWorker.Interval * time.Minute
	workerInit := NewWorkerMonitoring(ctx, cfg)
	workerInit.Monitoring()
	utils.RunWithTicker(ctx, interval, func() {
		workerInit.Monitoring()
	})
}
