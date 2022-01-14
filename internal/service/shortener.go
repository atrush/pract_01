package service

import (
	"errors"

	"github.com/google/uuid"

	"github.com/atrush/pract_01.git/internal/storage"
)

type Shortener struct {
	repository *storage.Repository
}

func NewShortener(repository *storage.Repository) (*Shortener, error) {
	if repository == nil || repository.URLStorer == nil {
		return nil, errors.New("ошибка инициализации хранилища")
	}

	return &Shortener{
		repository: repository,
	}, nil
}

func (sh *Shortener) GetURL(shortID string) (string, error) {
	longURL, err := sh.repository.GetURL(shortID)
	if err != nil {

		return "", err
	}

	return longURL, nil
}

func (sh *Shortener) SaveURL(srcURL string) (string, error) {
	shortID, err := sh.genShortURL(string(srcURL))
	if err != nil {

		return "", err
	}

	_, err = sh.repository.SaveURL(shortID, string(srcURL))
	if err != nil {

		return "", err
	}

	return shortID, nil
}

func (sh *Shortener) genShortURL(srcURL string) (string, error) {
	shortID, err := sh.iterShortURLGenerator(string(srcURL), 0, "")
	if err != nil {

		return "", err
	}

	return shortID, nil
}

func (sh *Shortener) iterShortURLGenerator(srcURL string, iterationCount int, salt string) (string, error) {
	shortID := GenerateShortLink(srcURL, salt)
	if !sh.repository.IsAvailableID(shortID) {
		iterationCount++
		salt := uuid.New().String()

		var err error
		shortID, err = sh.iterShortURLGenerator(srcURL, iterationCount, salt)
		if err != nil || iterationCount > 10 {

			return "", errors.New("ошибка генерации короткой ссылки")
		}

		return shortID, nil
	}

	return shortID, nil
}
