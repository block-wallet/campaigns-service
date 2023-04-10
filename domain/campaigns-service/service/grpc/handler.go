package campaignsgrpcservice

import (
	"context"

	campaignservicev1service "github.com/block-wallet/campaigns-service/protos/src/campaignsservicev1/campaigns"

	campaignsconverter "github.com/block-wallet/campaigns-service/domain/campaigns-service/converter"

	campaignsservice "github.com/block-wallet/campaigns-service/domain/campaigns-service/service"
	campaignsservicevalidator "github.com/block-wallet/campaigns-service/domain/campaigns-service/validator"
	"github.com/block-wallet/campaigns-service/utils/auth"
	"github.com/block-wallet/campaigns-service/utils/errors"
	"github.com/block-wallet/campaigns-service/utils/logger"
	"github.com/block-wallet/campaigns-service/utils/signatures"
)

type Handler struct {
	service       campaignsservice.Service
	validator     campaignsservicevalidator.Validator
	converter     campaignsconverter.Converter
	authenticator auth.Auth
}

func NewHandler(service campaignsservice.Service, validator campaignsservicevalidator.Validator, converter campaignsconverter.Converter, authenticator auth.Auth) *Handler {
	return &Handler{service: service, validator: validator, converter: converter, authenticator: authenticator}
}

func (h *Handler) GetCampaigns(ctx context.Context, req *campaignservicev1service.GetCampaignsMsg) (*campaignservicev1service.GetCampaignsReply, error) {
	err := h.validator.ValidateGetCampaignsRequest(req)
	if err != nil {
		logger.Sugar.WithCtx(ctx).Errorf("Error validating GetCampaigns request: %s - Req: %v", err.Error(), req)
		return nil, err.ToGRPCError()
	}

	modelFilters, _err := h.converter.ConvertFromProtoCampaignsFiltersToModelCampaignFilters(req.GetFilters())

	if _err != nil {
		logger.Sugar.WithCtx(ctx).Errorf("Error parsing campaign filters: %s", err.Error())
		return nil, err
	}

	modelCampaigns, err := h.service.GetCampaigns(ctx, modelFilters)
	if err != nil {
		logger.Sugar.WithCtx(ctx).Errorf("Error getting campaigns: %s", err.Error())
		return nil, err.ToGRPCError()
	}

	protoCampaigns := make([]*campaignservicev1service.Campaign, 0, len(*modelCampaigns))
	for _, modelCampaign := range *modelCampaigns {
		protoCampaign := h.converter.ConvertFromModelCampaignToProtoCampaign(&modelCampaign)
		protoCampaigns = append(protoCampaigns, protoCampaign)
	}

	return &campaignservicev1service.GetCampaignsReply{Campaigns: protoCampaigns}, nil
}

func (h *Handler) GetCampaignById(ctx context.Context, req *campaignservicev1service.GetCampaignByIdMsg) (*campaignservicev1service.GetCampaignByIdReply, error) {
	err := h.validator.ValidateGetCampaignByIdRequest(req)
	if err != nil {
		logger.Sugar.WithCtx(ctx).Errorf("Error validating GetCampaigns request: %s - Req: %v", err.Error(), req)
		return nil, err.ToGRPCError()
	}

	modelCampaign, err := h.service.GetCampaignById(ctx, req.GetId())
	if err != nil {
		logger.Sugar.WithCtx(ctx).Errorf("Error getting campaign: %s with id: %s", err.Error(), req.GetId())
		return nil, err.ToGRPCError()
	}
	if modelCampaign == nil {
		logger.Sugar.WithCtx(ctx).Errorf("Campaign with id: %s not found.", err.Error(), req.GetId())
		return nil, errors.NewNotFound(req.GetId()).ToGRPCError()
	}

	protoCampaign := h.converter.ConvertFromModelCampaignToProtoCampaign(modelCampaign)
	return &campaignservicev1service.GetCampaignByIdReply{Campaign: protoCampaign}, nil
}

func (h *Handler) GetCampaignAccounts(ctx context.Context, req *campaignservicev1service.GetCampaignByIdMsg) (*campaignservicev1service.GetCampaignAccountsReply, error) {
	campaign, err := h.GetCampaignById(ctx, req)
	if err != nil {
		return nil, err
	}
	return &campaignservicev1service.GetCampaignAccountsReply{
		Accounts: campaign.Campaign.Accounts,
	}, nil
}

func (h *Handler) GetCampaignEnrollMessage(ctx context.Context, req *campaignservicev1service.GetCampaignByIdMsg) (*campaignservicev1service.GetCampaignEnrollMessageReply, error) {
	campaign, err := h.GetCampaignById(ctx, req)
	if err != nil {
		return nil, err
	}
	return &campaignservicev1service.GetCampaignEnrollMessageReply{
		Message: campaign.Campaign.EnrollMessage,
	}, nil
}

func (h *Handler) CreateCampaign(ctx context.Context, req *campaignservicev1service.CreateCampaignMsg) (*campaignservicev1service.CreateCampaignReply, error) {
	if authErr := h.authenticator.AuthenticateUsingContext(ctx); authErr != nil {
		return nil, authErr.ToGRPCError()
	}
	err := h.validator.ValidateCreateCampaignRequest(req)
	if err != nil {
		logger.Sugar.WithCtx(ctx).Errorf("Error validating CreateCampaign request: %s - Req: %v", err.Error(), req)
		return nil, err.ToGRPCError()
	}
	input, _err := h.converter.ConvertFromProtoCreateCampaignToModelCreateCampaign(req)
	if _err != nil {
		logger.Sugar.WithCtx(ctx).Errorf("Error converting CreateCampaign input", err.Error())
		return nil, errors.NewInternal(_err.Error())
	}
	modelCampaign, err := h.service.CreateCampaign(ctx, input)
	if err != nil {
		logger.Sugar.WithCtx(ctx).Errorf("Error creating campaign.", err.Error())
		return nil, err
	}

	protoCampaign := h.converter.ConvertFromModelCampaignToProtoCampaign(modelCampaign)

	return &campaignservicev1service.CreateCampaignReply{Campaign: protoCampaign}, nil
}

func (h *Handler) UpdateCampaign(ctx context.Context, req *campaignservicev1service.UpdateCampaignMsg) (*campaignservicev1service.UpdateCampaignReply, error) {
	if authErr := h.authenticator.AuthenticateUsingContext(ctx); authErr != nil {
		return nil, authErr.ToGRPCError()
	}
	err := h.validator.ValidateUpdateCampaignRequest(req)
	if err != nil {
		logger.Sugar.WithCtx(ctx).Errorf("Error validating update campaign request", err.Error())
		return nil, err
	}

	updateInput := h.converter.ConvertFromProtoUpdateCampaignToModelUpdateCampaign(req)
	modelCamaping, err := h.service.UpdateCampaign(ctx, updateInput)
	if err != nil {
		return nil, err
	}

	protoCampaign := h.converter.ConvertFromModelCampaignToProtoCampaign(modelCamaping)
	return &campaignservicev1service.UpdateCampaignReply{Campaign: protoCampaign}, nil
}

func (h *Handler) EnrollInCampaign(ctx context.Context, req *campaignservicev1service.EnrollInCampaignMsg) (*campaignservicev1service.EnrollInCampaignReply, error) {
	campaign, err := h.service.GetCampaignById(ctx, req.CampaignId)
	if err != nil {
		return nil, err.ToGRPCError()
	}

	signatureVerfier := signatures.NewSimpleMessageVerifier(campaign.EnrollMessage)
	ok, signErr := signatureVerfier.Verify(req.Signature, req.AccountAddress)

	if signErr != nil {
		logger.Sugar.WithCtx(ctx).Errorf("Error checking user's signature: %s", signErr.Error())
		return nil, errors.NewInternal(signErr.Error()).ToGRPCError()
	}
	if !ok {
		logger.Sugar.WithCtx(ctx).Errorf("Invalid signature provided: %s for account %s", req.Signature, req.AccountAddress)
		return nil, errors.NewUnauthenticated("the signature you provided is invalid. Please make sure you are signing the correct campaing's enrollment message.").ToGRPCError()
	}

	enrollInput := h.converter.ConvertFromProtoEnrollInCampaignToModelEnrollInCampaign(req)
	enrolled, _err := h.service.EnrollInCampaign(ctx, enrollInput)
	if _err != nil {
		logger.Sugar.WithCtx(ctx).Errorf("Error registering account_address:%v in campaign: %v. Error: %v", req.AccountAddress, req.CampaignId, _err.Error())
		return nil, _err.ToGRPCError()
	}

	if !*enrolled {
		return nil, errors.NewInternal("unable to register user in campaign").ToGRPCError()
	}

	return &campaignservicev1service.EnrollInCampaignReply{}, nil
}

func (h *Handler) GetTokenById(ctx context.Context, in *campaignservicev1service.GetTokenByIdMsg) (*campaignservicev1service.GetTokenByIdReply, error) {
	t, err := h.service.GetTokenById(ctx, in.GetId())
	if err != nil {
		return nil, err.ToGRPCError()
	}

	protoToken := h.converter.ConvertFromModelMultichainTokenToProtoMultichainToken(t)
	return &campaignservicev1service.GetTokenByIdReply{Token: protoToken}, nil
}

func (h *Handler) GetTokens(ctx context.Context, in *campaignservicev1service.GetTokensMsg) (*campaignservicev1service.GetTokensReply, error) {
	modelTokens, err := h.service.GetAllTokens(ctx)
	if err != nil {
		return nil, err.ToGRPCError()
	}
	protoTokens := make([]*campaignservicev1service.MultichainToken, 0, len(*modelTokens))
	for _, t := range *modelTokens {
		protoTokens = append(protoTokens, h.converter.ConvertFromModelMultichainTokenToProtoMultichainToken(&t))
	}
	return &campaignservicev1service.GetTokensReply{
		Tokens: protoTokens,
	}, nil
}
