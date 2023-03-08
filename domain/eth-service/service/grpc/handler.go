package ethgrpcservice

import (
	"context"

	ethservicev1common "github.com/block-wallet/golang-service-template/protos/ethservicev1/src"
	ethservicev1service "github.com/block-wallet/golang-service-template/protos/ethservicev1/src/eth"

	ethserviceconverter "github.com/block-wallet/golang-service-template/domain/eth-service/converter"

	ethservice "github.com/block-wallet/golang-service-template/domain/eth-service/service"
	ethservicevalidator "github.com/block-wallet/golang-service-template/domain/eth-service/validator"
	"github.com/block-wallet/golang-service-template/utils/logger"
)

type Handler struct {
	service   ethservice.Service
	validator ethservicevalidator.Validator
	converter ethserviceconverter.Converter
}

func NewHandler(service ethservice.Service, validator ethservicevalidator.Validator, converter ethserviceconverter.Converter) *Handler {
	return &Handler{service: service, validator: validator, converter: converter}
}

func (h *Handler) GetEvents(ctx context.Context, req *ethservicev1service.GetEventsMsg) (*ethservicev1service.GetEventsReply, error) {
	logger.Sugar.WithCtx(ctx).Debug("GetEvents received")

	err := h.validator.ValidateGetEventsRequest(req)
	if err != nil {
		logger.Sugar.WithCtx(ctx).Errorf("Error validating GetEvents request: %s - Req: %v", err.Error(), req)
		return nil, err.ToGRPCError()
	}

	modelEvents, err := h.service.GetEvents(ctx, req.GetPair())
	if err != nil {
		logger.Sugar.WithCtx(ctx).Errorf("Error getting events: %s", err.Error())
		return nil, err.ToGRPCError()
	}

	protoEvents := make([]*ethservicev1service.Event, len(*modelEvents))
	for i, modelEvent := range *modelEvents {
		protoEvent, _err := h.converter.ConvertFromModelEventToProtoEvent(&modelEvent)
		if _err != nil {
			logger.Sugar.WithCtx(ctx).Errorf("Error parsing event: %s", err.Error())
			return nil, err
		}
		protoEvents[i] = protoEvent
	}

	return &ethservicev1service.GetEventsReply{Events: protoEvents}, nil
}

func (h *Handler) GetChains(ctx context.Context, _ *ethservicev1common.EmptyMsg) (*ethservicev1service.GetChainsReply, error) {
	logger.Sugar.WithCtx(ctx).Debug("GetChains received")

	modelChains, err := h.service.GetChains(ctx)
	if err != nil {
		logger.Sugar.WithCtx(ctx).Errorf("Error getting chains: %s", err.Error())
		return nil, err.ToGRPCError()
	}

	protoChains := make([]*ethservicev1service.Chain, len(*modelChains))
	for i, modelChain := range *modelChains {
		protoChain, _err := h.converter.ConvertFromModelChainToProtoChain(&modelChain)
		if _err != nil {
			logger.Sugar.WithCtx(ctx).Errorf("Error parsing event: %s", err.Error())
			return nil, err
		}
		protoChains[i] = protoChain
	}

	return &ethservicev1service.GetChainsReply{Chains: protoChains}, nil
}
