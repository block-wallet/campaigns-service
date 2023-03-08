package redis

import (
	"context"
	"testing"

	"github.com/Bose/minisentinel"
	"github.com/alicebob/miniredis/v2"
	"github.com/block-wallet/golang-service-template/utils/monitoring"

	"github.com/block-wallet/golang-service-template/storage/database/config"
	dberror "github.com/block-wallet/golang-service-template/storage/errors"
	"github.com/smartystreets/assertions"
	"github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
)

func TestDatabase_RetrieveFoundValues(t *testing.T) {
	redis := NewRedisSentinel(&config.RedisSentinelConfig{})
	convey.Convey("Given the values of redis mget result", t, func() {
		keys := []string{"key1", "key2", "key3"}
		values := []interface{}{"value1", "value2", nil}
		convey.Convey("When we retrieve the values found", func() {
			results, err := redis.retrieveFoundValues(keys, values)
			convey.Convey("Then we have two results value and no error", func() {
				convey.So(err, assertions.ShouldBeNil)
				convey.So(results, assertions.ShouldNotBeNil)
				convey.So(results, assertions.ShouldResemble, []interface{}{"value1", "value2"})
			})
		})
	})
}

func TestDatabase_RetrieveFoundValues_NoValidValues(t *testing.T) {
	redis := NewRedisSentinel(&config.RedisSentinelConfig{})
	convey.Convey("Given the values of redis mget result", t, func() {
		keys := []string{"key1", "key2", "key3"}
		values := []interface{}{nil, nil, nil}
		convey.Convey("When we retrieve the values found", func() {
			results, err := redis.retrieveFoundValues(keys, values)
			convey.Convey("Then we have notfound error", func() {
				convey.So(err, assertions.ShouldNotBeNil)
				convey.So(err, assertions.ShouldHaveSameTypeAs, &dberror.NotFound{})
				convey.So(results, assertions.ShouldBeNil)
			})
		})
	})
}

func TestRedis_Connect(t *testing.T) {
	redis, sentinel := getRedisSentinel(t)
	assert.NotNil(t, redis)

	defer disconnectAndCloseRedis(t, redis, sentinel)

	err := redis.Connect(context.Background())
	assert.Nil(t, err)
}

func TestRedis_Disconnect(t *testing.T) {
	redis, _ := getRedisSentinel(t)
	assert.NotNil(t, redis)

	err := redis.Connect(context.Background())
	assert.Nil(t, err)

	err = redis.Disconnect(context.Background())
	assert.Nil(t, err)
}

func TestRedis_Set(t *testing.T) {
	redis, sentinel := getRedisSentinel(t)
	assert.NotNil(t, redis)

	defer disconnectAndCloseRedis(t, redis, sentinel)

	err := redis.Connect(context.Background())
	assert.Nil(t, err)

	key := "key"
	field := "field"
	value := "value"

	err = redis.Set(context.Background(), key, field, value)

	assert.Nil(t, err)

	valueGet, errGet := redis.Get(context.Background(), key, field)
	assert.Nil(t, errGet)
	assert.EqualValues(t, value, valueGet)
}

func TestRedis_Get(t *testing.T) {
	redis, sentinel := getRedisSentinel(t)
	assert.NotNil(t, redis)

	defer disconnectAndCloseRedis(t, redis, sentinel)

	err := redis.Connect(context.Background())
	assert.Nil(t, err)

	_ = redis.Set(context.Background(), "key1", "field1", "value1")

	cases := []struct {
		name  string
		key   string
		field string
		value interface{}
		err   error
	}{
		{
			"should return not found error when key doesn't exist",
			"key",
			"",
			nil,
			dberror.NewNotFound("key"),
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
			value, err := redis.Get(context.Background(), c.key, c.field)

			// Validation
			if c.value != nil {
				assert.EqualValues(t, c.value, value.(string))
			} else {
				assert.EqualValues(t, c.value, value)

			}
			assert.EqualValues(t, c.err, err)
		})
	}
}

func TestDelete(t *testing.T) {
	// Initialization
	redis, sentinel := getRedisSentinel(t)
	assert.NotNil(t, redis)

	defer disconnectAndCloseRedis(t, redis, sentinel)

	err := redis.Connect(context.Background())
	assert.Nil(t, err)

	key := "key"
	field := "field"
	value := "value"

	// Setting data
	err = redis.Set(context.Background(), key, field, value)
	assert.Nil(t, err)

	// The data must exist
	valueGet, errGet := redis.Get(context.Background(), key, field)
	assert.Nil(t, errGet)
	assert.EqualValues(t, value, valueGet)

	// Deleting data
	err = redis.Delete(context.Background(), key, []string{field})
	assert.Nil(t, err)

	// The data must not exist
	valueGet, errGet = redis.Get(context.Background(), key, field)
	assert.Empty(t, valueGet)
	assert.EqualError(t, errGet, "Entity not found: key")
}

func TestBulkGet(t *testing.T) {
	// Initialization
	redis, sentinel := getRedisSentinel(t)
	assert.NotNil(t, redis)

	defer disconnectAndCloseRedis(t, redis, sentinel)

	err := redis.Connect(context.Background())
	assert.Nil(t, err)

	// Setting data
	key1 := "key_1"
	field1 := "field_1"
	value1 := "value_1"
	err = redis.Set(context.Background(), key1, field1, value1)
	assert.Nil(t, err)

	field2 := "field_2"
	value2 := "value_2"
	err = redis.Set(context.Background(), key1, field2, value2)
	assert.Nil(t, err)

	// Retreiving data
	results, errBulkGet := redis.BulkGet(context.Background(), key1, []string{field1, field2})
	assert.Nil(t, errBulkGet)
	assert.EqualValues(t, len(results), 2)
	assert.EqualValues(t, results[0], value1)
	assert.EqualValues(t, results[1], value2)
}

func getRedisSentinel(t *testing.T) (*RedisSentinel, *minisentinel.Sentinel) {
	m, err := miniredis.Run()
	assert.Nil(t, err)
	assert.NotNil(t, m)

	s := minisentinel.NewSentinel(m, minisentinel.WithReplica(m))
	err = s.Start()
	assert.Nil(t, err)

	redisSentinelConfig := config.NewRedisSentinelConfig([]string{s.Addr()},
		s.MasterInfo().Name,
		"",
		0,
		false,
		monitoring.RedisLatencyMetricSender)
	assert.NotNil(t, redisSentinelConfig)

	redis := NewRedisSentinel(redisSentinelConfig)
	assert.NotNil(t, redis)

	return redis, s
}

func disconnectAndCloseRedis(t *testing.T, r *RedisSentinel, s *minisentinel.Sentinel) {
	err := r.Disconnect(context.Background())
	assert.Nil(t, err)
	s.Master().Close()
	s.Close()
}
