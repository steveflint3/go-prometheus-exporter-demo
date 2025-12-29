package main

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/load"
)

var startTime = time.Now()

var cpuUsagePercent = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "custom_cpu_usage_percent",
		Help: "Host CPU usage percent (0-100)",
	},
)

var memUsagePercent = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "custom_memory_usage_percent",
		Help: "Host memory usage percent (0-100)",
	},
)

var rootDiskUsagePercent = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "custom_root_disk_usage_percent",
		Help: "Root filesystem disk usage percentage (0-100)",
	},
)

var load1 = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "custom_load1",
		Help: "System load average over 1 minute",
	},
)

var exporterUptime = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "custom_exporter_uptime_seconds",
		Help: "Uptime of the custom Go exporter in seconds",
	},
)

func main() {
	// Register the metric so Prometheus knows about it
	prometheus.MustRegister(cpuUsagePercent)
	prometheus.MustRegister(memUsagePercent)
	prometheus.MustRegister(exporterUptime)
	prometheus.MustRegister(rootDiskUsagePercent)
	prometheus.MustRegister(load1)

	go func() {
		for {
			// CPU usage (gopsutil returns % over an interval)
			// interval=0 means "since last call" on some platforms, so we use 1s for consistency
			cpuPercents, err := cpu.Percent(1*time.Second, false)
			if err != nil || len(cpuPercents) == 0 {
				log.Printf("CPU percent error: %v", err)
			} else {
				cpuUsagePercent.Set(cpuPercents[0])
			}

			// Memory usage
			vm, err := mem.VirtualMemory()
			if err != nil {
				log.Printf("Memory percent error: %v", err)
			} else {
				memUsagePercent.Set(vm.UsedPercent)
			}

			du, err := disk.Usage("/")
			if err != nil {
				log.Printf("Disk usage error: %v", err)
			} else {
				rootDiskUsagePercent.Set(du.UsedPercent)
			}

			avg, err := load.Avg()
			if err != nil {
				log.Printf("Load average error: %v", err)
			} else {
				load1.Set(avg.Load1)
			}

			// Exporter uptime — ALWAYS SET
			exporterUptime.Set(time.Since(startTime).Seconds())

			time.Sleep(5 * time.Second)
		}
	}()

	// Expose /metrics endpoint
	http.Handle("/metrics", promhttp.Handler())

	log.Println("Exporter running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
