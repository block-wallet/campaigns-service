package auth

import (
	"context"

	"github.com/block-wallet/campaigns-service/utils/errors"
)

type Auth interface {
	AuthenticateUsingContext(ctx context.Context) errors.RichError
}
