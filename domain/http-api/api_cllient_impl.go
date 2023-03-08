package httpapi

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/block-wallet/golang-service-template/domain/model"
	"github.com/block-wallet/golang-service-template/utils/http"

	"github.com/block-wallet/golang-service-template/utils/logger"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ApiClient interface {
	GetChains(ctx context.Context) (*[]model.Chain, error)
}

type ApiClientImpl struct {
	client   http.Client
	endpoint string
}

func NewApiClientImpl(client http.Client, protocol string, host string) *ApiClientImpl {
	return &ApiClientImpl{
		client:   client,
		endpoint: fmt.Sprintf("%s://%s", protocol, host),
	}
}

func (c *ApiClientImpl) GetChains(ctx context.Context) (*[]model.Chain, error) {
	result, err := c.client.Get(ctx, fmt.Sprintf("%s/chains.json", c.endpoint), nil)
	if err != nil {
		logger.Sugar.WithCtx(ctx).Errorf("request to chains api error - %s", err.Error())
		return nil, err
	}

	parsed, err := getResponseFromBody(result)
	if err != nil {
		logger.Sugar.WithCtx(ctx).Errorf("parsing chains api response error - %s", err.Error())
		return nil, err
	}

	return parsed, nil
}

func getResponseFromBody(body []byte) (*[]model.Chain, error) {
	resp := &[]model.Chain{}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, status.Error(codes.DataLoss, fmt.Sprintf("error unmarshaling response body: %v", err))
	}
	return resp, nil
}
