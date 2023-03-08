package ethrepository

import (
	"context"
	"encoding/json"
	"fmt"

	httpapi "github.com/block-wallet/golang-service-template/domain/http-api"
	"github.com/block-wallet/golang-service-template/domain/model"

	"github.com/block-wallet/golang-service-template/utils/logger"

	kvdb "github.com/block-wallet/golang-service-template/storage/database/kv"
	"github.com/block-wallet/golang-service-template/utils/errors"
)

type KVRepository struct {
	kvDatabase    kvdb.Database
	httpApiClient httpapi.ApiClient
}

func NewKVRepository(kvDatabase kvdb.Database, httpApiClient httpapi.ApiClient) *KVRepository {
	return &KVRepository{
		kvDatabase:    kvDatabase,
		httpApiClient: httpApiClient,
	}
}

func (k *KVRepository) GetEvents(ctx context.Context, pair string) (*[]model.Event, errors.RichError) {
	key := k.GetEventsKey(pair)
	eventsI, _ := k.kvDatabase.GetAll(ctx, key)

	var events []model.Event
	if eventsI != nil {
		for _, eventString := range eventsI.(map[string]string) {
			event := &model.Event{}
			if err := json.Unmarshal([]byte(eventString), event); err != nil {
				logger.Sugar.WithCtx(ctx).Errorf("Error getting events from kvDatabase: failed "+
					"unmarshalling - %s", err.Error())
				return nil, errors.NewInternal("Error getting events from kvDatabase: failed unmarshalling")
			}
			events = append(events, *event)
		}
	}

	return &events, nil
}
func (k *KVRepository) GetChains(ctx context.Context) (*[]model.Chain, errors.RichError) {
	chains, err := k.httpApiClient.GetChains(ctx)
	if err != nil {
		logger.Sugar.WithCtx(ctx).Errorf("Error getting chains from api - %s", err.Error())
		return nil, errors.NewInternal("Error getting chains from api")

	}
	return chains, nil
}

func (k *KVRepository) SetEvent(ctx context.Context, pair, eventId string, event *model.Event) errors.RichError {
	key := k.GetEventsKey(pair)
	eventString, err := dataToString(event)
	if err != nil {
		logger.Sugar.WithCtx(ctx).Errorf("Error parsing event - %s", err.Error())
		return errors.NewInternal("Error parsing event")
	}

	err = k.kvDatabase.Set(ctx, key, eventId, eventString)
	if err != nil {
		logger.Sugar.WithCtx(ctx).Errorf("Error saving event in kvDatabase: key %s %s - %s", key, eventId, err.Error())
		return errors.NewInternal("Error saving event in kvDatabase")
	}

	return nil
}

func (k *KVRepository) GetEventsKey(pair string) string {
	return fmt.Sprintf("events_%s", pair)
}

func dataToString(data interface{}) (string, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("could not marshal request data %v", err)
	}
	return string(b), nil
}
