package store

import (
	"fmt"
	"time"

	cache "github.com/patrickmn/go-cache"
)

// MemRepo holds a in-memory representation of a store backend
type MemRepo struct {
	c *cache.Cache
}

// NewMemRepository returns a newly initialized in-memory repository
func NewMemRepository(expiration, expiredPurge int) (*MemRepo, error) {
	return &MemRepo{
		c: cache.New(time.Duration(expiration)*time.Minute, time.Duration(expiredPurge)*time.Minute),
	}, nil
}

// Save stores a new value for a key in the k/v store
func (r *MemRepo) Save(key string, value string) error {
	r.c.Set(key, value, cache.DefaultExpiration)
	return nil
}

// Get retrieves a value from the k/v store
func (r *MemRepo) Get(key string) (string, error) {
	value, found := r.c.Get(key)
	if found {
		return value.(string), nil
	}
	return "", fmt.Errorf("no value found for key '%s'", key)
}
