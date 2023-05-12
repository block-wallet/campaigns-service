package campaignsservice

import (
	"context"
	"fmt"
	"time"

	"github.com/block-wallet/campaigns-service/domain/model"

	"github.com/block-wallet/campaigns-service/utils/errors"
	"github.com/block-wallet/campaigns-service/utils/logger"

	"github.com/block-wallet/campaigns-service/domain/campaigns-service/client"
	campaignsrepository "github.com/block-wallet/campaigns-service/domain/campaigns-service/repository"
)

type ServiceImpl struct {
	repository  campaignsrepository.Repository
	galxeClient client.GalxeClient
}

func NewServiceImpl(repository campaignsrepository.Repository, galxeClient client.GalxeClient) *ServiceImpl {
	return &ServiceImpl{
		repository:  repository,
		galxeClient: galxeClient,
	}
}

func (s *ServiceImpl) GetCampaigns(ctx context.Context, filters *model.GetCampaignsFilters) ([]*model.Campaign, errors.RichError) {
	if filters.Status == nil {
		// get only Active campaings by default
		filters.Status = &[]model.CampaignStatus{model.STATUS_ACTIVE}
	}
	campaigns, err := s.repository.GetCampaigns(ctx, filters)
	if err != nil {
		return nil, errors.NewInternal(err.Error())
	}
	return campaigns, nil
}

func (s *ServiceImpl) GetCampaignById(ctx context.Context, id string) (*model.Campaign, errors.RichError) {
	campaign, err := s.repository.GetCampaignById(ctx, id)
	if err != nil {
		return nil, errors.NewInternal(err.Error())
	}
	if campaign == nil {
		return nil, errors.NewNotFound(fmt.Sprintf("campaign with id: %v not found", id))
	}
	return campaign, nil
}

func (s *ServiceImpl) CreateCampaign(ctx context.Context, input *model.CreateCampaignInput) (*model.Campaign, errors.RichError) {
	//if id was specified, check that it exists
	if input.Rewards.Token.Id != nil {
		exists, err := s.repository.TokenExists(ctx, *input.Rewards.Token.Id)
		if err != nil {
			return nil, errors.NewInternal("unable to check token existance")
		}

		if !*exists {
			return nil, errors.NewInvalidArgument(fmt.Sprintf("token with id = %s does not exist.", *input.Rewards.Token.Id))
		}
	} else {
		tokenId, err := s.repository.NewToken(ctx, input.Rewards.Token.Create)
		if err != nil {
			return nil, errors.NewInternal(err.Error())
		}
		input.Rewards.Token.Id = tokenId
	}

	campaignId, err := s.repository.NewCampaign(ctx, input)
	if err != nil {
		return nil, errors.NewInternal(err.Error())
	}

	campaign, err := s.repository.GetCampaignById(ctx, *campaignId)
	if err != nil {
		return nil, errors.NewInternal(err.Error())
	}
	return campaign, nil
}

func (s *ServiceImpl) EnrollInCampaign(ctx context.Context, input *model.EnrollInCampaignInput) (bool, errors.RichError) {
	campaign, err := s.repository.GetCampaignById(ctx, input.CampaignId)
	if err != nil {
		return false, errors.NewInternal("error checking campaign existance")
	}

	if campaign == nil {
		return false, errors.NewNotFound(fmt.Sprintf("campaign with id: %v not found", input.CampaignId))
	}

	if campaign.Status != model.STATUS_ACTIVE && campaign.Status != model.STATUS_WAITLIST {
		return false, errors.NewFailedPrecondition("you can only enroll in an Active or Waitlist campaign.")
	}

	exists, err := s.repository.ParticipantExists(ctx, input.CampaignId, input.Adddress.String())
	if err != nil {
		return false, errors.NewInternal("error checking campaign participants")
	}

	if exists {
		logger.Sugar.WithCtx(ctx).Infof("Account: %s already enrolled in campaign: %v", input.Adddress, input.CampaignId)
		return exists, nil
	}

	input.EarlyEnrollment = model.STATUS_WAITLIST == campaign.Status

	ok, err := s.repository.EnrollInCampaign(ctx, input)
	if err != nil {
		return false, errors.NewInternal(err.Error())
	}

	//register user in galxe campaign
	if ok && campaign.Type == model.CAMPAIGN_TYPE_GALXE {
		populated, err := s.galxeClient.PopulateParticipant(ctx, client.PopulateParticipantsInput{
			Address:      input.Adddress,
			CredentialId: campaign.Metadata.GalxeMetadata.CredentialId,
		})

		if err != nil || !populated {
			if err != nil {
				logger.Sugar.WithCtx(ctx).Errorf("error populating participant to Galxe campaign. Error: %v ", err.Error())
			}

			s.UnenrollParticipant(ctx, &model.UnenrollFromCampaignInput{
				Adddress:   input.Adddress,
				CampaignId: input.CampaignId,
			})

			return false, errors.NewInternal("error populating participant to Galxe campaign.")
		}
	}

	return ok, nil
}

func (s *ServiceImpl) UnenrollParticipant(ctx context.Context, input *model.UnenrollFromCampaignInput) (bool, error) {
	ok, err := s.repository.UnenrollFromCampaign(ctx, input)
	if err != nil {
		logger.Sugar.WithCtx(ctx).Errorf("error unenrolling account from campaign. Error: %v", err.Error())
		return false, err
	}
	return ok, nil
}

func (s *ServiceImpl) UpdateCampaign(ctx context.Context, input *model.UpdateCampaignInput) (*model.Campaign, errors.RichError) {
	currentCampaign, err := s.repository.GetCampaignById(ctx, input.Id)
	if err != nil {
		logger.Sugar.WithCtx(ctx).Errorf("Error looking for original campaign with id: %v. Error: %v", input.Id, err.Error())
		return nil, errors.NewInternal("error getting original campaigns")
	}

	if currentCampaign == nil {
		logger.Sugar.WithCtx(ctx).Warnf("Campaign with id: %v does not exist", input.Id)
		return nil, errors.NewNotFound(fmt.Sprintf("campaign with id: %v does not exist", input.Id))
	}

	if _, _err := s.canUpdateCampaign(currentCampaign, input); _err != nil {
		return nil, _err
	}

	_, err = s.repository.UpdateCampaign(ctx, input)

	if err != nil {
		logger.Sugar.WithCtx(ctx).Errorf("Error updating campaign with id: %v. Error: %v", input.Id, err.Error())
		return nil, errors.NewInternal("error updating campaign")
	}

	return s.GetCampaignById(ctx, input.Id)
}

func (s *ServiceImpl) canUpdateCampaign(current *model.Campaign, updates *model.UpdateCampaignInput) (bool, errors.RichError) {
	if current.Status == model.STATUS_CANCELLED {
		return false, errors.NewFailedPrecondition("cannot update an already cancelled campaign")
	}

	switch current.Status {
	case model.STATUS_PENDING, model.STATUS_WAITLIST:
		{
			if updates.Stauts != nil {
				if *updates.Stauts == model.STATUS_ACTIVE {
					if current.StartDate.After(time.Now()) {
						return false, errors.NewFailedPrecondition(fmt.Sprintf("cannot activate a campaign that hasn't started yet. Campaign starts on: %v", current.StartDate))
					}
					if current.EndDate.Before(time.Now()) {
						return false, errors.NewFailedPrecondition(fmt.Sprintf("cannot activate a campaign that has already finished. Campaign ended on: %v", current.EndDate))
					}
				}
				if *updates.Stauts == model.STATUS_FINISHED {
					return false, errors.NewFailedPrecondition("cannot finalize a campaign that hasn't been active. You need to activate it first.")
				}
			}

		}
	case model.STATUS_ACTIVE:
		{
			if updates.Stauts != nil && (*updates.Stauts == model.STATUS_PENDING || *updates.Stauts == model.STATUS_WAITLIST) {
				return false, errors.NewFailedPrecondition("you can't set this campaign to PENDING or WAITLIST. You can only either CANCEL or FINISH it.")
			}
		}
	}

	if updates.Stauts != nil && *updates.Stauts == model.STATUS_FINISHED {
		if updates.EligibleAccounts != nil {
			winners := *updates.EligibleAccounts
			if len(winners) != len(current.Rewards.Amounts) && current.Rewards.Type == model.PODIUM_REWARD {
				return false, errors.NewInvalidArgument("winners length should match the rewards amounts length for a PODIUM like reward.")
			}
			participants := make(map[string]bool)
			for _, p := range current.Accounts {
				participants[p.String()] = true
			}
			for _, w := range winners {
				if !participants[w.String()] {
					return false, errors.NewInvalidArgument("all the winners should be registered in the campaign")
				}
			}
		}
	}
	return true, nil
}

func (s *ServiceImpl) GetTokenById(ctx context.Context, id string) (*model.MultichainToken, errors.RichError) {
	t, err := s.repository.GetTokenById(ctx, id)
	if err != nil {
		logger.Sugar.WithCtx(ctx).Errorf("Error getting token with id: %v. Error: %v", id, err.Error())
		return nil, errors.NewInternal("error getting token")
	}

	if t == nil {
		logger.Sugar.WithCtx(ctx).Errorf("token with id %v not found", id)
		return nil, errors.NewNotFound("token does not exist")
	}

	return t, nil
}

func (s *ServiceImpl) GetAllTokens(ctx context.Context) ([]*model.MultichainToken, errors.RichError) {
	tokens, err := s.repository.GetAllTokens(ctx)
	if err != nil {
		logger.Sugar.WithCtx(ctx).Errorf("error retrieving tokens. Error: %v", err.Error())
		return nil, errors.NewInternal("error getting tokens")
	}

	return tokens, nil
}
