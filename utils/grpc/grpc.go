package grpc

import (
	"context"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"google.golang.org/grpc"
)

type Service struct {
	RegisterFn     func(*grpc.Server, interface{})
	ServiceHandler interface{}
}

type EndpointHandlerFunc func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error

type Endpoint struct {
	Path    string
	Handler http.Handler
}

type MetricsOptions struct {
	Port             int
	MetricsEndpoints []Endpoint
}

type MarshalOptions struct {
	OrigName     bool
	EmitDefaults bool
}

type ServerOptions struct {
	Port                  int
	Name                  string
	StopTimeout           time.Duration
	HTTPEndpoints         []Endpoint
	MetricsOptions        *MetricsOptions
	GRPCOptions           []grpc.ServerOption
	GRPCServices          []Service
	EndpointHandlersFuncs []EndpointHandlerFunc
	MarshalOptions        *MarshalOptions
}
