package domain

import (
	"github.com/block-wallet/campaigns-service/utils/grpc"

	campaignsgrpcservice "github.com/block-wallet/campaigns-service/domain/campaigns-service/service/grpc"
	healthgrpcservice "github.com/block-wallet/campaigns-service/domain/health/service/grpc"
)

var HttpServiceEndpointsHandlersFuncs = []grpc.EndpointHandlerFunc{
	healthgrpcservice.HttpEndpointHandlerFunc,
	campaignsgrpcservice.HttpEndpointHandlerFunc,
}
