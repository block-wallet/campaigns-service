package srv

import (
	"database/sql"
	"fmt"
	nethttp "net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/block-wallet/campaigns-service/domain/campaigns-service/client"
	campaignsconverter "github.com/block-wallet/campaigns-service/domain/campaigns-service/converter"
	campaignsrepository "github.com/block-wallet/campaigns-service/domain/campaigns-service/repository"
	campaignsservice "github.com/block-wallet/campaigns-service/domain/campaigns-service/service"
	campaignsservicevalidator "github.com/block-wallet/campaigns-service/domain/campaigns-service/validator"
	sqldb "github.com/block-wallet/campaigns-service/storage/database/sql"
	"github.com/block-wallet/campaigns-service/utils/auth"
	monitoreddb "github.com/block-wallet/campaigns-service/utils/monitoring/monitored_db"

	"github.com/block-wallet/campaigns-service/utils/interceptors"

	"github.com/block-wallet/campaigns-service/domain"
	campaignsgrpcservice "github.com/block-wallet/campaigns-service/domain/campaigns-service/service/grpc"
	"github.com/block-wallet/campaigns-service/storage/database/config"
	"github.com/block-wallet/campaigns-service/utils/grpc"
	"github.com/block-wallet/campaigns-service/utils/logger"
	"github.com/block-wallet/campaigns-service/utils/monitoring"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	gogrpc "google.golang.org/grpc"
)

type Runnable struct{}

func NewRunnable() *Runnable {
	return &Runnable{}
}

type Args struct {
	LogLevel                string
	Version                 string
	Port                    int
	MetricsPort             int
	DBType                  config.DBType
	SQLConnectionString     string
	CampaignsPublicEndpoint string
}

func (r *Runnable) Cmd(version string) *cobra.Command {
	cmd.Run = func(_ *cobra.Command, _ []string) {
		server := r.Run(Args{
			LogLevel:                logLevelArg,
			Version:                 version,
			Port:                    port,
			MetricsPort:             metricsPort,
			DBType:                  dbType,
			SQLConnectionString:     sqlConnectionString,
			CampaignsPublicEndpoint: "CampaignsPublicEndpoint",
		})
		signal.Notify(server.StopChannel(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		_ = server.Start()
	}
	return cmd
}

func (r *Runnable) Run(args Args) *grpc.Server {
	serviceName := "campaignsservice"
	err := logger.Initialize(args.LogLevel, args.Version)
	if err != nil {
		panic(err)
	}

	var sqlDatabase *sql.DB
	var dbConfig *config.DBConfig = config.NewDBConfig(args.DBType, config.NewSQLConfig(args.SQLConnectionString, skipMigrations))

	sqlDatabase, err = r.getDatabase(dbConfig)

	if err != nil {
		panic(err)
	}

	if err = sqlDatabase.Ping(); err != nil {
		panic(err)
	}

	//register db metrics collector
	metricsCollector := monitoreddb.NewPrometheusDbMetricsCollector("campaigns-db", sqlDatabase)
	metricsCollector.Register()

	galxeClient := client.NewGalxeClient(galxeGraphQLEndpoint, galxeAccessToken)
	repository := campaignsrepository.NewSQLRepository(sqlDatabase)
	campaignsService := campaignsservice.NewServiceImpl(repository, galxeClient)
	campaignsServiceValidator := campaignsservicevalidator.NewRequestValidator()
	campaignsServiceConverter := campaignsconverter.NewConverterImpl()
	authenticator := auth.NewBasicAuth(adminUsername, adminPassword)

	grpcService := campaignsgrpcservice.GRPCService(campaignsgrpcservice.Options{
		CampaignsService: campaignsService,
		Validator:        campaignsServiceValidator,
		Converter:        campaignsServiceConverter,
		Authenticator:    authenticator,
	})

	return grpc.NewServer(
		grpc.ServerOptions{
			Port:        args.Port,
			Name:        serviceName,
			StopTimeout: 3 * time.Second,
			GRPCServices: []grpc.Service{
				grpcService,
			},
			EndpointHandlersFuncs: domain.HttpServiceEndpointsHandlersFuncs,
			GRPCOptions: []gogrpc.ServerOption{
				interceptors.UnaryInterceptors(
					monitoring.ServerPanicCounterMetricSender,
					monitoring.GRPCRequestLatencyMetricSender,
				),
				interceptors.StreamInterceptors(),
			},
			HTTPEndpoints: []grpc.Endpoint{
				{
					Path: "/ready",
					Handler: nethttp.HandlerFunc(func(writer nethttp.ResponseWriter, _ *nethttp.Request) {
						var _, _ = writer.Write([]byte("YES\n"))
					}),
				},
			},
			MetricsOptions: &grpc.MetricsOptions{
				Port: args.MetricsPort,
				MetricsEndpoints: []grpc.Endpoint{
					{Path: "/metrics", Handler: promhttp.Handler()},
					{Path: "/debug/pprof/profile", Handler: pprof.Handler("profile")},
					{Path: "/debug/pprof/trace", Handler: pprof.Handler("trace")},
					{Path: "/debug/pprof/heap", Handler: pprof.Handler("heap")},
				},
			},
			MarshalOptions: &grpc.MarshalOptions{
				EmitDefaults: false,
			},
		},
	)
}

func (r *Runnable) getDatabase(dbConfig *config.DBConfig) (*sql.DB, error) {
	sqlDatabase, err := sqldb.NewSQLDatabase(dbConfig)
	if err != nil {
		return nil, err
	}
	dbChannel := make(chan os.Signal, 1)
	signal.Notify(dbChannel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-dbChannel
		close(dbChannel)
		if sqlDatabase != nil {
			fmt.Printf("👋 Closing database connections...\n")
			sqlDatabase.Close()
		}
	}()
	return sqlDatabase, nil
}
