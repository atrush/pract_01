package infile

import (
	"sync"

	"github.com/atrush/pract_01.git/internal/storage/schema"
	"github.com/google/uuid"
)

//  cache stores records and indexes in memory.
type cache struct {
	sync.RWMutex
	urlCache    map[uuid.UUID]schema.ShortURL
	shortURLidx map[string]uuid.UUID
	srcURLidx   map[string]uuid.UUID
	userCache   map[uuid.UUID]uuid.UUID
}

//  newCache inits new cache.
func newCache() *cache {
	return &cache{
		urlCache:    make(map[uuid.UUID]schema.ShortURL),
		userCache:   make(map[uuid.UUID]uuid.UUID),
		shortURLidx: make(map[string]uuid.UUID),
		srcURLidx:   make(map[string]uuid.UUID),
	}
}
