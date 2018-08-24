package store

import (
	"fmt"

	"github.com/peterbourgon/diskv"
)

// DiskRepo holds a representation of a persistent store backend
type DiskRepo struct {
	d *diskv.Diskv
}

// NewDiskRepository returns a newly persistent store repository
func NewDiskRepository(path string) (*DiskRepo, error) {
	d := diskv.New(diskv.Options{
		BasePath:     path,
		CacheSizeMax: 1024 * 1024,
	})

	return &DiskRepo{
		d: d,
	}, nil
}

// Save stores a new value for a key in the k/v store
func (r *DiskRepo) Save(key string, value string) error {
	if err := r.d.Write(key, []byte(value)); err != nil {
		return err
	}
	return nil
}

// Get retrieves a value from the k/v store
func (r *DiskRepo) Get(key string) (string, error) {
	value, err := r.d.Read(key)
	if err != nil {
		return "", fmt.Errorf("no value found for key '%s', err: %s", key, err)
	}
	return string(value), nil
}
