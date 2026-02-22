package utils

import (
	"context"
	"time"
)

func RunWithTicker(ctx context.Context, interval time.Duration, job func()) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Run immediately first time
	job()

	for {
		select {
		case <-ticker.C:
			job()
		case <-ctx.Done():
			return
		}
	}
}
