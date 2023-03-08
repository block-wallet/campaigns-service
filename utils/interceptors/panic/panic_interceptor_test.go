package panic

import (
	"context"
	"testing"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/smartystreets/assertions"
	"github.com/smartystreets/goconvey/convey"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	mocks2 "github.com/block-wallet/golang-service-template/utils/monitoring/counter/mocks"
)

func TestPanicInterceptor_RecoverFromPanic(t *testing.T) {
	unaryHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		panic("some error")
	}
	convey.Convey("Panic interceptor", t, func() {
		counterMetricSender := &mocks2.CounterMetricSender{}
		var labels map[string]string
		counterMetricSender.On("Send", labels)
		interceptor := NewInterceptor(counterMetricSender)
		convey.Convey("When I execute the interceptor", func() {
			_, err := interceptor.UnaryInterceptor()(context.Background(), empty.Empty{}, nil, unaryHandler)
			convey.Convey("Then should recover from panic", func() {
				convey.So(err, assertions.ShouldNotBeNil)
				convey.So(err.Error(), assertions.ShouldEqual, "rpc error: code = Internal desc = server panic: some error")
				counterMetricSender.AssertCalled(t, "Send", labels)
			})
		})
	})
}

func TestPanicInterceptor_RecoverFromNilPanic(t *testing.T) {
	unaryHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		panic(nil)
	}
	convey.Convey("Panic interceptor", t, func() {
		counterMetricSender := &mocks2.CounterMetricSender{}
		var labels map[string]string
		counterMetricSender.On("Send", labels)
		interceptor := NewInterceptor(counterMetricSender)
		convey.Convey("When I execute the interceptor", func() {
			_, err := interceptor.UnaryInterceptor()(context.Background(), empty.Empty{}, nil, unaryHandler)
			convey.Convey("Then should recover from panic", func() {
				convey.So(err, assertions.ShouldNotBeNil)
				convey.So(err.Error(), assertions.ShouldEqual, "rpc error: code = Internal desc = server panic: <nil>")
				counterMetricSender.AssertCalled(t, "Send", labels)
			})
		})
	})
}

func TestPanicInterceptor_IgnoreIfNoPanic(t *testing.T) {
	unaryHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "ok response", nil
	}
	convey.Convey("Panic interceptor", t, func() {
		counterMetricSender := &mocks2.CounterMetricSender{}
		interceptor := NewInterceptor(counterMetricSender)
		convey.Convey("When I execute the interceptor", func() {
			resp, err := interceptor.UnaryInterceptor()(context.Background(), empty.Empty{}, nil, unaryHandler)
			convey.Convey("Then should return no error", func() {
				convey.So(err, assertions.ShouldBeNil)
				convey.So(resp.(string), assertions.ShouldEqual, "ok response")
				counterMetricSender.AssertNotCalled(t, "Send")
			})
		})
	})
}

func TestPanicInterceptor_IgnoreIfNoPanicAndError(t *testing.T) {
	unaryHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, status.Errorf(codes.DeadlineExceeded, "deadline exceeded")
	}
	convey.Convey("Panic interceptor", t, func() {
		counterMetricSender := &mocks2.CounterMetricSender{}
		interceptor := NewInterceptor(counterMetricSender)
		convey.Convey("When I execute the interceptor", func() {
			_, err := interceptor.UnaryInterceptor()(context.Background(), empty.Empty{}, nil, unaryHandler)
			convey.Convey("Then should return original error if no panic", func() {
				convey.So(err, assertions.ShouldNotBeNil)
				convey.So(err.Error(), assertions.ShouldEqual, "rpc error: code = DeadlineExceeded desc = deadline exceeded")
				counterMetricSender.AssertNotCalled(t, "Send")
			})
		})
	})
}
