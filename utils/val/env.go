package val

import (
	"os"
)

func GetEnvValWithDefault(key string, defaultValue string) string {
	value, found := os.LookupEnv(key)
	if !found {
		return defaultValue
	}

	return value
}
