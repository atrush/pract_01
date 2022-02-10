package service

import (
	st "github.com/atrush/pract_01.git/internal/storage"
	"github.com/google/uuid"
)

type Servicer interface {
	URL() URLShortener
	User() UserManager
}

type URLShortener interface {
	GetURL(shortID string) (string, error)
	GetUserURLList(userID uuid.UUID) ([]st.ShortURL, error)
	SaveURL(srcURL string, userID uuid.UUID) (string, error)
}

type UserManager interface {
	AddUser() (*st.User, error)
	Exist(id uuid.UUID) bool
}
