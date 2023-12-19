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
		Tags:            campaign.Tags,
		EnrollMessage:   campaign.EnrollMessage,
		EnrollmentMode:  campaign_FromModelEnrollmentModeToProtoEnrollmentMode(campaign.EnrollmentMode),
		CampaignType:    campaign_FromModelCampaignTypeToProtoCampaignType(campaign.Type),
		Participants:    campaign_FromModelParticipantsToProtoParticipants(campaign),
		CreatedAt:       campaign.CreatedAt.Format(model.CampaignTimeFormatLayout),
		UpdatedAt:       campaign.UpdatedAt.Format(model.CampaignTimeFormatLayout),
	}

	switch campaign.Type {
	case model.CAMPAIGN_TYPE_GALXE:
		ret.CampaignMetadata = &campaignservicev1service.Campaign_GalxeMetadata{
			GalxeMetadata: &campaignservicev1service.GalxeCampaignMetadata{
				CredentialId: campaign.Metadata.GalxeMetadata.CredentialId,
			},
		}
	case model.CAMPAIGN_TYPE_PARTNER_OFFERS:
		{
			ret.CampaignMetadata = &campaignservicev1service.Campaign_PartnerOffersMetadata{
				PartnerOffersMetadata: &campaignservicev1service.PartnerOffersCampaignMetadata{},
			}
		}
	}

	if campaign.Rewards != nil {
		ret.Rewards = &campaignservicev1service.Rewards{
			Type:    rewards_fromModelTypeToProtoType(campaign.Rewards.Type),
			Amounts: campaign_FromModelAmountsToProtoAmounts(campaign.Rewards.Amounts),
			Token:   campaign_FromModelTokenToProtoToken(campaign.Rewards.Token),
		}
	}

	if len(campaign.Participants) > 0 {
		accounts := make([]string, len(campaign.Participants))
		for i, p := range campaign.Participants {
			accounts[i] = p.AccountAddress.String()
		}
		ret.Accounts = accounts
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

	statuses := convertMultipleSlice(filters.Status, campaign_FromProtoStatusToModelStatus)

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
		EnrollmentMode:  campaign_FromProtoEnrollmentModeToModelEnrollmentMode(campaignInput.EnrollmentMode),
		Type:            campaign_FromProtoCampaignTypeToModelCampaignType(campaignInput.CampaignType),
	}

	campaignMetadata := model.CampaignMetadata{}

	switch createCampaignRet.Type {
	case model.CAMPAIGN_TYPE_GALXE:
		{
			campaignMetadata.GalxeMetadata = &model.GalxeCampaignMetadata{
				CredentialId: campaignInput.GetGalxeMetadata().CredentialId,
			}
		}
	case model.CAMPAIGN_TYPE_PARTNER_OFFERS:
		{
			campaignMetadata.PartnerOffersMetadata = &model.PartnerOffersMetadata{}
		}
	}

	createCampaignRet.Metadata = campaignMetadata

	return &createCampaignRet, nil
}

func (c *ConverterImpl) ConvertFromProtoUpdateCampaignToModelUpdateCampaign(campaignInput *campaignservicev1service.UpdateCampaignMsg) *model.UpdateCampaignInput {
	updateInput := &model.UpdateCampaignInput{
		Id: campaignInput.CampaignId,
	}

	if campaignInput.Status != campaignservicev1service.CampaignStatus_CAMPAIGN_STATUS_INVALID {
		status := campaign_FromProtoStatusToModelStatus(campaignInput.Status)
		updateInput.Status = &status
	}

	if len(campaignInput.EligibleAccounts) > 0 {
		elegibleAccounts := campaign_FromStringSliceToAddressesSlice(campaignInput.GetEligibleAccounts())
		updateInput.EligibleAccounts = &elegibleAccounts
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

func campaign_FromModelParticipantsToProtoParticipants(campaign *model.Campaign) []*campaignservicev1service.Participant {
	var protoParticipants = make([]*campaignservicev1service.Participant, len(campaign.Participants))

	for i, p := range campaign.Participants {
		participant := &campaignservicev1service.Participant{
			AccountAddress:  p.AccountAddress.String(),
			EarlyEnrollment: p.EarlyEnrollment,
		}

		if campaign.Status == model.STATUS_FINISHED {
			participant.Eligibility = &campaignservicev1service.Eligibility{
				IsEligible: p.Position != nil,
			}

			if participant.Eligibility.IsEligible && campaign.Rewards != nil {
				var rewardedAmount string
				switch campaign.Rewards.Type {
				case model.PODIUM_REWARD:
					{
						podiumPos := *p.Position - 1
						rewardedAmount = bigIntToString(campaign.Rewards.Amounts[podiumPos])
					}
				case model.SINGLE_REWARD:
					{
						rewardedAmount = bigIntToString(campaign.Rewards.Amounts[0])
					}
				}
				participant.Eligibility.RewardedAmount = rewardedAmount
			}
		}
		protoParticipants[i] = participant
	}

	return protoParticipants
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
	case campaignservicev1service.CampaignStatus_CAMPAIGN_STATUS_WAITLIST:
		return model.STATUS_WAITLIST
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
		return campaignservicev1service.CampaignStatus_CAMPAIGN_STATUS_PENDING
	case model.STATUS_WAITLIST:
		return campaignservicev1service.CampaignStatus_CAMPAIGN_STATUS_WAITLIST
	}
	return campaignservicev1service.CampaignStatus_CAMPAIGN_STATUS_INVALID
}

func campaign_FromModelEnrollmentModeToProtoEnrollmentMode(modelEnrollmentMode model.EnrollmentMode) campaignservicev1service.EnrollmentMode {
	switch modelEnrollmentMode {
	case model.INSTANCE_SINGLE_ENROLL:
		return campaignservicev1service.EnrollmentMode_INSTANCE_SINGLE_ENROLL
	case model.INSTANCE_UNLIMITED_ENROLL:
		return campaignservicev1service.EnrollmentMode_INSTANCE_UNLIMITED_ENROLL
	}
	return campaignservicev1service.EnrollmentMode_ENROLLMENT_MODE_INVALID
}

func campaign_FromProtoEnrollmentModeToModelEnrollmentMode(protoEnrollmentMode campaignservicev1service.EnrollmentMode) model.EnrollmentMode {
	switch protoEnrollmentMode {
	case campaignservicev1service.EnrollmentMode_INSTANCE_SINGLE_ENROLL:
		return model.INSTANCE_SINGLE_ENROLL
	case campaignservicev1service.EnrollmentMode_INSTANCE_UNLIMITED_ENROLL:
		return model.INSTANCE_UNLIMITED_ENROLL
	}
	return model.INSTANCE_SINGLE_ENROLL
}

func campaign_FromModelCampaignTypeToProtoCampaignType(modelCampaignType model.CampaignType) campaignservicev1service.CampaignType {
	switch modelCampaignType {
	case model.CAMPAIGN_TYPE_PARTNER_OFFERS:
		return campaignservicev1service.CampaignType_CAMPAIGN_TYPE_PARTNER_OFFERS
	case model.CAMPAIGN_TYPE_GALXE:
		return campaignservicev1service.CampaignType_CAMPAIGN_TYPE_GALXE
	case model.CAMPAIGN_TYPE_STAKING:
		return campaignservicev1service.CampaignType_CAMPAIGN_TYPE_STAKING
	}
	return campaignservicev1service.CampaignType_CAMPAIGN_TYPE_PARTNER_OFFERS
}

func campaign_FromProtoCampaignTypeToModelCampaignType(protoCampaignType campaignservicev1service.CampaignType) model.CampaignType {
	switch protoCampaignType {
	case campaignservicev1service.CampaignType_CAMPAIGN_TYPE_PARTNER_OFFERS:
		return model.CAMPAIGN_TYPE_PARTNER_OFFERS
	case campaignservicev1service.CampaignType_CAMPAIGN_TYPE_GALXE:
		return model.CAMPAIGN_TYPE_GALXE
	case campaignservicev1service.CampaignType_CAMPAIGN_TYPE_STAKING:
		return model.CAMPAIGN_TYPE_STAKING
	}
	return model.CAMPAIGN_TYPE_PARTNER_OFFERS
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
