package inmemory

import (
	"errors"
	"sync"

	"github.com/atrush/pract_01.git/internal/storage"
)

var _ storage.URLStorer = (*MapStorage)(nil)
var _ storage.UserStorer = (*MapStorage)(nil)

type MapStorage struct {
	urlMap  map[string]storage.ShortURL
	userMap map[string]string
	sync.RWMutex
}

func NewStorage() *MapStorage {

	return &MapStorage{
		urlMap:  make(map[string]storage.ShortURL),
		userMap: make(map[string]string),
	}
}

func (mp *MapStorage) AddUser(userID string) error {
	if userID == "" {
		return errors.New("нельзя использовать пустой id")
	}
	if !mp.IsAvailableID(userID) {
		return errors.New("id уже существует")
	}

	mp.Lock()
	mp.userMap[userID] = userID
	defer mp.Unlock()

	return nil
}

func (mp *MapStorage) IsAvailableUserID(userID string) bool {
	mp.RLock()
	defer mp.RUnlock()

	_, ok := mp.userMap[userID]

	return !ok
}

func (mp *MapStorage) GetURL(shortID string) (string, error) {
	if shortID == "" {
		return "", errors.New("нельзя использовать пустой id")
	}

	mp.RLock()
	item, ok := mp.urlMap[shortID]
	defer mp.RUnlock()

	if ok {
		return item.URL, nil
	}

	return "", nil
}
func (mp *MapStorage) SaveURL(shortID string, srcURL string, userID string) (string, error) {
	if shortID == "" {
		return "", errors.New("нельзя использовать пустой id")
	}
	if srcURL == "" {
		return "", errors.New("нельзя сохранить пустое значение")
	}
	if !mp.IsAvailableID(shortID) {
		return "", errors.New("id уже существует")
	}
	if userID != "" && mp.IsAvailableUserID(userID) {
		return "", errors.New("пользователь не найден")
	}

	mp.Lock()
	mp.urlMap[shortID] = storage.ShortURL{
		ShortID: shortID,
		URL:     srcURL,
		UserID:  userID,
	}
	defer mp.Unlock()

	return shortID, nil
}
func (mp *MapStorage) IsAvailableID(shortID string) bool {
	mp.RLock()
	defer mp.RUnlock()

	_, ok := mp.urlMap[shortID]

	return !ok
}
