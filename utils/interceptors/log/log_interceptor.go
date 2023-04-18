package log

import (
	"context"
	"time"

	"github.com/block-wallet/campaigns-service/utils/logger"
	gogrpc "google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type Interceptor struct{}

func NewInterceptor() *Interceptor {
	return &Interceptor{}
}

func (i *Interceptor) UnaryInterceptor() gogrpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *gogrpc.UnaryServerInfo,
		handler gogrpc.UnaryHandler) (response interface{}, err error) {
		start := time.Now()
		response, err = handler(ctx, req)
		if err != nil {
			if status.Code(err) >= 500 {
				duration := time.Since(start)
				statusCode := status.Code(err).String()

				reqLogger := logger.Sugar.WithCtx(
					ctx,
					"method", info.FullMethod,
					"duration", duration.String(),
					"status_code", statusCode,
				)

				reqLogger.Errorf("Incoming request: %v, Request Response: ERROR, Err: %s", req, err.Error())
			}
		}

		return response, err
	}
}
