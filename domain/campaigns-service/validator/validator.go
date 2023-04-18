package campaignsservicevalidator

import (
	campaignservicev1service "github.com/block-wallet/campaigns-service/protos/src/campaignsservicev1/campaigns"
	"github.com/block-wallet/campaigns-service/utils/errors"
)

type Validator interface {
	ValidateGetCampaignsRequest(req *campaignservicev1service.GetCampaignsMsg) errors.RichError
	ValidateGetCampaignByIdRequest(req *campaignservicev1service.GetCampaignByIdMsg) errors.RichError
	ValidateCreateCampaignRequest(req *campaignservicev1service.CreateCampaignMsg) errors.RichError
	ValidateUpdateCampaignRequest(req *campaignservicev1service.UpdateCampaignMsg) errors.RichError
}
