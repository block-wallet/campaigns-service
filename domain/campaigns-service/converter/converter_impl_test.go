package campaignsconverter

import (
	"math/big"
	"testing"
	"time"

	"github.com/block-wallet/campaigns-service/domain/model"
	campaignsservicev1 "github.com/block-wallet/campaigns-service/protos/src/campaignsservicev1/campaigns"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func Test_NewConverterImpl(t *testing.T) {
	conv := NewConverterImpl()
	assert.NotNil(t, conv)
}

func Test_ConvertFromModelCampaignToProtoCampaign(t *testing.T) {
	conv := NewConverterImpl()
	createdAt, _ := time.Parse(model.CampaignTimeFormatLayout, "2023-04-01T00:00:00Z")

	startDate, _ := time.Parse(model.CampaignTimeFormatLayout, "2023-04-01T00:00:00Z")
	endDate, _ := time.Parse(model.CampaignTimeFormatLayout, "2023-05-01T00:00:00Z")
	participans := []string{
		"0xf0F8B7C21e280b0167F14Af6db4B9F90430A6C22",
		"0x6B25a67345111737d5FBEA0dd1443F64F0ef17CB",
		"0xFBCFfCE77aad086216C8016A73faeec1fE8fFEcc",
	}
	amount1, _ := big.NewInt(0).SetString("1", 10)
	amount2, _ := big.NewInt(0).SetString("2", 10)
	amount3, _ := big.NewInt(0).SetString("3", 10)

	cases := []struct {
		name     string
		expected *campaignsservicev1.Campaign
		input    *model.Campaign
	}{
		{
			name:     "Nil convertion",
			input:    nil,
			expected: nil,
		},
		{
			name: "Should convert model.Campaign to proto Campaign",
			input: &model.Campaign{
				Id:              "123",
				SupportedChains: []uint32{1, 137},
				Name:            "Test campaign",
				Description:     "Description test campaign",
				Status:          model.STATUS_ACTIVE,
				StartDate:       startDate,
				EndDate:         endDate,
				CreatedAt:       createdAt,
				UpdatedAt:       createdAt,
				Participants: []model.CampaignParticipant{
					{
						AccountAddress: common.HexToAddress(participans[0]),
					},
					{
						AccountAddress: common.HexToAddress(participans[1]),
					},
					{
						AccountAddress: common.HexToAddress(participans[2]),
					},
				},
				Tags:          []string{"tag1", "tag2"},
				EnrollMessage: "Custom enroll message",
				Rewards: &model.Reward{
					Type: model.PODIUM_REWARD,
					Amounts: []*big.Int{
						amount1, amount2, amount3,
					},
					Token: &model.MultichainToken{
						Id:       "token-id-1",
						Name:     "GoBlank",
						Symbol:   "BLANK",
						Decimals: 18,
						ContractAddresses: map[string]common.Address{
							"1":   common.HexToAddress("0x41A3Dba3D677E573636BA691a70ff2D606c29666"),
							"137": common.HexToAddress("0xf4C83080E80AE530d6f8180572cBbf1Ac9D5d435"),
						},
					},
				},
				EnrollmentMode: model.INSTANCE_UNLIMITED_ENROLL,
				Type:           model.CAMPAIGN_TYPE_GALXE,
				Metadata: model.CampaignMetadata{
					GalxeMetadata: &model.GalxeCampaignMetadata{
						CredentialId: "123456",
					},
				},
			},
			expected: &campaignsservicev1.Campaign{
				Id:              "123",
				SupportedChains: []uint32{1, 137},
				Name:            "Test campaign",
				Description:     "Description test campaign",
				Status:          campaignsservicev1.CampaignStatus_CAMPAIGN_STATUS_ACTIVE,
				StartDate:       startDate.Format(model.CampaignTimeFormatLayout),
				EndDate:         endDate.Format(model.CampaignTimeFormatLayout),
				CreatedAt:       createdAt.Format(model.CampaignTimeFormatLayout),
				UpdatedAt:       createdAt.Format(model.CampaignTimeFormatLayout),
				Accounts:        participans,
				Participants: []*campaignsservicev1.Participant{
					{
						AccountAddress: participans[0],
					},
					{
						AccountAddress: participans[1],
					},
					{
						AccountAddress: participans[2],
					},
				},
				Tags:          []string{"tag1", "tag2"},
				EnrollMessage: "Custom enroll message",
				CampaignType:  campaignsservicev1.CampaignType_CAMPAIGN_TYPE_GALXE,
				CampaignMetadata: &campaignsservicev1.Campaign_GalxeMetadata{
					GalxeMetadata: &campaignsservicev1.GalxeCampaignMetadata{
						CredentialId: "123456",
					},
				},
				Rewards: &campaignsservicev1.Rewards{
					Type: campaignsservicev1.RewardType_REWARD_TYPE_PODIUM,
					Amounts: []string{
						"1", "2", "3",
					},
					Token: &campaignsservicev1.MultichainToken{
						Id:       "token-id-1",
						Name:     "GoBlank",
						Symbol:   "BLANK",
						Decimals: 18,
						ContractAddresses: map[string]string{
							"1":   common.HexToAddress("0x41A3Dba3D677E573636BA691a70ff2D606c29666").String(),
							"137": common.HexToAddress("0xf4C83080E80AE530d6f8180572cBbf1Ac9D5d435").String(),
						},
					},
				},
				EnrollmentMode: campaignsservicev1.EnrollmentMode_INSTANCE_UNLIMITED_ENROLL,
			},
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			result := conv.ConvertFromModelCampaignToProtoCampaign(c.input)
			assert.Equal(t, c.expected.String(), result.String())
		})
	}
}

func Test_ConvertFromProtoCampaignFiltersToModelCampaignFilters(t *testing.T) {
	conv := NewConverterImpl()

	fromDateStr, toDateStr := "2023-04-01T00:00:00Z", "2023-05-01T00:00:00Z"
	fromDate, _ := time.Parse(model.CampaignTimeFormatLayout, fromDateStr)
	toDate, _ := time.Parse(model.CampaignTimeFormatLayout, toDateStr)
	cases := []struct {
		name     string
		input    *campaignsservicev1.GetCampaignsFilters
		expected *model.GetCampaignsFilters
		expecErr bool
	}{
		{
			name:     "Should convert nil to empty struct",
			input:    nil,
			expected: &model.GetCampaignsFilters{},
		},
		{
			name: "Should convert a single filter",
			input: &campaignsservicev1.GetCampaignsFilters{
				Tags: []string{"test1"},
			},
			expected: &model.GetCampaignsFilters{
				Tags: &[]string{"test1"},
			},
		},
		{
			name: "Should convert all the filters",
			input: &campaignsservicev1.GetCampaignsFilters{
				Status: []campaignsservicev1.CampaignStatus{
					campaignsservicev1.CampaignStatus_CAMPAIGN_STATUS_ACTIVE,
					campaignsservicev1.CampaignStatus_CAMPAIGN_STATUS_CANCELLED,
					campaignsservicev1.CampaignStatus_CAMPAIGN_STATUS_FINISHED,
					campaignsservicev1.CampaignStatus_CAMPAIGN_STATUS_PENDING,
				},
				FromDate: &wrapperspb.StringValue{
					Value: fromDateStr,
				},
				ToDate: &wrapperspb.StringValue{
					Value: toDateStr,
				},
				Tags:     []string{"tag1", "tag2"},
				ChainIds: []uint32{1, 137},
			},
			expected: &model.GetCampaignsFilters{
				Status: &[]model.CampaignStatus{
					model.STATUS_ACTIVE,
					model.STATUS_CANCELLED,
					model.STATUS_FINISHED,
					model.STATUS_PENDING,
				},
				FromDate: &fromDate,
				ToDate:   &toDate,
				Tags:     &[]string{"tag1", "tag2"},
				ChainIds: &[]uint32{1, 137},
			},
		},
		{
			name: "should throw error on invalid date format",
			input: &campaignsservicev1.GetCampaignsFilters{
				ToDate: &wrapperspb.StringValue{
					Value: "this is an invalid date",
				},
			},
			expecErr: true,
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			result, err := conv.ConvertFromProtoCampaignsFiltersToModelCampaignFilters(c.input)
			if c.expecErr {
				assert.NotNil(t, err)
			} else if c.expected != nil {
				assert.EqualValues(t, c.expected, result)
			}
		})
	}
}

func Test_ConvertFromProtoCreateCampaignToModelCreateCampaign(t *testing.T) {
	conv := NewConverterImpl()
	tokenId := "token-test-1"
	startDateStr, endDateStr := "2023-04-01T00:00:00Z", "2023-05-01T00:00:00Z"
	cases := []struct {
		name     string
		input    *campaignsservicev1.CreateCampaignMsg
		expected *model.CreateCampaignInput
	}{
		{
			name: "should parse campaign input and set status active",
			input: &campaignsservicev1.CreateCampaignMsg{
				Campaign: &campaignsservicev1.CreateCampaignMsg_CreateCampaignInput{
					Name:        "campaign name",
					Description: "campaign description",
					IsActive:    true,
					StartDate:   startDateStr,
					EndDate:     endDateStr,
					Rewards: &campaignsservicev1.CreateCampaignMsg_CreateCampaignInput_CreateCampaignRewardInput{
						Token: &campaignsservicev1.CreateCampaignMsg_CreateCampaignInput_CreateCampaignRewardInput_CreateCampaignTokenInput{
							Id: tokenId,
						},
						Amounts: []string{"1", "2", "3"},
						Type:    campaignsservicev1.RewardType_REWARD_TYPE_PODIUM,
					},
					Tags:            []string{"tag1"},
					SupportedChains: []uint32{1, 137},
					EnrollMessage:   "custom enroll message",
					EnrollmentMode:  campaignsservicev1.EnrollmentMode_INSTANCE_SINGLE_ENROLL,
					CampaignType:    campaignsservicev1.CampaignType_CAMPAIGN_TYPE_GALXE,
					Metadata: &campaignsservicev1.CreateCampaignMsg_CreateCampaignInput_GalxeMetadata{
						GalxeMetadata: &campaignsservicev1.GalxeCampaignMetadata{
							CredentialId: "123456",
						},
					},
				},
			},
			expected: &model.CreateCampaignInput{
				Name:            "campaign name",
				Description:     "campaign description",
				StartDate:       startDateStr,
				EndDate:         endDateStr,
				Status:          model.STATUS_ACTIVE,
				Tags:            []string{"tag1"},
				SupportedChains: []uint32{1, 137},
				EnrollMessage:   "custom enroll message",
				EnrollmentMode:  model.INSTANCE_SINGLE_ENROLL,
				Type:            model.CAMPAIGN_TYPE_GALXE,
				Metadata: model.CampaignMetadata{
					GalxeMetadata: &model.GalxeCampaignMetadata{
						CredentialId: "123456",
					},
				},
				Rewards: model.CampaignRewardInput{
					Amounts: []string{"1", "2", "3"},
					Type:    model.PODIUM_REWARD,
					Token: model.CampaignRewardTokenInput{
						Id: &tokenId,
					},
				},
			},
		},
		{
			name: "should parse campaign input and set status pending",
			input: &campaignsservicev1.CreateCampaignMsg{
				Campaign: &campaignsservicev1.CreateCampaignMsg_CreateCampaignInput{
					Name:        "campaign name",
					Description: "campaign description",
					IsActive:    false,
					StartDate:   startDateStr,
					EndDate:     endDateStr,
					Rewards: &campaignsservicev1.CreateCampaignMsg_CreateCampaignInput_CreateCampaignRewardInput{
						Token: &campaignsservicev1.CreateCampaignMsg_CreateCampaignInput_CreateCampaignRewardInput_CreateCampaignTokenInput{
							Id: tokenId,
						},
						Amounts: []string{"1", "2", "3"},
						Type:    campaignsservicev1.RewardType_REWARD_TYPE_PODIUM,
					},
					Tags:            []string{"tag1"},
					SupportedChains: []uint32{1, 137},
					EnrollMessage:   "custom enroll message",
					EnrollmentMode:  campaignsservicev1.EnrollmentMode_INSTANCE_UNLIMITED_ENROLL,
					CampaignType:    campaignsservicev1.CampaignType_CAMPAIGN_TYPE_PARTNER_OFFERS,
					Metadata: &campaignsservicev1.CreateCampaignMsg_CreateCampaignInput_PartnerOffersMetadata{
						PartnerOffersMetadata: &campaignsservicev1.PartnerOffersCampaignMetadata{},
					},
				},
			},
			expected: &model.CreateCampaignInput{
				Name:            "campaign name",
				Description:     "campaign description",
				StartDate:       startDateStr,
				EndDate:         endDateStr,
				Status:          model.STATUS_PENDING,
				Tags:            []string{"tag1"},
				SupportedChains: []uint32{1, 137},
				EnrollMessage:   "custom enroll message",
				EnrollmentMode:  model.INSTANCE_UNLIMITED_ENROLL,
				Rewards: model.CampaignRewardInput{
					Amounts: []string{"1", "2", "3"},
					Type:    model.PODIUM_REWARD,
					Token: model.CampaignRewardTokenInput{
						Id: &tokenId,
					},
				},
				Type: model.CAMPAIGN_TYPE_PARTNER_OFFERS,
				Metadata: model.CampaignMetadata{
					PartnerOffersMetadata: &model.PartnerOffersMetadata{},
				},
			},
		},
		{
			name: "should parse campaign input and set default enroll message",
			input: &campaignsservicev1.CreateCampaignMsg{
				Campaign: &campaignsservicev1.CreateCampaignMsg_CreateCampaignInput{
					Name:        "campaign name",
					Description: "campaign description",
					IsActive:    false,
					StartDate:   startDateStr,
					EndDate:     endDateStr,
					Rewards: &campaignsservicev1.CreateCampaignMsg_CreateCampaignInput_CreateCampaignRewardInput{
						Token: &campaignsservicev1.CreateCampaignMsg_CreateCampaignInput_CreateCampaignRewardInput_CreateCampaignTokenInput{
							Id: tokenId,
						},
						Amounts: []string{"1", "2", "3"},
						Type:    campaignsservicev1.RewardType_REWARD_TYPE_PODIUM,
					},
					Tags:            []string{"tag1"},
					SupportedChains: []uint32{1, 137},
					EnrollmentMode:  campaignsservicev1.EnrollmentMode_INSTANCE_UNLIMITED_ENROLL,
					CampaignType:    campaignsservicev1.CampaignType_CAMPAIGN_TYPE_PARTNER_OFFERS,
					Metadata: &campaignsservicev1.CreateCampaignMsg_CreateCampaignInput_PartnerOffersMetadata{
						PartnerOffersMetadata: &campaignsservicev1.PartnerOffersCampaignMetadata{},
					},
				},
			},
			expected: &model.CreateCampaignInput{
				Name:            "campaign name",
				Description:     "campaign description",
				StartDate:       startDateStr,
				EndDate:         endDateStr,
				Status:          model.STATUS_PENDING,
				Tags:            []string{"tag1"},
				SupportedChains: []uint32{1, 137},
				EnrollmentMode:  model.INSTANCE_UNLIMITED_ENROLL,
				EnrollMessage:   "Sign this message to enroll in campaign name",
				Rewards: model.CampaignRewardInput{
					Amounts: []string{"1", "2", "3"},
					Type:    model.PODIUM_REWARD,
					Token: model.CampaignRewardTokenInput{
						Id: &tokenId,
					},
				},
				Type: model.CAMPAIGN_TYPE_PARTNER_OFFERS,
				Metadata: model.CampaignMetadata{
					PartnerOffersMetadata: &model.PartnerOffersMetadata{},
				},
			},
		},
		{
			name: "should parse campaign input and map the new token to create",
			input: &campaignsservicev1.CreateCampaignMsg{
				Campaign: &campaignsservicev1.CreateCampaignMsg_CreateCampaignInput{
					Name:        "campaign name",
					Description: "campaign description",
					IsActive:    true,
					StartDate:   startDateStr,
					EndDate:     endDateStr,
					Rewards: &campaignsservicev1.CreateCampaignMsg_CreateCampaignInput_CreateCampaignRewardInput{
						Token: &campaignsservicev1.CreateCampaignMsg_CreateCampaignInput_CreateCampaignRewardInput_CreateCampaignTokenInput{
							Create: &campaignsservicev1.MultichainToken{
								Name:     "TokenCampaigns",
								Decimals: 18,
								Symbol:   "TOKCA",
								ContractAddresses: map[string]string{
									"1":   "0xf0F8B7C21e280b0167F14Af6db4B9F90430A6C22",
									"137": "0x6B25a67345111737d5FBEA0dd1443F64F0ef17CB",
								},
							},
						},
						Amounts: []string{"1", "2", "3"},
						Type:    campaignsservicev1.RewardType_REWARD_TYPE_PODIUM,
					},
					Tags:            []string{"tag1"},
					SupportedChains: []uint32{1, 137},
					EnrollMessage:   "custom enroll message",
					EnrollmentMode:  campaignsservicev1.EnrollmentMode_INSTANCE_UNLIMITED_ENROLL,
					CampaignType:    campaignsservicev1.CampaignType_CAMPAIGN_TYPE_PARTNER_OFFERS,
					Metadata: &campaignsservicev1.CreateCampaignMsg_CreateCampaignInput_PartnerOffersMetadata{
						PartnerOffersMetadata: &campaignsservicev1.PartnerOffersCampaignMetadata{},
					},
				},
			},
			expected: &model.CreateCampaignInput{
				Name:            "campaign name",
				Description:     "campaign description",
				StartDate:       startDateStr,
				EndDate:         endDateStr,
				Status:          model.STATUS_ACTIVE,
				Tags:            []string{"tag1"},
				SupportedChains: []uint32{1, 137},
				EnrollMessage:   "custom enroll message",
				EnrollmentMode:  model.INSTANCE_UNLIMITED_ENROLL,
				Rewards: model.CampaignRewardInput{
					Amounts: []string{"1", "2", "3"},
					Type:    model.PODIUM_REWARD,
					Token: model.CampaignRewardTokenInput{
						Create: &model.MultichainToken{
							Name:     "TokenCampaigns",
							Decimals: 18,
							Symbol:   "TOKCA",
							ContractAddresses: map[string]common.Address{
								"1":   common.HexToAddress("0xf0F8B7C21e280b0167F14Af6db4B9F90430A6C22"),
								"137": common.HexToAddress("0x6B25a67345111737d5FBEA0dd1443F64F0ef17CB"),
							},
						},
					},
				},
				Type: model.CAMPAIGN_TYPE_PARTNER_OFFERS,
				Metadata: model.CampaignMetadata{
					PartnerOffersMetadata: &model.PartnerOffersMetadata{},
				},
			},
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			result, _ := conv.ConvertFromProtoCreateCampaignToModelCreateCampaign(c.input)
			assert.EqualValues(t, c.expected, result)
		})
	}
}

func Test_ConvertFromProtoUpdateCampaignToModelUpdateCampaign(t *testing.T) {
	conv := NewConverterImpl()
	statusCancelled := model.STATUS_CANCELLED
	statusFinished := model.STATUS_FINISHED
	winners := []string{
		"0xf0F8B7C21e280b0167F14Af6db4B9F90430A6C22",
		"0x6B25a67345111737d5FBEA0dd1443F64F0ef17CB",
		"0xFBCFfCE77aad086216C8016A73faeec1fE8fFEcc",
	}
	cases := []struct {
		name     string
		input    *campaignsservicev1.UpdateCampaignMsg
		expected *model.UpdateCampaignInput
	}{
		{
			name: "should map status update only",
			input: &campaignsservicev1.UpdateCampaignMsg{
				CampaignId: "campaign-1",
				Status:     campaignsservicev1.CampaignStatus_CAMPAIGN_STATUS_CANCELLED,
			},
			expected: &model.UpdateCampaignInput{
				Id:     "campaign-1",
				Stauts: &statusCancelled,
			},
		},
		{
			name: "should map status and winners",
			input: &campaignsservicev1.UpdateCampaignMsg{
				CampaignId:       "campaign-1",
				Status:           campaignsservicev1.CampaignStatus_CAMPAIGN_STATUS_FINISHED,
				EligibleAccounts: winners,
			},
			expected: &model.UpdateCampaignInput{
				Id:     "campaign-1",
				Stauts: &statusFinished,
				EligibleAccounts: &[]common.Address{
					common.HexToAddress(winners[0]),
					common.HexToAddress(winners[1]),
					common.HexToAddress(winners[2]),
				},
			},
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			result := conv.ConvertFromProtoUpdateCampaignToModelUpdateCampaign(c.input)
			assert.EqualValues(t, c.expected, result)
		})
	}
}
