package scraper

import (
	"fmt"
	"sync"
)

type factory struct {
	mu       sync.RWMutex
	builders map[string]Builder
}

func NewFactory() Factory {
	return &factory{
		builders: make(map[string]Builder),
	}
}

func (f *factory) Register(name string, b Builder) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.builders[name] = b
}

func (f *factory) Get(name string) (Scraper, error) {
	f.mu.RLock()
	b, ok := f.builders[name]
	f.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("aucun scraper défini pour %s", name)
	}
	return b(), nil
}
