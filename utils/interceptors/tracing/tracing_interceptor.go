package tracing

import (
	"context"

	"github.com/block-wallet/campaigns-service/utils/logger"
	"github.com/google/uuid"
	gogrpc "google.golang.org/grpc"
)

type Interceptor struct {
	messageIDField logger.ContextKey
}

func NewInterceptor(messageIDField logger.ContextKey) *Interceptor {
	return &Interceptor{
		messageIDField: messageIDField,
	}
}

func (i *Interceptor) UnaryInterceptor() gogrpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *gogrpc.UnaryServerInfo,
		handler gogrpc.UnaryHandler) (response interface{}, err error) {
		messageID := uuid.New().String()
		ctxWithMessageID := context.WithValue(ctx, i.messageIDField, messageID)

		return handler(ctxWithMessageID, req)
	}
}
