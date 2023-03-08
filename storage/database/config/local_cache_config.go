package config

import "time"

type LocalCacheConfig struct {
	DefaultExpiration time.Duration
	CleanupInterval   time.Duration
}

func NewLocalCacheConfig(defaultExpiration, cleanupInterval time.Duration) *LocalCacheConfig {
	return &LocalCacheConfig{
		DefaultExpiration: defaultExpiration,
		CleanupInterval:   cleanupInterval,
	}
}
