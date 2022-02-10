package service

import (
	"errors"

	"github.com/atrush/pract_01.git/internal/storage"
	"github.com/google/uuid"
)

var _ URLShortener = (*ShortURLService)(nil)

type ShortURLService struct {
	db storage.Storage
}

// Return new URL service
func NewShortURLService(db storage.Storage) (*ShortURLService, error) {
	if db == nil {
		return nil, errors.New("ошибка инициализации хранилища")
	}

	return &ShortURLService{
		db: db,
	}, nil
}

// Return array stored URLs by user UUID
func (sh *ShortURLService) GetUserURLList(userID uuid.UUID) ([]storage.ShortURL, error) {
	list, err := sh.db.URL().GetUserURLList(userID)
	if err != nil {
		return nil, err
	}

	return list, nil
}

// Return stored URL by shortID
func (sh *ShortURLService) GetURL(shortID string) (string, error) {
	longURL, err := sh.db.URL().GetURL(shortID)
	if err != nil {
		return "", err
	}

	return longURL, nil
}

// Save URL for user, return shortID
func (sh *ShortURLService) SaveURL(srcURL string, userID uuid.UUID) (string, error) {
	shortID, err := sh.genShortURL(string(srcURL))
	if err != nil {
		return "", err
	}

	sht := storage.ShortURL{
		ID:      uuid.New(),
		ShortID: shortID,
		UserID:  userID,
		URL:     srcURL,
	}

	if err = sh.db.URL().SaveURL(&sht); err != nil {
		return "", err
	}

	return sht.ShortID, nil
}

// Generate unique ShortID
func (sh *ShortURLService) genShortURL(srcURL string) (string, error) {
	shortID, err := sh.iterShortURLGenerator(string(srcURL), 0, "")
	if err != nil {
		return "", err
	}

	return shortID, nil
}

// Iter gen shortID, if exist try one more time, if maxIterate times throw error
func (sh *ShortURLService) iterShortURLGenerator(srcURL string, iterationCount int, salt string) (string, error) {
	maxIterate := 10
	shortID := GenerateShortLink(srcURL, salt)
	if !sh.db.URL().IsAvailableID(shortID) {
		iterationCount++
		salt := uuid.New().String()

		shortID, err := sh.iterShortURLGenerator(srcURL, iterationCount, salt)
		if err != nil || iterationCount > maxIterate {
			return "", errors.New("ошибка генерации короткой ссылки")
		}

		return shortID, nil
	}

	return shortID, nil
}
