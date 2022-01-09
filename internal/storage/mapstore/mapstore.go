package mapstore

import (
	"errors"
	"strconv"

	"github.com/atrush/pract_01.git/internal/shortener"
	"github.com/atrush/pract_01.git/internal/storage"
)

var _ storage.URLStorer = (*MapStorage)(nil)

//в чем разница с  = Mapstorage{} ???
// Storage keeps storage repository dependencies.
type MapStorage struct {
	urlMap map[string]string
}

// NewStorage creates a new Storage instance.
func NewStorage() *MapStorage {
	initURLMap := make(map[string]string)
	return &MapStorage{
		urlMap: initURLMap,
	}
}

func (mp *MapStorage) GetURL(shortURL string) (string, error) {
	if shortURL == "" {
		return "", errors.New("короткая ссылка пустая")
		//почему нельзя с заглавной бууквы????
	}
	longURL, ok := mp.urlMap[shortURL]
	if ok {
		return longURL, nil
	}
	return "", nil
}
func (mp *MapStorage) SaveURL(srcURL string) (string, error) {
	if srcURL == "" {
		return "", errors.New("ссылка не может быть пустой")
		//почему нельзя с заглавной буквы????
	}
	shortURL, err := mp.genShortURL(srcURL, len(mp.urlMap))
	if err != nil {
		return "", err
	}
	mp.urlMap[shortURL] = srcURL
	return shortURL, nil
}

func (mp *MapStorage) genShortURL(srcURL string, saltCount int) (string, error) {
	shortURL := shortener.GenerateShortLink(srcURL, strconv.Itoa(saltCount))
	_, ok := mp.urlMap[shortURL]
	if ok {
		saltCount++
		shortURL, err := mp.genShortURL(srcURL, saltCount)
		if err != nil || (saltCount-len(mp.urlMap) > 10) {
			return "", errors.New("ошибка генерации короткой ссылки")
		}
		return shortURL, nil
	}
	return shortURL, nil
}
