package histogram

import (
	"time"

	"github.com/block-wallet/campaigns-service/utils/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type PrometheusLatencyMetricSender struct {
	name         string
	histogramVec *prometheus.HistogramVec
}

func NewPrometheusLatencyMetricSender(namespace, name, help string, buckets []float64, labelNames []string) *PrometheusLatencyMetricSender {
	return &PrometheusLatencyMetricSender{
		name: name,
		histogramVec: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      name,
				Help:      help,
				Buckets:   buckets,
			},
			labelNames,
		),
	}
}

func (p *PrometheusLatencyMetricSender) Send(start, end time.Time, labels map[string]string) {
	observer, err := p.histogramVec.GetMetricWith(labels)
	if err != nil {
		logger.Sugar.Errorf("Error sending metric %s: %s, labels: %v", p.name, err.Error(), labels)
		return
	}
	observer.Observe(float64(end.Sub(start).Milliseconds()))
}
