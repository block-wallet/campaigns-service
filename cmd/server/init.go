package server

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	kvdb "github.com/block-wallet/golang-service-template/storage/database/kv"
	"github.com/block-wallet/golang-service-template/utils/val"

	"github.com/spf13/cobra"
)

const (
	defaultLogLevel                  = "debug"
	defaultPort                      = 8080
	defaultMetricsPort               = 9008
	defaultKVType                    = kvdb.LocalLRU
	defaultRedisSentinelMasterName   = "mymaster"
	defaultRedisSentinelHosts        = "localhost:26379"
	defaultRedisReadOnly             = "false"
	defaultRedisPassword             = ""
	defaultRedisDB                   = 0
	defaultLocalCacheExpiration      = 0 * time.Second
	defaultLocalCacheCleanUpInterval = 0 * time.Second
	defaultSomeHTTPEndpoint          = "chainid.network"
	defaultSomeHTTPProtocol          = "https"
	defaultSomeHTTPTimeout           = 60 * time.Second
	defaultETHEndpoint               = "wss://goerli-node.blockwallet.io/ws"
)

// Args for this cmd
var (
	logLevelArg               string
	port                      int
	metricsPort               int
	kvType                    string
	redisSentinelHosts        []string
	redisSentinelMasterName   string
	redisReadOnly             bool
	redisPassword             string
	redisDB                   int
	localCacheExpiration      time.Duration
	localCacheCleanUpInterval time.Duration
	someHTTPEndpoint          string
	someHTTPProtocol          string
	someHTTPTimeout           time.Duration
	ETHEndpoint               string
)

var cmd = &cobra.Command{
	Use:   "serve",
	Short: "Runs the gRPC Server",
	Long:  `Runs the gRPC Server`,
}

func init() {
	logLevelArg = val.GetEnvValWithDefault("LOG_LEVEL", defaultLogLevel)

	port, _ = strconv.Atoi(val.GetEnvValWithDefault("PORT", fmt.Sprint(defaultPort)))
	metricsPort, _ = strconv.Atoi(val.GetEnvValWithDefault("METRICS_PORT", fmt.Sprint(defaultMetricsPort)))

	kvType = val.GetEnvValWithDefault("KV_TYPE", defaultKVType)
	redisSentinelHosts = strings.Split(val.GetEnvValWithDefault("REDIS_SENTINEL_HOSTS", defaultRedisSentinelHosts), ",")
	redisSentinelMasterName = val.GetEnvValWithDefault("REDIS_SENTINEL_MASTER_NAME", defaultRedisSentinelMasterName)
	redisReadOnly, _ = strconv.ParseBool(val.GetEnvValWithDefault("REDIS_READ_ONLY", defaultRedisReadOnly))
	redisPassword = val.GetEnvValWithDefault("REDIS_PASSWORD", defaultRedisPassword)
	redisDB, _ = strconv.Atoi(val.GetEnvValWithDefault("REDIS_DB", fmt.Sprint(defaultRedisDB)))

	someHTTPEndpoint = val.GetEnvValWithDefault("SOME_HTTP_ENDPOINT", defaultSomeHTTPEndpoint)
	someHTTPProtocol = val.GetEnvValWithDefault("SOME_HTTP_PROTOCOL", defaultSomeHTTPProtocol)
	someHTTPTimeout, _ = time.ParseDuration(val.GetEnvValWithDefault("SOME_HTTP_TIMEOUT", fmt.Sprint(defaultSomeHTTPTimeout)))

	localCacheExpiration, _ = time.ParseDuration(val.GetEnvValWithDefault("LOCAL_CACHE_EXPIRATION", fmt.Sprint(defaultLocalCacheExpiration)))
	localCacheCleanUpInterval, _ = time.ParseDuration(val.GetEnvValWithDefault("LOCAL_CACHE_CLEAN_UP_INTERVAL", fmt.Sprint(defaultLocalCacheCleanUpInterval)))

	ETHEndpoint = val.GetEnvValWithDefault("ETH_ENDPOINT", defaultETHEndpoint)
}
