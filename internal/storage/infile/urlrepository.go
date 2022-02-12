package infile

import (
	"errors"
	"fmt"
	"sync"

	st "github.com/atrush/pract_01.git/internal/storage"
	"github.com/google/uuid"
)

var _ st.URLRepository = (*shortURLRepository)(nil)

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

func (r *shortURLRepository) SaveURLBuff(sht *st.ShortURL) error {
	return r.SaveURL(sht)
}

// Empty imitate flush
func (r *shortURLRepository) SaveURLBuffFlush() error {
	return nil
}

// Save URL
func (r *shortURLRepository) SaveURL(sht *st.ShortURL) error {
	if err := sht.Validate(); err != nil {
		return err
	}

	exist, _ := r.Exist(sht.ShortID)
	if exist {
		return errors.New("shortID уже существует")
	}

	_, userExist := r.cache.userCache[sht.UserID]
	if sht.UserID != uuid.Nil && !userExist {
		return errors.New("пользователь не найден")
	}

	if r.fileName != "" {
		r.Lock()
		if err := r.writeToFile(*sht); err != nil {
			return err
		}
		defer r.Unlock()
	}

	r.cache.urlCache[sht.ID] = *sht
	r.cache.shortURLidx[sht.ShortID] = sht.ID

	return nil
}

// Return stored URL by shortID
func (r *shortURLRepository) GetURL(shortID string) (string, error) {
	if shortID == "" {
		return "", errors.New("нельзя использовать пустой id")
	}

	r.RLock()
	idx, ok := r.cache.shortURLidx[shortID]
	if ok {

		item, ok := r.cache.urlCache[idx]
		if ok {
			return item.URL, nil
		}
		defer r.RUnlock()
	}
	return "", nil
}

// Get array of URL for user
func (r *shortURLRepository) GetUserURLList(userID uuid.UUID, limit int) ([]st.ShortURL, error) {
	if len(r.cache.urlCache) == 0 {
		return nil, nil
	}

	userURLs := make([]st.ShortURL, 0, limit)
	for _, v := range r.cache.urlCache {
		if v.UserID != uuid.Nil && v.UserID == userID {
			userURLs = append(userURLs, v)
			if len(userURLs) == limit {
				break
			}
		}
	}

	if len(userURLs) == 0 {
		return nil, nil
	}

	return userURLs, nil
}

// Check shortID not exist
func (r *shortURLRepository) Exist(shortID string) (bool, error) {
	r.RLock()
	_, ok := r.cache.shortURLidx[shortID]
	defer r.RUnlock()

	return ok, nil
}

// Write item to file
func (r *shortURLRepository) writeToFile(sht st.ShortURL) error {
	fileWriter, err := newFileWriter(r.fileName)
	if err != nil {
		return fmt.Errorf("ошибка записи в хранилище: %w", err)
	}

	defer fileWriter.Close()
	if err := fileWriter.WriteURL(sht); err != nil {
		return fmt.Errorf("ошибка записи в хранилище: %w", err)
	}

	return nil
}
