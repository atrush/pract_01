package storage

import (
	"github.com/google/uuid"
)

type Storage interface {
	URL() URLRepository
	User() UserRepository
	Close()
	Ping() error
}

type URLRepository interface {
	GetURL(shortID string) (string, error)
	GetUserURLList(userID uuid.UUID, limit int) ([]ShortURL, error)
	SaveURL(*ShortURL) error
	Exist(shortID string) (bool, error)
}

type UserRepository interface {
	AddUser(*User) error
	Exist(userID uuid.UUID) (bool, error)
}
