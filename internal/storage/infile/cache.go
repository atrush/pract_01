package infile

import (
	"github.com/atrush/pract_01.git/internal/storage"
	"github.com/google/uuid"
)

type cache struct {
	urlCache    map[uuid.UUID]storage.ShortURL
	shortURLidx map[string]uuid.UUID
	srcURLidx   map[string]uuid.UUID
	userCache   map[uuid.UUID]uuid.UUID
}

// Init new cahe
func newCache() *cache {
	return &cache{
		urlCache:    make(map[uuid.UUID]storage.ShortURL),
		userCache:   make(map[uuid.UUID]uuid.UUID),
		shortURLidx: make(map[string]uuid.UUID),
		srcURLidx:   make(map[string]uuid.UUID),
	}
}
