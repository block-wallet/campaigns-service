package grpc

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/common-nighthawk/go-figure"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	defaultOrigName     = true
	defaultEmitDefaults = true
)

type Server struct {
	grpcServer      *grpc.Server
	gatewayServer   *http.Server
	metricsServer   *http.Server
	httpHandlerFunc http.HandlerFunc
	options         ServerOptions
	stopChannel     chan os.Signal
	isRunning       bool
}

func NewServer(options ServerOptions) *Server {
	s := &Server{
		grpcServer:  grpc.NewServer(options.GRPCOptions...),
		options:     applyDefaultOptions(options),
		stopChannel: make(chan os.Signal, 1),
		isRunning:   false,
	}

	s.registerMetrics()
	s.registerServices()
	s.registerHttpHandlerFunc()
	return s
}

func applyDefaultOptions(options ServerOptions) ServerOptions {
	if options.MarshalOptions == nil {
		options.MarshalOptions = &MarshalOptions{
			OrigName:     defaultOrigName,
			EmitDefaults: defaultEmitDefaults,
		}
	}

	return options
}

func (s *Server) StopChannel() chan os.Signal {
	return s.stopChannel
}

func (s *Server) registerServices() {
	for _, registerProto := range s.options.GRPCServices {
		registerProto.RegisterFn(s.grpcServer, registerProto.ServiceHandler)
	}
	reflection.Register(s.grpcServer)
}

func (s *Server) registerEndpointHandlersFuncs(endpointsHandlersFuncs []EndpointHandlerFunc) (*runtime.ServeMux, error) {
	if len(endpointsHandlersFuncs) == 0 {
		return nil, nil
	}
	ctx := context.Background()
	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames:   s.options.MarshalOptions.OrigName,
				EmitUnpopulated: true,
			},
		}),
		runtime.WithIncomingHeaderMatcher(runtime.DefaultHeaderMatcher),
		// runtime.WithOutgoingHeaderMatcher(runtime.DefaultHeaderMatcher),
		runtime.WithOutgoingHeaderMatcher(func(key string) (string, bool) {
			return "", false
		}),
	)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	endpoint := fmt.Sprintf(":%d", s.options.Port)
	for _, endpointsHandlersFunc := range endpointsHandlersFuncs {
		if err := endpointsHandlersFunc(ctx, mux, endpoint, opts); err != nil {
			return nil, err
		}
	}
	return mux, nil
}

func (s *Server) registerHttpHandlerFunc() {
	handlersMux, err := s.registerEndpointHandlersFuncs(s.options.EndpointHandlersFuncs)
	if err != nil {
		return
	}

	mux := http.NewServeMux()
	if handlersMux != nil {
		mux.Handle("/", handlersMux)
	}

	for _, endpoint := range s.options.HTTPEndpoints {
		mux.Handle(endpoint.Path, endpoint.Handler)
	}

	if handlersMux != nil {
		s.httpHandlerFunc = func(w http.ResponseWriter, req *http.Request) {
			if req.ProtoMajor == 2 && strings.Contains(req.Header.Get("Content-Type"), "application/grpc") {
				s.grpcServer.ServeHTTP(w, req)
				return
			}
			mux.ServeHTTP(w, req)
		}
	} else {
		s.httpHandlerFunc = func(w http.ResponseWriter, req *http.Request) {
			mux.ServeHTTP(w, req)
		}
	}
}

func (s *Server) registerMetrics() {
	if s.options.MetricsOptions == nil {
		return
	}
	endpoints := s.options.MetricsOptions.MetricsEndpoints
	help := " Metrics endpoints\n"
	help += "---------------------------------------------\n"
	if len(endpoints) == 0 {
		help += "~> <N/A>\n"
	}
	for _, endpoint := range endpoints {
		help += fmt.Sprintf("~> %s\n", endpoint.Path)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, help)
	})
	for _, endpoint := range endpoints {
		mux.Handle(endpoint.Path, endpoint.Handler)
	}
	s.metricsServer = &http.Server{
		Handler: mux,
		Addr:    fmt.Sprintf(":%d", s.options.MetricsOptions.Port),
	}
}

func (s *Server) startMetrics() {
	if s.options.MetricsOptions == nil {
		return
	}
	go func() {
		fmt.Printf("📈 Metrics server is listening at port %d...\n", s.options.MetricsOptions.Port)
		if err := s.metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("💀 Error on metrics server at port %d: '%s'\n", s.options.MetricsOptions.Port, err.Error())
		}
	}()
}

func (s *Server) stopMetrics(ctx context.Context) {
	if s.options.MetricsOptions == nil {
		return
	}
	fmt.Printf("👋 Metrics server is shutting down...\n")
	err := s.metricsServer.Shutdown(ctx)
	if err != nil {
		fmt.Printf("💀 Error while stopping metrics server at %d: '%s'\n", s.options.MetricsOptions.Port, err)
	}
}

func (s *Server) Start() error {
	if s.IsRunning() {
		return nil
	}

	address := fmt.Sprintf(":%d", s.options.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Printf("💀 Error while listening at port %d: '%s'\n", s.options.Port, err)
		return err
	}

	s.isRunning = true

	s.printBanner()
	s.stopOnSignal()
	s.startMetrics()

	if s.options.Port > 0 {
		fmt.Printf("🌐 gRPC server is listening at port %d...\n", s.options.Port)
	}

	if s.httpHandlerFunc == nil {
		return s.grpcServer.Serve(listener)
	}

	s.gatewayServer = &http.Server{
		Handler: h2c.NewHandler(allowCors(s.httpHandlerFunc), &http2.Server{}),
	}

	return s.gatewayServer.Serve(listener)
}

func (s *Server) stop(ctx context.Context) error {
	if s.IsRunning() {
		ctx, cancel := context.WithTimeout(ctx, s.options.StopTimeout)
		defer cancel()

		s.stopMetrics(ctx)

		if s.options.Port > 0 {
			fmt.Printf("👋 gRPC server is shutting down...\n")
		}

		if s.gatewayServer != nil {
			s.isRunning = false
			return s.gatewayServer.Shutdown(ctx)
		} else {
			s.isRunning = false
			s.grpcServer.GracefulStop()
			return nil
		}
	}

	return nil
}

func (s *Server) stopOnSignal() {
	go func() {
		<-s.stopChannel
		close(s.stopChannel)
		var _ = s.stop(context.Background())
	}()
}

func (s *Server) IsRunning() bool {
	return s.isRunning
}

func (s *Server) printBanner() {
	fmt.Println("\033[32m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\u001B[0m")
	figure.NewColorFigure(" "+s.options.Name, "doom", "white", false).Print()
	fmt.Println("")
}
