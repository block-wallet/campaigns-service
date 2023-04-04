package counter

import (
	"github.com/block-wallet/campaigns-service/utils/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type PrometheusCounterMetricSender struct {
	name       string
	counterVec *prometheus.CounterVec
}

func NewPrometheusCounterMetricSender(namespace, name, help string, labelNames []string) *PrometheusCounterMetricSender {
	return &PrometheusCounterMetricSender{
		counterVec: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      name,
				Help:      help,
			},
			labelNames,
		),
	}
}

func (p *PrometheusCounterMetricSender) Send(labels map[string]string) {
	counter, err := p.counterVec.GetMetricWith(labels)
	if err != nil {
		logger.Sugar.Errorf("Error sending metric %s: %s, labels: %v", p.name, err.Error(), labels)
		return
	}
	counter.Inc()
}
