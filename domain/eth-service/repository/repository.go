package ethrepository

import (
	"context"

	"github.com/block-wallet/golang-service-template/domain/model"

	"github.com/block-wallet/golang-service-template/utils/errors"
)

type Repository interface {
	GetEvents(ctx context.Context, pair string) (*[]model.Event, errors.RichError)
	GetChains(ctx context.Context) (*[]model.Chain, errors.RichError)
	SetEvent(ctx context.Context, pair, eventId string, event *model.Event) errors.RichError
}
