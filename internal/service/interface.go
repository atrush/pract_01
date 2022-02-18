package service

import (
	st "github.com/atrush/pract_01.git/internal/storage"
	"github.com/google/uuid"
)

type URLShortener interface {
	GetURL(shortID string) (string, error)
	GetUserURLList(userID uuid.UUID) ([]st.ShortURL, error)
	SaveURL(srcURL string, userID uuid.UUID) (string, error)
	SaveURLList(srcArr map[string]string, userID uuid.UUID) error
	Ping() error
}

type UserManager interface {
	AddUser() (*st.User, error)
	Exist(id uuid.UUID) (bool, error)
}
