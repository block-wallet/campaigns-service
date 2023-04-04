package interceptors

import (
	"github.com/block-wallet/campaigns-service/utils/interceptors/log"
	"github.com/block-wallet/campaigns-service/utils/interceptors/monitoring"
	"github.com/block-wallet/campaigns-service/utils/interceptors/panic"
	"github.com/block-wallet/campaigns-service/utils/monitoring/counter"
	"github.com/block-wallet/campaigns-service/utils/monitoring/histogram"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	gogrpc "google.golang.org/grpc"
)

func UnaryInterceptors(
	serverPanicCounterMetricSender counter.MetricSender,
	grpcRequestLatencyMetricSender histogram.RequestLatencyMetricSender) gogrpc.ServerOption {
	return grpcmiddleware.WithUnaryServerChain(
		panic.NewInterceptor(serverPanicCounterMetricSender).UnaryInterceptor(),
		monitoring.NewGRPCRequestLatencyMetricInterceptor(grpcRequestLatencyMetricSender).UnaryInterceptor(),
		log.NewInterceptor().UnaryInterceptor(),
	)
}

func StreamInterceptors() gogrpc.ServerOption {
	return grpcmiddleware.WithStreamServerChain()
}
