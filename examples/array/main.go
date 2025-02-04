package main

import (
    "fmt"
    "log"
    "net/http"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "github.com/vzahanych/data-structures/array"  
)

var (
    // Define Prometheus metrics
    arrayAppendCounter = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "array_append_total",
            Help: "Total number of append operations",
        },
        []string{"status"},
    )
    arrayGetCounter = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "array_get_total",
            Help: "Total number of get operations",
        },
        []string{"status"},
    )
    arrayDeleteCounter = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "array_delete_total",
            Help: "Total number of delete operations",
        },
        []string{"status"},
    )
)

func init() {
    // Register Prometheus metrics
    prometheus.MustRegister(arrayAppendCounter)
    prometheus.MustRegister(arrayGetCounter)
    prometheus.MustRegister(arrayDeleteCounter)
}

func main() {
    config := array.ArrayConfig{MetricsEnabled: true}
    arr := array.NewArray[int](1000, config)

    // Simulate operations and increment metrics
    arr.Append(1)
    arrayAppendCounter.WithLabelValues("success").Inc()

    arr.Get(0)
    arrayGetCounter.WithLabelValues("success").Inc()

    arr.Delete(0)
    arrayDeleteCounter.WithLabelValues("success").Inc()

    // Expose Prometheus metrics at /metrics endpoint
    http.Handle("/metrics", promhttp.Handler())
    fmt.Println("Starting server on :8080...")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
