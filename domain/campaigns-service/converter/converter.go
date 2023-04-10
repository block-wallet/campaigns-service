package campaignsconverter

import (
	"github.com/block-wallet/campaigns-service/domain/model"
	campaignservicev1service "github.com/block-wallet/campaigns-service/protos/src/campaignsservicev1/campaigns"
)

type Converter interface {
	ConvertFromModelCampaignToProtoCampaign(campaign *model.Campaign) *campaignservicev1service.Campaign
	ConvertFromProtoCampaignsFiltersToModelCampaignFilters(filters *campaignservicev1service.GetCampaignsFilters) (*model.GetCampaignsFilters, error)
	ConvertFromProtoCreateCampaignToModelCreateCampaign(campaignInput *campaignservicev1service.CreateCampaignMsg) (*model.CreateCampaignInput, error)
	ConvertFromProtoEnrollInCampaignToModelEnrollInCampaign(campaignInput *campaignservicev1service.EnrollInCampaignMsg) *model.EnrollInCampaignInput
	ConvertFromProtoUpdateCampaignToModelUpdateCampaign(campaignInput *campaignservicev1service.UpdateCampaignMsg) *model.UpdateCampaignInput
	ConvertFromModelMultichainTokenToProtoMultichainToken(t *model.MultichainToken) *campaignservicev1service.MultichainToken
}
