package infile

import (
	"github.com/atrush/pract_01.git/internal/storage"
)

type cache struct {
	urlCache    map[string]storage.ShortURL
	shortURLidx map[string]string
	userCache   map[string]string
}

// Init new cahe
func newCache() *cache {
	return &cache{
		urlCache:    make(map[string]storage.ShortURL),
		userCache:   make(map[string]string),
		shortURLidx: make(map[string]string),
	}
}
