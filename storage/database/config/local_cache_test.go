package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewLocalCacheConfig(t *testing.T) {
	cases := []struct {
		name              string
		cleanUpInterval   time.Duration
		defaultExpiration time.Duration
	}{
		{
			"defaultExpiration and cleanupInterval 0",
			0,
			0,
		},
		{
			"defaultExpiration 0 and cleanupInterval with value",
			0,
			200,
		},
		{
			"defaultExpiration with value and cleanupInterval 0",
			300,
			0,
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			// Operation
			localDbConfig := NewLocalCacheConfig(c.defaultExpiration, c.cleanUpInterval)

			// Validation
			assert.NotNil(t, localDbConfig)
			assert.Equal(t, localDbConfig.DefaultExpiration, c.defaultExpiration)
			assert.Equal(t, localDbConfig.CleanupInterval, c.cleanUpInterval)
		})
	}
}
