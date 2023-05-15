package campaignsservicevalidator

import (
	"fmt"
	"math/big"
	"time"

	"github.com/block-wallet/campaigns-service/domain/model"
	campaignservicev1service "github.com/block-wallet/campaigns-service/protos/src/campaignsservicev1/campaigns"
	"github.com/block-wallet/campaigns-service/utils/errors"
	"github.com/google/uuid"
)

// RequestValidator implements the Validator interface
type RequestValidator struct {
}

func NewRequestValidator() *RequestValidator {
	return &RequestValidator{}
}

func (r *RequestValidator) ValidateGetCampaignsRequest(req *campaignservicev1service.GetCampaignsMsg) errors.RichError {
	if req.Filters != nil {
		fromDateP := req.Filters.FromDate
		toDateP := req.Filters.ToDate
		if fromDate := fromDateP.GetValue(); fromDate != "" {
			_, err := time.Parse(model.CampaignTimeFormatLayout, fromDate)
			if err != nil {
				return errors.NewInvalidArgument(fmt.Sprintf("Invalid from_date format. Please make sure it has the following format: %v", model.CampaignTimeFormatLayout))
			}
		}
		if toDate := toDateP.GetValue(); toDate != "" {
			_, err := time.Parse(model.CampaignTimeFormatLayout, toDate)
			if err != nil {
				return errors.NewInvalidArgument(fmt.Sprintf("Invalid to_date format. Please make sure it has the following format: %v", model.CampaignTimeFormatLayout))
			}
		}
	}
	return nil
}

func (r *RequestValidator) ValidateGetCampaignByIdRequest(req *campaignservicev1service.GetCampaignByIdMsg) errors.RichError {
	if req.Id == "" {
		return errors.NewInvalidArgument("id cannot be empty.")
	}
	if !IsValidUUID(req.Id) {
		return errors.NewInvalidArgument("invalid id format.")
	}
	return nil
}

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func (r *RequestValidator) ValidateCreateCampaignRequest(req *campaignservicev1service.CreateCampaignMsg) errors.RichError {
	campaignReq := req.Campaign
	if campaignReq.Name == "" {
		return errors.NewInvalidArgument("campaign name cannot be empty.")
	}
	if campaignReq.Description == "" {
		return errors.NewInvalidArgument("campaign description cannot be empty.")
	}

	if len(campaignReq.SupportedChains) == 0 {
		return errors.NewInvalidArgument("campaigns should target at least one chain.")
	}

	if campaignReq.CampaignType == campaignservicev1service.CampaignType_CAMPAIGN_TYPE_INVALID {
		return errors.NewInvalidArgument("you should specify the campaign type.")
	}

	if campaignReq.CampaignType == campaignservicev1service.CampaignType_CAMPAIGN_TYPE_GALXE {
		galxeMetadata := campaignReq.GetGalxeMetadata()
		if galxeMetadata == nil || galxeMetadata.GetCredentialId() == "" {
			return errors.NewInvalidArgument("galxe campaigns should specify the campaign's credential id.")
		}
	}

	if campaignReq.Rewards == nil {
		return errors.NewInvalidArgument("you must specify rewards for the campaign.")
	} else {
		if campaignReq.Rewards.Type == campaignservicev1service.RewardType_REWARD_TYPE_INVALID {
			return errors.NewInvalidArgument("you need to specify the rewards type of this campaign.")
		}

		if campaignReq.Rewards.Token == nil || (campaignReq.Rewards.Token.Id == "" && campaignReq.Rewards.Token.Create == nil) {
			return errors.NewInvalidArgument("campaign rewards should specify the token.")
		}

		if campaignReq.Rewards.Token.Id != "" && !IsValidUUID(campaignReq.Rewards.Token.Id) {
			return errors.NewInvalidArgument("invalid token id format. It should be a valid UUID")
		}

		if campaignReq.Rewards.Token.Create != nil {
			tokenCreate := campaignReq.Rewards.Token.Create
			if tokenCreate.Symbol == "" || tokenCreate.Name == "" || tokenCreate.Decimals == 0 {
				return errors.NewInvalidArgument("missing rewards token information. Symbol, Name and Decimals are required fileds.")
			}
		}

		rewardsType := campaignReq.Rewards.Type
		amounts := campaignReq.Rewards.Amounts
		switch rewardsType {
		case campaignservicev1service.RewardType_REWARD_TYPE_DYNAMIC:
			if len(amounts) > 0 {
				return errors.NewInvalidArgument("you can't specify amounts for a DYNAMIC reward type campaign")
			}
		case campaignservicev1service.RewardType_REWARD_TYPE_SINGLE:
			if len(amounts) != 1 {
				return errors.NewInvalidArgument("you need to specify only one amount for SINGLE reward type campaign.")
			}
		case campaignservicev1service.RewardType_REWARD_TYPE_PODIUM:
			if len(amounts) == 0 {
				return errors.NewInvalidArgument("campaign rewards amounts cannot be empty for either PODIUM rewards type campaign.")
			}
		}

		var prev *big.Int
		zero := big.NewInt(0)
		for _, amount := range campaignReq.Rewards.Amounts {
			if amount == "" {
				return errors.NewInvalidArgument("campaign rewards amounts cannot be empty.")
			}
			parsedAmount := new(big.Int)
			parsedAmount, ok := parsedAmount.SetString(amount, 10)
			if !ok {
				return errors.NewInvalidArgument(fmt.Sprintf("invalid amount number: %v", amount))
			}
			if parsedAmount.Cmp(zero) <= 0 {
				return errors.NewInvalidArgument("campaign rewards amounts should be bigger than 0.")
			}
			if prev != nil && parsedAmount.Cmp(prev) > 0 {
				return errors.NewInvalidArgument("campaign rewards amounts are not sorted.")
			}
			prev = parsedAmount
		}

	}

	//dates validations
	if campaignReq.StartDate == "" || campaignReq.EndDate == "" {
		return errors.NewInvalidArgument("campaign dates cannot be empty.")
	} else {
		startDate, err := time.Parse(model.CampaignTimeFormatLayout, campaignReq.StartDate)
		if err != nil {
			return errors.NewInvalidArgument(fmt.Sprintf("invalid start_date format. Please use the following format: %s", model.CampaignTimeFormatLayout))
		}
		endDate, err := time.Parse(model.CampaignTimeFormatLayout, campaignReq.EndDate)
		if err != nil {
			return errors.NewInvalidArgument(fmt.Sprintf("invalid end_date format. Please use the following format: %s", model.CampaignTimeFormatLayout))
		}

		if endDate.Before(startDate) {
			return errors.NewInvalidArgument("campaign's start_date cannot be after than end_date.")
		}

		if endDate.Before(time.Now()) {
			return errors.NewInvalidArgument("cannot create and already finished campaign.")
		}

		if startDate.After(time.Now()) && campaignReq.IsActive {
			return errors.NewInvalidArgument("cannot activate a campaign that hasn't started yet.")
		}
	}
	return nil
}

func (r *RequestValidator) ValidateUpdateCampaignRequest(req *campaignservicev1service.UpdateCampaignMsg) errors.RichError {
	if !IsValidUUID(req.CampaignId) {
		return errors.NewInvalidArgument("invalid campaign_id format")
	}
	status := req.Status
	if status != campaignservicev1service.CampaignStatus_CAMPAIGN_STATUS_INVALID {
		switch status {
		case campaignservicev1service.CampaignStatus_CAMPAIGN_STATUS_FINISHED:
			if len(req.EligibleAccounts) == 0 {
				return errors.NewInvalidArgument("you must specify the elegible accounts to finish a campaign")
			}
		case campaignservicev1service.CampaignStatus_CAMPAIGN_STATUS_ACTIVE, campaignservicev1service.CampaignStatus_CAMPAIGN_STATUS_CANCELLED:
			if len(req.EligibleAccounts) > 0 {
				return errors.NewInvalidArgument("you can only set campaign elegible accounts by updating its status to FINISHED")
			}
		}
	} else {
		if len(req.EligibleAccounts) > 0 {
			return errors.NewInvalidArgument("you need to update status to FINISHED to update the campaign elegible accounts")
		}
	}

	return nil
}
