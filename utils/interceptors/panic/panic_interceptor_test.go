package panic

import (
	"context"
	"testing"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/block-wallet/campaigns-service/utils/monitoring/counter/mocks"
	"github.com/stretchr/testify/assert"
)

func TestPanicInterceptor_RecoverFromPanic(t *testing.T) {
	unaryHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		panic("some error")
	}
	var labels map[string]string
	counterMetricSender := &mocks.CounterMetricSender{}
	counterMetricSender.On("Send", labels)
	interceptor := NewInterceptor(counterMetricSender)
	_, err := interceptor.UnaryInterceptor()(context.Background(), empty.Empty{}, nil, unaryHandler)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "rpc error: code = Internal desc = server panic: some error")
	counterMetricSender.AssertCalled(t, "Send", labels)
}

func TestPanicInterceptor_RecoverFromNilPanic(t *testing.T) {
	unaryHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		panic(nil)
	}

	var labels map[string]string
	counterMetricSender := &mocks.CounterMetricSender{}
	counterMetricSender.On("Send", labels)
	interceptor := NewInterceptor(counterMetricSender)
	_, err := interceptor.UnaryInterceptor()(context.Background(), empty.Empty{}, nil, unaryHandler)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "rpc error: code = Internal desc = server panic: <nil>")
	counterMetricSender.AssertCalled(t, "Send", labels)
}

func TestPanicInterceptor_IgnoreIfNoPanic(t *testing.T) {
	unaryHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "ok response", nil
	}
	counterMetricSender := &mocks.CounterMetricSender{}
	interceptor := NewInterceptor(counterMetricSender)
	resp, err := interceptor.UnaryInterceptor()(context.Background(), empty.Empty{}, nil, unaryHandler)

	assert.NoError(t, err)
	assert.Equal(t, resp.(string), "ok response")
	counterMetricSender.AssertNotCalled(t, "Send")
}

func TestPanicInterceptor_IgnoreIfNoPanicAndError(t *testing.T) {
	unaryHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, status.Errorf(codes.DeadlineExceeded, "deadline exceeded")
	}

	counterMetricSender := &mocks.CounterMetricSender{}
	interceptor := NewInterceptor(counterMetricSender)
	_, err := interceptor.UnaryInterceptor()(context.Background(), empty.Empty{}, nil, unaryHandler)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "rpc error: code = DeadlineExceeded desc = deadline exceeded")
	counterMetricSender.AssertNotCalled(t, "Send")
}
