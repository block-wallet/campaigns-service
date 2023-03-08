package http

import (
	"context"
)

type Client interface {
	Get(ctx context.Context, url string, headers map[string]string) ([]byte, error)
}
