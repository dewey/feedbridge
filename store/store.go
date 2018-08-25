package store

import "github.com/dewey/feedbridge/config"

// StorageRepository defines the interface for a k/v backend
type StorageRepository interface {
	// Save stores an item in the store
	Save(string, string) error

	// Get item from store
	Get(string) (string, error)
}

// NewStoreBackend is a factory to return a store implementation based on the config choosen
func NewStoreBackend(cfg config.Config) (*StorageRepository, error) {
	var storageRepo *StorageRepository
	switch cfg.StorageBackend {
	case "memory":
		memory, err := NewMemRepository(config.CacheExpiration, config.CacheExpiredPurge)
		if err != nil {
			return nil, err
		}
		storageRepo = memory
	case "persistent":
		disk, err := NewDiskRepository(config.StoragePath)
		if err != nil {
			return nil, err
		}
		storageRepo = disk
	}
	return storageRepo, nil
}
