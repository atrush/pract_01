package service

import "github.com/atrush/pract_01.git/internal/storage"

type URLShortener interface {
	GetURL(shortID string) (string, error)
	GetUserURLList(userID string) ([]storage.ShortURL, error)
	SaveURL(srcURL string, userID string) (string, error)
}

type UserManager interface {
	AddUser() (string, error)
}
