package api

import "github.com/dewey/feedbridge/store"

// Service provides access to the serving functions
type Service interface {
	// ServeFile serves a feed based on the plugin
	ServeFeed(format string, plugin string) (string, error)
}

type service struct {
	Storage *store.MemRepo
}

// NewService initializes a new store service
func NewService(sr *store.MemRepo) *service {
	return &service{
		Storage: sr,
	}
}
