package storage

import (
	"github.com/google/uuid"
)

type Storage interface {
	URL() URLRepository
	User() UserRepository
	Close()
}

type URLRepository interface {
	GetURL(shortID string) (string, error)
	GetUserURLList(userID uuid.UUID) ([]ShortURL, error)
	SaveURL(*ShortURL) error
	IsAvailableID(shortID string) bool
}

type UserRepository interface {
	AddUser(*User) error
	Exist(userID uuid.UUID) bool
}
