package campaignsservicevalidator

import (
	"testing"
	"time"

	"github.com/block-wallet/campaigns-service/domain/model"
	campaignsservicev1 "github.com/block-wallet/campaigns-service/protos/src/campaignsservicev1/campaigns"
	"github.com/block-wallet/campaigns-service/utils/errors"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func validateTest(t *testing.T, expected errors.RichError, result errors.RichError) {
	if expected != nil {
		assert.NotNil(t, result)
		validatorErrCode := status.Code(result.ToGRPCError())
		expectedErrorCode := status.Code(expected.ToGRPCError())
		assert.Equal(t, expectedErrorCode, validatorErrCode)
	} else {
		assert.Nil(t, result)
	}
}

func Test_ValidateGetCampaigns(t *testing.T) {
	validator := NewRequestValidator()
	cases := []struct {
		name     string
		filters  *campaignsservicev1.GetCampaignsFilters
		expected errors.RichError
	}{
		{
			name: "should return error if from_date format is not valid",
			filters: &campaignsservicev1.GetCampaignsFilters{
				FromDate: &wrapperspb.StringValue{
					Value: time.Now().AddDate(0, 1, 0).Format(time.RFC1123),
				},
			},
			expected: errors.NewInvalidArgument("invalid format"),
		},
		{
			name: "should return error if end_date format is not valid",
			filters: &campaignsservicev1.GetCampaignsFilters{
				ToDate: &wrapperspb.StringValue{
					Value: time.Now().AddDate(0, 1, 0).Format(time.RFC1123),
				},
			},
			expected: errors.NewInvalidArgument("invalid format"),
		},
		{
			name:     "should not return error if there dates are empty.",
			filters:  &campaignsservicev1.GetCampaignsFilters{},
			expected: nil,
		},
		{
			name: "should not return error if from_date and to_date are correct.",
			filters: &campaignsservicev1.GetCampaignsFilters{
				FromDate: &wrapperspb.StringValue{
					Value: "2023-04-01T00:00:00Z",
				},
				ToDate: &wrapperspb.StringValue{
					Value: "2023-04-01T00:00:00Z",
				},
			},
			expected: nil,
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			result := validator.ValidateGetCampaignsRequest(&campaignsservicev1.GetCampaignsMsg{
				Filters: c.filters,
			})
			validateTest(t, c.expected, result)
		})
	}
}

func Test_ValidateGetCampaignByIdRequest(t *testing.T) {
	validator := NewRequestValidator()
	cases := []struct {
		name     string
		id       string
		expected errors.RichError
	}{
		{
			name:     "should return error if id is empty",
			id:       "",
			expected: errors.NewInvalidArgument("missing id"),
		},
		{
			name:     "should return error if id is not an UUID",
			id:       "invalid_id",
			expected: errors.NewInvalidArgument("invalid id"),
		},
		{
			name:     "should return empty errors if id is a valid UUID",
			id:       "0af4cde5-9a1b-4d3e-a1b9-d0bce781f61e",
			expected: nil,
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			result := validator.ValidateGetCampaignByIdRequest(&campaignsservicev1.GetCampaignByIdMsg{
				Id: c.id,
			})
			validateTest(t, c.expected, result)
		})
	}
}

func getCreateCampaignBase() *campaignsservicev1.CreateCampaignMsg_CreateCampaignInput {
	return &campaignsservicev1.CreateCampaignMsg_CreateCampaignInput{
		Name:            "test",
		Description:     "description test",
		IsActive:        true,
		StartDate:       time.Now().AddDate(0, -1, 0).Format(model.CampaignTimeFormatLayout),
		EndDate:         time.Now().AddDate(0, 1, 0).Format(model.CampaignTimeFormatLayout),
		Tags:            []string{"tag1", "tag2"},
		EnrollMessage:   "enroll message",
		SupportedChains: []uint32{1, 137},
		Rewards: &campaignsservicev1.CreateCampaignMsg_CreateCampaignInput_CreateCampaignRewardInput{
			Token: &campaignsservicev1.CreateCampaignMsg_CreateCampaignInput_CreateCampaignRewardInput_CreateCampaignTokenInput{
				Id: "0af4cde5-9a1b-4d3e-a1b9-d0bce781f61e",
			},
			Type:    campaignsservicev1.RewardType_REWARD_TYPE_PODIUM,
			Amounts: []string{"3", "2", "1"},
		},
		EnrollmentMode: campaignsservicev1.EnrollmentMode_INSTANCE_SINGLE_ENROLL,
		CampaignType:   campaignsservicev1.CampaignType_CAMPAIGN_TYPE_PARTNER_OFFERS,
	}
}

func Test_ValidateCreateCampaignRequest(t *testing.T) {
	validator := NewRequestValidator()
	type inputType = campaignsservicev1.CreateCampaignMsg_CreateCampaignInput

	cases := []struct {
		name      string
		generator func() *inputType
		expected  errors.RichError
	}{
		{
			name: "should fail on empty name",
			generator: func() *inputType {
				new := getCreateCampaignBase()
				new.Name = ""
				return new
			},
			expected: errors.NewInvalidArgument("invalid"),
		},
		{
			name: "should fail on empty description",
			generator: func() *inputType {
				new := getCreateCampaignBase()
				new.Description = ""
				return new
			},
			expected: errors.NewInvalidArgument("invalid"),
		},
		{
			name: "should fail on empty supported_chains",
			generator: func() *inputType {
				new := getCreateCampaignBase()
				new.SupportedChains = []uint32{}
				return new
			},
			expected: errors.NewInvalidArgument("invalid"),
		},
		{
			name: "should fail on empty rewards",
			generator: func() *inputType {
				new := getCreateCampaignBase()
				new.Rewards = nil
				return new
			},
			expected: errors.NewInvalidArgument("invalid"),
		},
		{
			name: "should fail on invalid rewards missing type",
			generator: func() *inputType {
				new := getCreateCampaignBase()
				new.Rewards.Type = campaignsservicev1.RewardType_REWARD_TYPE_INVALID
				return new
			},
			expected: errors.NewInvalidArgument("invalid"),
		},
		{
			name: "should fail on invalid rewards missing token id and create",
			generator: func() *inputType {
				new := getCreateCampaignBase()
				new.Rewards.Token.Id = ""
				new.Rewards.Token.Create = nil
				return new
			},
			expected: errors.NewInvalidArgument("invalid"),
		},
		{
			name: "should fail on invalid rewards invalid token id",
			generator: func() *inputType {
				new := getCreateCampaignBase()
				new.Rewards.Token.Id = "invalid_id"
				new.Rewards.Token.Create = nil
				return new
			},
			expected: errors.NewInvalidArgument("invalid"),
		},
		{
			name: "should fail on invalid rewards invalid token create data",
			generator: func() *inputType {
				new := getCreateCampaignBase()
				new.Rewards.Token.Id = ""
				new.Rewards.Token.Create = &campaignsservicev1.MultichainToken{
					Name:   "a token name",
					Symbol: "SYMBOL",
					//missing decimals
				}
				return new
			},
			expected: errors.NewInvalidArgument("invalid"),
		},
		{
			name: "should fail on invalid rewards amounts for dynamic reward type",
			generator: func() *inputType {
				new := getCreateCampaignBase()
				new.Rewards.Amounts = []string{"3", "2", "1"}
				new.Rewards.Type = campaignsservicev1.RewardType_REWARD_TYPE_DYNAMIC
				return new
			},
			expected: errors.NewInvalidArgument("invalid"),
		},
		{
			name: "should fail on invalid rewards amounts for single reward type",
			generator: func() *inputType {
				new := getCreateCampaignBase()
				new.Rewards.Amounts = []string{"3", "2", "1"}
				new.Rewards.Type = campaignsservicev1.RewardType_REWARD_TYPE_SINGLE
				return new
			},
			expected: errors.NewInvalidArgument("invalid"),
		},
		{
			name: "should fail on invalid rewards amounts for single reward type 2",
			generator: func() *inputType {
				new := getCreateCampaignBase()
				new.Rewards.Amounts = []string{}
				new.Rewards.Type = campaignsservicev1.RewardType_REWARD_TYPE_SINGLE
				return new
			},
			expected: errors.NewInvalidArgument("invalid"),
		},
		{
			name: "should fail on invalid rewards amounts for podium reward type",
			generator: func() *inputType {
				new := getCreateCampaignBase()
				new.Rewards.Amounts = []string{}
				new.Rewards.Type = campaignsservicev1.RewardType_REWARD_TYPE_PODIUM
				return new
			},
			expected: errors.NewInvalidArgument("invalid"),
		},
		{
			name: "should fail on empty amount",
			generator: func() *inputType {
				new := getCreateCampaignBase()
				new.Rewards.Amounts = []string{"3", "", "2"}
				new.Rewards.Type = campaignsservicev1.RewardType_REWARD_TYPE_PODIUM
				return new
			},
			expected: errors.NewInvalidArgument("invalid"),
		},
		{
			name: "should fail on invalid amount number",
			generator: func() *inputType {
				new := getCreateCampaignBase()
				new.Rewards.Amounts = []string{"3", "invalid_number", "2"}
				new.Rewards.Type = campaignsservicev1.RewardType_REWARD_TYPE_PODIUM
				return new
			},
			expected: errors.NewInvalidArgument("invalid"),
		},
		{
			name: "should fail on amount lower than 0",
			generator: func() *inputType {
				new := getCreateCampaignBase()
				new.Rewards.Amounts = []string{"3", "-1", "2"}
				new.Rewards.Type = campaignsservicev1.RewardType_REWARD_TYPE_PODIUM
				return new
			},
			expected: errors.NewInvalidArgument("invalid"),
		},
		{
			name: "should fail on unsorted amounts",
			generator: func() *inputType {
				new := getCreateCampaignBase()
				new.Rewards.Amounts = []string{"2", "1", "0"}
				new.Rewards.Type = campaignsservicev1.RewardType_REWARD_TYPE_PODIUM
				return new
			},
			expected: errors.NewInvalidArgument("invalid"),
		},
		{
			name: "should fail on empty start_date",
			generator: func() *inputType {
				new := getCreateCampaignBase()
				new.StartDate = ""
				return new
			},
			expected: errors.NewInvalidArgument("invalid"),
		},
		{
			name: "should fail on empty end_date",
			generator: func() *inputType {
				new := getCreateCampaignBase()
				new.EndDate = ""
				return new
			},
			expected: errors.NewInvalidArgument("invalid"),
		},
		{
			name: "should fail on invalid start_date format",
			generator: func() *inputType {
				new := getCreateCampaignBase()
				new.StartDate = time.Now().AddDate(0, -1, 0).Format(time.RFC1123)
				return new
			},
			expected: errors.NewInvalidArgument("invalid"),
		},
		{
			name: "should fail on invalid end_date format",
			generator: func() *inputType {
				new := getCreateCampaignBase()
				new.EndDate = time.Now().AddDate(0, 1, 0).Format(time.RFC1123)
				return new
			},
			expected: errors.NewInvalidArgument("invalid"),
		},
		{
			name: "should fail on end_date before now",
			generator: func() *inputType {
				new := getCreateCampaignBase()
				new.EndDate = time.Now().AddDate(0, -1, 0).Format(model.CampaignTimeFormatLayout)
				return new
			},
			expected: errors.NewInvalidArgument("invalid"),
		},
		{
			name: "should fail on end_date before start_date",
			generator: func() *inputType {
				new := getCreateCampaignBase()
				new.IsActive = false
				new.EndDate = time.Now().AddDate(0, 1, 0).Format(model.CampaignTimeFormatLayout)
				new.StartDate = time.Now().AddDate(0, 2, 0).Format(model.CampaignTimeFormatLayout)
				return new
			},
			expected: errors.NewInvalidArgument("invalid"),
		},
		{
			name: "should fail on future start_date and active campaign",
			generator: func() *inputType {
				new := getCreateCampaignBase()
				new.EndDate = time.Now().AddDate(0, 3, 0).Format(model.CampaignTimeFormatLayout)
				new.StartDate = time.Now().AddDate(0, 2, 0).Format(model.CampaignTimeFormatLayout)
				new.IsActive = true
				return new
			},
			expected: errors.NewInvalidArgument("invalid"),
		},
		{
			name: "should fail if campaign type is invalid",
			generator: func() *inputType {
				new := getCreateCampaignBase()
				new.CampaignType = campaignsservicev1.CampaignType_CAMPAIGN_TYPE_INVALID
				return new
			},
			expected: errors.NewInvalidArgument("invalid"),
		},
		{
			name: "should fail if campaign type is galxe and the credential_id is not specified",
			generator: func() *inputType {
				new := getCreateCampaignBase()
				new.CampaignType = campaignsservicev1.CampaignType_CAMPAIGN_TYPE_GALXE
				return new
			},
			expected: errors.NewInvalidArgument("invalid"),
		},
		{
			name:      "should pass all the validations",
			generator: func() *inputType { return getCreateCampaignBase() },
			expected:  nil,
		},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			input := c.generator()
			result := validator.ValidateCreateCampaignRequest(&campaignsservicev1.CreateCampaignMsg{
				Campaign: input,
			})
			validateTest(t, c.expected, result)
		})
	}
}

func Test_ValidateUpdateCampaignRequest(t *testing.T) {
	validator := NewRequestValidator()
	cases := []struct {
		name     string
		input    *campaignsservicev1.UpdateCampaignMsg
		expected errors.RichError
	}{
		{
			name:     "should fail on empty id",
			input:    &campaignsservicev1.UpdateCampaignMsg{},
			expected: errors.NewInvalidArgument("invalid argument"),
		},
		{
			name: "should fail on invalid id format",
			input: &campaignsservicev1.UpdateCampaignMsg{
				CampaignId: "invalid",
			},
			expected: errors.NewInvalidArgument("invalid argument"),
		},
		{
			name: "should fail if winners length is bigger than 0 and status is not specified",
			input: &campaignsservicev1.UpdateCampaignMsg{
				CampaignId:       "35ce94ce-8135-46f0-a5b1-9a09df0b6a75",
				EligibleAccounts: []string{"0x6B25a67345111737d5FBEA0dd1443F64F0ef17CB"},
			},
			expected: errors.NewInvalidArgument("invalid argument"),
		},
		{
			name: "should fail if winners length is bigger than 0 and status is active",
			input: &campaignsservicev1.UpdateCampaignMsg{
				CampaignId:       "35ce94ce-8135-46f0-a5b1-9a09df0b6a75",
				Status:           campaignsservicev1.CampaignStatus_CAMPAIGN_STATUS_ACTIVE,
				EligibleAccounts: []string{"0x6B25a67345111737d5FBEA0dd1443F64F0ef17CB"},
			},
			expected: errors.NewInvalidArgument("invalid argument"),
		},
		{
			name: "should fail if winners length is 0 and status is finished",
			input: &campaignsservicev1.UpdateCampaignMsg{
				CampaignId: "35ce94ce-8135-46f0-a5b1-9a09df0b6a75",
				Status:     campaignsservicev1.CampaignStatus_CAMPAIGN_STATUS_FINISHED,
			},
			expected: errors.NewInvalidArgument("invalid argument"),
		},
		{
			name: "shuld work updating status to active",
			input: &campaignsservicev1.UpdateCampaignMsg{
				CampaignId: "35ce94ce-8135-46f0-a5b1-9a09df0b6a75",
				Status:     campaignsservicev1.CampaignStatus_CAMPAIGN_STATUS_ACTIVE,
			},
		},
		{
			name: "shuld work updating status to waitlist",
			input: &campaignsservicev1.UpdateCampaignMsg{
				CampaignId: "35ce94ce-8135-46f0-a5b1-9a09df0b6a75",
				Status:     campaignsservicev1.CampaignStatus_CAMPAIGN_STATUS_WAITLIST,
			},
		},
		{
			name: "shuld work updating status to finished",
			input: &campaignsservicev1.UpdateCampaignMsg{
				CampaignId:       "35ce94ce-8135-46f0-a5b1-9a09df0b6a75",
				Status:           campaignsservicev1.CampaignStatus_CAMPAIGN_STATUS_FINISHED,
				EligibleAccounts: []string{"0x6B25a67345111737d5FBEA0dd1443F64F0ef17CB"},
			},
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			validateTest(t, c.expected, validator.ValidateUpdateCampaignRequest(c.input))
		})
	}
}
