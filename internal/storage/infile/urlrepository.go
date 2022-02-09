package infile

import (
	"errors"
	"fmt"
	"sync"

	"github.com/atrush/pract_01.git/internal/storage"
)

var _ storage.URLRepository = (*shortURLRepository)(nil)

type shortURLRepository struct {
	cache    cache
	fileName string
	sync.RWMutex
}

// Init new repository
func newShortURLRepository(c *cache, fileName string) (*shortURLRepository, error) {
	if c == nil {
		return nil, errors.New("cant init repository cache not init")
	}

	return &shortURLRepository{
		cache:    *c,
		fileName: fileName,
	}, nil
}

//Save URL
func (r *shortURLRepository) SaveURL(shortID string, srcURL string, userID string) (string, error) {
	if shortID == "" {
		return "", errors.New("нельзя использовать id")
	}
	if srcURL == "" {
		return "", errors.New("нельзя сохранить пустое значение")
	}
	if !r.IsAvailableID(shortID) {
		return "", errors.New("shortID уже существует")
	}

	_, userExist := r.cache.userCache[userID]
	if userID != "" && !userExist {
		return "", errors.New("пользователь не найден")
	}

	if r.fileName != "" {
		r.Lock()
		if err := r.writeToFile(shortID, srcURL, userID); err != nil {
			return "", err
		}
		defer r.Unlock()
	}

	r.cache.urlCache[shortID] = storage.ShortURL{
		ShortID: shortID,
		URL:     srcURL,
		UserID:  userID,
	}
	r.cache.shortURLidx[shortID] = shortID

	return shortID, nil
}

func (r *shortURLRepository) GetURL(shortID string) (string, error) {
	if shortID == "" {
		return "", errors.New("нельзя использовать пустой id")
	}

	r.RLock()
	item, ok := r.cache.urlCache[shortID]
	if ok {
		return item.URL, nil
	}
	defer r.RUnlock()

	return "", nil
}

// Get array of URL for user
func (r *shortURLRepository) GetUserURLList(userID string) ([]storage.ShortURL, error) {
	if userID == "" {
		return nil, errors.New("нельзя использовать пустой id")
	}

	if len(r.cache.urlCache) == 0 {
		return nil, nil
	}

	userURLs := make([]storage.ShortURL, 0, len(r.cache.urlCache))
	for _, v := range r.cache.urlCache {
		if v.UserID == userID {
			userURLs = append(userURLs, v)
		}
	}

	if len(userURLs) == 0 {
		return nil, nil
	}

	return userURLs, nil
}

func (r *shortURLRepository) IsAvailableID(shortID string) bool {
	r.RLock()
	_, ok := r.cache.shortURLidx[shortID]
	defer r.RUnlock()

	return !ok
}

// Write item to file
func (r *shortURLRepository) writeToFile(shortID string, srcURL string, userID string) error {
	fileWriter, err := newFileWriter(r.fileName)
	if err != nil {
		return fmt.Errorf("ошибка записи в хранилище: %w", err)
	}

	defer fileWriter.Close()
	if err := fileWriter.WriteURL(shortID, srcURL, userID); err != nil {
		return fmt.Errorf("ошибка записи в хранилище: %w", err)
	}

	return nil
}
