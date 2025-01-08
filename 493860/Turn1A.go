// main.go
package main

import (
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"log"
	"time"
)

func main() {
	flag.Parse()
	networkSetup()
	startBenchmark()
	prometheus.StartServer(8080)
	prometheus.WaitForExit()
}

// Define metrics
var latencyMetric = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: "network",
	Name:      "latency_ms",
	Help:      "Latency in milliseconds",
}, []string{"test_id"})

var throughputMetric = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: "network",
	Name:      "throughput_mbps",
	Help:      "Throughput in megabits per second",
}, []string{"test_id"})

func initMetrics() {
	// Register metrics
	prometheus.MustRegister(latencyMetric)
	prometheus.MustRegister(throughputMetric)
	prometheus.MustRegister(collectors.NewBuildInfoCollector())
	prometheus.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
}

func startBenchmark() {
	client := netperf.Client{}
	client.Addr = "example.com"
	client.Port = 12345

	var runs int
	for runs < 10 {
		tests, err := netperf.AvailableTests(&client)
		if err != nil {
			log.Fatalf("Failed to retrieve available tests: %v", err)
		}

		for _, t := range tests {
			result, err := netperf.RunTest(&client, t, netperf.TestOptions{})
			if err != nil {
				log.Fatalf("Failed to run test %d: %v", t, err)
			}
			fmt.Printf("Test %d: %#v\n", t, result)
		}
		time.Sleep(5 * time.Second)
		runs++
	}
}
