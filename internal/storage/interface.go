package storage

import (
	"context"
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
	GetURL(ctx context.Context, shortID string) (model.ShortURL, error)
	GetUserURLList(ctx context.Context, userID uuid.UUID, limit int) ([]model.ShortURL, error)
	SaveURL(ctx context.Context, shURL *model.ShortURL) error
	Exist(shortID string) (bool, error)
	SaveURLBuff(shURL *model.ShortURL) error
	SaveURLBuffFlush() error
	DeleteURLBatch(userID uuid.UUID, shortIDList ...string) error
}

type UserRepository interface {
	AddUser(ctx context.Context, user *model.User) error
	Exist(userID uuid.UUID) (bool, error)
}
