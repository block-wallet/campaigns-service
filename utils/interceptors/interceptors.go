package interceptors

import (
	"github.com/block-wallet/campaigns-service/utils/interceptors/log"
	"github.com/block-wallet/campaigns-service/utils/interceptors/monitoring"
	"github.com/block-wallet/campaigns-service/utils/interceptors/panic"
	"github.com/block-wallet/campaigns-service/utils/monitoring/counter"
	"github.com/block-wallet/campaigns-service/utils/monitoring/histogram"
	gogrpc "google.golang.org/grpc"
)

func UnaryInterceptors(
	serverPanicCounterMetricSender counter.MetricSender,
	grpcRequestLatencyMetricSender histogram.RequestLatencyMetricSender) gogrpc.ServerOption {
	return gogrpc.ChainUnaryInterceptor(
		panic.NewInterceptor(serverPanicCounterMetricSender).UnaryInterceptor(),
		monitoring.NewGRPCRequestLatencyMetricInterceptor(grpcRequestLatencyMetricSender).UnaryInterceptor(),
		log.NewInterceptor().UnaryInterceptor(),
	)
}

func StreamInterceptors() gogrpc.ServerOption {
	return gogrpc.ChainUnaryInterceptor()
}
