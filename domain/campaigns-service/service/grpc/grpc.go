package campaignsgrpcservice

import (
	"context"

	campaignsservicev1service "github.com/block-wallet/campaigns-service/protos/src/campaignsservicev1/campaigns"

	campaignsconverter "github.com/block-wallet/campaigns-service/domain/campaigns-service/converter"

	campaignsservice "github.com/block-wallet/campaigns-service/domain/campaigns-service/service"
	campaignsvalidator "github.com/block-wallet/campaigns-service/domain/campaigns-service/validator"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	"github.com/block-wallet/campaigns-service/utils/auth"
	grpcServer "github.com/block-wallet/campaigns-service/utils/grpc"
)

type Options struct {
	CampaignsService campaignsservice.Service
	Validator        campaignsvalidator.Validator
	Converter        campaignsconverter.Converter
	Authenticator    auth.Auth
}

func GRPCService(options Options) grpcServer.Service {
	return grpcServer.Service{
		RegisterFn: func(server *grpc.Server, self interface{}) {
			campaignsservicev1service.RegisterCampaignsSerivceServer(server, self.(campaignsservicev1service.CampaignsSerivceServer))
		},
		ServiceHandler: NewHandler(
			options.CampaignsService,
			options.Validator,
			options.Converter,
			options.Authenticator,
		),
	}
}

func HttpEndpointHandlerFunc(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	return campaignsservicev1service.RegisterCampaignsSerivceHandlerFromEndpoint(ctx, mux, endpoint, opts)
}
