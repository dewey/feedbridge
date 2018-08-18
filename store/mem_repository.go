package store

import (
	"fmt"

	cache "github.com/patrickmn/go-cache"
)

// MemRepo holds a in-memory representation of a store backend
type MemRepo struct {
	c *cache.Cache
}

// NewMemRepository returns a newly initialized in-memory repository
func NewMemRepository(c *cache.Cache) (*MemRepo, error) {
	return &MemRepo{
		c: c,
	}, nil
}

// Save stores a new value for a key in the k/v store
func (r *MemRepo) Save(key string, value string) {
	r.c.Set(key, value, cache.DefaultExpiration)
}

// Get retrieves a value from the k/v store
func (r *MemRepo) Get(key string) (string, error) {
	value, found := r.c.Get(key)
	if found {
		return value.(string), nil
	}
	return "", fmt.Errorf("no value found for key '%s'", key)
}
