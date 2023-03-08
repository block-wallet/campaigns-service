package redis

import (
	"context"
	"errors"
	"strings"

	"github.com/block-wallet/golang-service-template/storage/database/config"
	dberror "github.com/block-wallet/golang-service-template/storage/errors"
	hookredis "github.com/block-wallet/golang-service-template/utils/interceptors/hook/redis"
	"github.com/block-wallet/golang-service-template/utils/logger"
	"github.com/block-wallet/golang-service-template/utils/monitoring/histogram"
	"github.com/go-redis/redis/v8"
)

type RedisSentinel struct {
	options             *redis.FailoverOptions
	client              *redis.Client
	latencyMetricSender histogram.LatencyMetricSender
}

func NewRedisSentinel(sentinelConfig *config.RedisSentinelConfig) *RedisSentinel {
	return &RedisSentinel{
		options: &redis.FailoverOptions{
			MasterName:    sentinelConfig.MasterName,
			SentinelAddrs: sentinelConfig.Hosts,
			Password:      sentinelConfig.Password,
			DB:            sentinelConfig.DB,
			SlaveOnly:     sentinelConfig.ReadOnly,
		},
		latencyMetricSender: sentinelConfig.LatencyMetricSender,
	}
}

func (r *RedisSentinel) Connect(ctx context.Context) error {
	r.client = redis.NewFailoverClient(r.options)
	err := r.client.Ping(ctx).Err()
	if err != nil {
		logger.Sugar.WithCtx(ctx).Errorf("Could not connect to redis db using sentinels '%s': %s", r.options.SentinelAddrs, err)
		return err
	}
	logger.Sugar.WithCtx(ctx).Debugf("Connection established with redis db '%s'", r.client.Options().Addr)
	r.client.AddHook(hookredis.NewMetricsHook(r.latencyMetricSender))
	return err
}

func (r *RedisSentinel) Disconnect(ctx context.Context) error {
	err := r.client.Close()
	if err != nil {
		logger.Sugar.WithCtx(ctx).Errorf("disconnected from redis db: %s", r.client.Options().Addr)
	}
	return err
}

func (r *RedisSentinel) Get(ctx context.Context, key, field string) (interface{}, error) {
	value, err := r.client.HGet(ctx, key, field).Result()
	if errors.Is(err, redis.Nil) {
		return nil, dberror.NewNotFound(key)
	}
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (r *RedisSentinel) GetAll(ctx context.Context, key string) (interface{}, error) {
	value, err := r.client.HGetAll(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, dberror.NewNotFound(key)
	}
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (r *RedisSentinel) Set(ctx context.Context, key, field string, value interface{}) error {
	return r.client.HSet(ctx, key, field, value).Err()
}

func (r *RedisSentinel) Delete(ctx context.Context, key string, fields []string) error {
	deleted, err := r.client.HDel(ctx, key, fields...).Result()
	if err != nil {
		return err
	}
	if deleted == 0 {
		return dberror.NewNotFound(key)
	}

	return nil
}

func (r *RedisSentinel) BulkGet(ctx context.Context, key string, fields []string) ([]interface{}, error) {
	values, err := r.client.HMGet(ctx, key, fields...).Result()
	if err != nil {
		return nil, err
	}
	return r.retrieveFoundValues(fields, values)
}

func (r *RedisSentinel) retrieveFoundValues(keys []string, values []interface{}) ([]interface{}, error) {
	results := make([]interface{}, 0)
	for _, val := range values {
		if val != nil {
			results = append(results, val)
		}
	}
	if len(results) == 0 {
		return nil, dberror.NewNotFound(strings.Join(keys, " - "))
	}
	return results, nil
}
