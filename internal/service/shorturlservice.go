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

func NewShortURLService(db storage.Storage) (*ShortURLService, error) {
	if db == nil {
		return nil, errors.New("ошибка инициализации хранилища")
	}

	return &ShortURLService{
		db: db,
	}, nil
}

func (sh *ShortURLService) GetUserURLList(userID string) ([]storage.ShortURL, error) {
	list, err := sh.db.URL().GetUserURLList(userID)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (sh *ShortURLService) GetURL(shortID string) (string, error) {
	longURL, err := sh.db.URL().GetURL(shortID)
	if err != nil {
		return "", err
	}

	return longURL, nil
}

func (sh *ShortURLService) SaveURL(srcURL string, userID string) (string, error) {
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

func (sh *ShortURLService) genShortURL(srcURL string) (string, error) {
	shortID, err := sh.iterShortURLGenerator(string(srcURL), 0, "")
	if err != nil {
		return "", err
	}

	return shortID, nil
}

func (sh *ShortURLService) iterShortURLGenerator(srcURL string, iterationCount int, salt string) (string, error) {
	shortID := GenerateShortLink(srcURL, salt)
	if !sh.db.URL().IsAvailableID(shortID) {
		iterationCount++
		salt := uuid.New().String()

		shortID, err := sh.iterShortURLGenerator(srcURL, iterationCount, salt)
		if err != nil || iterationCount > 10 {
			return "", errors.New("ошибка генерации короткой ссылки")
		}

		return shortID, nil
	}

	return shortID, nil
}
