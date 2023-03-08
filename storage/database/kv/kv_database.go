package kv

import (
	"context"

	"github.com/block-wallet/golang-service-template/storage/database/kv/redis"

	"github.com/block-wallet/golang-service-template/storage/database/config"
	"github.com/block-wallet/golang-service-template/storage/database/kv/local"
)

const RedisDb = "redis"
const LocalLRU = "local"

type Database interface {
	Connect(ctx context.Context) error
	Get(ctx context.Context, key, field string) (interface{}, error)
	GetAll(ctx context.Context, key string) (interface{}, error)
	Set(ctx context.Context, key, field string, value interface{}) error
	Delete(ctx context.Context, key string, fields []string) error
	BulkGet(ctx context.Context, key string, fields []string) ([]interface{}, error)
}

func NewKVDatabase(dbConfig *config.DBConfig) Database {
	switch dbConfig.DBType {
	case RedisDb:
		return redis.NewRedisSentinel(dbConfig.RedisSentinelConfig)
	case LocalLRU:
		return local.NewLocalCache(dbConfig.LocalCacheConfig)
	default:
		return local.NewLocalCache(dbConfig.LocalCacheConfig)
	}
}
