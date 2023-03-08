package local

import (
	"context"
	"strings"

	"github.com/block-wallet/golang-service-template/storage/database/config"
	dberror "github.com/block-wallet/golang-service-template/storage/errors"

	"github.com/patrickmn/go-cache"
)

type Cache struct {
	cache *cache.Cache
}

// NewLocalCache creates a new in local database.
// If the expiration duration is less than one (or NoExpiration),
// the items in the cache never expire (by default), and must be deleted
// manually. If the cleanup interval is less than one, expired items are not
// deleted from the cache before calling c.DeleteExpired().
func NewLocalCache(localCacheConfig *config.LocalCacheConfig) *Cache {
	return &Cache{
		cache: cache.New(localCacheConfig.DefaultExpiration, localCacheConfig.CleanupInterval),
	}
}

func (c *Cache) Connect(_ context.Context) error {
	return nil
}

func (c *Cache) Disconnect(_ context.Context) error {
	return nil
}

func (c *Cache) Get(_ context.Context, key, field string) (interface{}, error) {
	value, found := c.cache.Get(key + field)
	if !found {
		return "", dberror.NewNotFound(key)
	}

	return value, nil
}

func (c *Cache) GetAll(_ context.Context, key string) (interface{}, error) {
	values := map[string]string{}

	for _key := range c.cache.Items() {
		if strings.Contains(_key, key) {
			value, found := c.cache.Get(_key)
			if found {
				values[key] = value.(string)
			}
		}
	}

	if len(values) == 0 {
		return nil, dberror.NewNotFound("")
	}

	return values, nil
}

func (c *Cache) Set(_ context.Context, key, field string, value interface{}) error {
	// c.cache is already thread-safe
	c.cache.Set(key+field, value, 0)
	return nil
}

func (c *Cache) Delete(_ context.Context, key string, fields []string) error {
	for i := range fields {
		if _, found := c.cache.Get(key + fields[i]); found {
			c.cache.Delete(key + fields[i])
		}
	}
	return nil
}

func (c *Cache) BulkGet(_ context.Context, key string, fields []string) ([]interface{}, error) {
	results := make([]interface{}, 0)
	for i := range fields {
		if value, found := c.cache.Get(key + fields[i]); found {
			results = append(results, value)
		}
	}
	if len(results) == 0 {
		return nil, dberror.NewNotFound(strings.Join(fields, " - "))
	}
	return results, nil
}
