package log

import (
	"context"
	"time"

	"github.com/block-wallet/golang-service-template/utils/logger"
	gogrpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
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
		duration := time.Since(start)

		statusCode := codes.OK.String()
		if err != nil {
			statusCode = status.Code(err).String()
		}

		reqLogger := logger.Sugar.WithCtx(
			ctx,
			"method", info.FullMethod,
			"duration", duration.String(),
			"status_code", statusCode,
		)

		if err != nil {
			reqLogger.Errorf("Incoming request: %v, Request Response: ERROR, Err: %v",
				req, err)
		} else {
			reqLogger.Debugf("Incoming request: %v, Request Response: OK",
				req)
		}

		return response, err
	}
}
