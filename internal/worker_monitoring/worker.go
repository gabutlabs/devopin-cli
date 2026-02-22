package worker_monitoring

import (
	"context"
	"fmt"
	"gabutlabs/devopin-cli/internal/config"
	"gabutlabs/devopin-cli/internal/notification"
	"strings"

	"github.com/coreos/go-systemd/v22/dbus"
)

type DiscoveredService struct {
	Name   string
	Status string // Raw status: "active", "started", "failed", "error", etc.
}

type WorkerMonitoring struct {
	ctx context.Context
	cfg *config.Config
}

func NewWorkerMonitoring(ctx context.Context, cfg *config.Config) *WorkerMonitoring {
	return &WorkerMonitoring{
		ctx: ctx,
		cfg: cfg,
	}
}

func (wm *WorkerMonitoring) discoverWorker() ([]DiscoveredService, error) {
	commonServices := map[string]bool{"nginx.service": true, "postgresql.service": true}

	conn, err := dbus.NewSystemConnectionContext(wm.ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	units, err := conn.ListUnitsByPatternsContext(wm.ctx, []string{}, []string{"*.service"})
	if err != nil {
		return nil, err
	}
	var discoveredServices []DiscoveredService
	found := make(map[string]bool)
	for _, unit := range units {
		if wm.contains(wm.cfg.MonitorWorker.ExcludeWorkers, unit.Name) {
			continue
		}
		properties, _ := conn.GetUnitPropertiesContext(wm.ctx, unit.Name)
		fragmentPath, _ := properties["FragmentPath"].(string)
		isCustomService := strings.HasPrefix(fragmentPath, "/etc/systemd/system/")
		isCommonService := commonServices[unit.Name]
		if (isCustomService || isCommonService) && !found[unit.Name] {
			discoveredServices = append(discoveredServices, DiscoveredService{
				Name: unit.Name, Status: unit.ActiveState,
			})
			found[unit.Name] = true
		}
	}
	return discoveredServices, nil
}

func (wm *WorkerMonitoring) Monitoring() {
	workers, err := wm.discoverWorker()
	if err != nil {
		panic(err)
	}
	notif := notification.NewNotification(wm.ctx, wm.cfg)
	var inactiveWorkers []string
	for _, w := range workers {
		if w.Status != "active" {
			inactiveWorkers = append(inactiveWorkers, w.Name)
		}
	}
	msg := notif.FormatMonitorWorkerAlertMessage(wm.cfg.Server.Host, inactiveWorkers)
	notif.SendTelegramAlert(msg)
}
func (wm *WorkerMonitoring) contains(list []string, item string) bool {
	for _, v := range list {
		if strings.Contains(item, ".service") {
			if fmt.Sprintf("%s.service", v) == item {
				return true
			}
		}
		if v == item {
			return true
		}
	}
	return false
}
