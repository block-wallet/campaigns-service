package campaignsrepository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/block-wallet/campaigns-service/domain/model"
	"github.com/block-wallet/campaigns-service/utils/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
)

type SQLRepository struct {
	db *sql.DB
}

func NewSQLRepository(sqlDatabase *sql.DB) Repository {
	return &SQLRepository{
		db: sqlDatabase,
	}
}

func (r *SQLRepository) GetCampaigns(ctx context.Context, filters *model.GetCampaignsFilters) ([]*model.Campaign, error) {
	campaignsQueryBuilder := NewCampaignsQueryBuilder(filters)
	q, p := campaignsQueryBuilder.Query(ctx)
	rows, err := r.db.Query(q, p...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	campaigns := make([]*model.Campaign, 0)

	for rows.Next() {
		modelCampaign, err := campaignsQueryBuilder.Parse(ctx, rows)
		if err != nil {
			continue
		}
		campaigns = append(campaigns, modelCampaign)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return campaigns, nil
}

func (r *SQLRepository) GetCampaignById(ctx context.Context, id string) (*model.Campaign, error) {
	filters := &model.GetCampaignsFilters{
		Id: &id,
	}
	campaignsQueryBuilder := NewCampaignsQueryBuilder(filters)
	q, p := campaignsQueryBuilder.Query(ctx)
	rows, err := r.db.Query(q, p...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var modelCampaign *model.Campaign
	for rows.Next() {
		if modelCampaign != nil {
			return nil, fmt.Errorf("more than 1 campaign record found for id: %s", id)
		}
		modelCampaign, err = campaignsQueryBuilder.Parse(ctx, rows)
		if err != nil {
			continue
		}
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return modelCampaign, nil
}

func (r *SQLRepository) NewCampaign(ctx context.Context, input *model.CreateCampaignInput) (*string, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	campaignId := uuid.NewString()

	_, err = tx.ExecContext(ctx, "INSERT INTO campaigns (id,name,description,status,start_date,end_date,enroll_message) VALUES ($1,$2,$3,$4,$5,$6,$7)", campaignId, input.Name, input.Description, input.Status, input.StartDate, input.EndDate, input.EnrollMessage)
	if err != nil {
		return nil, err
	}
	for _, chainId := range input.SupportedChains {
		if _, err = tx.ExecContext(ctx, "INSERT INTO campaigns_supported_chains (campaign_id,chain_id) VALUES ($1,$2)", campaignId, chainId); err != nil {
			return nil, err
		}
	}
	for _, tag := range input.Tags {
		if _, err = tx.ExecContext(ctx, "INSERT INTO campaigns_tags (campaign_id,tag) VALUES ($1,$2)", campaignId, tag); err != nil {
			return nil, err
		}
	}
	rewardId := uuid.NewString()
	amounts := strings.Join(input.Rewards.Amounts, ",")
	_, err = tx.ExecContext(ctx, "INSERT INTO rewards (reward_id,campaign_id,token_id,type,amounts) VALUES ($1,$2,$3,$4,$5)", rewardId, campaignId, *input.Rewards.Token.Id, input.Rewards.Type, amounts)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &campaignId, nil
}

func (r *SQLRepository) EnrollInCampaign(ctx context.Context, input *model.EnrollInCampaignInput) (*bool, error) {
	ok := true
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO participants (campaign_id,account_address) VALUES ($1,$2)", input.CampaignId, input.Adddress.String())
	if err != nil {
		if _err := tx.Rollback(); _err != nil {
			logger.Sugar.WithCtx(ctx).Warnf("error applying transaction rollback: %v", _err.Error())
			return nil, _err
		}
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &ok, nil
}

func (r *SQLRepository) UpdateCampaign(ctx context.Context, updates *model.UpdateCampaignInput) (*bool, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	ok := true
	if err != nil {
		return nil, err
	}
	params := make([]any, 0)
	updatesVariables := make([]string, 0)
	if updates.Stauts != nil {
		params = append(params, string(*updates.Stauts))
		updatesVariables = append(updatesVariables, "status = $1")
	}

	if len(updatesVariables) > 0 {
		joinedVariables := strings.Join(updatesVariables, ", ")
		params = append(params, updates.Id)
		q := fmt.Sprintf("UPDATE campaigns SET %v WHERE id = $%v;", joinedVariables, len(params))
		_, err := tx.ExecContext(ctx, q, params...)
		if err != nil {
			if _err := tx.Rollback(); _err != nil {
				logger.Sugar.WithCtx(ctx).Warnf("error applying transaction rollback: %v", _err.Error())
				return nil, _err
			}
			return nil, err
		}
	}

	if updates.Winners != nil && len(*updates.Winners) > 0 {
		_, err := tx.ExecContext(ctx, "UPDATE participants SET position = NULL where campaign_id = $1;", updates.Id)
		if err != nil {
			if _err := tx.Rollback(); _err != nil {
				logger.Sugar.WithCtx(ctx).Warnf("error applying transaction rollback: %v", _err.Error())
				return nil, _err
			}
			return nil, err
		}
		winners := *updates.Winners
		for i := 0; i < len(winners); i++ {
			winnerAddress := winners[i].String()
			_, err := tx.ExecContext(ctx, "UPDATE participants SET position = $1 WHERE campaign_id = $2 and account_address = $3;", i+1, updates.Id, winnerAddress)
			if err != nil {
				if _err := tx.Rollback(); _err != nil {
					logger.Sugar.WithCtx(ctx).Warnf("error applying transaction rollback: %v", _err.Error())
					return nil, _err
				}
				return nil, err
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &ok, nil
}

func (r *SQLRepository) ParticipantExists(ctx context.Context, campaignId string, accountAddress string) (*bool, error) {
	statement, err := r.db.Prepare("SELECT campaign_id from participants p where p.campaign_id = $1 and p.account_address = $2;")
	if err != nil {
		return nil, err
	}
	var id string
	err = statement.QueryRow(campaignId, accountAddress).Scan(&id)
	ret := true
	if err != nil {
		ret = false
		//An error occurred.
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	if id != campaignId {
		ret = false
	}

	return &ret, nil
}

func (r *SQLRepository) GetTokenById(ctx context.Context, id string) (*model.MultichainToken, error) {
	statement, err := r.db.Prepare("SELECT t.id, t.name, t.symbol, t.decimals from tokens t where t.id = $1;")
	if err != nil {
		return nil, err
	}
	row := statement.QueryRow(id)

	if err = row.Err(); err != nil {
		return nil, err
	}

	tokenRow := tokenrow{}
	err = row.Scan(&tokenRow.id, &tokenRow.name, &tokenRow.description, &tokenRow.decimals)
	if err != nil {
		return nil, err
	}

	contracts, err := r.getTokenContracts(ctx, id)
	if err != nil {
		return nil, err
	}

	return &model.MultichainToken{
		Name:              tokenRow.name,
		Symbol:            tokenRow.symbol,
		Decimals:          uint8(tokenRow.decimals),
		ContractAddresses: *contracts,
	}, nil
}

func (r *SQLRepository) getTokenContracts(ctx context.Context, id string) (*map[string]common.Address, error) {
	contracts := map[string]common.Address{}
	contractsRows, err := r.db.Query("SELECT tc.chain_id, tc.address FROM tokens_contracts tc WHERE tc.token_id = $1;", id)

	if err != nil {
		return nil, err
	}

	defer contractsRows.Close()

	for contractsRows.Next() {
		var chainid, address string
		err := contractsRows.Scan(&chainid, &address)
		if err != nil {
			return nil, err
		}
		contracts[chainid] = common.HexToAddress(address)
	}
	return &contracts, nil
}

func (r *SQLRepository) TokenExists(ctx context.Context, id string) (*bool, error) {
	statement, err := r.db.Prepare("SELECT t.id from tokens t where t.id = $1;")
	if err != nil {
		return nil, err
	}
	err = statement.QueryRow(id).Scan(&id)
	ret := true
	if err != nil {
		ret = false
		//An error occurred.
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	return &ret, nil
}

func (r *SQLRepository) NewToken(ctx context.Context, token *model.MultichainToken) (*string, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	tokenId := uuid.NewString()

	_, err = tx.ExecContext(ctx, "INSERT INTO tokens (id,name,symbol,decimals) VALUES ($1,$2,$3,$4)", tokenId, token.Name, token.Symbol, token.Decimals)

	if err != nil {
		if _err := tx.Rollback(); _err != nil {
			logger.Sugar.WithCtx(ctx).Warnf("error applying transaction rollback: %v", _err.Error())
			return nil, _err
		}
		return nil, err
	}

	for chain, addrr := range token.ContractAddresses {
		_, err := tx.ExecContext(ctx, "INSERT INTO tokens_contracts (token_id,chain_id,address) VALUES ($1,$2,$3)", tokenId, chain, addrr.String())
		if err != nil {
			if _err := tx.Rollback(); _err != nil {
				logger.Sugar.WithCtx(ctx).Warnf("error applying transaction rollback: %v", _err.Error())
				return nil, _err
			}
			return nil, err
		}
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return &tokenId, nil
}

func (r *SQLRepository) GetAllTokens(ctx context.Context) ([]*model.MultichainToken, error) {
	rows, err := r.db.Query("SELECT t.id, t.name, t.symbol, t.decimals from tokens t;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tokens := make([]*model.MultichainToken, 0)

	for rows.Next() {
		tokenRow := tokenrow{}
		err := rows.Scan(&tokenRow.id, &tokenRow.name, &tokenRow.description, &tokenRow.decimals)
		if err != nil {
			continue
		}
		contracts, err := r.getTokenContracts(ctx, tokenRow.id)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, &model.MultichainToken{
			Id:                tokenRow.id,
			Name:              tokenRow.name,
			Symbol:            tokenRow.symbol,
			Decimals:          uint8(tokenRow.decimals),
			ContractAddresses: *contracts,
		})
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return tokens, nil
}
