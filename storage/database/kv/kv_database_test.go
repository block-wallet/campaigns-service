package kv

import (
	"testing"

	"github.com/block-wallet/golang-service-template/storage/database/kv/redis"

	"github.com/block-wallet/golang-service-template/storage/database/config"
	"github.com/block-wallet/golang-service-template/storage/database/kv/local"
	"github.com/stretchr/testify/assert"
)

func TestNewKVDatabaseLocal(t *testing.T) {
	localDbConfig := config.NewLocalCacheConfig(0, 0)
	assert.NotNil(t, localDbConfig)

	dbConfig := config.NewDBConfig("local", localDbConfig, nil)
	assert.NotNil(t, dbConfig)

	database := NewKVDatabase(dbConfig)
	assert.NotNil(t, database)
	assert.NotNil(t, database.(*local.Cache))
}

func TestNewKVDatabaseRedisSentinel(t *testing.T) {
	redisConfig := config.NewRedisSentinelConfig([]string{""}, "", "", 0, true, nil)
	assert.NotNil(t, redisConfig)

	dbConfig := config.NewDBConfig(RedisDb, nil, redisConfig)
	assert.NotNil(t, dbConfig)

	database := NewKVDatabase(dbConfig)
	assert.NotNil(t, database)
	assert.NotNil(t, database.(*redis.RedisSentinel))
}

func TestNewKVDatabaseEmpty(t *testing.T) {
	localDbConfig := config.NewLocalCacheConfig(0, 0)
	assert.NotNil(t, localDbConfig)

	dbConfig := config.NewDBConfig("", localDbConfig, nil)
	assert.NotNil(t, dbConfig)

	database := NewKVDatabase(dbConfig)
	assert.NotNil(t, database)
	assert.NotNil(t, database.(*local.Cache))
}
