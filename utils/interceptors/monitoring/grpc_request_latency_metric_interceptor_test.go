package monitoring

import (
	"context"
	"testing"

	"github.com/block-wallet/campaigns-service/utils/errors"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/block-wallet/campaigns-service/utils/monitoring/histogram/mocks"
	"github.com/stretchr/testify/assert"
	gogrpc "google.golang.org/grpc"
)

func TestServerMetricInterceptor_Fn(t *testing.T) {
	unaryInfo := &gogrpc.UnaryServerInfo{
		FullMethod: "method_name",
	}
	unaryHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, nil
	}

	requestLatencyMetricSender := &mocks.RequestLatencyMetricSender{}
	requestLatencyMetricSender.On("Send", mock.AnythingOfType("time.Time"),
		mock.AnythingOfType("time.Time"), "method_name", status.Error(codes.OK, ""))
	interceptor := NewGRPCRequestLatencyMetricInterceptor(requestLatencyMetricSender)
	_, err := interceptor.UnaryInterceptor()(context.Background(), empty.Empty{}, unaryInfo, unaryHandler)

	assert.NoError(t, err)
	requestLatencyMetricSender.AssertCalled(t, "Send", mock.AnythingOfType("time.Time"),
		mock.AnythingOfType("time.Time"), "method_name", status.Error(codes.OK, ""))
}

func TestServerMetricInterceptor_Fn_WithError(t *testing.T) {
	err := errors.NewInvalidArgument("invalid").ToGRPCError()
	unaryInfo := &gogrpc.UnaryServerInfo{
		FullMethod: "method_name",
	}
	unaryHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, err
	}

	requestLatencyMetricSender := &mocks.RequestLatencyMetricSender{}
	requestLatencyMetricSender.On("Send", mock.AnythingOfType("time.Time"),
		mock.AnythingOfType("time.Time"), "method_name", status.Error(codes.InvalidArgument, "invalid"))
	interceptor := NewGRPCRequestLatencyMetricInterceptor(requestLatencyMetricSender)
	_, err = interceptor.UnaryInterceptor()(context.Background(), empty.Empty{}, unaryInfo, unaryHandler)

	assert.Error(t, err)
	requestLatencyMetricSender.AssertCalled(t, "Send", mock.AnythingOfType("time.Time"),
		mock.AnythingOfType("time.Time"), "method_name", status.Error(codes.InvalidArgument, "invalid"))
}
