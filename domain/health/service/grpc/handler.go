package healthgrpcservice

import (
	"context"

	"github.com/block-wallet/campaigns-service/utils/logger"
	"github.com/golang/protobuf/ptypes/empty"

	campaignsservicev1health "github.com/block-wallet/campaigns-service/protos/src/campaignsservicev1/health"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Status(ctx context.Context, _ *empty.Empty) (*campaignsservicev1health.StatusReply, error) {
	logger.Sugar.WithCtx(ctx).Debug("Status request received")
	return &campaignsservicev1health.StatusReply{
		Status: campaignsservicev1health.HealthStatus_HEALTH_STATUS_ALIVE,
	}, nil
}
