package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDBConfig(t *testing.T) {
	cases := []struct {
		name         string
		dbType       DBType
		sqliteConfig *SQLConfig
	}{
		{
			"all empty",
			"",
			nil,
		},
		{
			"db type SQLite",
			SQLiteDBType,
			NewSQLConfig("test.db", false),
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			// Operation
			dbConfig := NewDBConfig(c.dbType, c.sqliteConfig)

			// Validation
			assert.NotNil(t, dbConfig)
			assert.Equal(t, dbConfig.DBType, c.dbType)
			assert.Equal(t, dbConfig.SQLConfig, c.sqliteConfig)
		})
	}
}
