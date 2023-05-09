package campaignsservice

import (
	"context"

	"github.com/block-wallet/campaigns-service/domain/model"

	"github.com/block-wallet/campaigns-service/utils/errors"
)

type Service interface {
	GetCampaigns(ctx context.Context, filters *model.GetCampaignsFilters) ([]*model.Campaign, errors.RichError)
	GetCampaignById(ctx context.Context, id string) (*model.Campaign, errors.RichError)
	CreateCampaign(ctx context.Context, input *model.CreateCampaignInput) (*model.Campaign, errors.RichError)
	EnrollInCampaign(ctx context.Context, input *model.EnrollInCampaignInput) (bool, errors.RichError)
	UpdateCampaign(ctx context.Context, input *model.UpdateCampaignInput) (*model.Campaign, errors.RichError)

	//tokens
	GetTokenById(ctx context.Context, id string) (*model.MultichainToken, errors.RichError)
	GetAllTokens(ctx context.Context) ([]*model.MultichainToken, errors.RichError)
}
