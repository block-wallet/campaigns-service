package panic

import (
	"context"

	"github.com/block-wallet/golang-service-template/utils/logger"
	"github.com/block-wallet/golang-service-template/utils/monitoring/counter"
	gogrpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Interceptor struct {
	counterMetricSender counter.MetricSender
}

func NewInterceptor(counterMetricSender counter.MetricSender) *Interceptor {
	return &Interceptor{
		counterMetricSender: counterMetricSender,
	}
}

func (i *Interceptor) UnaryInterceptor() gogrpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *gogrpc.UnaryServerInfo,
		handler gogrpc.UnaryHandler) (response interface{}, err error) {
		panicked := true

		defer func() {
			if r := recover(); r != nil || panicked {
				logger.Sugar.Errorf("server panic: %v from request %+v", r, req)
				i.counterMetricSender.Send(nil)
				response = nil
				err = status.Errorf(codes.Internal, "server panic: %v", r)
			}
		}()
		response, err = handler(ctx, req)
		panicked = false
		return response, err
	}
}
