package monitoring

import (
	"context"
	"time"

	"github.com/block-wallet/golang-service-template/utils/monitoring/histogram"
	gogrpc "google.golang.org/grpc"
)

type GRPCRequestLatencyMetricInterceptor struct {
	requestLatencyMetricSender histogram.RequestLatencyMetricSender
}

func NewGRPCRequestLatencyMetricInterceptor(
	requestLatencyMetricSender histogram.RequestLatencyMetricSender) *GRPCRequestLatencyMetricInterceptor {
	return &GRPCRequestLatencyMetricInterceptor{requestLatencyMetricSender: requestLatencyMetricSender}
}

func (g *GRPCRequestLatencyMetricInterceptor) UnaryInterceptor() gogrpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *gogrpc.UnaryServerInfo,
		handler gogrpc.UnaryHandler) (response interface{}, err error) {
		start := time.Now()
		response, err = handler(ctx, req)
		end := time.Now()

		g.requestLatencyMetricSender.Send(start, end, info.FullMethod, req, err)
		return response, err
	}
}
