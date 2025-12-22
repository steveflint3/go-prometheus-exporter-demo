package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

// A gauge metric we update periodically
var randomGauge = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "demo_random_value",
		Help: "A random value updated every second",
	},
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
	// Seed local random number generator so values change each run
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Register the metric so Prometheus knows about it
	prometheus.MustRegister(randomGauge)
	prometheus.MustRegister(cpuUsagePercent)
	prometheus.MustRegister(memUsagePercent)

	// Goroutine: update the gauge every second
	go func() {
		for {
			val := r.Float64() * 100
			randomGauge.Set(val)
			log.Println("Setting random value:", val)

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
