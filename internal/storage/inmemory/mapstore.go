package inmemory

import (
	"errors"

	"github.com/atrush/pract_01.git/internal/storage"
)

var _ storage.URLStorer = (*MapStorage)(nil)

//в чем разница с  = Mapstorage{} ???
// Storage keeps storage repository dependencies.
type MapStorage struct {
	urlMap map[string]string
}

// NewStorage creates a new Storage instance.
func NewStorage() *MapStorage {
	initURLMap := make(map[string]string)
	return &MapStorage{
		urlMap: initURLMap,
	}
}

func (mp *MapStorage) GetURL(shortID string) (string, error) {
	if shortID == "" {
		return "", errors.New("нельзя использовать пустой id")
	}
	longURL, ok := mp.urlMap[shortID]
	if ok {
		return longURL, nil
	}
	return "", nil
}
func (mp *MapStorage) SaveURL(shortID string, srcURL string) (string, error) {
	if shortID == "" {
		return "", errors.New("нельзя использовать пустой id")
	}
	if srcURL == "" {
		return "", errors.New("нельзя сохранить пустое значение")
	}
	if !mp.IsAvailableID(shortID) {
		return "", errors.New("id уже существует")
	}
	mp.urlMap[shortID] = srcURL
	return shortID, nil
}
func (mp *MapStorage) IsAvailableID(shortID string) bool {
	_, ok := mp.urlMap[shortID]
	return !ok
}
