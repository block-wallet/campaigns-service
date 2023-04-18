package healthgrpcservice

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	campaignsservicev1health "github.com/block-wallet/campaigns-service/protos/src/campaignsservicev1/health"
	grpcServer "github.com/block-wallet/campaigns-service/utils/grpc"
)

func GRPCService() grpcServer.Service {
	return grpcServer.Service{
		RegisterFn: func(server *grpc.Server, self interface{}) {
			campaignsservicev1health.RegisterHealthServer(server, self.(campaignsservicev1health.HealthServer))
		},
		ServiceHandler: NewHandler(),
	}
}

func HttpEndpointHandlerFunc(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return campaignsservicev1health.RegisterHealthHandlerFromEndpoint(ctx, mux, endpoint, opts)
}
