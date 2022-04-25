package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/atrush/pract_01.git/internal/model"
	"github.com/atrush/pract_01.git/internal/storage"
	"github.com/google/uuid"
)

var _ URLShortener = (*ShortURLService)(nil)

//  ShortURLService implements URLShortener interface, provides operations with urls.
type ShortURLService struct {
	db storage.Storage
}

//  NewShortURLService inits and returns new URL service.
func NewShortURLService(db storage.Storage) (*ShortURLService, error) {
	if db == nil {
		return nil, errors.New("ошибка инициализации хранилища")
	}

	return &ShortURLService{
		db: db,
	}, nil
}

//  DeleteURLList marks list of short urls as deleted.
func (sh *ShortURLService) DeleteURLList(userID uuid.UUID, shotIDList ...string) error {
	return sh.db.URL().DeleteURLBatch(userID, shotIDList...)
}

//  SaveURLList saves map[external_id]URL to storage, updates URL to ShortID in map.
func (sh *ShortURLService) SaveURLList(src map[string]string, userID uuid.UUID) (map[string]string, error) {

	//  map of new shortURL with incoming IDs
	toAdd := make(map[string]model.ShortURL, len(src))

	//  map for cheking new shortID for unique
	checkShortID := make(map[string]string, len(src))

	//  generate new shortURLs and send to save db
	for k, v := range src {

		sht := model.NewShortURL(v, userID)

		shortID, err := sh.genShortURL(v, sht.ID, checkShortID)
		if err != nil {
			return nil, err
		}

		sht.ShortID = shortID
		if err := sh.db.URL().SaveURLBuff(&sht); err != nil {
			return nil, err
		}
		toAdd[k] = sht
	}

	//  flush buffer
	if err := sh.db.URL().SaveURLBuffFlush(); err != nil {
		return nil, err
	}

	resMap := make(map[string]string)
	for k, v := range toAdd {
		resMap[k] = v.ShortID
	}

	return resMap, nil
}

//  GetUserURLList returns array of stored urlss by user id.
func (sh *ShortURLService) GetUserURLList(ctx context.Context, userID uuid.UUID) ([]model.ShortURL, error) {
	list, err := sh.db.URL().GetUserURLList(ctx, userID, 100)
	if err != nil {
		return nil, err
	}

	return list, nil
}

//  GetURL returns stored url by shortID.
func (sh *ShortURLService) GetURL(ctx context.Context, shortID string) (model.ShortURL, error) {
	longURL, err := sh.db.URL().GetURL(ctx, shortID)
	if err != nil {
		return model.ShortURL{}, err
	}

	return longURL, nil
}

//  SaveURL saves url for user, return shortID.
func (sh *ShortURLService) SaveURL(ctx context.Context, srcURL string, userID uuid.UUID) (string, error) {

	sht := model.NewShortURL(srcURL, userID)

	var err error
	if sht.ShortID, err = sh.genShortURL(srcURL, sht.ID, nil); err != nil {
		return "", err
	}

	sht, err = sh.db.URL().SaveURL(ctx, sht)
	if err != nil {
		return "", err
	}

	return sht.ShortID, nil
}

//  Ping checks storage connection.
func (sh *ShortURLService) Ping(ctx context.Context) error {
	return sh.db.Ping()
}

//  GetCount returns count of stored, not deleted urls.
func (sh *ShortURLService) GetCount() (int, error) {
	return sh.db.URL().GetCount()
}

//  genShortURL generate unique ShortID, for generating multiple shortIDs use generatedCheck.
func (sh *ShortURLService) genShortURL(srcURL string, id uuid.UUID, generatedCheck map[string]string) (string, error) {
	shortID, err := sh.iterShortURLGenerator(string(srcURL), 0, id.String(), generatedCheck)
	if err != nil {
		return "", err
	}

	return shortID, nil
}

//  iterShortURLGenerator iter generates shortID.
//  If generated short id is exist, tryes one more time, if maxIterate times throws error.
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
