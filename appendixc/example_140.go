// Example 140
// internal/monitoring/alerting/manager.go
type AlertManager struct {
    client  *alertmanager.Client
    metrics *MetricsCollector
    logger  *StructuredLogger
}

type AlertRule struct {
    Name        string
    Query       string
    Duration    time.Duration
    Severity    string
    Annotations map[string]string
}

func (am *AlertManager) ConfigureAlerts(rules []AlertRule) error {
    for _, rule := range rules {
        if err := am.client.CreateAlertRule(rule); err != nil {
            return fmt.Errorf("creating alert rule %s: %w", rule.Name, err)
        }
    }
    return nil
}

// Example alert rules
var defaultAlerts = []AlertRule{
    {
        Name:     "HighErrorRate",
        Query:    `rate(error_total[5m]) > 0.1`,
        Duration: 5 * time.Minute,
        Severity: "critical",
        Annotations: map[string]string{
            "summary": "High error rate detected",
            "description": "Error rate exceeded 10% in the last 5 minutes",
        },
    },
    {
        Name:     "HighLatency",
        Query:    `histogram_quantile(0.95, rate(request_duration_seconds_bucket[5m])) > 1`,
        Duration: 5 * time.Minute,
        Severity: "warning",
        Annotations: map[string]string{
            "summary": "High request latency detected",
            "description": "95th percentile latency exceeded 1s in the last 5 minutes",
        },
    },
}