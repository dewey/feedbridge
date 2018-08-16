package repository

import (
	"errors"
	"fmt"
	"sync"
)

// MemRepo is a memory-backed repo
type MemRepo struct {
	plugins map[string]Plugin
	sync.RWMutex
}

// NewMemRepo creates a memory-backed repo
func NewMemRepo() *MemRepo {
	return &MemRepo{
		plugins: map[string]Plugin{},
	}
}

// Install installs a new plugin
func (mr *MemRepo) Install(p Plugin) {
	if p == nil {
		panic("no plugin func, or source is empty")
	}

	mr.Lock()
	defer mr.Unlock()

	if _, ok := mr.plugins[p.String()]; ok {
		panic(fmt.Sprintf("plugin with name '%s' already exists", p.String()))
	}

	mr.plugins[p.String()] = p
}

// Find finds a plugin by name
func (mr *MemRepo) Find(name string) (Plugin, error) {
	mr.RLock()
	defer mr.RUnlock()
	pf, ok := mr.plugins[name]
	if !ok {
		return nil, errors.New("no plugin found")
	}
	return pf, nil
}

// All returns all the available plugins
func (mr *MemRepo) All() []Plugin {
	mr.RLock()
	defer mr.RUnlock()
	plugins := []Plugin{}
	for _, p := range mr.plugins {
		plugins = append(plugins, p)
	}

	return plugins
}
