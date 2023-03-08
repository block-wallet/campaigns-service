package local

import (
	"context"
	"testing"

	"github.com/block-wallet/golang-service-template/storage/database/config"
	"github.com/block-wallet/golang-service-template/storage/errors"

	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	localCache := NewLocalCache(config.NewLocalCacheConfig(-1, -1))

	err := localCache.Connect(context.Background())
	assert.Nil(t, err)
}

func TestDisconnect(t *testing.T) {
	localCache := NewLocalCache(config.NewLocalCacheConfig(-1, -1))

	err := localCache.Disconnect(context.Background())
	assert.Nil(t, err)
}

func TestGet(t *testing.T) {
	localCache := NewLocalCache(config.NewLocalCacheConfig(-1, -1))
	_ = localCache.Set(context.Background(), "key1", "field1", "value1")

	cases := []struct {
		name  string
		key   string
		field string
		value string
		err   error
	}{
		{
			"should return not found error when key doesn't exist",
			"key",
			"field",
			"",
			errors.NewNotFound("key"),
		},
		{
			"should return the value when key exists on local cache",
			"key1",
			"field1",
			"value1",
			nil,
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			// Operation
			value, err := localCache.Get(context.Background(), c.key, c.field)

			// Validation
			assert.EqualValues(t, c.value, value.(string))
			assert.EqualValues(t, c.err, err)
		})
	}
}

func TestSet(t *testing.T) {
	// Initialization
	localCache := NewLocalCache(config.NewLocalCacheConfig(-1, -1))
	key := "key"
	field := "field"
	value := "value"

	// Operation
	err := localCache.Set(context.Background(), key, field, value)

	// Validation
	assert.Nil(t, err)

	valueGet, errGet := localCache.Get(context.Background(), key, field)
	assert.Nil(t, errGet)
	assert.EqualValues(t, value, valueGet)
}
