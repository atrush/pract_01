package mapstore

import (
	"errors"
	"strconv"

	"github.com/atrush/pract_01.git/internal/shortener"
	"github.com/atrush/pract_01.git/internal/storage"
)

var _ storage.UrlStorer = (*MapStorage)(nil)

//в чем разница с  = Mapstorage{} ???
// Storage keeps storage repository dependencies.
type MapStorage struct {
	urlMap map[string]string
}

// NewStorage creates a new Storage instance.
func NewStorage() *MapStorage {
	initUrlMap := make(map[string]string)
	return &MapStorage{
		urlMap: initUrlMap,
	}
}

func (mp *MapStorage) GetUrl(shortUrl string) (string, error) {
	if shortUrl == "" {
		return "", errors.New("короткая ссылка пустая")
		//почему нельзя с заглавной бууквы????
	}
	longUrl, ok := mp.urlMap[shortUrl]
	if ok {
		return longUrl, nil
	}
	return "", nil
}
func (mp *MapStorage) SaveUrl(srcUrl string) (string, error) {
	if srcUrl == "" {
		return "", errors.New("ссылка не может быть пустой")
		//почему нельзя с заглавной буквы????
	}
	shortUrl, err := mp.genShortUrl(srcUrl, len(mp.urlMap))
	if err != nil {
		return "", err
	}
	mp.urlMap[shortUrl] = srcUrl
	return shortUrl, nil
}

func (mp *MapStorage) genShortUrl(srcUrl string, saltCount int) (string, error) {
	shortUrl := shortener.GenerateShortLink(srcUrl, strconv.Itoa(saltCount))
	_, ok := mp.urlMap[shortUrl]
	if ok {
		saltCount++
		shortUrl, err := mp.genShortUrl(srcUrl, saltCount)
		if err != nil || (saltCount-len(mp.urlMap) > 10) {
			return "", errors.New("ошибка генерации короткой ссылки")
		}
		return shortUrl, nil
	}
	return shortUrl, nil
}
