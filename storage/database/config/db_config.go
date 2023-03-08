package config

type DBConfig struct {
	DBType              string
	LocalCacheConfig    *LocalCacheConfig
	RedisSentinelConfig *RedisSentinelConfig
}

func NewDBConfig(DBType string, localCacheConfig *LocalCacheConfig, redisSentinelConfig *RedisSentinelConfig) *DBConfig {
	return &DBConfig{
		DBType:              DBType,
		LocalCacheConfig:    localCacheConfig,
		RedisSentinelConfig: redisSentinelConfig}
}
