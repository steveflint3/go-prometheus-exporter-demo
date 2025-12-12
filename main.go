package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// A gauge metric we update periodically
var randomGauge = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "demo_random_value",
		Help: "A random value updated every second",
	},
)

func main() {
	// Seed the RNG so values change each run
	rand.Seed(time.Now().UnixNano())

	// Register the metric so Prometheus knows about it
	prometheus.MustRegister(randomGauge)

	// Goroutine: update the gauge every second
	go func() {
		for {
			val := rand.Float64() * 100
			randomGauge.Set(val)
			fmt.Println("Setting random value:", val)
			time.Sleep(5 * time.Second)
		}
	}()

	// Expose /metrics endpoint
	http.Handle("/metrics", promhttp.Handler())

	fmt.Println("Exporter running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting HTTP server:", err)
	}
}
