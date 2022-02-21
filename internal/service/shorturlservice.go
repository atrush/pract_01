package service

import (
	"errors"
	"fmt"

	"github.com/atrush/pract_01.git/internal/model"
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

func (sh *ShortURLService) DeleteURLList(userID uuid.UUID, shotIDList ...string) error {
	return sh.db.URL().DeleteURLBatch(userID, shotIDList...)
}

// Save map[id]URL to db, updates URL to ShortID in map
func (sh *ShortURLService) SaveURLList(src map[string]string, userID uuid.UUID) error {

	//map of new shortURL with incoming IDs
	toAdd := make(map[string]model.ShortURL, len(src))

	//map for cheking new shortID for unique
	checkShortID := make(map[string]string, len(src))

	//generate new shortURLs and send to save db
	for k, v := range src {
		sht := model.ShortURL{
			ID:        uuid.New(),
			UserID:    userID,
			URL:       v,
			IsDeleted: false,
		}

		shortID, err := sh.genShortURL(v, sht.ID, checkShortID)
		if err != nil {
			return err
		}

		sht.ShortID = shortID
		if err := sh.db.URL().SaveURLBuff(&sht); err != nil {
			return err
		}
		toAdd[k] = sht
	}

	//flush buffer
	if err := sh.db.URL().SaveURLBuffFlush(); err != nil {
		return err
	}

	//update incoming map URLs fo shortIDs
	for k, v := range toAdd {
		src[k] = v.ShortID
	}

	return nil
}

// Return array stored URLs by user UUID
func (sh *ShortURLService) GetUserURLList(userID uuid.UUID) ([]model.ShortURL, error) {
	list, err := sh.db.URL().GetUserURLList(userID, 100)
	if err != nil {
		return nil, err
	}

	return list, nil
}

// Return stored URL by shortID
func (sh *ShortURLService) GetURL(shortID string) (model.ShortURL, error) {
	longURL, err := sh.db.URL().GetURL(shortID)
	if err != nil {
		return model.ShortURL{}, err
	}

	return longURL, nil
}

// Save URL for user, return shortID
func (sh *ShortURLService) SaveURL(srcURL string, userID uuid.UUID) (string, error) {

	sht := model.ShortURL{
		ID:        uuid.New(),
		UserID:    userID,
		URL:       srcURL,
		IsDeleted: false,
	}

	var err error
	if sht.ShortID, err = sh.genShortURL(srcURL, sht.ID, nil); err != nil {
		return "", err
	}

	if err := sh.db.URL().SaveURL(&sht); err != nil {
		return "", err
	}

	return sht.ShortID, nil
}

// Check db connection
func (sh *ShortURLService) Ping() error {
	return sh.db.Ping()
}

// Generate unique ShortID, for generating multiple shortIDs use generatedCheck
func (sh *ShortURLService) genShortURL(srcURL string, id uuid.UUID, generatedCheck map[string]string) (string, error) {
	shortID, err := sh.iterShortURLGenerator(string(srcURL), 0, id.String(), generatedCheck)
	if err != nil {
		return "", err
	}

	return shortID, nil
}

// Iter gen shortID, if exist try one more time, if maxIterate times throw error
func (sh *ShortURLService) iterShortURLGenerator(srcURL string, iterationCount int, salt string, generatedCheck map[string]string) (string, error) {
	maxIterate := 10
	shortID := GenerateShortLink(srcURL, salt)

	existInCheck := false
	if generatedCheck != nil {
		_, existInCheck = generatedCheck[shortID]
	}
	exist, err := sh.db.URL().Exist(shortID)
	if err != nil {
		return "", fmt.Errorf("ошибка генерации короткой ссылки:%w", err)
	}
	if exist || existInCheck {
		iterationCount++
		if iterationCount > maxIterate {
			return "", fmt.Errorf("ошибка генерации короткой ссылки, число попыток:%v", maxIterate)
		}

		salt := uuid.New().String()

		shortID, err := sh.iterShortURLGenerator(srcURL, iterationCount, salt, generatedCheck)
		if err != nil {
			return "", err
		}

		return shortID, nil
	}

	if generatedCheck != nil {
		generatedCheck[shortID] = ""
	}

	return shortID, nil
}
