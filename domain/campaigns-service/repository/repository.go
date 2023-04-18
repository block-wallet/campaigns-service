package campaignsrepository

import (
	"context"

	"github.com/block-wallet/campaigns-service/domain/model"
)

type Repository interface {
	//Campaigns management
	GetCampaigns(ctx context.Context, filters *model.GetCampaignsFilters) ([]*model.Campaign, error)
	GetCampaignById(ctx context.Context, id string) (*model.Campaign, error)
	NewCampaign(ctx context.Context, input *model.CreateCampaignInput) (*string, error)
	EnrollInCampaign(ctx context.Context, input *model.EnrollInCampaignInput) (*bool, error)
	ParticipantExists(ctx context.Context, campaignId string, accountAddress string) (*bool, error)
	UpdateCampaign(ctx context.Context, updates *model.UpdateCampaignInput) (*bool, error)

	//Token managemeent
	GetTokenById(ctx context.Context, id string) (*model.MultichainToken, error)
	GetAllTokens(ctx context.Context) ([]*model.MultichainToken, error)
	TokenExists(ctx context.Context, id string) (*bool, error)
	NewToken(ctx context.Context, token *model.MultichainToken) (*string, error)
}
