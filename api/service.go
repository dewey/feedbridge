package api

import (
	"github.com/dewey/feedbridge/repository"
	"github.com/dewey/feedbridge/store"
)

// Service provides access to the functions by the public API
type Service interface {
	// ServeFeed serves a feed based on the plugin and format
	ServeFeed(format string, plugin string) (string, error)

	// ListFeeds lists all available feed plugins
	ListFeeds() []string
}

type service struct {
	StorageRepository *store.MemRepo
	PluginRepository  *repository.MemRepo
}

// NewService initializes a new API service
func NewService(sr *store.MemRepo, pr *repository.MemRepo) *service {
	return &service{
		StorageRepository: sr,
		PluginRepository:  pr,
	}
}
