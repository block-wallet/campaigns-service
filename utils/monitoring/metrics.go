package monitoring

import (
	"github.com/block-wallet/campaigns-service/utils/monitoring/counter"
	"github.com/block-wallet/campaigns-service/utils/monitoring/histogram"
)

const NAMESPACE = "campaignsservice"

var (
	defBucketsMs = []float64{5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000}

	GRPCRequestLatencyMetricSender = histogram.NewGRPCLatencyMetricSender(
		histogram.NewPrometheusLatencyMetricSender(
			NAMESPACE,
			"grpc_request_latency_ms",
			"GRPC request latency histogram, labeled by method and status",
			defBucketsMs,
			[]string{MethodLabel, StatusLabel},
		),
	)
	HTTPRequestLatencyMetricSender = histogram.NewPrometheusLatencyMetricSender(
		NAMESPACE,
		"http_request_latency_ms",
		"HTTP request latency histogram, labeled by method and status",
		defBucketsMs,
		[]string{MethodLabel, StatusLabel},
	)
	ServerPanicCounterMetricSender = counter.NewPrometheusCounterMetricSender(
		NAMESPACE,
		"server_panic",
		"Count server panic",
		[]string{},
	)
	RedisLatencyMetricSender = histogram.NewPrometheusLatencyMetricSender(
		NAMESPACE,
		"redis_latency_ms",
		"Redis latency histogram, labeled by method and status",
		defBucketsMs,
		[]string{MethodLabel, StatusLabel},
	)
)
