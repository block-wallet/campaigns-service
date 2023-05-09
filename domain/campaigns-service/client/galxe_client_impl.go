package client

import (
	"context"
	"time"

	"github.com/block-wallet/campaigns-service/utils/logger"
	"github.com/hasura/go-graphql-client"
)

const GALXE_ACCES_TOKEN_HEADER_KEY = "access-token"

type GalxeClientImpl struct {
	client     *graphql.Client
	maxRetries int8
}

func NewGalxeClient(graphQLEndpoint, accessToken string) *GalxeClientImpl {
	return &GalxeClientImpl{
		client:     graphql.NewClient(graphQLEndpoint, NewAuthenticatedClient(accessToken, GALXE_ACCES_TOKEN_HEADER_KEY)),
		maxRetries: 3,
	}
}

func (gc *GalxeClientImpl) PopulateParticipant(ctx context.Context, input PopulateParticipantsInput) (bool, error) {
	addressesStr := []string{
		input.Address.String(),
	}

	type MutateCredItemInput struct {
		CredId    string   `json:"credId"`
		Operation string   `json:"operation"`
		Items     []string `json:"items"`
	}

	var mutation struct {
		CredentialItems struct {
			Name string
		} `graphql:"credentialItems(input:$input)"`
	}

	variables := map[string]interface{}{
		"input": MutateCredItemInput{
			CredId:    input.CredentialId,
			Operation: "APPEND",
			Items:     addressesStr,
		},
	}

	err := gc.execWithRetries(ctx, func() error {
		return gc.client.Mutate(ctx, &mutation, variables, graphql.OperationName("credentialItems"))
	}, 0)

	if err != nil {
		logger.Sugar.WithCtx(ctx).Errorf("error populing participant = %v to galxe. Error: %v", input.Address, err.Error())
		return false, err
	}

	return true, nil
}

func (gc *GalxeClientImpl) execWithRetries(ctx context.Context, execute func() error, retryCount int8) error {
	err := execute()
	if err != nil && retryCount < gc.maxRetries {
		nextRetry := retryCount + 1
		logger.Sugar.WithCtx(ctx).Warnf("error populating account to galxe. Error: %v", err.Error())
		logger.Sugar.WithCtx(ctx).Warnf("retrying #%d starts in 1 second...", retryCount)
		time.Sleep(1 * time.Second)
		return gc.execWithRetries(ctx, execute, nextRetry)
	}
	return err
}
