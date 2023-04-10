package campaignsconverter

import (
	"math/big"
	"time"

	"github.com/block-wallet/campaigns-service/domain/model"
	campaignservicev1service "github.com/block-wallet/campaigns-service/protos/src/campaignsservicev1/campaigns"
	"github.com/block-wallet/campaigns-service/utils/signatures"
	"github.com/ethereum/go-ethereum/common"
)

// ConverterImpl has the responsibility to convert from grpc objects to our domain objects and vice versa.
type ConverterImpl struct {
}

func NewConverterImpl() *ConverterImpl {
	return &ConverterImpl{}
}

func (c *ConverterImpl) ConvertFromModelCampaignToProtoCampaign(campaign *model.Campaign) *campaignservicev1service.Campaign {
	if campaign == nil {
		return nil
	}
	ret := &campaignservicev1service.Campaign{
		Id:              campaign.Id,
		SupportedChains: campaign.SupportedChains,
		Name:            campaign.Name,
		Description:     campaign.Description,
		Status:          campaign_FromModelStatusToProtoStatus(campaign.Status),
		StartDate:       campaign.StartDate.Format(model.CampaignTimeFormatLayout),
		EndDate:         campaign.EndDate.Format(model.CampaignTimeFormatLayout),
		Accounts:        campaign_FromAddressesSliceToStringSlice(campaign.Accounts),
		Winners:         campaign_FromAddressesSliceToStringSlice(campaign.Winners),
		Tags:            campaign.Tags,
		EnrollMessage:   campaign.EnrollMessage,
	}

	if campaign.Rewards != nil {
		ret.Rewards = &campaignservicev1service.Rewards{
			Type:    rewards_fromModelTypeToProtoType(campaign.Rewards.Type),
			Amounts: campaign_FromModelAmountsToProtoAmounts(campaign.Rewards.Amounts),
			Token:   campaign_FromModelTokenToProtoToken(campaign.Rewards.Token),
		}
	}

	return ret
}
func (c *ConverterImpl) ConvertFromProtoCampaignsFiltersToModelCampaignFilters(filters *campaignservicev1service.GetCampaignsFilters) (*model.GetCampaignsFilters, error) {
	if filters == nil {
		return &model.GetCampaignsFilters{}, nil
	}
	campaignFilters := &model.GetCampaignsFilters{}
	if len(filters.Tags) > 0 {
		campaignFilters.Tags = &filters.Tags
	}

	if len(filters.ChainIds) > 0 {
		campaignFilters.ChainIds = &filters.ChainIds
	}

	if filters.FromDate != nil {
		fromDate, e := time.Parse(model.CampaignTimeFormatLayout, filters.FromDate.Value)
		if e != nil {
			return nil, e
		}
		campaignFilters.FromDate = &fromDate
	}
	if filters.ToDate != nil {
		toDate, e := time.Parse(model.CampaignTimeFormatLayout, filters.ToDate.Value)
		if e != nil {
			return nil, e
		}
		campaignFilters.ToDate = &toDate
	}

	statuses := convertMultipleSlice(filters.Statuses, campaign_FromProtoStatusToModelStatus)

	if len(statuses) > 0 {
		campaignFilters.Status = &statuses
	}

	return campaignFilters, nil
}

func (c *ConverterImpl) ConvertFromProtoCreateCampaignToModelCreateCampaign(input *campaignservicev1service.CreateCampaignMsg) (*model.CreateCampaignInput, error) {
	campaignInput := input.Campaign
	status := model.STATUS_PENDING
	if campaignInput.IsActive {
		status = model.STATUS_ACTIVE
	}
	var rewardToken model.CampaignRewardTokenInput

	if campaignInput.Rewards.Token.Id != "" {
		rewardToken.Id = &campaignInput.Rewards.Token.Id
	} else {
		rewardToken.Create = campaign_FromProtoTokenToModelToken(campaignInput.Rewards.Token.Create)
	}
	rewards := model.CampaignRewardInput{
		Amounts: campaignInput.Rewards.Amounts,
		Token:   rewardToken,
		Type:    rewards_fromProtoTypeToModelType(input.Campaign.Rewards.Type),
	}

	enrollMessage := input.Campaign.EnrollMessage

	if enrollMessage == "" {
		enrollMessage = signatures.GenerateDefaultCampaignVerificationMessage(input.Campaign.Name)
	}

	createCampaignRet := model.CreateCampaignInput{
		Name:            campaignInput.Name,
		Description:     campaignInput.Description,
		StartDate:       campaignInput.StartDate,
		EndDate:         campaignInput.EndDate,
		Status:          status,
		Rewards:         rewards,
		Tags:            campaignInput.Tags,
		SupportedChains: campaignInput.SupportedChains,
		EnrollMessage:   enrollMessage,
	}

	return &createCampaignRet, nil
}

func (c *ConverterImpl) ConvertFromProtoUpdateCampaignToModelUpdateCampaign(campaignInput *campaignservicev1service.UpdateCampaignMsg) *model.UpdateCampaignInput {
	updateInput := &model.UpdateCampaignInput{
		Id: campaignInput.CampaignId,
	}

	if campaignInput.Status != campaignservicev1service.CampaignStatus_CAMPAIGN_STATUS_INVALID {
		status := campaign_FromProtoStatusToModelStatus(campaignInput.Status)
		updateInput.Stauts = &status
	}

	if len(campaignInput.Winners) > 0 {
		winners := campaign_FromStringSliceToAddressesSlice(campaignInput.GetWinners())
		updateInput.Winners = &winners
	}

	return updateInput
}

func (c *ConverterImpl) ConvertFromProtoEnrollInCampaignToModelEnrollInCampaign(input *campaignservicev1service.EnrollInCampaignMsg) *model.EnrollInCampaignInput {
	return &model.EnrollInCampaignInput{
		Adddress:   common.HexToAddress(input.AccountAddress),
		CampaignId: input.CampaignId,
	}
}

func (c *ConverterImpl) ConvertFromModelMultichainTokenToProtoMultichainToken(t *model.MultichainToken) *campaignservicev1service.MultichainToken {
	return campaign_FromModelTokenToProtoToken(t)
}

func campaign_FromModelTokenToProtoToken(modelToken *model.MultichainToken) *campaignservicev1service.MultichainToken {
	return &campaignservicev1service.MultichainToken{
		Id:                modelToken.Id,
		Name:              modelToken.Name,
		Decimals:          int32(modelToken.Decimals),
		Symbol:            modelToken.Symbol,
		ContractAddresses: campaign_FromModelContractAddressesToProtoContractAddresses(modelToken.ContractAddresses),
	}
}

func campaign_FromProtoTokenToModelToken(protoToken *campaignservicev1service.MultichainToken) *model.MultichainToken {
	if protoToken == nil {
		return nil
	}
	return &model.MultichainToken{
		Name:              protoToken.Name,
		Symbol:            protoToken.Symbol,
		Decimals:          uint8(protoToken.Decimals),
		ContractAddresses: campaign_FromProtoContractAddressesToModelContractAddresses(protoToken.ContractAddresses),
	}
}

func campaign_FromProtoStatusToModelStatus(protoStatus campaignservicev1service.CampaignStatus) model.CampaignStatus {
	switch protoStatus {
	case campaignservicev1service.CampaignStatus_CAMPAIGN_STATUS_ACTIVE:
		return model.STATUS_ACTIVE
	case campaignservicev1service.CampaignStatus_CAMPAIGN_STATUS_CANCELLED:
		return model.STATUS_CANCELLED
	case campaignservicev1service.CampaignStatus_CAMPAIGN_STATUS_FINISHED:
		return model.STATUS_FINISHED
	case campaignservicev1service.CampaignStatus_CAMPAIGN_STATUS_PENDING:
		return model.STATUS_PENDING
	}
	return model.STATUS_UNKNOWN
}

func campaign_FromModelStatusToProtoStatus(modelStatus model.CampaignStatus) campaignservicev1service.CampaignStatus {
	switch modelStatus {
	case model.STATUS_ACTIVE:
		return campaignservicev1service.CampaignStatus_CAMPAIGN_STATUS_ACTIVE
	case model.STATUS_CANCELLED:
		return campaignservicev1service.CampaignStatus_CAMPAIGN_STATUS_CANCELLED
	case model.STATUS_FINISHED:
		return campaignservicev1service.CampaignStatus_CAMPAIGN_STATUS_FINISHED
	case model.STATUS_PENDING:
		return campaignservicev1service.CampaignStatus_CAMPAIGN_STATUS_ACTIVE
	}
	return campaignservicev1service.CampaignStatus_CAMPAIGN_STATUS_INVALID
}

func rewards_fromModelTypeToProtoType(rewardType model.RewardType) campaignservicev1service.RewardType {
	switch rewardType {
	case model.DYNAMIC_REWARD:
		return campaignservicev1service.RewardType_REWARD_TYPE_DYNAMIC
	case model.SINGLE_REWARD:
		return campaignservicev1service.RewardType_REWARD_TYPE_SINGLE
	case model.PODIUM_REWARD:
		return campaignservicev1service.RewardType_REWARD_TYPE_PODIUM
	}
	return campaignservicev1service.RewardType_REWARD_TYPE_INVALID
}

func rewards_fromProtoTypeToModelType(rewardType campaignservicev1service.RewardType) model.RewardType {
	switch rewardType {
	case campaignservicev1service.RewardType_REWARD_TYPE_DYNAMIC:
		return model.DYNAMIC_REWARD
	case campaignservicev1service.RewardType_REWARD_TYPE_SINGLE:
		return model.SINGLE_REWARD
	case campaignservicev1service.RewardType_REWARD_TYPE_PODIUM:
		return model.PODIUM_REWARD
	}
	return model.PODIUM_REWARD
}

func campaign_FromModelAmountsToProtoAmounts(amounts []*big.Int) []string {
	return convertMultipleSlice(amounts, bigIntToString)
}

func campaign_FromModelContractAddressesToProtoContractAddresses(modelContractAddresses map[string]common.Address) map[string]string {
	return convertMapValues(modelContractAddresses, func(a common.Address) string { return a.String() })
}

func campaign_FromProtoContractAddressesToModelContractAddresses(modelContractAddresses map[string]string) map[string]common.Address {
	return convertMapValues(modelContractAddresses, func(s string) common.Address { return common.HexToAddress(s) })
}

func campaign_FromAddressesSliceToStringSlice(addresses []common.Address) []string {
	return convertMultipleSlice(addresses, func(a common.Address) string { return a.String() })
}

func campaign_FromStringSliceToAddressesSlice(addresses []string) []common.Address {
	return convertMultipleSlice(addresses, func(s string) common.Address { return common.HexToAddress(s) })
}

func convertMultipleSlice[K *big.Int | []byte | string | common.Address | campaignservicev1service.CampaignStatus, V *big.Int | []byte | string | common.Address | model.CampaignStatus](input []K, converter func(i K) V) []V {
	output := make([]V, 0, len(input))
	for _, i := range input {
		output = append(output, converter(i))
	}
	return output
}

func convertMapValues[I string | common.Address, O string | common.Address](input map[string]I, converter func(i I) O) map[string]O {
	output := make(map[string]O)
	for k, i := range input {
		output[k] = converter(i)
	}
	return output
}

func bigIntToString(value *big.Int) string {
	return value.String()
}
