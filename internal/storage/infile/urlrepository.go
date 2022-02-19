package infile

import (
	"errors"
	"fmt"
	"sync"

	"github.com/atrush/pract_01.git/internal/model"
	"github.com/atrush/pract_01.git/internal/shterrors"
	st "github.com/atrush/pract_01.git/internal/storage"
	"github.com/atrush/pract_01.git/internal/storage/schema"
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

func (r *shortURLRepository) SaveURLBuff(sht *model.ShortURL) error {
	return r.SaveURL(sht)
}

// Empty imitate flush
func (r *shortURLRepository) SaveURLBuffFlush() error {
	return nil
}

// Save URL
func (r *shortURLRepository) SaveURL(sht *model.ShortURL) error {
	dbObj, err := schema.NewURLFromCanonical(*sht)
	if err != nil {
		return fmt.Errorf("ошибка хранилица:%w", err)
	}

	exist, _ := r.Exist(dbObj.ShortID)
	if exist {
		return errors.New("shortID уже существует")
	}

	existSrcURL, _ := r.ExistSrcURL(dbObj.URL)
	if existSrcURL {

		return &shterrors.ErrorConflictSaveURL{
			Err:           errors.New("конфлит добавления записи, URL уже существует"),
			ExistShortURL: r.GetShortURLBySrcURL(dbObj.URL),
		}
	}

	_, userExist := r.cache.userCache[dbObj.UserID]
	if dbObj.UserID != uuid.Nil && !userExist {
		return errors.New("пользователь не найден")
	}

	if r.fileName != "" {
		r.Lock()
		if err := r.writeToFile(dbObj); err != nil {
			return err
		}
		defer r.Unlock()
	}

	r.cache.urlCache[dbObj.ID] = dbObj
	r.cache.shortURLidx[dbObj.ShortID] = dbObj.ID
	r.cache.srcURLidx[dbObj.URL] = dbObj.ID

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

// Return stored shortID by srcURL
func (r *shortURLRepository) GetShortURLBySrcURL(url string) string {
	r.RLock()
	id, ok := r.cache.srcURLidx[url]
	if ok {
		sht, okSht := r.cache.urlCache[id]
		if okSht {
			return sht.ShortID
		}
	}
	defer r.RUnlock()

	return ""
}

// Get array of URL for user
func (r *shortURLRepository) GetUserURLList(userID uuid.UUID, limit int) ([]model.ShortURL, error) {
	if len(r.cache.urlCache) == 0 {
		return nil, nil
	}
	var dbURLs schema.URLList
	dbURLs = make([]schema.ShortURL, 0, limit)
	for _, v := range r.cache.urlCache {
		if v.UserID != uuid.Nil && v.UserID == userID {
			dbURLs = append(dbURLs, v)
			if len(dbURLs) == limit {
				break
			}
		}
	}

	if len(dbURLs) == 0 {
		return nil, nil
	}

	return dbURLs.ToCanonical()
}

// Check shortID not exist
func (r *shortURLRepository) Exist(shortID string) (bool, error) {
	r.RLock()
	_, ok := r.cache.shortURLidx[shortID]
	defer r.RUnlock()

	return ok, nil
}

// Check shortID not exist
func (r *shortURLRepository) ExistSrcURL(url string) (bool, error) {
	r.RLock()
	_, ok := r.cache.srcURLidx[url]
	defer r.RUnlock()

	return ok, nil
}

// Write item to file
func (r *shortURLRepository) writeToFile(sht schema.ShortURL) error {
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
