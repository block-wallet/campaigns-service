package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDBConfig(t *testing.T) {
	cases := []struct {
		name             string
		dbType           string
		localCacheConfig *LocalCacheConfig
		redisConfig      *RedisSentinelConfig
	}{
		{
			"all empty",
			"",
			nil,
			nil,
		},
		{
			"db type local",
			"local",
			NewLocalCacheConfig(0, 0),
			nil,
		},
		{
			"db type redis sentinel",
			"redis",
			nil,
			NewRedisSentinelConfig([]string{"host"}, "master", "password", 0, true, nil),
		},
		{
			"db type redis sentinel with local",
			"redis",
			NewLocalCacheConfig(0, 0),
			NewRedisSentinelConfig([]string{"host"}, "master", "password", 0, true, nil),
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			// Operation
			dbConfig := NewDBConfig(c.dbType, c.localCacheConfig, c.redisConfig)

			// Validation
			assert.NotNil(t, dbConfig)
			assert.Equal(t, dbConfig.DBType, c.dbType)
			assert.Equal(t, dbConfig.LocalCacheConfig, c.localCacheConfig)
			assert.Equal(t, dbConfig.RedisSentinelConfig, c.redisConfig)
		})
	}
}
