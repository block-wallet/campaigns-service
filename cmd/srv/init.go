package srv

import (
	"fmt"
	"strconv"

	"github.com/block-wallet/campaigns-service/storage/database/config"
	"github.com/block-wallet/campaigns-service/utils/logger"
	"github.com/block-wallet/campaigns-service/utils/val"

	"github.com/spf13/cobra"
)

const (
	defaultLogLevel             = "debug"
	defaultPort                 = 8080
	defaultMetricsPort          = 9008
	defaultDbType               = config.PostgresDBType
	defaultSQLConnectionString  = "postgresql://localhost:5432/postgres?user=postgres&password=admin&sslmode=disable"
	defaultAdminUsername        = "blockwallet"
	defaultAdminPassword        = "password123"
	defaultSkipMigrations       = false
	defaultGalxeGraphQLEndpoint = "https://graphigo.prd.galaxy.eco/query"
	defaultGalxeAccessToken     = ""
)

// Args for this cmd
var (
	logLevelArg          string
	port                 int
	metricsPort          int
	dbType               config.DBType
	sqlConnectionString  string
	adminUsername        string
	adminPassword        string
	skipMigrations       bool
	galxeGraphQLEndpoint string
	galxeAccessToken     string
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

	dbType = parseDBType(val.GetEnvValWithDefault("DB_TYPE", string(defaultDbType)))
	sqlConnectionString = val.GetEnvValWithDefault("SQL_CONNECTION", defaultSQLConnectionString)
	adminUsername = val.GetEnvValWithDefault("ADMIN_USERNAME", defaultAdminUsername)
	adminPassword = val.GetEnvValWithDefault("ADMIN_PASSWORD", defaultAdminPassword)
	skipMigrations = val.GetBoolEnvValWithDefault("SKIP_MIGRATIONS", defaultSkipMigrations)
	galxeGraphQLEndpoint = val.GetEnvValWithDefault("GALXE_GRAPHQL_ENDPOINT", defaultGalxeGraphQLEndpoint)
	galxeAccessToken = val.GetEnvValWithDefault("GALXE_ACCESS_TOKEN", defaultGalxeAccessToken)
}

func parseDBType(dbType string) config.DBType {
	switch dbType {
	case string(config.PostgresDBType):
		return config.PostgresDBType
	}
	logger.Sugar.Warnf("Unknown specified db type: %s. Using default: %s", dbType, config.PostgresDBType)
	return config.PostgresDBType
}
