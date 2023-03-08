package healthgrpcservice

import (
	"context"

	"github.com/block-wallet/golang-service-template/utils/logger"
	"github.com/golang/protobuf/ptypes/empty"

	ethservicev1health "github.com/block-wallet/golang-service-template/protos/ethservicev1/src/health"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Status(ctx context.Context, _ *empty.Empty) (*ethservicev1health.StatusReply, error) {
	logger.Sugar.WithCtx(ctx).Debug("Status request received")
	return &ethservicev1health.StatusReply{
		Status: ethservicev1health.HealthStatus_HEALTH_STATUS_ALIVE,
	}, nil
}
