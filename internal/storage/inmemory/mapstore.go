package inmemory

import (
	"errors"
	"sync"

	"github.com/atrush/pract_01.git/internal/storage"
)

var _ storage.URLStorer = (*MapStorage)(nil)

type MapStorage struct {
	urlMap map[string]string
	mutex  *sync.RWMutex
}

func NewStorage() *MapStorage {

	return &MapStorage{
		urlMap: make(map[string]string),
		mutex:  &sync.RWMutex{},
	}
}

func (mp *MapStorage) GetURL(shortID string) (string, error) {
	if shortID == "" {
		return "", errors.New("нельзя использовать пустой id")
	}

	mp.mutex.RLock()
	longURL, ok := mp.urlMap[shortID]
	defer mp.mutex.RUnlock()

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

	mp.mutex.Lock()
	mp.urlMap[shortID] = srcURL
	defer mp.mutex.Unlock()

	return shortID, nil
}
func (mp *MapStorage) IsAvailableID(shortID string) bool {
	_, ok := mp.urlMap[shortID]

	return !ok
}
