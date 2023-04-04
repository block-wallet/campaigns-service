package tracing

import (
	"context"
	"testing"

	"github.com/block-wallet/campaigns-service/utils/logger"
	"github.com/stretchr/testify/assert"
)

// nolint:staticcheck
func TestTracingInterceptor_UnaryInterceptor(t *testing.T) {
	unaryHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, nil
	}
	ctx := context.WithValue(context.Background(), "key", "val")
	requestIDField := "request_id"
	_, err := NewInterceptor(logger.ContextKey(requestIDField)).UnaryInterceptor()(ctx, nil, nil, unaryHandler)

	assert.NoError(t, err)
	assert.Empty(t, ctx.Value(requestIDField))
}
