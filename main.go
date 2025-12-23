package main

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

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

func main() {
	// Register the metric so Prometheus knows about it
	prometheus.MustRegister(cpuUsagePercent)
	prometheus.MustRegister(memUsagePercent)

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

			time.Sleep(5 * time.Second)
		}
	}()

	// Expose /metrics endpoint
	http.Handle("/metrics", promhttp.Handler())

	log.Println("Exporter running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
