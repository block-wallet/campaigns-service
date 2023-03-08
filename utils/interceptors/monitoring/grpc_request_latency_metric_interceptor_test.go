package monitoring

import (
	"context"
	"testing"

	"github.com/block-wallet/golang-service-template/utils/errors"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/smartystreets/assertions"
	"github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	mocks2 "github.com/block-wallet/golang-service-template/utils/monitoring/histogram/mocks"
	gogrpc "google.golang.org/grpc"
)

func TestServerMetricInterceptor_Fn(t *testing.T) {
	unaryInfo := &gogrpc.UnaryServerInfo{
		FullMethod: "method_name",
	}
	unaryHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, nil
	}
	convey.Convey("Server metric interceptor", t, func() {
		requestLatencyMetricSender := &mocks2.RequestLatencyMetricSender{}
		requestLatencyMetricSender.On("Send", mock.AnythingOfType("time.Time"),
			mock.AnythingOfType("time.Time"), "method_name", empty.Empty{},
			status.Error(codes.OK, ""))
		interceptor := NewGRPCRequestLatencyMetricInterceptor(requestLatencyMetricSender)
		convey.Convey("When I execute the interceptor", func() {
			_, err := interceptor.UnaryInterceptor()(context.Background(), empty.Empty{}, unaryInfo, unaryHandler)
			convey.Convey("Then should send metrics", func() {
				convey.So(err, assertions.ShouldBeNil)
				requestLatencyMetricSender.AssertCalled(t, "Send", mock.AnythingOfType("time.Time"),
					mock.AnythingOfType("time.Time"), "method_name", empty.Empty{},
					status.Error(codes.OK, ""))
			})
		})
	})
}

func TestServerMetricInterceptor_Fn_WithError(t *testing.T) {
	err := errors.NewInvalidArgument("invalid").ToGRPCError()
	unaryInfo := &gogrpc.UnaryServerInfo{
		FullMethod: "method_name",
	}
	unaryHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, err
	}
	convey.Convey("Server metric interceptor", t, func() {
		requestLatencyMetricSender := &mocks2.RequestLatencyMetricSender{}
		requestLatencyMetricSender.On("Send", mock.AnythingOfType("time.Time"),
			mock.AnythingOfType("time.Time"), "method_name", empty.Empty{},
			status.Error(codes.InvalidArgument, "invalid"))
		interceptor := NewGRPCRequestLatencyMetricInterceptor(requestLatencyMetricSender)
		convey.Convey("When I execute the interceptor", func() {
			_, err := interceptor.UnaryInterceptor()(context.Background(), empty.Empty{}, unaryInfo, unaryHandler)
			convey.Convey("Then should send metrics", func() {
				convey.So(err, assertions.ShouldNotBeNil)
				requestLatencyMetricSender.AssertCalled(t, "Send", mock.AnythingOfType("time.Time"),
					mock.AnythingOfType("time.Time"), "method_name", empty.Empty{},
					status.Error(codes.InvalidArgument, "invalid"))
			})
		})
	})
}
