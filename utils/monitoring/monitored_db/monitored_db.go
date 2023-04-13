package monitoreddb

import (
	"database/sql"

	"github.com/dlmiddlecote/sqlstats"
	"github.com/prometheus/client_golang/prometheus"
)

type PrometheusDbMetricCollector struct {
	db   *sql.DB
	name string
}

func NewPrometheusDbMetricsCollector(name string, db *sql.DB) *PrometheusDbMetricCollector {
	return &PrometheusDbMetricCollector{
		db:   db,
		name: name,
	}
}

func (mc *PrometheusDbMetricCollector) Register() {
	// Create a new collector, the name will be used as a label on the metrics
	collector := sqlstats.NewStatsCollector(mc.name, mc.db)

	// Register it with Prometheus
	prometheus.MustRegister(collector)
}
