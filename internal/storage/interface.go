package storage

import (
	"github.com/atrush/pract_01.git/internal/model"
	"github.com/google/uuid"
)

type Storage interface {
	URL() URLRepository
	User() UserRepository
	Close()
	Ping() error
}

type URLRepository interface {
	GetURL(shortID string) (model.ShortURL, error)
	GetUserURLList(userID uuid.UUID, limit int) ([]model.ShortURL, error)
	SaveURL(*model.ShortURL) error
	SaveURLBuff(*model.ShortURL) error
	SaveURLBuffFlush() error
	DeleteURLBatch(userID uuid.UUID, shortIDList ...string) error
	Exist(shortID string) (bool, error)
}

type UserRepository interface {
	AddUser(*model.User) error
	Exist(userID uuid.UUID) (bool, error)
}
