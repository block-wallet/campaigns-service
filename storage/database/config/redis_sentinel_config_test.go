package config

import (
	"testing"

	"github.com/block-wallet/golang-service-template/utils/monitoring"
	histogrammonitoring "github.com/block-wallet/golang-service-template/utils/monitoring/histogram"
	"github.com/stretchr/testify/assert"
)

func TestNewSentinelRedisConfig(t *testing.T) {
	cases := []struct {
		name                string
		hosts               []string
		masterName          string
		password            string
		DB                  int
		readOnly            bool
		LatencyMetricSender histogrammonitoring.LatencyMetricSender
	}{
		{
			"one host",
			[]string{"host:1234"},
			"",
			"",
			0,
			true,
			nil,
		},
		{"multiple hosts",
			[]string{"host:1234", "host:1245"},
			"localhost:1234",
			"password",
			0,
			false,
			monitoring.RedisLatencyMetricSender,
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			redisConfig := NewRedisSentinelConfig(c.hosts, c.masterName, c.password, c.DB, c.readOnly, c.LatencyMetricSender)

			assert.NotNil(t, redisConfig)
			assert.Equal(t, redisConfig.Hosts, c.hosts)
			assert.Equal(t, redisConfig.Password, c.password)
			assert.Equal(t, redisConfig.DB, c.DB)
			assert.Equal(t, redisConfig.LatencyMetricSender, c.LatencyMetricSender)
		})
	}
}
