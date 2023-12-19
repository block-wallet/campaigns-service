package campaignsservice

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/block-wallet/campaigns-service/domain/campaigns-service/client"
	campaignsservicemocks "github.com/block-wallet/campaigns-service/domain/campaigns-service/mocks"
	"github.com/block-wallet/campaigns-service/domain/model"
	"github.com/block-wallet/campaigns-service/utils/errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/status"
)

func getActiveCampaign() model.Campaign {
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
	return model.Campaign{
		Id:              "123",
		SupportedChains: []uint32{1, 137},
		Name:            "Test campaign",
		Description:     "Description test campaign",
		Status:          model.STATUS_ACTIVE,
		StartDate:       startDate,
		EndDate:         endDate,
		Tags:            []string{"tag1", "tag2"},
		EnrollMessage:   "Custom enroll message",
		EnrollmentMode:  model.INSTANCE_SINGLE_ENROLL,
		Type:            model.CAMPAIGN_TYPE_PARTNER_OFFERS,
		Metadata: model.CampaignMetadata{
			PartnerOffersMetadata: &model.PartnerOffersMetadata{},
		},
		Participants: []*model.CampaignParticipant{
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
	}
}

func Test_GetCampaigns(t *testing.T) {
	activeCampaign := getActiveCampaign()
	cancelledCampaign := getActiveCampaign()
	cancelledCampaign.Status = model.STATUS_CANCELLED

	statusesActiveDefault := []model.CampaignStatus{model.STATUS_ACTIVE}
	statusesCancelled := []model.CampaignStatus{model.STATUS_CANCELLED}

	activeCampaigns := []*model.Campaign{
		&activeCampaign,
	}
	cancelledCampaigns := []*model.Campaign{
		&cancelledCampaign,
	}
	cases := []struct {
		name                 string
		input                *model.GetCampaignsFilters
		repositoryInput      *model.GetCampaignsFilters
		repositoryResponse   []*model.Campaign
		expected             []*model.Campaign
		expectedRespotoryErr error
		expectedServiceErr   errors.RichError
	}{
		{
			name:               "should fetch active campaings if filters.status is not provided",
			input:              &model.GetCampaignsFilters{},
			repositoryInput:    &model.GetCampaignsFilters{Status: &statusesActiveDefault},
			expected:           activeCampaigns,
			repositoryResponse: activeCampaigns,
		},
		{
			name:               "should not add the active status filter if at least one filter.status is provided",
			input:              &model.GetCampaignsFilters{Status: &statusesCancelled},
			repositoryInput:    &model.GetCampaignsFilters{Status: &statusesCancelled},
			expected:           cancelledCampaigns,
			repositoryResponse: cancelledCampaigns,
		},
		{
			name:                 "should return internal error if the repository fails",
			input:                &model.GetCampaignsFilters{Status: &statusesCancelled},
			repositoryInput:      &model.GetCampaignsFilters{Status: &statusesCancelled},
			expectedRespotoryErr: fmt.Errorf("internal error"),
			expectedServiceErr:   errors.NewInternal("internal error"),
			repositoryResponse:   nil,
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			repository := new(campaignsservicemocks.Repository)
			galxeClient := new(campaignsservicemocks.GalxeClient)
			service := NewServiceImpl(repository, galxeClient)
			repository.On("GetCampaigns", mock.Anything, c.repositoryInput).Return(c.repositoryResponse, c.expectedRespotoryErr)
			result, err := service.GetCampaigns(context.Background(), c.input)
			if err != nil {
				assert.NotNil(t, c.expectedServiceErr)
				serviceErrCode := status.Code(err.ToGRPCError())
				expectedServErrorCode := status.Code(c.expectedServiceErr.ToGRPCError())
				assert.Equal(t, expectedServErrorCode, serviceErrCode)
			} else {
				assert.EqualValues(t, c.expected, result)
			}
		})
	}
}

func Test_GetCampaignById(t *testing.T) {
	activeCampaign := getActiveCampaign()
	cases := []struct {
		name                 string
		id                   string
		repositoryResponse   *model.Campaign
		expected             *model.Campaign
		expectedRespotoryErr error
		expectedServiceErr   errors.RichError
	}{
		{
			name:                 "should return not found if the campaign does not exists",
			id:                   "456",
			expected:             nil,
			repositoryResponse:   nil,
			expectedRespotoryErr: nil,
			expectedServiceErr:   errors.NewNotFound("not found"),
		},
		{
			name:               "should return the campaign with id=123",
			id:                 "123",
			expected:           &activeCampaign,
			repositoryResponse: &activeCampaign,
		},
		{
			name:                 "should return internal error if the repository fails",
			id:                   "id_will_fail",
			expectedRespotoryErr: fmt.Errorf("internal error"),
			expectedServiceErr:   errors.NewInternal("internal error"),
			repositoryResponse:   nil,
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			repository := new(campaignsservicemocks.Repository)
			galxeClient := new(campaignsservicemocks.GalxeClient)
			service := NewServiceImpl(repository, galxeClient)
			repository.On("GetCampaignById", mock.Anything, c.id).Return(c.repositoryResponse, c.expectedRespotoryErr)
			result, err := service.GetCampaignById(context.Background(), c.id)
			if err != nil {
				assert.NotNil(t, c.expectedServiceErr)
				serviceErrCode := status.Code(err.ToGRPCError())
				expectedServErrorCode := status.Code(c.expectedServiceErr.ToGRPCError())
				assert.Equal(t, expectedServErrorCode, serviceErrCode)
			} else {
				assert.EqualValues(t, c.expected, result)
			}
		})
	}
}

func Test_EnrollInCampaignTest(t *testing.T) {
	galxeCredential := "galxe-credential-123"
	activeCampaign := getActiveCampaign()
	activeGalxeCampaign := getActiveCampaign()
	activeGalxeCampaign.Type = model.CAMPAIGN_TYPE_GALXE
	activeGalxeCampaign.Metadata = model.CampaignMetadata{
		GalxeMetadata: &model.GalxeCampaignMetadata{
			CredentialId: galxeCredential,
		},
	}

	cancelledCampaign := getActiveCampaign()
	cancelledCampaign.Status = model.STATUS_CANCELLED
	type repositoryMock struct {
		getByIdExpectedResponse    *model.Campaign
		getByIdExpectedErr         error
		enrollExpectedResponse     bool
		enrollExpectedErr          error
		participantsExistsResponse bool
		participantsExistsErr      error
		shouldCallUnenroll         bool
		unenrollExpectedResponse   bool
		unenrollExpectedErr        error
	}
	type galxeMock struct {
		populateParticipantExpectedInput     client.PopulateParticipantsInput
		populateParticipantsExpectedResponse bool
		populateParticipantsExpectedErr      error
	}
	cases := []struct {
		name               string
		expectedServiceRes bool
		expectedServiceErr errors.RichError
		repository         repositoryMock
		id                 string
		accountAddress     string
		galxeMock          *galxeMock
	}{
		{
			name:               "should return internal error if the database fails getting campaign by id",
			id:                 "123",
			accountAddress:     "0xf0F8B7C21e280b0167F14Af6db4B9F90430A6C22",
			expectedServiceRes: false,
			expectedServiceErr: errors.NewInternal("internal error"),
			repository: repositoryMock{
				getByIdExpectedErr: fmt.Errorf("internal error"),
			},
		},
		{
			name:               "should return not found if the campaign does not exists",
			expectedServiceRes: false,
			id:                 "890",
			accountAddress:     "0xf0F8B7C21e280b0167F14Af6db4B9F90430A6C22",
			expectedServiceErr: errors.NewNotFound("campaign not found"),
			repository: repositoryMock{
				getByIdExpectedResponse: nil,
			},
		},
		{
			name:               "should return validation error if the campaign is not pending or active",
			expectedServiceRes: false,
			id:                 "123",
			accountAddress:     "0xf0F8B7C21e280b0167F14Af6db4B9F90430A6C22",
			expectedServiceErr: errors.NewFailedPrecondition("cannot enroll"),
			repository: repositoryMock{
				getByIdExpectedResponse: &cancelledCampaign,
			},
		},
		{
			name:               "should return internal error if participants check fails",
			expectedServiceErr: errors.NewInternal("internal error"),
			id:                 "123",
			accountAddress:     "0xf0F8B7C21e280b0167F14Af6db4B9F90430A6C22",
			repository: repositoryMock{
				getByIdExpectedResponse: &activeCampaign,
				participantsExistsErr:   fmt.Errorf("error checking participants"),
			},
		},
		{
			name:               "should return true if the account is already enrolled",
			expectedServiceRes: true,
			id:                 "123",
			accountAddress:     "0xf0F8B7C21e280b0167F14Af6db4B9F90430A6C22",
			repository: repositoryMock{
				getByIdExpectedResponse:    &activeCampaign,
				participantsExistsResponse: true,
			},
		},
		{
			name:               "should return internal error if database fails registering account in the campaign",
			expectedServiceErr: errors.NewInternal("internal error"),
			id:                 "123",
			accountAddress:     "0xf0F8B7C21e280b0167F14Af6db4B9F90430A6C22",
			repository: repositoryMock{
				getByIdExpectedResponse:    &activeCampaign,
				participantsExistsResponse: false,
				enrollExpectedErr:          fmt.Errorf("error"),
			},
		},
		{
			name:               "should enroll the account in the requested campaign",
			expectedServiceRes: true,
			id:                 "123",
			accountAddress:     "0xf0F8B7C21e280b0167F14Af6db4B9F90430A6C22",
			repository: repositoryMock{
				getByIdExpectedResponse:    &activeCampaign,
				participantsExistsResponse: false,
				enrollExpectedResponse:     true,
			},
		},
		{
			name:               "should unenroll the account if galxe participant population fails",
			expectedServiceRes: false,
			expectedServiceErr: errors.NewInternal("failed"),
			id:                 "101112",
			accountAddress:     "0xf0F8B7C21e280b0167F14Af6db4B9F90430A6C22",
			repository: repositoryMock{
				getByIdExpectedResponse:    &activeGalxeCampaign,
				participantsExistsResponse: false,
				enrollExpectedResponse:     true,
				shouldCallUnenroll:         true,
				unenrollExpectedResponse:   true,
				unenrollExpectedErr:        nil,
			},
			galxeMock: &galxeMock{
				populateParticipantExpectedInput: client.PopulateParticipantsInput{
					Address:      common.HexToAddress("0xf0F8B7C21e280b0167F14Af6db4B9F90430A6C22"),
					CredentialId: activeGalxeCampaign.Metadata.GalxeMetadata.CredentialId,
				},
				populateParticipantsExpectedResponse: false,
				populateParticipantsExpectedErr:      fmt.Errorf("galxe population err"),
			},
		},
		{
			name:               "should return ok if enroll and galxe population works",
			expectedServiceRes: true,
			id:                 "101112",
			accountAddress:     "0xf0F8B7C21e280b0167F14Af6db4B9F90430A6C22",
			repository: repositoryMock{
				getByIdExpectedResponse:    &activeGalxeCampaign,
				participantsExistsResponse: false,
				enrollExpectedResponse:     true,
				shouldCallUnenroll:         false,
			},
			galxeMock: &galxeMock{
				populateParticipantExpectedInput: client.PopulateParticipantsInput{
					Address:      common.HexToAddress("0xf0F8B7C21e280b0167F14Af6db4B9F90430A6C22"),
					CredentialId: activeGalxeCampaign.Metadata.GalxeMetadata.CredentialId,
				},
				populateParticipantsExpectedResponse: true,
			},
		},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			repository := new(campaignsservicemocks.Repository)
			galxeClient := new(campaignsservicemocks.GalxeClient)
			service := NewServiceImpl(repository, galxeClient)
			input := &model.EnrollInCampaignInput{
				Adddress:   common.HexToAddress(c.accountAddress),
				CampaignId: c.id,
			}
			repository.On("GetCampaignById", mock.Anything, c.id).Return(c.repository.getByIdExpectedResponse, c.repository.getByIdExpectedErr)
			repository.On("ParticipantExists", mock.Anything, c.id, input.Adddress.String()).Return(c.repository.participantsExistsResponse, c.repository.participantsExistsErr)
			repository.On("EnrollInCampaign", mock.Anything, input).Return(c.repository.enrollExpectedResponse, c.repository.enrollExpectedErr)
			repository.On("UnenrollFromCampaign", mock.Anything, &model.UnenrollFromCampaignInput{
				Adddress:   common.HexToAddress(c.accountAddress),
				CampaignId: c.id,
			}).Return(c.repository.unenrollExpectedResponse, c.repository.unenrollExpectedErr)

			if c.galxeMock != nil {
				galxeClient.On("PopulateParticipant", mock.Anything, c.galxeMock.populateParticipantExpectedInput).Return(c.galxeMock.populateParticipantsExpectedResponse, c.repository.getByIdExpectedErr)
			}

			result, err := service.EnrollInCampaign(context.Background(), input)
			if err != nil {
				assert.NotNil(t, c.expectedServiceErr)
				serviceErrCode := status.Code(err.ToGRPCError())
				expectedServErrorCode := status.Code(c.expectedServiceErr.ToGRPCError())
				assert.Equal(t, expectedServErrorCode, serviceErrCode)
			} else {
				assert.EqualValues(t, c.expectedServiceRes, result)
				unenrollTimes := 0
				if c.repository.shouldCallUnenroll {
					unenrollTimes = 1
				}
				repository.AssertNumberOfCalls(t, "UnenrollFromCampaign", unenrollTimes)
			}
		})
	}
}

func Test_CreateCampaign(t *testing.T) {
	tokenId1 := "token-123"
	tokenId2 := "token-456"
	rewardWithTokenId := model.CampaignRewardTokenInput{
		Id: &tokenId1,
	}
	rewardWithTokenData := model.CampaignRewardTokenInput{
		Create: &model.MultichainToken{
			Name:     "testToken",
			Symbol:   "TESTTOKEN",
			Decimals: 18,
		},
	}

	createCampaignInputWithTokenId := model.CreateCampaignInput{
		Name:            "test",
		Description:     "test desc",
		StartDate:       "2023-04-01T00:00:00Z",
		EndDate:         "2023-05-01T00:00:00Z",
		Tags:            []string{"tag1"},
		SupportedChains: []uint32{1, 137},
		Status:          model.STATUS_ACTIVE,
		EnrollMessage:   "custom",
		Rewards: model.CampaignRewardInput{
			Amounts: []string{"1", "2", "3"},
			Token:   rewardWithTokenId,
			Type:    model.PODIUM_REWARD,
		},
	}

	createCampaignInputWithTokenData := createCampaignInputWithTokenId
	createCampaignInputWithTokenData.Rewards.Token = rewardWithTokenData

	createCampaignInputWithTokenDataAndId := createCampaignInputWithTokenData
	createCampaignInputWithTokenDataAndId.Rewards.Token.Id = &tokenId2

	activeCampaign := getActiveCampaign()
	trueBool := true
	falseBool := false
	type repositoryMock struct {
		tokenExistsRes     *bool
		tokenExistsErr     error
		newTokenRes        *string
		newTokenErr        error
		newCampaignInput   *model.CreateCampaignInput
		newCampaignRes     *string
		newCampaignErr     error
		getCampaignByIdRes *model.Campaign
		getCampaignByIdErr error
	}
	cases := []struct {
		name               string
		input              *model.CreateCampaignInput
		expectedServiceRes *model.Campaign
		expectedServiceErr errors.RichError
		repository         repositoryMock
	}{
		//failure cases
		{
			name:               "should fail if token existance check fails",
			input:              &createCampaignInputWithTokenId,
			expectedServiceErr: errors.NewInternal("internal error"),
			expectedServiceRes: nil,
			repository: repositoryMock{
				tokenExistsErr: fmt.Errorf("error"),
			},
		},
		{
			name:               "should return not found if the token provided in the reward does not exists",
			input:              &createCampaignInputWithTokenId,
			expectedServiceErr: errors.NewInvalidArgument("token not found"),
			expectedServiceRes: nil,
			repository: repositoryMock{
				tokenExistsRes: &falseBool,
			},
		},
		{
			name:               "should return internal error if token creation fails",
			input:              &createCampaignInputWithTokenData,
			expectedServiceErr: errors.NewInternal("create token failed"),
			expectedServiceRes: nil,
			repository: repositoryMock{
				newTokenErr: fmt.Errorf("error creating token"),
			},
		},
		{
			name:               "should return error if campaign creation fails",
			input:              &createCampaignInputWithTokenId,
			expectedServiceErr: errors.NewInternal("create campaign failed"),
			expectedServiceRes: nil,
			repository: repositoryMock{
				newCampaignInput: &createCampaignInputWithTokenId,
				tokenExistsRes:   &trueBool,
				newCampaignErr:   fmt.Errorf("i have failed you"),
			},
		},
		{
			name:               "should return error if getCampaignById fails after campaign creation",
			input:              &createCampaignInputWithTokenId,
			expectedServiceErr: errors.NewInternal("get campaign failed"),
			expectedServiceRes: nil,
			repository: repositoryMock{
				tokenExistsRes:     &trueBool,
				newCampaignInput:   &createCampaignInputWithTokenId,
				newCampaignRes:     &activeCampaign.Id,
				getCampaignByIdErr: fmt.Errorf("error getting campaign"),
			},
		},
		//ok cases
		{
			name:               "should create and return the new campaign",
			input:              &createCampaignInputWithTokenId,
			expectedServiceRes: &activeCampaign,
			repository: repositoryMock{
				newCampaignInput:   &createCampaignInputWithTokenId,
				tokenExistsRes:     &trueBool,
				newCampaignRes:     &activeCampaign.Id,
				getCampaignByIdRes: &activeCampaign,
			},
		},
		{
			name:               "should create token and return the new campaign",
			input:              &createCampaignInputWithTokenData,
			expectedServiceRes: &activeCampaign,
			repository: repositoryMock{
				newTokenRes:        &tokenId2,
				newCampaignInput:   &createCampaignInputWithTokenDataAndId,
				newCampaignRes:     &activeCampaign.Id,
				getCampaignByIdRes: &activeCampaign,
			},
		},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			repository := new(campaignsservicemocks.Repository)
			galxeClient := new(campaignsservicemocks.GalxeClient)
			service := NewServiceImpl(repository, galxeClient)
			if c.input.Rewards.Token.Id != nil {
				repository.On("TokenExists", mock.Anything, *c.input.Rewards.Token.Id).Return(c.repository.tokenExistsRes, c.repository.tokenExistsErr)
			}
			repository.On("NewToken", mock.Anything, c.input.Rewards.Token.Create).Return(c.repository.newTokenRes, c.repository.newTokenErr)
			repository.On("NewCampaign", mock.Anything, c.repository.newCampaignInput).Return(c.repository.newCampaignRes, c.repository.newCampaignErr)
			if c.repository.newCampaignRes != nil {
				repository.On("GetCampaignById", mock.Anything, *c.repository.newCampaignRes).Return(c.repository.getCampaignByIdRes, c.repository.getCampaignByIdErr)
			}

			result, err := service.CreateCampaign(context.Background(), c.input)

			if err != nil {
				assert.NotNil(t, c.expectedServiceErr)
				serviceErrCode := status.Code(err.ToGRPCError())
				expectedServErrorCode := status.Code(c.expectedServiceErr.ToGRPCError())
				assert.Equal(t, expectedServErrorCode, serviceErrCode)
			} else {
				assert.EqualValues(t, c.expectedServiceRes, result)
			}
		})
	}
}

func Test_UpdateCampaign(t *testing.T) {
	type repositoryMock struct {
		getCampaignByIdRes *model.Campaign
		getCampaignByIdErr error
		updateCampaignRes  *bool
		updateCampaignErr  error
	}
	statusFinished := model.STATUS_FINISHED

	statusActive := model.STATUS_ACTIVE
	statusPending := model.STATUS_PENDING

	activeCampaign := getActiveCampaign()

	cancelledCampaign := getActiveCampaign()
	cancelledCampaign.Status = model.STATUS_CANCELLED

	pendingCampaign := getActiveCampaign()
	pendingCampaign.Status = model.STATUS_PENDING

	pendingCampaignFutureStartDate := pendingCampaign
	pendingCampaignFutureStartDate.StartDate = time.Now().AddDate(0, 1, 0)

	pendingCampaignPastEndDate := pendingCampaign
	pendingCampaignPastEndDate.EndDate = time.Now().AddDate(0, -1, 0)

	trueBool := true

	cases := []struct {
		name               string
		input              *model.UpdateCampaignInput
		expectedServiceRes *model.Campaign
		expectedServiceErr errors.RichError
		repository         repositoryMock
	}{
		//basic errors
		{
			name:               "should return internal error if database fails",
			expectedServiceErr: errors.NewInternal("internal"),
			input: &model.UpdateCampaignInput{
				Id:     cancelledCampaign.Id,
				Status: &statusActive,
			},
			repository: repositoryMock{
				getCampaignByIdErr: fmt.Errorf("err"),
			},
		},
		{
			name: "should return not found error if campaign does not exists",
			input: &model.UpdateCampaignInput{
				Id:     cancelledCampaign.Id,
				Status: &statusActive,
			},
			expectedServiceErr: errors.NewNotFound("campaign not found"),
			repository: repositoryMock{
				getCampaignByIdRes: nil,
			},
		},
		//validations
		{
			name: "cannot update a cancelled campaign",
			input: &model.UpdateCampaignInput{
				Id:     cancelledCampaign.Id,
				Status: &statusActive,
			},
			expectedServiceErr: errors.NewFailedPrecondition("cannot update cancelled campaign"),
			repository: repositoryMock{
				getCampaignByIdRes: &cancelledCampaign,
			},
		},
		{
			name: "cannot set status finished to a non-active campaign",
			input: &model.UpdateCampaignInput{
				Id:     pendingCampaign.Id,
				Status: &statusFinished,
			},
			expectedServiceErr: errors.NewFailedPrecondition("invalid"),
			repository: repositoryMock{
				getCampaignByIdRes: &pendingCampaign,
			},
		},
		{
			name: "cannot activate a campaign which start_date is in the future",
			input: &model.UpdateCampaignInput{
				Id:     pendingCampaignFutureStartDate.Id,
				Status: &statusActive,
			},
			expectedServiceErr: errors.NewFailedPrecondition("invalid"),
			repository: repositoryMock{
				getCampaignByIdRes: &pendingCampaignFutureStartDate,
			},
		},
		{
			name: "cannot activate a campaign which end_date is in the past",
			input: &model.UpdateCampaignInput{
				Id:     pendingCampaignPastEndDate.Id,
				Status: &statusActive,
			},
			expectedServiceErr: errors.NewFailedPrecondition("invalid"),
			repository: repositoryMock{
				getCampaignByIdRes: &pendingCampaignPastEndDate,
			},
		},
		{
			name: "cannot set to pending an active campaign",
			input: &model.UpdateCampaignInput{
				Id:     activeCampaign.Id,
				Status: &statusPending,
			},
			expectedServiceErr: errors.NewFailedPrecondition("invalid"),
			repository: repositoryMock{
				getCampaignByIdRes: &activeCampaign,
			},
		},
		{
			name: "should return error if the number of elegible accounts does not match the amount of rewards in a podium like campaign",
			input: &model.UpdateCampaignInput{
				Id:               activeCampaign.Id,
				Status:           &statusFinished,
				EligibleAccounts: &[]common.Address{activeCampaign.Participants[0].AccountAddress},
			},
			expectedServiceErr: errors.NewInvalidArgument("invalid"),
			repository: repositoryMock{
				getCampaignByIdRes: &activeCampaign,
			},
		},
		{
			name: "should return error if one of the elegible accounts is not registered in the campaign",
			input: &model.UpdateCampaignInput{
				Id:               activeCampaign.Id,
				Status:           &statusFinished,
				EligibleAccounts: &[]common.Address{common.HexToAddress("0xB1e8eB3bd367095F1eD945ba8bf67cc698D958c9")},
			},
			expectedServiceErr: errors.NewInvalidArgument("invalid"),
			repository: repositoryMock{
				getCampaignByIdRes: &activeCampaign,
			},
		},
		//validations ok
		{
			name: "should return the updated campaign",
			input: &model.UpdateCampaignInput{
				Id:     activeCampaign.Id,
				Status: &statusFinished,
				EligibleAccounts: &[]common.Address{
					activeCampaign.Participants[0].AccountAddress,
					activeCampaign.Participants[1].AccountAddress,
					activeCampaign.Participants[2].AccountAddress,
				},
			},
			expectedServiceRes: &activeCampaign,
			repository: repositoryMock{
				getCampaignByIdRes: &activeCampaign,
				updateCampaignRes:  &trueBool,
			},
		},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			repository := new(campaignsservicemocks.Repository)
			galxeClient := new(campaignsservicemocks.GalxeClient)
			service := NewServiceImpl(repository, galxeClient)
			repository.On("GetCampaignById", mock.Anything, c.input.Id).Return(c.repository.getCampaignByIdRes, c.repository.getCampaignByIdErr)
			repository.On("UpdateCampaign", mock.Anything, c.input).Return(c.repository.updateCampaignRes, c.repository.updateCampaignErr)

			result, err := service.UpdateCampaign(context.Background(), c.input)
			if err != nil {
				assert.NotNil(t, c.expectedServiceErr)
				serviceErrCode := status.Code(err.ToGRPCError())
				expectedServErrorCode := status.Code(c.expectedServiceErr.ToGRPCError())
				assert.Equal(t, expectedServErrorCode, serviceErrCode)
			} else {
				assert.EqualValues(t, c.expectedServiceRes, result)
			}
		})
	}
}

func Test_GetTokenById(t *testing.T) {
	token := &model.MultichainToken{
		Id:       "token-123",
		Name:     "test-token",
		Symbol:   "TESTTOKEN",
		Decimals: 18,
	}
	type repositoryMock struct {
		getTokenByIdRes *model.MultichainToken
		getTokenByIdErr error
	}

	cases := []struct {
		name               string
		id                 string
		expectedServiceRes *model.MultichainToken
		expectedServiceErr errors.RichError
		repository         repositoryMock
	}{
		{
			name:               "should return internal error if database fails",
			expectedServiceErr: errors.NewInternal("internal"),
			repository: repositoryMock{
				getTokenByIdErr: fmt.Errorf("err"),
			},
		},
		{
			name:               "should return not found if token does not exists",
			expectedServiceErr: errors.NewNotFound("token not found"),
			repository: repositoryMock{
				getTokenByIdRes: nil,
				getTokenByIdErr: nil,
			},
		},
		{
			name:               "should return the requested token",
			expectedServiceRes: token,
			repository: repositoryMock{
				getTokenByIdRes: token,
			},
		},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			repository := new(campaignsservicemocks.Repository)
			galxeClient := new(campaignsservicemocks.GalxeClient)
			service := NewServiceImpl(repository, galxeClient)
			repository.On("GetTokenById", mock.Anything, c.id).Return(c.repository.getTokenByIdRes, c.repository.getTokenByIdErr)
			result, err := service.GetTokenById(context.Background(), c.id)
			if err != nil {
				assert.NotNil(t, c.expectedServiceErr)
				serviceErrCode := status.Code(err.ToGRPCError())
				expectedServErrorCode := status.Code(c.expectedServiceErr.ToGRPCError())
				assert.Equal(t, expectedServErrorCode, serviceErrCode)
			} else {
				assert.EqualValues(t, c.expectedServiceRes, result)
			}
		})
	}
}

func Test_GetAllTokens(t *testing.T) {
	token := model.MultichainToken{
		Id:       "token-123",
		Name:     "test-token",
		Symbol:   "TESTTOKEN",
		Decimals: 18,
	}
	token2 := model.MultichainToken{
		Id:       "token-456",
		Name:     "second-token",
		Symbol:   "SECTOK",
		Decimals: 6,
	}
	allTokens := []*model.MultichainToken{&token, &token2}
	type repositoryMock struct {
		getAllTokensRes []*model.MultichainToken
		getAllTokensErr error
	}
	cases := []struct {
		name               string
		expectedServiceRes []*model.MultichainToken
		expectedServiceErr errors.RichError
		repository         repositoryMock
	}{
		{
			name:               "should return internal error if database fails",
			expectedServiceErr: errors.NewInternal("internal err"),
			repository: repositoryMock{
				getAllTokensErr: fmt.Errorf("error"),
			},
		},
		{
			name:               "should return all the tokens",
			expectedServiceRes: allTokens,
			repository: repositoryMock{
				getAllTokensRes: allTokens,
			},
		},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			repository := new(campaignsservicemocks.Repository)
			galxeClient := new(campaignsservicemocks.GalxeClient)
			service := NewServiceImpl(repository, galxeClient)
			repository.On("GetAllTokens", mock.Anything).Return(c.repository.getAllTokensRes, c.repository.getAllTokensErr)
			result, err := service.GetAllTokens(context.Background())
			if err != nil {
				assert.NotNil(t, c.expectedServiceErr)
				serviceErrCode := status.Code(err.ToGRPCError())
				expectedServErrorCode := status.Code(c.expectedServiceErr.ToGRPCError())
				assert.Equal(t, expectedServErrorCode, serviceErrCode)
			} else {
				assert.EqualValues(t, c.expectedServiceRes, result)
			}
		})
	}
}
