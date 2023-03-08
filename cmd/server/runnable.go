package server

import (
	"context"
	nethttp "net/http"
	"net/http/pprof"
	"os/signal"
	"syscall"
	"time"

	"github.com/block-wallet/golang-service-template/domain/eth"

	ethserviceconverter "github.com/block-wallet/golang-service-template/domain/eth-service/converter"
	"github.com/block-wallet/golang-service-template/utils/grpc/converter"

	"github.com/block-wallet/golang-service-template/utils/http"

	httpapi "github.com/block-wallet/golang-service-template/domain/http-api"

	"github.com/block-wallet/golang-service-template/utils/interceptors"

	ethgrpcservice "github.com/block-wallet/golang-service-template/domain/eth-service/service/grpc"

	"github.com/block-wallet/golang-service-template/domain"
	ethrepository "github.com/block-wallet/golang-service-template/domain/eth-service/repository"
	ethservice "github.com/block-wallet/golang-service-template/domain/eth-service/service"
	ethservicevalidator "github.com/block-wallet/golang-service-template/domain/eth-service/validator"
	"github.com/block-wallet/golang-service-template/storage/database/config"
	kvdb "github.com/block-wallet/golang-service-template/storage/database/kv"
	"github.com/block-wallet/golang-service-template/utils/grpc"
	"github.com/block-wallet/golang-service-template/utils/logger"
	"github.com/block-wallet/golang-service-template/utils/monitoring"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	gogrpc "google.golang.org/grpc"
)

type Runnable struct{}

func NewRunnable() *Runnable {
	return &Runnable{}
}

type Args struct {
	LogLevel                  string
	Version                   string
	Port                      int
	MetricsPort               int
	KVType                    string
	RedisSentinelHosts        []string
	RedisSentinelMasterName   string
	RedisHost                 string
	RedisPassword             string
	RedisReadOnly             bool
	RedisDB                   int
	LocalCacheExpiration      time.Duration
	LocalCacheCleanUpInterval time.Duration
	SomeHTTPEndpoint          string
	SomeHTTPProtocol          string
	SomeHTTPTimeout           time.Duration
	ETHEndpoint               string
}

func (r *Runnable) Cmd(version string) *cobra.Command {
	cmd.Run = func(_ *cobra.Command, _ []string) {
		server := r.Run(Args{
			LogLevel:                  logLevelArg,
			Version:                   version,
			Port:                      port,
			MetricsPort:               metricsPort,
			KVType:                    kvType,
			RedisSentinelHosts:        redisSentinelHosts,
			RedisSentinelMasterName:   redisSentinelMasterName,
			RedisPassword:             redisPassword,
			RedisDB:                   redisDB,
			RedisReadOnly:             redisReadOnly,
			LocalCacheExpiration:      localCacheExpiration,
			LocalCacheCleanUpInterval: localCacheCleanUpInterval,
			SomeHTTPEndpoint:          someHTTPEndpoint,
			SomeHTTPProtocol:          someHTTPProtocol,
			SomeHTTPTimeout:           someHTTPTimeout,
			ETHEndpoint:               ETHEndpoint,
		})
		signal.Notify(server.StopChannel(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

		_ = server.Start()
	}
	return cmd
}

func (r *Runnable) Run(args Args) *grpc.Server {
	serviceName := "ethservice"
	ctx := context.Background()

	err := logger.Initialize(args.LogLevel, args.Version)
	if err != nil {
		panic(err)
	}

	var kvDatabase kvdb.Database
	var dbConfig *config.DBConfig

	if args.KVType == kvdb.RedisDb {
		dbConfig = config.NewDBConfig(args.KVType, nil, config.NewRedisSentinelConfig(args.RedisSentinelHosts, args.RedisSentinelMasterName, args.RedisPassword, args.RedisDB, args.RedisReadOnly, monitoring.RedisLatencyMetricSender))
	} else {
		dbConfig = config.NewDBConfig(args.KVType, config.NewLocalCacheConfig(args.LocalCacheExpiration, args.LocalCacheCleanUpInterval), nil)
	}

	kvDatabase = kvdb.NewKVDatabase(dbConfig)
	err = kvDatabase.Connect(ctx)
	if err != nil {
		panic(err)
	}

	httpClient := http.NewClientImpl(args.SomeHTTPTimeout, monitoring.HTTPRequestLatencyMetricSender)
	httpApiClient := httpapi.NewApiClientImpl(httpClient, args.SomeHTTPProtocol, args.SomeHTTPEndpoint)
	repository := ethrepository.NewKVRepository(kvDatabase, httpApiClient)
	ethService := ethservice.NewServiceImpl(repository)
	ethServiceValidator := ethservicevalidator.NewRequestValidator()
	ethServiceConverter := ethserviceconverter.NewConverterImpl(converter.NewGRPCConverter())
	subscriptions := eth.NewSubscriptions(ethService, repository, args.ETHEndpoint)

	grpcService := ethgrpcservice.GRPCService(ethgrpcservice.Options{
		ETHService: ethService,
		Validator:  ethServiceValidator,
		Converter:  ethServiceConverter,
	})

	subscriptions.BackgroundStart(ctx)

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
				interceptors.UnaryInterceptors(logger.Sugar.GetMessageIDField(), monitoring.ServerPanicCounterMetricSender, monitoring.GRPCRequestLatencyMetricSender),
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
