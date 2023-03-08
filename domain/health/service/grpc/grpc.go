package healthgrpcservice

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"

	ethservicev1health "github.com/block-wallet/golang-service-template/protos/ethservicev1/src/health"
	grpcServer "github.com/block-wallet/golang-service-template/utils/grpc"
)

func GRPCService() grpcServer.Service {
	return grpcServer.Service{
		RegisterFn: func(server *grpc.Server, self interface{}) {
			ethservicev1health.RegisterHealthServer(server, self.(ethservicev1health.HealthServer))
		},
		ServiceHandler: NewHandler(),
	}
}

func HttpEndpointHandlerFunc(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return ethservicev1health.RegisterHealthHandlerFromEndpoint(ctx, mux, endpoint, opts)
}
