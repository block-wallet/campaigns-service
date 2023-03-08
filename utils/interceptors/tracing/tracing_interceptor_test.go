package tracing

import (
	"context"
	"testing"

	"github.com/block-wallet/golang-service-template/utils/logger"
	"github.com/smartystreets/assertions"
	"github.com/smartystreets/goconvey/convey"
)

// nolint:staticcheck
func TestTracingInterceptor_UnaryInterceptor(t *testing.T) {
	convey.Convey("Unary tracing interceptor", t, func() {
		unaryHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
			return nil, nil
		}
		ctx := context.WithValue(context.Background(), "key", "val")
		requestIDField := "request_id"
		convey.Convey("When I execute the interceptor", func() {
			_, err := NewInterceptor(logger.ContextKey(requestIDField)).UnaryInterceptor()(ctx, nil, nil, unaryHandler)
			convey.Convey("Then should add request id on context", func() {
				convey.So(err, assertions.ShouldBeNil)
				convey.So(ctx.Value(requestIDField), assertions.ShouldBeEmpty)
			})
		})
	})
}
