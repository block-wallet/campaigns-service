package config

import histogrammonitoring "github.com/block-wallet/golang-service-template/utils/monitoring/histogram"

type RedisSentinelConfig struct {
	Hosts               []string
	MasterName          string
	Password            string
	DB                  int
	ReadOnly            bool
	LatencyMetricSender histogrammonitoring.LatencyMetricSender
}

func NewRedisSentinelConfig(hosts []string, masterName,
	password string, DB int, readOnly bool,
	latencyMetricSender histogrammonitoring.LatencyMetricSender) *RedisSentinelConfig {
	return &RedisSentinelConfig{
		Hosts:               hosts,
		MasterName:          masterName,
		Password:            password,
		DB:                  DB,
		ReadOnly:            readOnly,
		LatencyMetricSender: latencyMetricSender,
	}
}
