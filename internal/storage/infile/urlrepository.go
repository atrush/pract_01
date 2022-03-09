package infile

import (
	"context"
	"errors"
	"fmt"

	"github.com/atrush/pract_01.git/internal/model"
	"github.com/atrush/pract_01.git/internal/shterrors"
	st "github.com/atrush/pract_01.git/internal/storage"
	"github.com/atrush/pract_01.git/internal/storage/schema"
	"github.com/google/uuid"
)

var _ st.URLRepository = (*shortURLRepository)(nil)

type shortURLRepository struct {
	cache    *cache
	fileName string
}

// Init new repository
func newShortURLRepository(c *cache, fileName string) (*shortURLRepository, error) {
	if c == nil {
		return nil, errors.New("cant init repository cache not init")
	}

	return &shortURLRepository{
		cache:    c,
		fileName: fileName,
	}, nil
}

// Delete batch of URLs for user
func (r *shortURLRepository) DeleteURLBatch(userID uuid.UUID, shortIDList ...string) error {
	if len(shortIDList) == 0 {
		return nil
	}
	for _, v := range shortIDList {
		sht, _ := r.GetURL(context.TODO(), v)
		if sht != (model.ShortURL{}) {
			if sht.UserID == userID && !sht.IsDeleted {
				sht.IsDeleted = true
				toAdd, err := schema.NewURLFromCanonical(sht)
				if err != nil {
					return fmt.Errorf("ошибка обновления запси: %w", err)
				}

				r.cache.Lock()
				r.cache.urlCache[sht.ID] = toAdd
				r.cache.Unlock()
			}

		}
	}
	return nil
}

func (r *shortURLRepository) SaveURLBuff(sht *model.ShortURL) error {
	if sht == nil {
		return errors.New("short URL is nil")
	}
	return r.SaveURL(context.TODO(), sht)
}

// Empty imitate flush
func (r *shortURLRepository) SaveURLBuffFlush() error {
	return nil
}

// Save URL
func (r *shortURLRepository) SaveURL(_ context.Context, sht *model.ShortURL) error {
	if sht == nil {
		return errors.New("short URL is nil")
	}

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

	if userExist := r.UserExist(sht.UserID); !userExist {
		return errors.New("пользователь не найден")
	}

	r.cache.Lock()
	defer r.cache.Unlock()
	if r.fileName != "" {

		if err := r.writeToFile(dbObj); err != nil {
			return err
		}
	}

	r.cache.urlCache[dbObj.ID] = dbObj
	r.cache.shortURLidx[dbObj.ShortID] = dbObj.ID
	r.cache.srcURLidx[dbObj.URL] = dbObj.ID

	return nil
}

// Return stored URL by shortID
func (r *shortURLRepository) GetURL(_ context.Context, shortID string) (model.ShortURL, error) {
	if shortID == "" {
		return model.ShortURL{}, errors.New("нельзя использовать пустой id")
	}

	r.cache.RLock()
	idx, ok := r.cache.shortURLidx[shortID]
	defer r.cache.RUnlock()

	if ok {
		item, ok := r.cache.urlCache[idx]
		if ok {
			return item.ToCanonical()
		}
	}

	return model.ShortURL{}, nil
}

// Return stored shortID by srcURL
func (r *shortURLRepository) GetShortURLBySrcURL(url string) string {
	r.cache.RLock()
	id, ok := r.cache.srcURLidx[url]
	defer r.cache.RUnlock()

	if ok {
		sht, okSht := r.cache.urlCache[id]
		if okSht {
			return sht.ShortID
		}
	}

	return ""
}

// Get array of URL for user
func (r *shortURLRepository) GetUserURLList(_ context.Context, userID uuid.UUID, limit int) ([]model.ShortURL, error) {
	r.cache.RLock()
	defer r.cache.RUnlock()

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

// Check user exist
func (r *shortURLRepository) UserExist(userID uuid.UUID) bool {
	r.cache.RLock()
	_, userExist := r.cache.userCache[userID]
	defer r.cache.RUnlock()
	return userExist
}

// Check shortID not exist
func (r *shortURLRepository) Exist(shortID string) (bool, error) {
	r.cache.RLock()
	_, ok := r.cache.shortURLidx[shortID]
	defer r.cache.RUnlock()

	return ok, nil
}

// Check shortID not exist
func (r *shortURLRepository) ExistSrcURL(url string) (bool, error) {
	r.cache.RLock()
	_, ok := r.cache.srcURLidx[url]
	defer r.cache.RUnlock()

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
