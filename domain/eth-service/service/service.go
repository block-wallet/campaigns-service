package ethservice

import (
	"context"

	"github.com/block-wallet/golang-service-template/domain/model"

	"github.com/block-wallet/golang-service-template/utils/errors"
)

type Service interface {
	GetEvents(ctx context.Context, pair string) (*[]model.Event, errors.RichError)
	GetChains(ctx context.Context) (*[]model.Chain, errors.RichError)
}
