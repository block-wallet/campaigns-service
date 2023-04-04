package val

import (
	"log"
	"os"
	"strconv"
)

func GetEnvValWithDefault(key string, defaultValue string) string {
	value, found := os.LookupEnv(key)
	if !found {
		return defaultValue
	}

	return value
}

func GetBoolEnvValWithDefault(key string, defaultValue bool) bool {
	value, found := os.LookupEnv(key)
	if !found {
		return defaultValue
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		log.Fatal(err)
	}

	return boolValue
}
