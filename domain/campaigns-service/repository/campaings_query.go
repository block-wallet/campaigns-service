package campaignsrepository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/block-wallet/campaigns-service/domain/model"
	"github.com/block-wallet/campaigns-service/utils/logger"
	"github.com/ethereum/go-ethereum/common"
)

const SPLIT_TOKEN = ","

type QueryBuilder[P interface{}, F interface{}] interface {
	Query(ctx context.Context) (string, []string)
	Parse(ctx context.Context, row *sql.Rows) (*P, error)
}

type CampaignsQueryBuilder struct {
	filters *model.GetCampaignsFilters
	params  []any
}

func NewCampaignsQueryBuilder(filters *model.GetCampaignsFilters) *CampaignsQueryBuilder {
	return &CampaignsQueryBuilder{
		filters: filters,
		params:  make([]any, 0),
	}
}

func (r *CampaignsQueryBuilder) Query(ctx context.Context) (string, []any) {
	campaignSelectFields := "c.id, c.name, c.description, c.status, c.start_date, c.end_date, c.enroll_message, c.enrollment_mode, c.campaign_type, c.external_campaign_id, c.created_at, c.updated_at"
	tagsSelectFields := "string_agg(distinct ct.tag, ',') as tags"
	supportedChainFields := "string_agg(distinct csc.chain_id, ',') as supported_chains"
	tokenSelectFields := "t.id as token_id, t.name as token_name, t.symbol as token_symbol, t.decimals as token_decimals"
	rewardsSelectFields := "r.reward_id as reward_id, r.amounts as reward_amounts, r.type as reward_type"
	participantsSelectStatement := "json_agg(distinct jsonb_build_object('account_address',p.account_address,'position',p.position,'early_enrollment', p.early_enrollment)) participants"
	campaignsSelect := fmt.Sprintf("SELECT %s, %s, %s, %s, %s, %s", campaignSelectFields, tagsSelectFields, supportedChainFields, tokenSelectFields, rewardsSelectFields, participantsSelectStatement)
	fromAndJoinStatements := `from campaigns c 
	LEFT JOIN rewards r on c.id = r.campaign_id  
	LEFT JOIN tokens t on t.id = r.token_id
	LEFT JOIN participants p on p.campaign_id = c.id
	LEFT JOIN campaigns_tags ct on ct.campaign_id = c.id
	LEFT JOIN campaigns_supported_chains csc on csc.campaign_id = c.id
	`
	groupStatement := "GROUP by c.id, t.id,r.reward_id"
	q := fmt.Sprintf("%s %s", campaignsSelect, fromAndJoinStatements)
	queryWithFilters := r.withFilters(ctx, q)
	q = fmt.Sprintf("%s %s;", queryWithFilters, groupStatement)
	return q, r.params
}

func (r *CampaignsQueryBuilder) withFilters(ctx context.Context, query string) string {
	extraJoin := ""
	filtersQuery := "WHERE 1=1"
	if r.filters != nil {
		//first check filters that will need an extra join
		if r.filters.ChainIds != nil && len(*r.filters.ChainIds) > 0 {
			chains := r.filters.ChainIds
			chainsInStr := ""
			for _, chain := range *chains {
				queryArgument := r.newQueryArgument(chain)
				if chainsInStr != "" {
					chainsInStr = fmt.Sprintf("%s,%s", chainsInStr, queryArgument)
				} else {
					chainsInStr = queryArgument
				}
			}
			//extra join for filtering by tags
			extraJoin = fmt.Sprintf("%s LEFT JOIN campaigns_supported_chains cscf on cscf.campaign_id = c.id", extraJoin)
			filtersQuery = fmt.Sprintf("%s AND cscf.chain_id IN (%s)", filtersQuery, chainsInStr)
		}

		if r.filters.Tags != nil && len(*r.filters.Tags) > 0 {
			tags := r.filters.Tags
			tagsInStr := ""
			for _, tag := range *tags {
				queryArg := r.newQueryArgument(string(tag))
				if tagsInStr != "" {
					tagsInStr = fmt.Sprintf("%s,%s", tagsInStr, queryArg)
				} else {
					tagsInStr = queryArg
				}
			}
			//extra join for filtering by tags
			extraJoin = fmt.Sprintf("%s LEFT JOIN campaigns_tags ctf on ctf.campaign_id = c.id", extraJoin)
			filtersQuery = fmt.Sprintf("%s AND ctf.tag IN (%s)", filtersQuery, tagsInStr)
		}

		if r.filters.Status != nil {
			status := r.filters.Status
			statusInStr := ""
			for _, status := range *status {
				queryArg := r.newQueryArgument(string(status))
				if statusInStr != "" {
					statusInStr = fmt.Sprintf("%s,%s", statusInStr, queryArg)
				} else {
					statusInStr = queryArg
				}
			}
			filtersQuery = fmt.Sprintf("%s AND c.status IN (%s)", filtersQuery, statusInStr)
		}
		if r.filters.FromDate != nil {
			queryArg := r.newQueryArgument(r.filters.FromDate.Format(model.CampaignTimeFormatLayout))
			filtersQuery = fmt.Sprintf("%s AND c.start_date::date >= %s::date", filtersQuery, queryArg)
		}

		if r.filters.ToDate != nil {
			queryArg := r.newQueryArgument(r.filters.ToDate.Format(model.CampaignTimeFormatLayout))
			filtersQuery = fmt.Sprintf("%s AND stc.end_date::date <= %s::date", filtersQuery, queryArg)
		}
		if r.filters.Id != nil {
			queryArg := r.newQueryArgument(*r.filters.Id)
			filtersQuery = fmt.Sprintf("%s AND c.id = %s", filtersQuery, queryArg)
		}
	}

	return fmt.Sprintf("%s %s %s", query, extraJoin, filtersQuery)
}

func (r *CampaignsQueryBuilder) Parse(ctx context.Context, rows *sql.Rows) (*model.Campaign, error) {
	var parsedRow campaignrow
	err := rows.Scan(&parsedRow.id, &parsedRow.name, &parsedRow.description, &parsedRow.status,
		&parsedRow.startDate, &parsedRow.endDate, &parsedRow.enrollMessage, &parsedRow.enrollmentMode, &parsedRow.campaignType, &parsedRow.externalCampaignId, &parsedRow.createdAt, &parsedRow.updatedAt, &parsedRow.tags, &parsedRow.supportedChains,
		&parsedRow.tokenId, &parsedRow.tokenName, &parsedRow.tokenSymbol, &parsedRow.decimals, &parsedRow.rewardId, &parsedRow.amounts, &parsedRow.rewardType, &parsedRow.participants)
	if err != nil {
		return nil, err
	}
	return r.rowToCampaign(parsedRow)
}

func (r *CampaignsQueryBuilder) newQueryArgument(arg any) string {
	r.params = append(r.params, arg)
	return fmt.Sprintf("$%v", len(r.params))
}

func (r *CampaignsQueryBuilder) rowToCampaign(row campaignrow) (*model.Campaign, error) {
	startDate, err := time.Parse(model.CampaignTimeFormatLayout, row.startDate)

	if err != nil {
		return nil, err
	}
	endDate, err := time.Parse(model.CampaignTimeFormatLayout, row.endDate)

	if err != nil {
		return nil, err
	}
	campaign := model.Campaign{
		Id:             row.id,
		Name:           row.name,
		Description:    row.description,
		Type:           model.CampaignType(row.campaignType),
		Status:         model.CampaignStatus(row.status),
		StartDate:      startDate,
		EndDate:        endDate,
		EnrollMessage:  *row.enrollMessage,
		EnrollmentMode: model.EnrollmentMode(row.enrollmentMode),
		CreatedAt:      row.createdAt,
		UpdatedAt:      row.updatedAt,
	}

	var campaignMetadata model.CampaignMetadata

	switch campaign.Type {
	case model.CAMPAIGN_TYPE_GALXE:
		{
			campaignMetadata.GalxeMetadata = &model.GalxeCampaignMetadata{
				CredentialId: *row.externalCampaignId,
			}
		}
	case model.CAMPAIGN_TYPE_PARTNER_OFFERS:
		{
			campaignMetadata.PartnerOffersMetadata = &model.PartnerOffersMetadata{}
		}
	}

	campaign.Metadata = campaignMetadata

	if row.supportedChains != "" {
		chainsStr := strings.Split(row.supportedChains, SPLIT_TOKEN)
		supportedChains := make([]uint32, 0, len(chainsStr))
		for _, chainStr := range chainsStr {
			chainId, err := strconv.Atoi(chainStr)
			if err != nil {
				return nil, err
			}
			supportedChains = append(supportedChains, uint32(int32(chainId)))
		}
		campaign.SupportedChains = supportedChains
	}

	if row.tokenName != nil && row.amounts != nil {
		amountsStr := strings.Split(*row.amounts, SPLIT_TOKEN)
		amounts := make([]*big.Int, 0, len(amountsStr))
		for _, amount := range amountsStr {
			bigAmount := new(big.Int)
			bigAmount, ok := bigAmount.SetString(amount, 10)
			if !ok {
				return nil, fmt.Errorf("invalid reward amount %v", amount)
			}
			amounts = append(amounts, bigAmount)
		}
		campaign.Rewards = &model.Reward{
			Token: &model.MultichainToken{
				Name:              *row.tokenName,
				Symbol:            *row.tokenSymbol,
				Decimals:          uint8(*row.decimals),
				ContractAddresses: map[string]common.Address{},
			},
			Amounts: amounts,
			Type:    model.RewardType(*row.rewardType),
		}
	}

	if row.tags != nil {
		campaign.Tags = strings.Split(*row.tags, SPLIT_TOKEN)
	}

	if row.participants != nil {
		var data = make([]*participantJSONRow, 0)
		participants := make([]*model.CampaignParticipant, 0)
		if err = json.Unmarshal(*row.participants, &data); err != nil {
			logger.Sugar.Warnf("unable to parse participants data for campaign: %v. Error: %v", row.id, err.Error())
		} else {
			for _, d := range data {
				if d.AccountAddress != "" {
					participants = append(participants, &model.CampaignParticipant{
						AccountAddress:  common.HexToAddress(d.AccountAddress),
						EarlyEnrollment: d.EarlyEnrollment,
						Position:        d.Position,
					})
				}
			}
		}

		campaign.Participants = participants
	}
	return &campaign, nil
}
