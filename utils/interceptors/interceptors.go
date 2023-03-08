package interceptors

import (
	"github.com/block-wallet/golang-service-template/utils/interceptors/log"
	"github.com/block-wallet/golang-service-template/utils/interceptors/monitoring"
	"github.com/block-wallet/golang-service-template/utils/interceptors/panic"
	"github.com/block-wallet/golang-service-template/utils/interceptors/tracing"
	"github.com/block-wallet/golang-service-template/utils/logger"
	"github.com/block-wallet/golang-service-template/utils/monitoring/counter"
	"github.com/block-wallet/golang-service-template/utils/monitoring/histogram"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	gogrpc "google.golang.org/grpc"
)

func UnaryInterceptors(messageIDField logger.ContextKey,
	serverPanicCounterMetricSender counter.MetricSender,
	grpcRequestLatencyMetricSender histogram.RequestLatencyMetricSender) gogrpc.ServerOption {
	return grpcmiddleware.WithUnaryServerChain(
		panic.NewInterceptor(serverPanicCounterMetricSender).UnaryInterceptor(),
		monitoring.NewGRPCRequestLatencyMetricInterceptor(grpcRequestLatencyMetricSender).UnaryInterceptor(),
		tracing.NewInterceptor(messageIDField).UnaryInterceptor(),
		log.NewInterceptor().UnaryInterceptor(),
	)
}

func StreamInterceptors() gogrpc.ServerOption {
	return grpcmiddleware.WithStreamServerChain()
}
