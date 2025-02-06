// Example 137
// internal/monitoring/metrics/custom.go
type BusinessMetrics struct {
    orderProcessingTime   *prometheus.HistogramVec
    orderStatusCount     *prometheus.CounterVec
    activeUsers          prometheus.Gauge
    paymentSuccess       *prometheus.CounterVec
    inventoryLevel       *prometheus.GaugeVec
}

func (bm *BusinessMetrics) RecordOrderProcessing(duration time.Duration, status string) {
    bm.orderProcessingTime.WithLabelValues(status).Observe(duration.Seconds())
    bm.orderStatusCount.WithLabelValues(status).Inc()
}

func (bm *BusinessMetrics) TrackInventory(productID string, quantity int) {
    bm.inventoryLevel.WithLabelValues(productID).Set(float64(quantity))
}

// Resource utilization metrics
type ResourceMetrics struct {
    cpuUsage    prometheus.Gauge
    memoryUsage prometheus.Gauge
    goroutines  prometheus.Gauge
    gcDuration  prometheus.Histogram
}

func (rm *ResourceMetrics) CollectRuntimeMetrics() {
    var memStats runtime.MemStats
    runtime.ReadMemStats(&memStats)
    
    rm.memoryUsage.Set(float64(memStats.Alloc))
    rm.goroutines.Set(float64(runtime.NumGoroutine()))
    rm.gcDuration.Observe(float64(memStats.PauseTotalNs) / 1e9)
}