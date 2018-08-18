package store

// StorageRepository defines the interface for a k/v backend
type StorageRepository interface {
	// Save stores an item in the store
	Save(string, string) error

	// Get item from store
	Get(string) (string, error)
}
