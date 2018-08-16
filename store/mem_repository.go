package store

import (
	"errors"

	cache "github.com/patrickmn/go-cache"
)

type MemRepo struct {
	c *cache.Cache
}

// NewMemRepository returns a newly initialized in-memory repository
func NewMemRepository(c *cache.Cache) (*MemRepo, error) {
	return &MemRepo{
		c: c,
	}, nil
}

func (r *MemRepo) Save(key string, value string) {
	r.c.Set(key, value, cache.DefaultExpiration)
}

func (r *MemRepo) Get(key string) (string, error) {
	value, found := r.c.Get(key)
	if found {
		return value.(string), nil
	}
	return "", errors.New("no value found for this key")
}
