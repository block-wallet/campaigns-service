package ethgrpcservice

import (
	"context"

	ethservicev1service "github.com/block-wallet/golang-service-template/protos/ethservicev1/src/eth"

	ethserviceconverter "github.com/block-wallet/golang-service-template/domain/eth-service/converter"

	ethservice "github.com/block-wallet/golang-service-template/domain/eth-service/service"
	ethservicevalidator "github.com/block-wallet/golang-service-template/domain/eth-service/validator"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"

	grpcServer "github.com/block-wallet/golang-service-template/utils/grpc"
)

type Options struct {
	ETHService ethservice.Service
	Validator  ethservicevalidator.Validator
	Converter  ethserviceconverter.Converter
}

func GRPCService(options Options) grpcServer.Service {
	return grpcServer.Service{
		RegisterFn: func(server *grpc.Server, self interface{}) {
			ethservicev1service.RegisterETHServiceServer(server, self.(ethservicev1service.ETHServiceServer))
		},
		ServiceHandler: NewHandler(
			options.ETHService,
			options.Validator,
			options.Converter,
		),
	}
}

func HttpEndpointHandlerFunc(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return ethservicev1service.RegisterETHServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
}
