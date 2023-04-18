package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSQLConfig(t *testing.T) {
	cases := []struct {
		name     string
		fileName string
	}{
		{
			"Production database named Campaigns",
			"campaigns.db",
		},
		{
			"Test database named Test",
			"test.db",
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			// Operation
			sqliteDbConfig := NewSQLConfig(c.fileName, false)

			// Validation
			assert.NotNil(t, sqliteDbConfig)
			assert.Equal(t, sqliteDbConfig.Connection, c.fileName)
		})
	}
}
