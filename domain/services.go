package domain

import (
	"github.com/block-wallet/golang-service-template/utils/grpc"

	ethgrpcservice "github.com/block-wallet/golang-service-template/domain/eth-service/service/grpc"
	healthgrpcservice "github.com/block-wallet/golang-service-template/domain/health/service/grpc"
)

var HttpServiceEndpointsHandlersFuncs = []grpc.EndpointHandlerFunc{
	healthgrpcservice.HttpEndpointHandlerFunc,
	ethgrpcservice.HttpEndpointHandlerFunc,
}
