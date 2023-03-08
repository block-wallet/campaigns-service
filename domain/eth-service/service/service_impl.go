package ethservice

import (
	"context"

	ethrepository "github.com/block-wallet/golang-service-template/domain/eth-service/repository"
	"github.com/block-wallet/golang-service-template/domain/model"

	"github.com/block-wallet/golang-service-template/utils/errors"
)

type ServiceImpl struct {
	repository ethrepository.Repository
}

func NewServiceImpl(repository ethrepository.Repository) *ServiceImpl {
	return &ServiceImpl{
		repository: repository,
	}
}

func (s *ServiceImpl) GetEvents(ctx context.Context, pair string) (*[]model.Event, errors.RichError) {
	return s.repository.GetEvents(ctx, pair)
}
func (s *ServiceImpl) GetChains(ctx context.Context) (*[]model.Chain, errors.RichError) {
	return s.repository.GetChains(ctx)
}
