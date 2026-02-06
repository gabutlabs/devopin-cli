package resource_alert

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"
)

// SystemStats menyimpan semua statistik sumber daya sistem
type SystemStats struct {
	CPU       CPUStats     `json:"cpu"`
	Memory    MemoryStats  `json:"memory"`
	Disk      DiskStats    `json:"disk"`
	Network   NetworkStats `json:"network"`
	Timestamp time.Time    `json:"timestamp"`
}

// CPUStats menyimpan statistik CPU
type CPUStats struct {
	UsagePercent float64 `json:"usage_percent"`
	Cores        int     `json:"cores"`
	ModelName    string  `json:"model_name"`
}

// MemoryStats menyimpan statistik memori
type MemoryStats struct {
	Total       uint64  `json:"total"`
	Available   uint64  `json:"available"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"used_percent"`
}

// DiskStats menyimpan statistik disk
type DiskStats struct {
	Total       uint64  `json:"total"`
	Free        uint64  `json:"free"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"used_percent"`
	Path        string  `json:"path"`
}

// NetworkStats menyimpan statistik jaringan
type NetworkStats struct {
	BytesSent   uint64 `json:"bytes_sent"`
	BytesRecv   uint64 `json:"bytes_recv"`
	PacketsSent uint64 `json:"packets_sent"`
	PacketsRecv uint64 `json:"packets_recv"`
}

// Collector menyediakan metode untuk mengumpulkan metrik sistem
type Collector struct{}

// NewCollector membuat instance baru dari system metrics collector
func collector() *Collector {
	return &Collector{}
}

// GetCPUUsage mengambil statistik penggunaan CPU
func (c *Collector) GetCPUUsage() (CPUStats, error) {
	var cpuStats CPUStats

	// Mengambil persentase penggunaan CPU secara keseluruhan selama 1 detik
	percentages, err := cpu.Percent(time.Second, false)
	if err != nil {
		return cpuStats, err
	}
	if len(percentages) > 0 {
		cpuStats.UsagePercent = percentages[0]
	}

	// Mengambil informasi CPU (jumlah core dan model)
	info, err := cpu.Info()
	if err != nil {
		return cpuStats, err
	}
	if len(info) > 0 {
		cpuStats.Cores = int(info[0].Cores)
		cpuStats.ModelName = info[0].ModelName
	}

	return cpuStats, nil
}

// GetMemoryUsage mengambil statistik penggunaan memori
func (c *Collector) GetMemoryUsage() (MemoryStats, error) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return MemoryStats{}, err
	}

	return MemoryStats{
		Total:       vmStat.Total,
		Available:   vmStat.Available,
		Used:        vmStat.Used,
		UsedPercent: vmStat.UsedPercent,
	}, nil
}

// GetDiskUsage mengambil statistik penggunaan disk untuk path tertentu
func (c *Collector) GetDiskUsage(path string) (DiskStats, error) {
	usage, err := disk.Usage(path)
	if err != nil {
		return DiskStats{}, err
	}

	return DiskStats{
		Total:       usage.Total,
		Free:        usage.Free,
		Used:        usage.Used,
		UsedPercent: usage.UsedPercent,
		Path:        path,
	}, nil
}

// GetNetworkUsage mengambil statistik penggunaan jaringan (agregat semua interface)
func (c *Collector) GetNetworkUsage() (NetworkStats, error) {
	// 'false' berarti mengambil total dari semua interface jaringan
	ioCounters, err := net.IOCounters(false)
	if err != nil {
		return NetworkStats{}, err
	}

	if len(ioCounters) == 0 {
		return NetworkStats{}, fmt.Errorf("tidak ada network interface yang ditemukan")
	}

	// Ambil statistik agregat pertama
	stats := ioCounters[0]
	return NetworkStats{
		BytesSent:   stats.BytesSent,
		BytesRecv:   stats.BytesRecv,
		PacketsSent: stats.PacketsSent,
		PacketsRecv: stats.PacketsRecv,
	}, nil
}

// GetSystemStats mengambil semua statistik sumber daya sistem secara keseluruhan
func (c *Collector) GetSystemStats() (*SystemStats, error) {
	cpuStats, err := c.GetCPUUsage()
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data CPU: %w", err)
	}

	memoryStats, err := c.GetMemoryUsage()
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data memori: %w", err)
	}

	// Mengambil data disk untuk root path "/"
	diskStats, err := c.GetDiskUsage("/")
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data disk: %w", err)
	}

	networkStats, err := c.GetNetworkUsage()
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data jaringan: %w", err)
	}

	stats := &SystemStats{
		CPU:       cpuStats,
		Memory:    memoryStats,
		Disk:      diskStats,
		Network:   networkStats,
		Timestamp: time.Now(),
	}

	return stats, nil
}
